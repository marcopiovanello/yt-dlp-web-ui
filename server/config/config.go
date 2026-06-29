package config

import (
	"path/filepath"
	"sync"
	"time"
)

type Config struct {
	Server         ServerConfig   `mapstructure:"server"`
	Logging        LoggingConfig  `mapstructure:"logging"`
	Paths          PathsConfig    `mapstructure:"paths"`
	Authentication AuthConfig     `mapstructure:"authentication"`
	OpenId         OpenIdConfig   `mapstructure:"openid"`
	Frontend       FrontendConfig `mapstructure:"frontend"`
	AutoArchive    bool           `mapstructure:"auto_archive"`
	Twitch         TwitchConfig   `mapstructure:"twitch"`
	path           string
}

type ServerConfig struct {
	BaseURL   string `mapstructure:"base_url"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	QueueSize int    `mapstructure:"queue_size"`
}

type LoggingConfig struct {
	LogPath           string `mapstructure:"log_path"`
	EnableFileLogging bool   `mapstructure:"enable_file_logging"`
}

type PathsConfig struct {
	DownloadPath      string `mapstructure:"download_path"`
	DownloaderPath    string `mapstructure:"downloader_path"`
	LocalDatabasePath string `mapstructure:"local_database_path"`
	JSRuntimePath     string `mapstructure:"js_runtime_path"`
}

type AuthConfig struct {
	RequireAuth  bool   `mapstructure:"require_auth"`
	Username     string `mapstructure:"username"`
	PasswordHash string `mapstructure:"password_hash"`
}

type OpenIdConfig struct {
	UseOpenId      bool     `mapstructure:"use_openid"`
	ProviderURL    string   `mapstructure:"provider_url"`
	ClientId       string   `mapstructure:"client_id"`
	ClientSecret   string   `mapstructure:"client_secret"`
	RedirectURL    string   `mapstructure:"redirect_url"`
	EmailWhitelist []string `mapstructure:"email_whitelist"`
}

type FrontendConfig struct {
	FrontendPath string `mapstructure:"frontend_path"`
}

type TwitchConfig struct {
	ClientId      string        `mapstructure:"client_id"`
	ClientSecret  string        `mapstructure:"client_secret"`
	CheckInterval time.Duration `mapstructure:"check_interval"`
}

var (
	instance     *Config
	instanceOnce sync.Once
)

func Instance() *Config {
	if instance == nil {
		instanceOnce.Do(func() {
			instance = &Config{}
			instance.Twitch.CheckInterval = time.Minute * 5
		})
	}
	return instance
}

// Path of the directory containing the config file
func (c *Config) Dir() string { return filepath.Dir(c.path) }

// Absolute path of the config file
func (c *Config) Path() string { return c.path }
