package domain

import "time"

type Location struct {
	ID    string
	Name  string
	Email string
}

type ClubTransferRow struct {
	MemberId       string
	FobNumber      string
	FirstName      string
	LastName       string
	MembershipType string
	HomeClub       string
	TargetClub     string
}

type ClubTransferData struct {
	MemberId       string
	FobNumber      string
	FirstName      string
	LastName       string
	MembershipType string
	HomeClub       string
	TargetClub     string
	TransferType   string
	TransferDate   time.Time
}
