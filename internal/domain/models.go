package domain

import "time"

type Location struct {
	ID    string
	Name  string
	Email string
}

type ClubTransferRow struct {
	MemberID       string
	FobNumber      string
	FirstName      string
	LastName       string
	MembershipType string
	HomeClub       string
	TargetClub     string
}

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
