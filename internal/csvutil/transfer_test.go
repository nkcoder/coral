package csvutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"coral.daniel-guo.com/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CSVUtilTestSuite struct {
	suite.Suite
	tempDir string
}

func (suite *CSVUtilTestSuite) SetupTest() {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "csvutil_test")
	assert.NoError(suite.T(), err)
	suite.tempDir = tempDir
}

func (suite *CSVUtilTestSuite) TearDownTest() {
	// Clean up temporary directory
	_ = os.RemoveAll(suite.tempDir)
}

func (suite *CSVUtilTestSuite) createTestCSVFile(filename, content string) string {
	filePath := filepath.Join(suite.tempDir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(suite.T(), err)
	return filePath
}

func (suite *CSVUtilTestSuite) TestReadClubTransferCSVSuccess() {
	csvContent := `Member Id,Fob Number,First Name,Last Name,Membership Type,Home Club,Target Club
12345,FOB001,John,Doe,Premium,CLUB A,CLUB B
67890,FOB002,Jane,Smith,Standard,club c,club d`

	filePath := suite.createTestCSVFile("test.csv", csvContent)

	result, err := ReadClubTransferCSV(filePath)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	// Check first row
	assert.Equal(suite.T(), "12345", result[0].MemberID)
	assert.Equal(suite.T(), "FOB001", result[0].FobNumber)
	assert.Equal(suite.T(), "John", result[0].FirstName)
	assert.Equal(suite.T(), "Doe", result[0].LastName)
	assert.Equal(suite.T(), "Premium", result[0].MembershipType)
	assert.Equal(suite.T(), "CLUB A", result[0].HomeClub)
	assert.Equal(suite.T(), "CLUB B", result[0].TargetClub)

	// Check second row (should be uppercase)
	assert.Equal(suite.T(), "67890", result[1].MemberID)
	assert.Equal(suite.T(), "CLUB C", result[1].HomeClub)
	assert.Equal(suite.T(), "CLUB D", result[1].TargetClub)
}

func (suite *CSVUtilTestSuite) TestReadClubTransferCSVFileNotFound() {
	_, err := ReadClubTransferCSV("nonexistent.csv")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to open file")
}

func (suite *CSVUtilTestSuite) TestReadClubTransferCSVEmptyFile() {
	filePath := suite.createTestCSVFile("empty.csv", "")

	_, err := ReadClubTransferCSV(filePath)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "no records found")
}

func (suite *CSVUtilTestSuite) TestReadClubTransferCSVMissingColumns() {
	csvContent := `Member Id,First Name,Last Name
12345,John,Doe`

	filePath := suite.createTestCSVFile("missing_cols.csv", csvContent)

	_, err := ReadClubTransferCSV(filePath)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "column")
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *CSVUtilTestSuite) TestReadClubTransferCSVInvalidCSV() {
	csvContent := `Member Id,Fob Number,First Name,Last Name,Membership Type,Home Club,Target Club
12345,FOB001,John,"Doe,Premium,CLUB A,CLUB B`

	filePath := suite.createTestCSVFile("invalid.csv", csvContent)

	_, err := ReadClubTransferCSV(filePath)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to read file")
}

func (suite *CSVUtilTestSuite) TestGenerateCSVContentSuccess() {
	transferDate := time.Date(2023, 12, 15, 10, 30, 0, 0, time.UTC)

	data := []model.ClubTransferData{
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
			TargetClub:     "CLUB D",
			TransferType:   "TRANSFER OUT",
			TransferDate:   transferDate,
		},
	}

	result, err := GenerateCSVContent(data)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)

	// Convert to string for easier testing
	csvString := string(result)
	lines := strings.Split(strings.TrimSpace(csvString), "\n")

	// Check header
	expectedHeader := "Member Id,Fob Number,First Name,Last Name,Membership Type,Home Club,Target Club,Transfer Type,Transfer Date"
	assert.Equal(suite.T(), expectedHeader, lines[0])

	// Check first data row
	expectedFirstRow := "12345,FOB001,John,Doe,Premium,CLUB A,CLUB B,TRANSFER IN,2023-12-15"
	assert.Equal(suite.T(), expectedFirstRow, lines[1])

	// Check second data row
	expectedSecondRow := "67890,FOB002,Jane,Smith,Standard,CLUB C,CLUB D,TRANSFER OUT,2023-12-15"
	assert.Equal(suite.T(), expectedSecondRow, lines[2])
}

func (suite *CSVUtilTestSuite) TestGenerateCSVContentEmpty() {
	data := []model.ClubTransferData{}

	result, err := GenerateCSVContent(data)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)

	// Should still have header
	csvString := string(result)
	lines := strings.Split(strings.TrimSpace(csvString), "\n")
	assert.Len(suite.T(), lines, 1) // Only header

	expectedHeader := "Member Id,Fob Number,First Name,Last Name,Membership Type,Home Club,Target Club,Transfer Type,Transfer Date"
	assert.Equal(suite.T(), expectedHeader, lines[0])
}

func (suite *CSVUtilTestSuite) TestGenerateCSVContentSpecialCharacters() {
	transferDate := time.Date(2023, 12, 15, 10, 30, 0, 0, time.UTC)

	data := []model.ClubTransferData{
		{
			MemberID:       "12345",
			FobNumber:      "FOB001",
			FirstName:      "John,Jr",
			LastName:       "O'Doe",
			MembershipType: "Premium \"VIP\"",
			HomeClub:       "CLUB A",
			TargetClub:     "CLUB B",
			TransferType:   "TRANSFER IN",
			TransferDate:   transferDate,
		},
	}

	result, err := GenerateCSVContent(data)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)

	// CSV should properly escape special characters
	csvString := string(result)
	assert.Contains(suite.T(), csvString, "\"John,Jr\"")
	assert.Contains(suite.T(), csvString, "O'Doe")
	assert.Contains(suite.T(), csvString, "\"Premium \"\"VIP\"\"\"")
}

func TestCSVUtilSuite(t *testing.T) {
	suite.Run(t, new(CSVUtilTestSuite))
}
