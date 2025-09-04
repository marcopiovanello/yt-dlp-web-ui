package config

import (
	"path/filepath"
	"sync"
	"time"
)

type Config struct {
	Server         ServerConfig   `yaml:"server"`
	Logging        LoggingConfig  `yaml:"logging"`
	Paths          PathsConfig    `yaml:"paths"`
	Authentication AuthConfig     `yaml:"authentication"`
	OpenId         OpenIdConfig   `yaml:"openid"`
	Frontend       FrontendConfig `yaml:"frontend"`
	AutoArchive    bool           `yaml:"auto_archive"`
	Twitch         TwitchConfig   `yaml:"twitch"`
	path           string
}

type ServerConfig struct {
	BaseURL   string `yaml:"base_url"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	QueueSize int    `yaml:"queue_size"`
}

type LoggingConfig struct {
	LogPath           string `yaml:"log_path"`
	EnableFileLogging bool   `yaml:"enable_file_logging"`
}

type PathsConfig struct {
	DownloadPath      string `yaml:"download_path"`
	DownloaderPath    string `yaml:"downloader_path"`
	LocalDatabasePath string `yaml:"local_database_path"`
}

type AuthConfig struct {
	RequireAuth  bool   `yaml:"require_auth"`
	Username     string `yaml:"username"`
	PasswordHash string `yaml:"password"`
}

type OpenIdConfig struct {
	UseOpenId      bool     `yaml:"use_openid"`
	ProviderURL    string   `yaml:"openid_provider_url"`
	ClientId       string   `yaml:"openid_client_id"`
	ClientSecret   string   `yaml:"openid_client_secret"`
	RedirectURL    string   `yaml:"openid_redirect_url"`
	EmailWhitelist []string `yaml:"openid_email_whitelist"`
}

type FrontendConfig struct {
	FrontendPath string `yaml:"frontend_path"`
}

type TwitchConfig struct {
	ClientId      string        `yaml:"client_id"`
	ClientSecret  string        `yaml:"client_secret"`
	CheckInterval time.Duration `yaml:"check_interval"`
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
