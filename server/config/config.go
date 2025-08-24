package config

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogPath              string   `yaml:"log_path"`
	EnableFileLogging    bool     `yaml:"enable_file_logging"`
	BaseURL              string   `yaml:"base_url"`
	Host                 string   `yaml:"host"`
	Port                 int      `yaml:"port"`
	DownloadPath         string   `yaml:"downloadPath"`
	DownloaderPath       string   `yaml:"downloaderPath"`
	RequireAuth          bool     `yaml:"require_auth"`
	Username             string   `yaml:"username"`
	Password             string   `yaml:"password"`
	QueueSize            int      `yaml:"queue_size"`
	LocalDatabasePath    string   `yaml:"local_database_path"`
	SessionFilePath      string   `yaml:"session_file_path"`
	path                 string   // private
	UseOpenId            bool     `yaml:"use_openid"`
	OpenIdProviderURL    string   `yaml:"openid_provider_url"`
	OpenIdClientId       string   `yaml:"openid_client_id"`
	OpenIdClientSecret   string   `yaml:"openid_client_secret"`
	OpenIdRedirectURL    string   `yaml:"openid_redirect_url"`
	OpenIdEmailWhitelist []string `yaml:"openid_email_whitelist"`
	FrontendPath         string   `yaml:"frontend_path"`
	AutoArchive          bool     `yaml:"auto_archive"`
	Twitch               struct {
		ClientId     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	} `yaml:"twitch"`
}

var (
	instance     *Config
	instanceOnce sync.Once
)

func Instance() *Config {
	if instance == nil {
		instanceOnce.Do(func() {
			instance = &Config{}
		})
	}
	return instance
}

// Initialises the Config struct given its config file
func (c *Config) LoadFile(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}

	c.path = filename

	if err := yaml.NewDecoder(fd).Decode(c); err != nil {
		return err
	}

	return nil
}

// Path of the directory containing the config file
func (c *Config) Dir() string { return filepath.Dir(c.path) }

// Absolute path of the config file
func (c *Config) Path() string { return c.path }
