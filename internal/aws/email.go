package aws

import (
	"bytes"
	"fmt"
	"mime"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"

	"encoding/base64"

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

	// Extract filename from path
	filename := filepath.Base(attachmentPath)

	// Create a buffer for the message
	var buf bytes.Buffer

	// Set up email headers
	headers := textproto.MIMEHeader{}
	headers.Set("From", sender)
	headers.Set("To", recipient)
	headers.Set("Subject", subject)
	headers.Set("MIME-Version", "1.0")

	// Create multipart writer
	writer := multipart.NewWriter(&buf)
	headers.Set("Content-Type", "multipart/mixed; boundary="+writer.Boundary())

	// Write headers
	for k, vv := range headers {
		for _, v := range vv {
			fmt.Fprintf(&buf, "%s: %s\r\n", k, v)
		}
	}
	buf.WriteString("\r\n")

	// Create alternative part for text/html versions
	altWriter := multipart.NewWriter(&buf)
	fmt.Fprintf(&buf, "--%s\r\n", writer.Boundary())
	fmt.Fprintf(&buf, "Content-Type: multipart/alternative; boundary=%s\r\n\r\n", altWriter.Boundary())

	// Add text part
	fmt.Fprintf(&buf, "--%s\r\n", altWriter.Boundary())
	fmt.Fprintf(&buf, "Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	buf.WriteString(StripHTML(body))
	buf.WriteString("\r\n")

	// Add HTML part
	fmt.Fprintf(&buf, "--%s\r\n", altWriter.Boundary())
	fmt.Fprintf(&buf, "Content-Type: text/html; charset=UTF-8\r\n\r\n")
	buf.WriteString(body)
	buf.WriteString("\r\n")

	// Close alternative part
	fmt.Fprintf(&buf, "--%s--\r\n", altWriter.Boundary())

	// Add attachment
	fileExt := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(fileExt)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	fmt.Fprintf(&buf, "--%s\r\n", writer.Boundary())
	fmt.Fprintf(&buf, "Content-Type: %s\r\n", mimeType)
	fmt.Fprintf(&buf, "Content-Disposition: attachment; filename=%s\r\n", filename)
	fmt.Fprintf(&buf, "Content-Transfer-Encoding: base64\r\n\r\n")

	// Base64 encode the attachment
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	_, err = encoder.Write(fileContent)
	if err != nil {
		return fmt.Errorf("failed to encode attachment: %w", err)
	}
	encoder.Close()
	buf.WriteString("\r\n")

	// Close the multipart message
	fmt.Fprintf(&buf, "--%s--\r\n", writer.Boundary())

	// Create a new AWS SES session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	})
	if err != nil {
		return fmt.Errorf("failed to create SES session: %w", err)
	}

	// Create a new SES client
	svc := ses.New(sess)

	// Send the raw email
	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: buf.Bytes(),
		},
		Source: aws.String(sender),
		Destinations: []*string{
			aws.String(recipient),
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
