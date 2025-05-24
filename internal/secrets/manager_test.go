package secrets

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	expectedRegion := "ap-southeast-2"
	if config.Region != expectedRegion {
		t.Errorf("DefaultConfig().Region = %v, want %v", config.Region, expectedRegion)
	}
}

func TestConfig_StructFields(t *testing.T) {
	// Test that all struct fields are accessible and have correct types
	config := Config{
		Region: "us-east-1",
	}

	// Verify field types and accessibility
	var _ string = config.Region

	// Verify field can be set
	config.Region = "us-west-2"
	if config.Region != "us-west-2" {
		t.Errorf("Config.Region = %v, want %v", config.Region, "us-west-2")
	}
}

func TestNewManager(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected Config
	}{
		{
			name: "default config",
			config: Config{
				Region: "ap-southeast-2",
			},
			expected: Config{
				Region: "ap-southeast-2",
			},
		},
		{
			name: "custom region",
			config: Config{
				Region: "us-east-1",
			},
			expected: Config{
				Region: "us-east-1",
			},
		},
		{
			name: "empty region",
			config: Config{
				Region: "",
			},
			expected: Config{
				Region: "",
			},
		},
		{
			name: "different region",
			config: Config{
				Region: "eu-west-1",
			},
			expected: Config{
				Region: "eu-west-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager(tt.config)

			// Verify manager is not nil
			if manager == nil {
				t.Fatal("NewManager returned nil")
			}

			// Verify config is set correctly
			if manager.config.Region != tt.expected.Region {
				t.Errorf("Manager.config.Region = %v, want %v", manager.config.Region, tt.expected.Region)
			}
		})
	}
}

func TestNewManager_WithDefaultConfig(t *testing.T) {
	// Test using DefaultConfig
	defaultConfig := DefaultConfig()
	manager := NewManager(defaultConfig)

	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if manager.config.Region != "ap-southeast-2" {
		t.Errorf("Manager.config.Region = %v, want %v", manager.config.Region, "ap-southeast-2")
	}
}

func TestManager_StructFields(t *testing.T) {
	// Test that Manager struct fields are accessible
	config := Config{Region: "test-region"}
	manager := NewManager(config)

	// Verify field types and accessibility
	var _ Config = manager.config

	// Verify the config is stored correctly
	if manager.config.Region != "test-region" {
		t.Errorf("Manager.config.Region = %v, want %v", manager.config.Region, "test-region")
	}
}

// Note: GetSecret method makes actual AWS API calls, so these tests would require:
// 1. Mock AWS SDK (complex setup)
// 2. Integration test environment with actual AWS credentials
// 3. Test doubles/stubs

// For demonstration, here are some basic unit tests that test the input validation
// and structure without making AWS calls:

// Benchmark tests
func BenchmarkDefaultConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DefaultConfig()
	}
}

func BenchmarkNewManager(b *testing.B) {
	config := DefaultConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		NewManager(config)
	}
}

// Test mocks for TestableManager
type mockSecretsManagerClient struct {
	GetSecretValueFunc func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

func (m *mockSecretsManagerClient) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return m.GetSecretValueFunc(input)
}

type mockSessionCreator struct {
	NewSessionFunc func(config *aws.Config) (*session.Session, error)
}

func (m *mockSessionCreator) NewSession(config *aws.Config) (*session.Session, error) {
	return m.NewSessionFunc(config)
}

func TestNewTestableManager(t *testing.T) {
	config := Config{Region: "us-east-1"}

	tests := []struct {
		name           string
		sessionCreator SessionCreator
		clientFactory  func(*aws.Config) SecretsManagerAPI
		expectDefaults bool
	}{
		{
			name:           "with nil dependencies - should use defaults",
			sessionCreator: nil,
			clientFactory:  nil,
			expectDefaults: true,
		},
		{
			name: "with custom session creator",
			sessionCreator: &mockSessionCreator{
				NewSessionFunc: func(config *aws.Config) (*session.Session, error) {
					return nil, nil
				},
			},
			clientFactory:  nil,
			expectDefaults: false,
		},
		{
			name:           "with custom client factory",
			sessionCreator: nil,
			clientFactory: func(config *aws.Config) SecretsManagerAPI {
				return &mockSecretsManagerClient{}
			},
			expectDefaults: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewTestableManager(config, tt.sessionCreator, tt.clientFactory)

			if manager == nil {
				t.Fatal("NewTestableManager returned nil")
			}

			if manager.config.Region != config.Region {
				t.Errorf("Manager.config.Region = %v, want %v", manager.config.Region, config.Region)
			}

			// Test that dependencies are set (even if default)
			if manager.sessionCreator == nil {
				t.Error("sessionCreator should not be nil")
			}
			if manager.clientFactory == nil {
				t.Error("clientFactory should not be nil")
			}
		})
	}
}

func TestNewTestableManagerWithDefaults(t *testing.T) {
	config := Config{Region: "ap-southeast-2"}
	manager := NewTestableManagerWithDefaults(config)

	if manager == nil {
		t.Fatal("NewTestableManagerWithDefaults returned nil")
	}

	if manager.config.Region != config.Region {
		t.Errorf("Manager.config.Region = %v, want %v", manager.config.Region, config.Region)
	}
}

func TestTestableManager_GetSecret(t *testing.T) {
	tests := []struct {
		name           string
		secretName     string
		mockResponse   *secretsmanager.GetSecretValueOutput
		mockError      error
		expectedValue  string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:       "successful retrieval",
			secretName: "test-secret",
			mockResponse: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String("secret-value"),
			},
			expectedValue: "secret-value",
			expectError:   false,
		},
		{
			name:           "api error",
			secretName:     "missing-secret",
			mockError:      fmt.Errorf("ResourceNotFoundException: Secret not found"),
			expectError:    true,
			expectedErrMsg: "failed to get secret",
		},
		{
			name:       "nil secret string",
			secretName: "nil-secret",
			mockResponse: &secretsmanager.GetSecretValueOutput{
				SecretString: nil,
			},
			expectError:    true,
			expectedErrMsg: "secret value is nil",
		},
		{
			name:       "empty secret name",
			secretName: "",
			mockResponse: &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(""),
			},
			expectedValue: "",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockSecretsManagerClient{
				GetSecretValueFunc: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
					// Verify the input
					if *input.SecretId != tt.secretName {
						t.Errorf("GetSecretValue called with wrong secret name: got %v, want %v", *input.SecretId, tt.secretName)
					}
					return tt.mockResponse, tt.mockError
				},
			}

			clientFactory := func(config *aws.Config) SecretsManagerAPI {
				// Verify the config is passed correctly
				if config.Region == nil || *config.Region == "" {
					t.Error("AWS config should have a region set")
				}
				return mockClient
			}

			config := Config{Region: "us-east-1"}
			manager := NewTestableManager(config, nil, clientFactory)

			result, err := manager.GetSecret(tt.secretName)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if tt.expectedErrMsg != "" && !contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("Expected error message to contain %v, got: %v", tt.expectedErrMsg, err.Error())
				}
				if result != "" {
					t.Errorf("Expected empty result on error, got: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expectedValue {
					t.Errorf("GetSecret() = %v, want %v", result, tt.expectedValue)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Benchmark for TestableManager
func BenchmarkNewTestableManager(b *testing.B) {
	config := DefaultConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		NewTestableManager(config, nil, nil)
	}
}
