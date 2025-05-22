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
}

// NewAppConfig creates a new application configuration with default values
func NewAppConfig(environment string) *AppConfig {
	return &AppConfig{
		Environment:   environment,
		Email:         email.DefaultConfig(),
		Secrets:       secrets.DefaultConfig(),
		DefaultSender: "no-reply@the-hub.ai",
	}
}

// WithTestEmail sets a test email address
func (c *AppConfig) WithTestEmail(email string) *AppConfig {
	c.TestEmail = email
	return c
}

// WithSender sets the default sender email
func (c *AppConfig) WithSender(sender string) *AppConfig {
	c.DefaultSender = sender
	return c
}
