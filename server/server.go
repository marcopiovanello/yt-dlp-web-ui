// a stupid package name...
package server

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archiver"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/dbutil"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/filebrowser"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/livestream"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/logging"
	middlewares "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/middleware"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/openid"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/rest"
	ytdlpRPC "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/rpc"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/status"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/task"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/user"

	_ "modernc.org/sqlite"
)

type RunConfig struct {
	App     fs.FS
	Swagger fs.FS
}

type serverConfig struct {
	frontend fs.FS
	swagger  fs.FS
	mdb      *internal.MemoryDB
	db       *sql.DB
	mq       *internal.MessageQueue
	lm       *livestream.Monitor
}

// TODO: change scope
var observableLogger = logging.NewObservableLogger()

func RunBlocking(rc *RunConfig) {
	mdb := internal.NewMemoryDB()

	// ---- LOGGING ---------------------------------------------------
	logWriters := []io.Writer{
		os.Stdout,
		observableLogger, // for web-ui
	}

	conf := config.Instance()

	// file based logging
	if conf.EnableFileLogging {
		logger, err := logging.NewRotableLogger(conf.LogPath)
		if err != nil {
			panic(err)
		}

		defer logger.Rotate()

		go func() {
			for {
				time.Sleep(time.Hour * 24)
				logger.Rotate()
			}
		}()

		logWriters = append(logWriters, logger)
	}

	logger := slog.New(slog.NewTextHandler(io.MultiWriter(logWriters...), &slog.HandlerOptions{
		Level: slog.LevelInfo, // TODO: detect when launched in debug mode -> slog.LevelDebug
	}))

	// make the new logger the default one with all the new writers
	slog.SetDefault(logger)
	// ----------------------------------------------------------------

	db, err := sql.Open("sqlite", conf.LocalDatabasePath)
	if err != nil {
		slog.Error("failed to open database", slog.String("err", err.Error()))
	}

	if err := dbutil.Migrate(context.Background(), db); err != nil {
		slog.Error("failed to init database", slog.String("err", err.Error()))
	}

	mq, err := internal.NewMessageQueue()
	if err != nil {
		panic(err)
	}
	mq.SetupConsumers()
	go mdb.Restore(mq)
	go mdb.EventListener()

	lm := livestream.NewMonitor(mq, mdb)
	go lm.Schedule()
	go lm.Restore()

	srv := newServer(serverConfig{
		frontend: rc.App,
		swagger:  rc.Swagger,
		mdb:      mdb,
		mq:       mq,
		db:       db,
		lm:       lm,
	})

	go gracefulShutdown(srv, mdb)
	go autoPersist(time.Minute*5, mdb, lm)

	var (
		network = "tcp"
		address = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	)

	// support unix sockets
	if strings.HasPrefix(conf.Host, "/") {
		network = "unix"
		address = conf.Host
	}

	listener, err := net.Listen(network, address)
	if err != nil {
		slog.Error("failed to listen", slog.String("err", err.Error()))
		return
	}

	slog.Info("yt-dlp-webui started", slog.String("address", address))

	if err := srv.Serve(listener); err != nil {
		slog.Warn("http server stopped", slog.String("err", err.Error()))
	}
}

func newServer(c serverConfig) *http.Server {
	archiver.Register(c.db)

	cronTaskRunner := task.NewCronTaskRunner(c.mq, c.mdb)
	go cronTaskRunner.Spawner(context.TODO())

	service := ytdlpRPC.Container(c.mdb, c.mq, c.lm)
	rpc.Register(service)

	r := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	r.Use(corsMiddleware.Handler)
	// use in dev
	// r.Use(middleware.Logger)

	baseUrl := config.Instance().BaseURL
	r.Mount(baseUrl+"/", http.StripPrefix(baseUrl, http.FileServerFS(c.frontend)))

	// swagger
	r.Mount("/openapi", http.FileServerFS(c.swagger))

	// Filebrowser routes
	r.Route("/filebrowser", func(r chi.Router) {
		if config.Instance().RequireAuth {
			r.Use(middlewares.Authenticated)
		}
		if config.Instance().UseOpenId {
			r.Use(openid.Middleware)
		}
		r.Post("/downloaded", filebrowser.ListDownloaded)
		r.Post("/delete", filebrowser.DeleteFile)
		r.Get("/d/{id}", filebrowser.DownloadFile)
		r.Get("/v/{id}", filebrowser.SendFile)
		r.Get("/bulk", filebrowser.BulkDownload(c.mdb))
	})


	// Filebrowser routes
	r.Route("/public", func(r chi.Router) {
		r.Get("/{id}", filebrowser.SendFile)
	})

	// Archive routes
	r.Route("/archive", archive.ApplyRouter(c.db))

	// Authentication routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", user.Login)
		r.Get("/logout", user.Logout)

		r.Route("/openid", func(r chi.Router) {
			r.Get("/login", openid.Login)
			r.Get("/signin", openid.SingIn)
			r.Get("/logout", openid.Logout)
		})
	})

	// RPC handlers
	r.Route("/rpc", ytdlpRPC.ApplyRouter())

	// REST API handlers
	r.Route("/api/v1", rest.ApplyRouter(&rest.ContainerArgs{
		DB:  c.db,
		MDB: c.mdb,
		MQ:  c.mq,
	}))

	// Logging
	r.Route("/log", logging.ApplyRouter(observableLogger))

	// Status
	r.Route("/status", status.ApplyRouter(c.mdb))

	// Subscriptions
	r.Route("/subscriptions", subscription.Container(c.db, cronTaskRunner).ApplyRouter())

	return &http.Server{Handler: r}
}

func gracefulShutdown(srv *http.Server, db *internal.MemoryDB) {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()
		slog.Info("shutdown signal received")

		defer func() {
			db.Persist()

			stop()
			srv.Shutdown(context.Background())
		}()
	}()
}

func autoPersist(d time.Duration, db *internal.MemoryDB, lm *livestream.Monitor) {
	for {
		if err := db.Persist(); err != nil {
			slog.Warn("failed to persisted session", slog.Any("err", err))
		}
		if err := lm.Persist(); err != nil {
			slog.Warn(
				"failed to persisted livestreams monitor session", slog.Any("err", err.Error()))
		}
		slog.Debug("sucessfully persisted session")
		time.Sleep(d)
	}
}
