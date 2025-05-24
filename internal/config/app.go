// Package config provides application configuration
package config

import (
	"coral.daniel-guo.com/internal/email"
	"coral.daniel-guo.com/internal/secrets"
)

// AppConfig holds the application configuration
type AppConfig struct {
	// Environment (dev, staging, prod)
	Environment string

	// Email sender configuration
	Email email.Config

	// Secrets manager configuration
	Secrets secrets.Config

	// Default sender email address
	DefaultSender string

	// Test email (if set, all emails go here)
	TestEmail string

	// Worker pool configuration
	WorkerPoolSize int
	WorkerDelayMs  int
}

// NewAppConfig creates a new application configuration with default values
func NewAppConfig(environment string, testEmail string, sender string) *AppConfig {
	cfg := &AppConfig{
		Environment:    environment,
		Email:          email.DefaultConfig(),
		Secrets:        secrets.DefaultConfig(),
		DefaultSender:  "no-reply@the-hub.ai",
		TestEmail:      "",
		WorkerPoolSize: 5,
		WorkerDelayMs:  1000,
	}
	if testEmail != "" {
		cfg.TestEmail = testEmail
	}
	if sender != "" {
		cfg.DefaultSender = sender
	}
	return cfg
}
