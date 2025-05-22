// Package secrets provides access to secure configuration values
package secrets

import (
	"fmt"

	"coral.daniel-guo.com/internal/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// Config contains AWS Secrets Manager configuration
type Config struct {
	Region string
}

// DefaultConfig returns a default secrets configuration
func DefaultConfig() Config {
	return Config{
		Region: "ap-southeast-2",
	}
}

// Manager handles retrieving secrets from AWS Secrets Manager
type Manager struct {
	config Config
}

// NewManager creates a new secrets manager with the given configuration
func NewManager(config Config) *Manager {
	return &Manager{
		config: config,
	}
}

// GetSecret gets a secret from AWS Secrets Manager
func (m *Manager) GetSecret(secretName string) (string, error) {
	logger.Info("Getting secret: %s", secretName)

	// Create a new AWS session with the configuration
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(m.config.Region),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	// Create a new Secrets Manager client
	svc := secretsmanager.New(sess)

	// Create a request to get the secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	// Get the secret value
	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	if result.SecretString == nil {
		return "", fmt.Errorf("secret value is nil")
	}

	secretString := *result.SecretString

	return secretString, nil
}
