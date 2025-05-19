package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// GetSecret gets a secret from AWS Secrets Manager
func GetSecret(secretName string) (string, error) {
	fmt.Printf("Getting secret: %s\n", secretName)

	// Create a new AWS session with the default configuration
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	// Create a new Secrets Manager client
	svc := secretsmanager.New(session)

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
