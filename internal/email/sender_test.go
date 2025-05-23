package email

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockSESAPI is a mock implementation of the SES API
type MockSESAPI struct {
	mock.Mock
}

func (m *MockSESAPI) SendRawEmail(input *ses.SendRawEmailInput) (*ses.SendRawEmailOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*ses.SendRawEmailOutput), args.Error(1)
}

type EmailSenderTestSuite struct {
	suite.Suite
	sender  *Sender
	mockSES *MockSESAPI
}

func (suite *EmailSenderTestSuite) SetupTest() {
	config := Config{Region: "us-east-1"}
	suite.sender = NewSender(config)
	suite.mockSES = new(MockSESAPI)
}

func (suite *EmailSenderTestSuite) TestNewSender() {
	config := Config{Region: "us-west-2"}
	sender := NewSender(config)

	assert.NotNil(suite.T(), sender)
	assert.Equal(suite.T(), "us-west-2", sender.config.Region)
}

func (suite *EmailSenderTestSuite) TestDefaultConfig() {
	config := DefaultConfig()
	assert.Equal(suite.T(), "ap-southeast-2", config.Region)
}

func (suite *EmailSenderTestSuite) TestStripHTML() {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple HTML",
			input:    "<p>Hello <b>World</b></p>",
			expected: "Hello World",
		},
		{
			name:     "Complex HTML",
			input:    "<html><head></head><body><p>Hello team,</p><p>Regards</p></body></html>",
			expected: "Hello team,Regards",
		},
		{
			name:     "No HTML",
			input:    "Plain text message",
			expected: "Plain text message",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "HTML with attributes",
			input:    `<div class="container"><span style="color: red;">Red text</span></div>`,
			expected: "Red text",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			result := StripHTML(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Note: Testing SendWithAttachment and SendWithAttachmentFile would require
// more complex mocking of the AWS session and SES service creation.
// For now, we'll focus on testing the components we can easily test.

func (suite *EmailSenderTestSuite) TestSendWithAttachmentValidation() {
	// Test that the method handles basic validation
	sender := "test@example.com"
	recipient := "recipient@example.com"
	subject := "Test Subject"
	body := "<p>Test Body</p>"
	attachmentName := "test.csv"
	attachmentContent := []byte("test,data\n1,2")

	// This will fail at AWS session creation, but we can test the input validation
	err := suite.sender.SendWithAttachment(sender, recipient, subject, body, attachmentName, attachmentContent)

	// We expect an error because we don't have AWS credentials in test environment
	// The error could be about AWS session, credentials, or region configuration
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to send email")
}

func (suite *EmailSenderTestSuite) TestSendWithAttachmentFileNotFound() {
	err := suite.sender.SendWithAttachmentFile(
		"test@example.com",
		"recipient@example.com",
		"Test Subject",
		"Test Body",
		"nonexistent.csv",
	)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to read file")
}

func TestEmailSenderSuite(t *testing.T) {
	suite.Run(t, new(EmailSenderTestSuite))
}
