package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	CompanyID    string `mapstructure:"company_id"`
	OutputFormat string `mapstructure:"output_format"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectPort int    `mapstructure:"redirect_port"`
}

// DefaultConfigDir returns ~/.config/bokio-cli.
func DefaultConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "bokio-cli")
}

// Load reads the config file and environment variables.
func Load(configDir string) (*Config, error) {
	if configDir == "" {
		configDir = DefaultConfigDir()
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	viper.SetDefault("output_format", "table")
	viper.SetDefault("redirect_port", 8585)

	viper.SetEnvPrefix("BOKIO")
	viper.AutomaticEnv()

	// Bind specific env vars
	viper.BindEnv("company_id", "BOKIO_COMPANY_ID")
	viper.BindEnv("output_format", "BOKIO_OUTPUT")
	viper.BindEnv("redirect_port", "BOKIO_REDIRECT_PORT")
	viper.BindEnv("client_id", "BOKIO_CLIENT_ID")
	viper.BindEnv("client_secret", "BOKIO_CLIENT_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the config to disk.
func Save(configDir string, cfg *Config) error {
	if configDir == "" {
		configDir = DefaultConfigDir()
	}
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	viper.Set("company_id", cfg.CompanyID)
	viper.Set("output_format", cfg.OutputFormat)
	viper.Set("client_id", cfg.ClientID)
	viper.Set("client_secret", cfg.ClientSecret)
	viper.Set("redirect_port", cfg.RedirectPort)

	return viper.WriteConfigAs(filepath.Join(configDir, "config.yaml"))
}
