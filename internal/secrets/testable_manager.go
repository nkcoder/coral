package secrets

import (
	"fmt"

	"coral.daniel-guo.com/internal/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// TestableManager is a version of Manager that accepts dependencies for easier testing
type TestableManager struct {
	config         Config
	sessionCreator SessionCreator
	clientFactory  func(*aws.Config) SecretsManagerAPI
}

// NewTestableManager creates a new testable secrets manager with dependency injection
func NewTestableManager(
	config Config,
	sessionCreator SessionCreator,
	clientFactory func(*aws.Config) SecretsManagerAPI,
) *TestableManager {
	if sessionCreator == nil {
		sessionCreator = &DefaultSessionCreator{}
	}
	if clientFactory == nil {
		clientFactory = func(cfg *aws.Config) SecretsManagerAPI {
			sess, _ := sessionCreator.NewSession(cfg)
			return secretsmanager.New(sess)
		}
	}

	return &TestableManager{
		config:         config,
		sessionCreator: sessionCreator,
		clientFactory:  clientFactory,
	}
}

// NewTestableManagerWithDefaults creates a testable manager with default dependencies
func NewTestableManagerWithDefaults(config Config) *TestableManager {
	return NewTestableManager(config, nil, nil)
}

// GetSecret gets a secret from AWS Secrets Manager using the injected dependencies
func (m *TestableManager) GetSecret(secretName string) (string, error) {
	logger.Info("Getting secret: %s", secretName)

	// Create AWS config
	awsConfig := &aws.Config{
		Region: aws.String(m.config.Region),
	}

	// Create client using the factory
	client := m.clientFactory(awsConfig)

	// Create a request to get the secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	// Get the secret value
	result, err := client.GetSecretValue(input)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	if result.SecretString == nil {
		return "", fmt.Errorf("secret value is nil")
	}

	secretString := *result.SecretString

	return secretString, nil
}
