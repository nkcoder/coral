package config

import (
	"testing"

	"coral.daniel-guo.com/internal/email"
	"coral.daniel-guo.com/internal/secrets"
)

func TestNewAppConfig(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		testEmail   string
		sender      string
		expected    *AppConfig
	}{
		{
			name:        "default values with empty parameters",
			environment: "dev",
			testEmail:   "",
			sender:      "",
			expected: &AppConfig{
				Environment:    "dev",
				Email:          email.DefaultConfig(),
				Secrets:        secrets.DefaultConfig(),
				DefaultSender:  "no-reply@the-hub.ai",
				TestEmail:      "",
				WorkerPoolSize: 5,
				WorkerDelayMs:  1000,
			},
		},
		{
			name:        "with test email",
			environment: "test",
			testEmail:   "test@example.com",
			sender:      "",
			expected: &AppConfig{
				Environment:    "test",
				Email:          email.DefaultConfig(),
				Secrets:        secrets.DefaultConfig(),
				DefaultSender:  "no-reply@the-hub.ai",
				TestEmail:      "test@example.com",
				WorkerPoolSize: 5,
				WorkerDelayMs:  1000,
			},
		},
		{
			name:        "with custom sender",
			environment: "prod",
			testEmail:   "",
			sender:      "custom@sender.com",
			expected: &AppConfig{
				Environment:    "prod",
				Email:          email.DefaultConfig(),
				Secrets:        secrets.DefaultConfig(),
				DefaultSender:  "custom@sender.com",
				TestEmail:      "",
				WorkerPoolSize: 5,
				WorkerDelayMs:  1000,
			},
		},
		{
			name:        "with both test email and custom sender",
			environment: "staging",
			testEmail:   "test@staging.com",
			sender:      "staging@sender.com",
			expected: &AppConfig{
				Environment:    "staging",
				Email:          email.DefaultConfig(),
				Secrets:        secrets.DefaultConfig(),
				DefaultSender:  "staging@sender.com",
				TestEmail:      "test@staging.com",
				WorkerPoolSize: 5,
				WorkerDelayMs:  1000,
			},
		},
		{
			name:        "empty environment",
			environment: "",
			testEmail:   "",
			sender:      "",
			expected: &AppConfig{
				Environment:    "",
				Email:          email.DefaultConfig(),
				Secrets:        secrets.DefaultConfig(),
				DefaultSender:  "no-reply@the-hub.ai",
				TestEmail:      "",
				WorkerPoolSize: 5,
				WorkerDelayMs:  1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAppConfig(tt.environment, tt.testEmail, tt.sender)

			// Verify all fields are set correctly
			if got.Environment != tt.expected.Environment {
				t.Errorf("Environment = %v, want %v", got.Environment, tt.expected.Environment)
			}
			if got.DefaultSender != tt.expected.DefaultSender {
				t.Errorf("DefaultSender = %v, want %v", got.DefaultSender, tt.expected.DefaultSender)
			}
			if got.TestEmail != tt.expected.TestEmail {
				t.Errorf("TestEmail = %v, want %v", got.TestEmail, tt.expected.TestEmail)
			}
			if got.WorkerPoolSize != tt.expected.WorkerPoolSize {
				t.Errorf("WorkerPoolSize = %v, want %v", got.WorkerPoolSize, tt.expected.WorkerPoolSize)
			}
			if got.WorkerDelayMs != tt.expected.WorkerDelayMs {
				t.Errorf("WorkerDelayMs = %v, want %v", got.WorkerDelayMs, tt.expected.WorkerDelayMs)
			}

			// Verify Email config is set with default values
			expectedEmailConfig := email.DefaultConfig()
			if got.Email.Region != expectedEmailConfig.Region {
				t.Errorf("Email.Region = %v, want %v", got.Email.Region, expectedEmailConfig.Region)
			}

			// Verify Secrets config is set with default values
			expectedSecretsConfig := secrets.DefaultConfig()
			if got.Secrets.Region != expectedSecretsConfig.Region {
				t.Errorf("Secrets.Region = %v, want %v", got.Secrets.Region, expectedSecretsConfig.Region)
			}
		})
	}
}

func TestAppConfig_FieldModification(t *testing.T) {
	// Test that fields can be modified directly
	config := NewAppConfig("test", "", "")

	// Modify fields directly
	config.TestEmail = "modified@test.com"
	config.DefaultSender = "modified@sender.com"

	// Verify both values were set
	if config.TestEmail != "modified@test.com" {
		t.Errorf("TestEmail = %v, want %v", config.TestEmail, "modified@test.com")
	}
	if config.DefaultSender != "modified@sender.com" {
		t.Errorf("DefaultSender = %v, want %v", config.DefaultSender, "modified@sender.com")
	}
}

func TestAppConfig_DefaultValues(t *testing.T) {
	// Test that default values are properly set
	config := NewAppConfig("test", "", "")

	// Test default values
	expectedDefaults := map[string]interface{}{
		"WorkerPoolSize": 5,
		"WorkerDelayMs":  1000,
		"DefaultSender":  "no-reply@the-hub.ai",
		"TestEmail":      "",
		"Environment":    "test",
	}

	actualValues := map[string]interface{}{
		"WorkerPoolSize": config.WorkerPoolSize,
		"WorkerDelayMs":  config.WorkerDelayMs,
		"DefaultSender":  config.DefaultSender,
		"TestEmail":      config.TestEmail,
		"Environment":    config.Environment,
	}

	for key, expected := range expectedDefaults {
		if actual := actualValues[key]; actual != expected {
			t.Errorf("Default value for %s = %v, want %v", key, actual, expected)
		}
	}
}

func TestAppConfig_DependencyConfigs(t *testing.T) {
	// Test that dependency configs are properly initialized
	config := NewAppConfig("test", "", "")

	// Test email config initialization
	expectedEmailRegion := "ap-southeast-2"
	if config.Email.Region != expectedEmailRegion {
		t.Errorf("Email.Region = %v, want %v", config.Email.Region, expectedEmailRegion)
	}

	// Test secrets config initialization
	expectedSecretsRegion := "ap-southeast-2"
	if config.Secrets.Region != expectedSecretsRegion {
		t.Errorf("Secrets.Region = %v, want %v", config.Secrets.Region, expectedSecretsRegion)
	}
}

func TestAppConfig_StructFields(t *testing.T) {
	// Test that all struct fields are accessible and have correct types
	config := NewAppConfig("test", "test@example.com", "test@sender.com")

	// Verify field types and accessibility
	var _ string = config.Environment
	var _ email.Config = config.Email
	var _ secrets.Config = config.Secrets
	var _ string = config.DefaultSender
	var _ string = config.TestEmail
	var _ int = config.WorkerPoolSize
	var _ int = config.WorkerDelayMs

	// This test ensures that if struct fields change, the test will break
	// and force us to update the tests accordingly
}

func BenchmarkNewAppConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewAppConfig("dev", "test@example.com", "sender@example.com")
	}
}

func BenchmarkFieldAssignment(b *testing.B) {
	config := NewAppConfig("dev", "", "")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		config.TestEmail = "test@example.com"
		config.DefaultSender = "sender@example.com"
	}
}
