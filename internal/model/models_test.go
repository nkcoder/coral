package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ModelTestSuite struct {
	suite.Suite
}

func (suite *ModelTestSuite) TestClubTransferDataCreation() {
	transferDate := time.Date(2023, 12, 15, 10, 30, 0, 0, time.UTC)

	data := ClubTransferData{
		MemberID:       "12345",
		FobNumber:      "FOB001",
		FirstName:      "John",
		LastName:       "Doe",
		MembershipType: "Premium",
		HomeClub:       "CLUB A",
		TargetClub:     "CLUB B",
		TransferType:   "TRANSFER IN",
		TransferDate:   transferDate,
	}

	assert.Equal(suite.T(), "12345", data.MemberID)
	assert.Equal(suite.T(), "FOB001", data.FobNumber)
	assert.Equal(suite.T(), "John", data.FirstName)
	assert.Equal(suite.T(), "Doe", data.LastName)
	assert.Equal(suite.T(), "Premium", data.MembershipType)
	assert.Equal(suite.T(), "CLUB A", data.HomeClub)
	assert.Equal(suite.T(), "CLUB B", data.TargetClub)
	assert.Equal(suite.T(), "TRANSFER IN", data.TransferType)
	assert.Equal(suite.T(), transferDate, data.TransferDate)
}

func (suite *ModelTestSuite) TestClubTransferRowCreation() {
	row := ClubTransferRow{
		MemberID:       "67890",
		FobNumber:      "FOB002",
		FirstName:      "Jane",
		LastName:       "Smith",
		MembershipType: "Standard",
		HomeClub:       "club c",
		TargetClub:     "club d",
	}

	assert.Equal(suite.T(), "67890", row.MemberID)
	assert.Equal(suite.T(), "FOB002", row.FobNumber)
	assert.Equal(suite.T(), "Jane", row.FirstName)
	assert.Equal(suite.T(), "Smith", row.LastName)
	assert.Equal(suite.T(), "Standard", row.MembershipType)
	assert.Equal(suite.T(), "club c", row.HomeClub)
	assert.Equal(suite.T(), "club d", row.TargetClub)
}

func (suite *ModelTestSuite) TestLocationCreation() {
	location := Location{
		ID:    "loc123",
		Name:  "Test Club",
		Email: "test@club.com",
	}

	assert.Equal(suite.T(), "loc123", location.ID)
	assert.Equal(suite.T(), "Test Club", location.Name)
	assert.Equal(suite.T(), "test@club.com", location.Email)
}

func (suite *ModelTestSuite) TestClubTransferDataZeroValues() {
	var data ClubTransferData

	assert.Empty(suite.T(), data.MemberID)
	assert.Empty(suite.T(), data.FobNumber)
	assert.Empty(suite.T(), data.FirstName)
	assert.Empty(suite.T(), data.LastName)
	assert.Empty(suite.T(), data.MembershipType)
	assert.Empty(suite.T(), data.HomeClub)
	assert.Empty(suite.T(), data.TargetClub)
	assert.Empty(suite.T(), data.TransferType)
	assert.True(suite.T(), data.TransferDate.IsZero())
}

func (suite *ModelTestSuite) TestLocationZeroValues() {
	var location Location

	assert.Empty(suite.T(), location.ID)
	assert.Empty(suite.T(), location.Name)
	assert.Empty(suite.T(), location.Email)
}

func (suite *ModelTestSuite) TestClubTransferRowZeroValues() {
	var row ClubTransferRow

	assert.Empty(suite.T(), row.MemberID)
	assert.Empty(suite.T(), row.FobNumber)
	assert.Empty(suite.T(), row.FirstName)
	assert.Empty(suite.T(), row.LastName)
	assert.Empty(suite.T(), row.MembershipType)
	assert.Empty(suite.T(), row.HomeClub)
	assert.Empty(suite.T(), row.TargetClub)
}

func (suite *ModelTestSuite) TestTransferTypes() {
	// Test common transfer types
	transferIn := ClubTransferData{TransferType: "TRANSFER IN"}
	transferOut := ClubTransferData{TransferType: "TRANSFER OUT"}

	assert.Equal(suite.T(), "TRANSFER IN", transferIn.TransferType)
	assert.Equal(suite.T(), "TRANSFER OUT", transferOut.TransferType)
}

func (suite *ModelTestSuite) TestMembershipTypes() {
	// Test common membership types
	premium := ClubTransferData{MembershipType: "Premium"}
	standard := ClubTransferData{MembershipType: "Standard"}
	basic := ClubTransferData{MembershipType: "Basic"}

	assert.Equal(suite.T(), "Premium", premium.MembershipType)
	assert.Equal(suite.T(), "Standard", standard.MembershipType)
	assert.Equal(suite.T(), "Basic", basic.MembershipType)
}

func TestModelSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}
