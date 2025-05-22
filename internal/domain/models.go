package domain

import "time"

// Location represents a club location with its contact information
type Location struct {
	ID    string
	Name  string
	Email string
}

// ClubTransferRow represents a raw row from the CSV input file
type ClubTransferRow struct {
	MemberID       string `csv:"Member Id"`
	FobNumber      string `csv:"Fob Number"`
	FirstName      string `csv:"First Name"`
	LastName       string `csv:"Last Name"`
	MembershipType string `csv:"Membership Type"`
	HomeClub       string `csv:"Home Club"`
	TargetClub     string `csv:"Target Club"`
}

// ClubTransferData represents processed transfer data ready for output
type ClubTransferData struct {
	MemberID       string
	FobNumber      string
	FirstName      string
	LastName       string
	MembershipType string
	HomeClub       string
	TargetClub     string
	TransferType   string
	TransferDate   time.Time
}
