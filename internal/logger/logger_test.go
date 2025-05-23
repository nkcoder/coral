package logger

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
	originalDebugLogger *log.Logger
	originalInfoLogger  *log.Logger
	originalWarnLogger  *log.Logger
	originalErrorLogger *log.Logger
	originalLevel       LogLevel
	debugBuf            *bytes.Buffer
	infoBuf             *bytes.Buffer
	warnBuf             *bytes.Buffer
	errorBuf            *bytes.Buffer
}

func (suite *LoggerTestSuite) SetupTest() {
	// Save original loggers and level
	suite.originalDebugLogger = debugLogger
	suite.originalInfoLogger = infoLogger
	suite.originalWarnLogger = warnLogger
	suite.originalErrorLogger = errorLogger
	suite.originalLevel = currentLevel

	// Create buffers for capturing output
	suite.debugBuf = &bytes.Buffer{}
	suite.infoBuf = &bytes.Buffer{}
	suite.warnBuf = &bytes.Buffer{}
	suite.errorBuf = &bytes.Buffer{}

	// Replace loggers with test loggers
	debugLogger = log.New(suite.debugBuf, "DEBUG: ", log.Ldate|log.Ltime)
	infoLogger = log.New(suite.infoBuf, "INFO:  ", log.Ldate|log.Ltime)
	warnLogger = log.New(suite.warnBuf, "WARN:  ", log.Ldate|log.Ltime)
	errorLogger = log.New(suite.errorBuf, "ERROR: ", log.Ldate|log.Ltime)

	// Reset to default level
	currentLevel = InfoLevel
}

func (suite *LoggerTestSuite) TearDownTest() {
	// Restore original loggers and level
	debugLogger = suite.originalDebugLogger
	infoLogger = suite.originalInfoLogger
	warnLogger = suite.originalWarnLogger
	errorLogger = suite.originalErrorLogger
	currentLevel = suite.originalLevel
}

func (suite *LoggerTestSuite) TestSetLevel() {
	SetLevel(DebugLevel)
	assert.Equal(suite.T(), DebugLevel, currentLevel)

	SetLevel(ErrorLevel)
	assert.Equal(suite.T(), ErrorLevel, currentLevel)
}

func (suite *LoggerTestSuite) TestDebugLogging() {
	SetLevel(DebugLevel)

	Debug("Test debug message")
	assert.Contains(suite.T(), suite.debugBuf.String(), "Test debug message")

	// Should not log when level is higher
	SetLevel(InfoLevel)
	suite.debugBuf.Reset()
	Debug("Should not appear")
	assert.Empty(suite.T(), suite.debugBuf.String())
}

func (suite *LoggerTestSuite) TestInfoLogging() {
	SetLevel(InfoLevel)

	Info("Test info message")
	assert.Contains(suite.T(), suite.infoBuf.String(), "Test info message")

	// Should not log when level is higher
	SetLevel(WarnLevel)
	suite.infoBuf.Reset()
	Info("Should not appear")
	assert.Empty(suite.T(), suite.infoBuf.String())
}

func (suite *LoggerTestSuite) TestWarnLogging() {
	SetLevel(WarnLevel)

	Warn("Test warn message")
	assert.Contains(suite.T(), suite.warnBuf.String(), "Test warn message")

	// Should not log when level is higher
	SetLevel(ErrorLevel)
	suite.warnBuf.Reset()
	Warn("Should not appear")
	assert.Empty(suite.T(), suite.warnBuf.String())
}

func (suite *LoggerTestSuite) TestErrorLogging() {
	SetLevel(ErrorLevel)

	Error("Test error message")
	assert.Contains(suite.T(), suite.errorBuf.String(), "Test error message")

	// Error should always log regardless of level
	SetLevel(DebugLevel)
	suite.errorBuf.Reset()
	Error("Should appear")
	assert.Contains(suite.T(), suite.errorBuf.String(), "Should appear")
}

func (suite *LoggerTestSuite) TestFormattedLogging() {
	SetLevel(DebugLevel)

	Debug("Debug: %s %d", "test", 123)
	Info("Info: %s %d", "test", 456)
	Warn("Warn: %s %d", "test", 789)
	Error("Error: %s %d", "test", 999)

	assert.Contains(suite.T(), suite.debugBuf.String(), "Debug: test 123")
	assert.Contains(suite.T(), suite.infoBuf.String(), "Info: test 456")
	assert.Contains(suite.T(), suite.warnBuf.String(), "Warn: test 789")
	assert.Contains(suite.T(), suite.errorBuf.String(), "Error: test 999")
}

func (suite *LoggerTestSuite) TestLogLevelHierarchy() {
	// Test that higher levels include lower levels
	SetLevel(DebugLevel)
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")

	assert.Contains(suite.T(), suite.debugBuf.String(), "debug")
	assert.Contains(suite.T(), suite.infoBuf.String(), "info")
	assert.Contains(suite.T(), suite.warnBuf.String(), "warn")
	assert.Contains(suite.T(), suite.errorBuf.String(), "error")

	// Reset buffers
	suite.debugBuf.Reset()
	suite.infoBuf.Reset()
	suite.warnBuf.Reset()
	suite.errorBuf.Reset()

	// Test that lower levels exclude higher levels
	SetLevel(WarnLevel)
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")

	assert.Empty(suite.T(), suite.debugBuf.String())
	assert.Empty(suite.T(), suite.infoBuf.String())
	assert.Contains(suite.T(), suite.warnBuf.String(), "warn")
	assert.Contains(suite.T(), suite.errorBuf.String(), "error")
}

func (suite *LoggerTestSuite) TestLogLevelConstants() {
	assert.Equal(suite.T(), LogLevel(0), DebugLevel)
	assert.Equal(suite.T(), LogLevel(1), InfoLevel)
	assert.Equal(suite.T(), LogLevel(2), WarnLevel)
	assert.Equal(suite.T(), LogLevel(3), ErrorLevel)
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
