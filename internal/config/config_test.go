package config

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) SetupTest() {
	// Reset viper before each test
	viper.Reset()
}

func (suite *ConfigTestSuite) TearDownTest() {
	// Clean up environment variables
	_ = os.Unsetenv("CORAL_ENV")
	_ = os.Unsetenv("CORAL_EMAIL_SENDER")
	_ = os.Unsetenv("CORAL_AWS_REGION")
	_ = os.Unsetenv("CORAL_WORKER_POOL_SIZE")
	_ = os.Unsetenv("CORAL_WORKER_DELAY_MS")
}

func (suite *ConfigTestSuite) TestSetDefaults() {
	setDefaults()

	assert.Equal(suite.T(), "dev", viper.GetString("env"))
	assert.Equal(suite.T(), "no-reply@the-hub.ai", viper.GetString("email.sender"))
	assert.Equal(suite.T(), "", viper.GetString("email.test"))
	assert.Equal(suite.T(), "ap-southeast-2", viper.GetString("aws.region"))
	assert.Equal(suite.T(), 5, viper.GetInt("worker.pool_size"))
	assert.Equal(suite.T(), 1000, viper.GetInt("worker.delay_ms"))
}

func (suite *ConfigTestSuite) TestLoadConfigWithDefaults() {
	setDefaults()

	cfg := LoadConfig()

	assert.Equal(suite.T(), "dev", cfg.Environment)
	assert.Equal(suite.T(), "no-reply@the-hub.ai", cfg.DefaultSender)
	assert.Equal(suite.T(), "", cfg.TestEmail)
	assert.Equal(suite.T(), "ap-southeast-2", cfg.Email.Region)
	assert.Equal(suite.T(), "ap-southeast-2", cfg.Secrets.Region)
	assert.Equal(suite.T(), 5, cfg.WorkerPoolSize)
	assert.Equal(suite.T(), 1000, cfg.WorkerDelayMs)
}

func (suite *ConfigTestSuite) TestLoadConfigWithEnvironmentVariables() {
	// Set environment variables
	_ = os.Setenv("CORAL_ENV", "prod")
	_ = os.Setenv("CORAL_EMAIL_SENDER", "test@example.com")
	_ = os.Setenv("CORAL_AWS_REGION", "us-west-2")
	_ = os.Setenv("CORAL_WORKER_POOL_SIZE", "10")
	_ = os.Setenv("CORAL_WORKER_DELAY_MS", "2000")

	// Initialize with environment variables and defaults
	setDefaults()
	viper.SetEnvPrefix("CORAL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	cfg := LoadConfig()

	assert.Equal(suite.T(), "prod", cfg.Environment)
	assert.Equal(suite.T(), "test@example.com", cfg.DefaultSender)
	assert.Equal(suite.T(), "us-west-2", cfg.Email.Region)
	assert.Equal(suite.T(), "us-west-2", cfg.Secrets.Region)
	assert.Equal(suite.T(), 10, cfg.WorkerPoolSize)
	assert.Equal(suite.T(), 2000, cfg.WorkerDelayMs)
}

func (suite *ConfigTestSuite) TestLoadConfigWithViperSettings() {
	setDefaults()

	// Override some settings directly in viper
	viper.Set("env", "staging")
	viper.Set("email.sender", "staging@example.com")
	viper.Set("email.test", "test@example.com")
	viper.Set("worker.pool_size", 3)

	cfg := LoadConfig()

	assert.Equal(suite.T(), "staging", cfg.Environment)
	assert.Equal(suite.T(), "staging@example.com", cfg.DefaultSender)
	assert.Equal(suite.T(), "test@example.com", cfg.TestEmail)
	assert.Equal(suite.T(), 3, cfg.WorkerPoolSize)
	assert.Equal(suite.T(), 1000, cfg.WorkerDelayMs) // Should remain default
}

func (suite *ConfigTestSuite) TestNewAppConfigWithDefaults() {
	cfg := NewAppConfig("test")

	assert.Equal(suite.T(), "test", cfg.Environment)
	assert.Equal(suite.T(), "no-reply@the-hub.ai", cfg.DefaultSender)
	assert.Equal(suite.T(), "", cfg.TestEmail)
	assert.Equal(suite.T(), "ap-southeast-2", cfg.Email.Region)
	assert.Equal(suite.T(), "ap-southeast-2", cfg.Secrets.Region)
	assert.Equal(suite.T(), 5, cfg.WorkerPoolSize)
	assert.Equal(suite.T(), 1000, cfg.WorkerDelayMs)
}

func (suite *ConfigTestSuite) TestAppConfigChaining() {
	cfg := NewAppConfig("test").
		WithSender("custom@example.com").
		WithTestEmail("test@example.com")

	assert.Equal(suite.T(), "test", cfg.Environment)
	assert.Equal(suite.T(), "custom@example.com", cfg.DefaultSender)
	assert.Equal(suite.T(), "test@example.com", cfg.TestEmail)
}

func (suite *ConfigTestSuite) TestInitWithoutConfigFile() {
	// This should not error even without a config file
	err := Init()
	assert.NoError(suite.T(), err)
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
