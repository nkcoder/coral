// Package config provides application configuration
package config

import (
	"fmt"
	"strings"

	"coral.daniel-guo.com/internal/email"
	"coral.daniel-guo.com/internal/logger"
	"coral.daniel-guo.com/internal/secrets"
	"github.com/spf13/viper"
)

// Init initializes the configuration with Viper
func Init() error {
	// Set defaults
	setDefaults()

	// Read environment variables
	viper.SetEnvPrefix("CORAL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read config file if present
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.coral")

	// Use config file from the flag if specified
	configFile := viper.GetString("config")
	if configFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(configFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		logger.Debug("No config file found, using defaults and environment variables")
	} else {
		logger.Info("Using config file: %s", viper.ConfigFileUsed())
	}

	return nil
}

// LoadConfig loads the configuration from Viper into an AppConfig struct
func LoadConfig() *AppConfig {
	cfg := &AppConfig{
		Environment:    viper.GetString("env"),
		DefaultSender:  viper.GetString("email.sender"),
		TestEmail:      viper.GetString("email.test"),
		WorkerPoolSize: viper.GetInt("worker.pool_size"),
		WorkerDelayMs:  viper.GetInt("worker.delay_ms"),
		Email: email.Config{
			Region: viper.GetString("aws.region"),
		},
		Secrets: secrets.Config{
			Region: viper.GetString("aws.region"),
		},
	}

	return cfg
}

// setDefaults sets default values for configuration
func setDefaults() {
	// General defaults
	viper.SetDefault("env", "dev")

	// Email defaults
	viper.SetDefault("email.sender", "no-reply@the-hub.ai")
	viper.SetDefault("email.test", "")

	// AWS region
	viper.SetDefault("aws.region", "ap-southeast-2")

	// Worker pool defaults
	viper.SetDefault("worker.pool_size", 5)
	viper.SetDefault("worker.delay_ms", 1000)
}
