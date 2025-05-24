package secrets

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SecretsManagerAPI defines the interface for AWS Secrets Manager operations
// This interface allows for mocking in tests
type SecretsManagerAPI interface {
	GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

// SessionCreator defines the interface for creating AWS sessions
// This interface allows for mocking session creation in tests
type SessionCreator interface {
	NewSession(config *aws.Config) (*session.Session, error)
}

// DefaultSessionCreator implements SessionCreator using the real AWS SDK
type DefaultSessionCreator struct{}

// NewSession creates a new AWS session using the real AWS SDK
func (d *DefaultSessionCreator) NewSession(config *aws.Config) (*session.Session, error) {
	return session.NewSession(config)
}
