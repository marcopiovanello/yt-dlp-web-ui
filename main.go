package main

import (
	"context"
	"embed"
	"flag"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/openid"

	"github.com/spf13/viper"
)

//go:embed frontend/dist/index.html
//go:embed frontend/dist/assets/*
var frontend embed.FS

//go:embed openapi/*
var swagger embed.FS

func main() {
	// Parse optional config path from flag
	var configFile string
	flag.StringVar(&configFile, "conf", "./config.yml", "Config file path")
	flag.Parse()

	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")

	// Defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 3033)
	v.SetDefault("server.queue_size", 2)
	v.SetDefault("paths.download_path", ".")
	v.SetDefault("paths.downloader_path", "yt-dlp")
	v.SetDefault("paths.local_database_path", ".")
	v.SetDefault("logging.log_path", "yt-dlp-webui.log")
	v.SetDefault("logging.enable_file_logging", false)
	v.SetDefault("authentication.require_auth", false)

	// Env binding
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	// Load YAML file if exists
	if err := v.ReadInConfig(); err != nil {
		slog.Debug("using defaults")
	}

	cfg := config.Instance()
	if err := v.Unmarshal(&cfg); err != nil {
		slog.Error("failed to load config", "error", err)
	}

	if cfg.Server.QueueSize <= 0 || runtime.NumCPU() <= 2 {
		cfg.Server.QueueSize = 2
	}

	// 6. Frontend FS
	var appFS fs.FS
	if fp := v.GetString("frontend_path"); fp != "" {
		appFS = os.DirFS(fp)
	} else {
		sub, err := fs.Sub(frontend, "frontend/dist")
		if err != nil {
			slog.Error("failed to load embedded frontend", "error", err)
			os.Exit(1)
		}
		appFS = sub
	}

	// Configure OpenID if needed
	openid.Configure()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("starting server",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
		"queue_size", cfg.Server.QueueSize,
	)

	if err := server.Run(ctx, &server.RunConfig{
		App:     appFS,
		Swagger: swagger,
	}); err != nil {
		slog.Error("server stopped with error", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited cleanly")
}
