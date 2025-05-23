// Package testutil provides testing utilities and helpers
package testutil

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"coral.daniel-guo.com/internal/config"
	"coral.daniel-guo.com/internal/email"
	"coral.daniel-guo.com/internal/model"
	"coral.daniel-guo.com/internal/secrets"
)

// CreateTempDir creates a temporary directory for testing
func CreateTempDir(t *testing.T, prefix string) string {
	tempDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tempDir
}

// CreateTestCSVFile creates a CSV file with the given content in the specified directory
func CreateTestCSVFile(t *testing.T, dir, filename, content string) string {
	filePath := filepath.Join(dir, filename)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}
	return filePath
}

// CleanupTempDir removes a temporary directory
func CleanupTempDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Logf("Warning: Failed to cleanup temp dir %s: %v", dir, err)
	}
}

// CreateTestConfig creates a test configuration
func CreateTestConfig() *config.AppConfig {
	return &config.AppConfig{
		Environment:    "test",
		DefaultSender:  "test@example.com",
		TestEmail:      "",
		WorkerPoolSize: 2,
		WorkerDelayMs:  100,
		Email: email.Config{
			Region: "us-east-1",
		},
		Secrets: secrets.Config{
			Region: "us-east-1",
		},
	}
}

// CreateTestClubTransferData creates test club transfer data
func CreateTestClubTransferData() []model.ClubTransferData {
	transferDate := time.Date(2023, 12, 15, 10, 30, 0, 0, time.UTC)

	return []model.ClubTransferData{
		{
			MemberID:       "12345",
			FobNumber:      "FOB001",
			FirstName:      "John",
			LastName:       "Doe",
			MembershipType: "Premium",
			HomeClub:       "CLUB A",
			TargetClub:     "CLUB B",
			TransferType:   "TRANSFER IN",
			TransferDate:   transferDate,
		},
		{
			MemberID:       "67890",
			FobNumber:      "FOB002",
			FirstName:      "Jane",
			LastName:       "Smith",
			MembershipType: "Standard",
			HomeClub:       "CLUB C",
			TargetClub:     "CLUB A",
			TransferType:   "TRANSFER OUT",
			TransferDate:   transferDate,
		},
	}
}

// CreateTestLocation creates a test location
func CreateTestLocation(id, name, email string) *model.Location {
	return &model.Location{
		ID:    id,
		Name:  name,
		Email: email,
	}
}

// SampleCSVContent provides sample CSV content for testing
const SampleCSVContent = `Member Id,Fob Number,First Name,Last Name,Membership Type,Home Club,Target Club
12345,FOB001,John,Doe,Premium,CLUB A,CLUB B
67890,FOB002,Jane,Smith,Standard,CLUB C,CLUB A
11111,FOB003,Bob,Johnson,Basic,CLUB B,CLUB C`

// InvalidCSVContent provides invalid CSV content for testing error cases
const InvalidCSVContent = `Member Id,First Name,Last Name
12345,John,Doe`

// EmptyCSVContent provides empty CSV content for testing
const EmptyCSVContent = ``
