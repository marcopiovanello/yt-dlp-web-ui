// a stupid package name...
package server

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/filebrowser"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/livestream"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/pipeline"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/logging"
	middlewares "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/middleware"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/openid"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/rest"
	ytdlpRPC "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/rpc"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/status"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/task"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/twitch"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/user"

	bolt "go.etcd.io/bbolt"
)

type RunConfig struct {
	App     fs.FS
	Swagger fs.FS
}

type serverConfig struct {
	frontend      fs.FS
	swagger       fs.FS
	mdb           *kv.Store
	db            *bolt.DB
	mq            *queue.MessageQueue
	lm            *livestream.Monitor
	taskRunner    task.TaskRunner
	twitchMonitor *twitch.Monitor
}

// TODO: change scope
var observableLogger = logging.NewObservableLogger()

func Run(ctx context.Context, rc *RunConfig) error {
	dbPath := filepath.Join(config.Instance().Paths.LocalDatabasePath, "bolt.db")

	boltdb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return err
	}

	mdb, err := kv.NewStore(boltdb, time.Second*15)
	if err != nil {
		return err
	}

	// ---- LOGGING ---------------------------------------------------
	logWriters := []io.Writer{
		os.Stdout,
		observableLogger, // for web-ui
	}

	conf := config.Instance()

	// file based logging
	if conf.Logging.EnableFileLogging {
		logger, err := logging.NewRotableLogger(conf.Logging.LogPath)
		if err != nil {
			return err
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

	mq, err := queue.NewMessageQueue()
	if err != nil {
		return err
	}
	mq.SetupConsumers()
	go mdb.Restore(mq)
	go mdb.EventListener()

	lm := livestream.NewMonitor(mq, mdb, boltdb)
	go lm.Schedule()
	go lm.Restore()

	tm := twitch.NewMonitor(
		twitch.NewAuthenticationManager(
			config.Instance().Twitch.ClientId,
			config.Instance().Twitch.ClientSecret,
		),
		boltdb,
	)
	go tm.Monitor(
		ctx,
		config.Instance().Twitch.CheckInterval,
		twitch.DEFAULT_DOWNLOAD_HANDLER(mdb, mq),
	)
	go tm.Restore()

	cronTaskRunner := task.NewCronTaskRunner(mq, mdb)
	go cronTaskRunner.Spawner(ctx)

	scfg := serverConfig{
		frontend:      rc.App,
		swagger:       rc.Swagger,
		mdb:           mdb,
		db:            boltdb,
		mq:            mq,
		lm:            lm,
		twitchMonitor: tm,
		taskRunner:    cronTaskRunner,
	}

	srv := newServer(scfg)

	go gracefulShutdown(ctx, srv, &scfg)

	var (
		network = "tcp"
		address = fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	)

	// support unix sockets
	if strings.HasPrefix(conf.Server.Host, "/") {
		network = "unix"
		address = conf.Server.Host
	}

	listener, err := net.Listen(network, address)
	if err != nil {
		slog.Error("failed to listen", slog.String("err", err.Error()))
		return err
	}

	slog.Info("yt-dlp-webui started", slog.String("address", address))

	if err := srv.Serve(listener); err != nil {
		slog.Warn("http server stopped", slog.String("err", err.Error()))
	}

	return nil
}

func newServer(c serverConfig) *http.Server {
	// archiver.Register(c.db)
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

	baseUrl := config.Instance().Server.BaseURL
	r.Mount(baseUrl+"/", http.StripPrefix(baseUrl, http.FileServerFS(c.frontend)))

	// swagger
	r.Mount("/openapi", http.FileServerFS(c.swagger))

	// Filebrowser routes
	r.Route("/filebrowser", func(r chi.Router) {
		r.Use(middlewares.ApplyAuthenticationByConfig)
		r.Post("/downloaded", filebrowser.ListDownloaded)
		r.Post("/delete", filebrowser.DeleteFile)
		r.Get("/d/{id}", filebrowser.DownloadFile)
		r.Get("/v/{id}", filebrowser.SendFile)
		r.Get("/bulk", filebrowser.BulkDownload(c.mdb))
	})

	// Archive routes
	// r.Route("/archive", archive.ApplyRouter(c.db))

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
		LM:  c.lm,
	}))

	// Logging
	r.Route("/log", logging.ApplyRouter(observableLogger))

	// Status
	r.Route("/status", status.ApplyRouter(c.mdb))

	// Subscriptions
	r.Route("/subscriptions", subscription.Container(c.db, c.taskRunner).ApplyRouter())

	// Twitch
	r.Route("/twitch", func(r chi.Router) {
		r.Use(middlewares.ApplyAuthenticationByConfig)
		r.Get("/users", twitch.GetMonitoredUsers(c.twitchMonitor))
		r.Post("/user", twitch.MonitorUserHandler(c.twitchMonitor))
		r.Delete("/user/{user}", twitch.DeleteUser(c.twitchMonitor))
	})

	// Pipelines
	r.Route("/pipelines", func(r chi.Router) {
		h := pipeline.NewRestHandler(c.db)
		r.Use(middlewares.ApplyAuthenticationByConfig)
		r.Get("/id/{id}", h.GetPipeline)
		r.Get("/all", h.GetAllPipelines)
		r.Post("/", h.SavePipeline)
		r.Delete("/id/{id}", h.DeletePipeline)
	})

	return &http.Server{Handler: r}
}

func gracefulShutdown(ctx context.Context, srv *http.Server, cfg *serverConfig) {
	<-ctx.Done()
	slog.Info("shutdown signal received")

	defer func() {
		cfg.db.Close()
		srv.Shutdown(context.Background())
	}()
}
