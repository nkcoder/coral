package aws

import (
	"bytes"
	"fmt"
	"mime"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SendEmailWithAttachment sends an email with an attachment using AWS SES
func SendEmailWithAttachment(sender, recipient, subject, body, attachmentPath string) error {
	// Read the file content
	fileContent, err := os.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create a new multipart message
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add an alternate part for the plain text
	textPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/plain; charset=UTF-8"},
	})
	if err != nil {
		return fmt.Errorf("failed to create text part: %w", err)
	}
	_, err = textPart.Write([]byte(StripHTML(body)))
	if err != nil {
		return fmt.Errorf("failed to write text part: %w", err)
	}

	// Add an alternate part for the HTML
	htmlPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/html; charset=UTF-8"},
	})
	if err != nil {
		return fmt.Errorf("failed to create HTML part: %w", err)
	}
	_, err = htmlPart.Write([]byte(body))
	if err != nil {
		return fmt.Errorf("failed to write HTML part: %w", err)
	}

	// Add attachment
	fileExt := filepath.Ext(attachmentPath)
	mimeType := mime.TypeByExtension(fileExt)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	attachmentPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":        []string{mimeType},
		"Content-Disposition": []string{fmt.Sprintf("attachment; filename=%s", attachmentPath)},
	})
	if err != nil {
		return fmt.Errorf("failed to create attachment part: %w", err)
	}
	_, err = attachmentPart.Write(fileContent)
	if err != nil {
		return fmt.Errorf("failed to write attachment part: %w", err)
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// create a new AWS SES session
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	})
	if err != nil {
		return fmt.Errorf("failed to create SES session: %w", err)
	}

	// create a new SES client
	svc := ses.New(session)
	rawMessage := buf.Bytes()

	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: rawMessage,
		},
	}

	_, err = svc.SendRawEmail(input)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("Email sent to: %s successfully\n", recipient)

	return nil
}

// Helper function to strip HTML tags for plain text version
func StripHTML(html string) string {
	var buf bytes.Buffer
	inTag := false

	for _, r := range html {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			buf.WriteRune(r)
		}
	}

	return buf.String()
}
