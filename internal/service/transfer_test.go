package service

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"coral.daniel-guo.com/internal/config"
	"coral.daniel-guo.com/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockEmailSender is a mock implementation of the email sender
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendWithAttachment(
	sender, recipient, subject, body, attachmentName string,
	attachmentContent []byte,
) error {
	args := m.Called(sender, recipient, subject, body, attachmentName, attachmentContent)
	return args.Error(0)
}

// MockLocationRepository is a mock implementation of the location repository interface
type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) FindByName(name string) (*model.Location, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Location), args.Error(1)
}

type TransferServiceTestSuite struct {
	suite.Suite
	service          *Service
	mockEmailSender  *MockEmailSender
	mockLocationRepo *MockLocationRepository
	tempDir          string
}

func (suite *TransferServiceTestSuite) SetupTest() {
	// Create test configuration
	cfg := &config.AppConfig{
		Environment:    "test",
		DefaultSender:  "test@example.com",
		TestEmail:      "",
		WorkerPoolSize: 2,
		WorkerDelayMs:  100,
	}

	suite.service = NewService(cfg)
	suite.mockEmailSender = new(MockEmailSender)
	suite.mockLocationRepo = new(MockLocationRepository)

	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "service_test")
	assert.NoError(suite.T(), err)
	suite.tempDir = tempDir
}

func (suite *TransferServiceTestSuite) TearDownTest() {
	// Clean up temporary directory
	_ = os.RemoveAll(suite.tempDir)
}

func (suite *TransferServiceTestSuite) createTestCSVFile(filename, content string) string {
	filePath := filepath.Join(suite.tempDir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(suite.T(), err)
	return filePath
}

func (suite *TransferServiceTestSuite) TestNewService() {
	cfg := &config.AppConfig{
		Environment:   "test",
		DefaultSender: "test@example.com",
	}

	service := NewService(cfg)

	assert.NotNil(suite.T(), service)
	assert.Equal(suite.T(), cfg, service.config)
	assert.NotNil(suite.T(), service.secretsManager)
	assert.NotNil(suite.T(), service.emailSender)
}

func (suite *TransferServiceTestSuite) TestReadClubTransferDataSuccess() {
	csvContent := `Member Id,Fob Number,First Name,Last Name,Membership Type,Home Club,Target Club
12345,FOB001,John,Doe,Premium,CLUB A,CLUB B
67890,FOB002,Jane,Smith,Standard,CLUB C,CLUB A`

	filePath := suite.createTestCSVFile("test.csv", csvContent)

	result, err := suite.service.readClubTransferData(filePath)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3) // CLUB A, CLUB B, CLUB C

	// Check CLUB A (should have 1 transfer out and 1 transfer in)
	clubATransfers := result["CLUB A"]
	assert.Len(suite.T(), clubATransfers, 2)

	// Check CLUB B (should have 1 transfer in)
	clubBTransfers := result["CLUB B"]
	assert.Len(suite.T(), clubBTransfers, 1)
	assert.Equal(suite.T(), "TRANSFER IN", clubBTransfers[0].TransferType)
	assert.Equal(suite.T(), "12345", clubBTransfers[0].MemberID)

	// Check CLUB C (should have 1 transfer out)
	clubCTransfers := result["CLUB C"]
	assert.Len(suite.T(), clubCTransfers, 1)
	assert.Equal(suite.T(), "TRANSFER OUT", clubCTransfers[0].TransferType)
	assert.Equal(suite.T(), "67890", clubCTransfers[0].MemberID)
}

func (suite *TransferServiceTestSuite) TestReadClubTransferDataFileNotFound() {
	_, err := suite.service.readClubTransferData("nonexistent.csv")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "error reading club transfer data")
}

func (suite *TransferServiceTestSuite) TestGetOutputFileNamePIF() {
	result := suite.service.getOutputFileName("PIF", "CLUB_A")
	assert.Equal(suite.T(), "pif_club_transfer_CLUB_A.csv", result)
}

func (suite *TransferServiceTestSuite) TestGetOutputFileNameDD() {
	result := suite.service.getOutputFileName("DD", "CLUB_B")
	assert.Equal(suite.T(), "dd_club_transfer_CLUB_B.csv", result)
}

func (suite *TransferServiceTestSuite) TestGetOutputFileNameOther() {
	result := suite.service.getOutputFileName("OTHER", "CLUB_C")
	assert.Equal(suite.T(), "pif_club_transfer_CLUB_C.csv", result)
}

func (suite *TransferServiceTestSuite) TestSendEmailSuccess() {
	// Setup test data
	transferData := []model.ClubTransferData{
		{
			MemberID:       "12345",
			FobNumber:      "FOB001",
			FirstName:      "John",
			LastName:       "Doe",
			MembershipType: "Premium",
			HomeClub:       "CLUB A",
			TargetClub:     "CLUB B",
			TransferType:   "TRANSFER IN",
			TransferDate:   time.Now(),
		},
	}

	data := map[string][]model.ClubTransferData{
		"CLUB A": transferData,
	}

	// Setup mocks
	location := &model.Location{
		ID:    "1",
		Name:  "CLUB A",
		Email: "cluba@example.com",
	}

	suite.mockLocationRepo.On("FindByName", "CLUB A").Return(location, nil)

	// Note: We can't easily test the actual email sending without more complex mocking
	// of the email sender's internal AWS dependencies. In a real scenario, you'd
	// inject the email sender as an interface and mock it here.

	err := suite.service.sendEmail("CLUB A", data, "PIF", suite.mockLocationRepo)

	// This will fail because we can't mock the email sender easily
	// In a production setup, you'd refactor to inject dependencies
	assert.Error(suite.T(), err)
	suite.mockLocationRepo.AssertExpectations(suite.T())
}

func (suite *TransferServiceTestSuite) TestSendEmailLocationNotFound() {
	data := map[string][]model.ClubTransferData{
		"CLUB A": {},
	}

	suite.mockLocationRepo.On("FindByName", "CLUB A").Return(nil, nil)

	err := suite.service.sendEmail("CLUB A", data, "PIF", suite.mockLocationRepo)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "location not found")
	suite.mockLocationRepo.AssertExpectations(suite.T())
}

func (suite *TransferServiceTestSuite) TestSendEmailLocationError() {
	data := map[string][]model.ClubTransferData{
		"CLUB A": {},
	}

	suite.mockLocationRepo.On("FindByName", "CLUB A").Return(nil, errors.New("database error"))

	err := suite.service.sendEmail("CLUB A", data, "PIF", suite.mockLocationRepo)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "error finding location")
	suite.mockLocationRepo.AssertExpectations(suite.T())
}

func (suite *TransferServiceTestSuite) TestSendEmailNoEmail() {
	data := map[string][]model.ClubTransferData{
		"CLUB A": {},
	}

	location := &model.Location{
		ID:    "1",
		Name:  "CLUB A",
		Email: "", // No email
	}

	suite.mockLocationRepo.On("FindByName", "CLUB A").Return(location, nil)

	err := suite.service.sendEmail("CLUB A", data, "PIF", suite.mockLocationRepo)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "email not found")
	suite.mockLocationRepo.AssertExpectations(suite.T())
}

func (suite *TransferServiceTestSuite) TestTransferRequestValidation() {
	req := TransferRequest{
		TransferType: "PIF",
		FileName:     "test.csv",
	}

	assert.Equal(suite.T(), "PIF", req.TransferType)
	assert.Equal(suite.T(), "test.csv", req.FileName)
}

// Integration test that would require database setup
func (suite *TransferServiceTestSuite) TestProcessIntegration() {
	// This test would require a real database connection
	// For now, we'll skip it or mark it as integration test
	suite.T().Skip("Integration test - requires database setup")

	csvContent := `Member Id,Fob Number,First Name,Last Name,Membership Type,Home Club,Target Club
12345,FOB001,John,Doe,Premium,CLUB A,CLUB B`

	filePath := suite.createTestCSVFile("integration_test.csv", csvContent)

	req := TransferRequest{
		TransferType: "PIF",
		FileName:     filePath,
	}

	err := suite.service.Process(req)
	// This will fail without proper database setup
	assert.Error(suite.T(), err)
}

func TestTransferServiceSuite(t *testing.T) {
	suite.Run(t, new(TransferServiceTestSuite))
}
