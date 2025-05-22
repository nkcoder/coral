package clubtransfer

import (
	"fmt"
	"time"

	"coral.daniel-guo.com/internal/aws"
	"coral.daniel-guo.com/internal/db"
	"coral.daniel-guo.com/internal/domain"
	"coral.daniel-guo.com/internal/logger"
)

// Config holds the configuration for the club transfer process
type Config struct {
	TransferType string
	FileName     string
	Sender       string
	Environment  string
	TestEmail    string // Email for testing - if set, sends to this address instead of club email
}

// Process handles the club transfer workflow
func Process(cfg Config) error {
	// Setup database connection pool
	db, err := db.NewPool(cfg.Environment)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	logger.Info("Starting club transfer process for type: %s", cfg.TransferType)

	// Read club transfer data from CSV file
	data, err := readClubTransferData(cfg.FileName)
	if err != nil {
		return fmt.Errorf("failed to read club transfer data: %w", err)
	}
	logger.Info("Successfully read club transfer data from %s", cfg.FileName)

	// Write club transfer data to CSV files for each club
	if err := writeClubTransferData(data, cfg.TransferType); err != nil {
		return fmt.Errorf("failed to write club transfer data: %w", err)
	}
	logger.Info("Successfully wrote club transfer data to individual files")

	// Send emails to clubs
	if err := sendEmailToClub(data, db, cfg); err != nil {
		return fmt.Errorf("failed to send emails to clubs: %w", err)
	}

	logger.Info("Club transfer process completed successfully")
	return nil
}

// readClubTransferData reads the club transfer data from the CSV file based on payment type
func readClubTransferData(fileName string) (map[string][]domain.ClubTransferData, error) {
	// Read CSV and parse data
	clubTransferRows, err := domain.ReadClubTransferCSV(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading club transfer data: %w", err)
	}

	transfers := make(map[string][]domain.ClubTransferData)
	for _, row := range clubTransferRows {
		transferIn := domain.ClubTransferData{
			MemberID:       row.MemberID,
			FobNumber:      row.FobNumber,
			FirstName:      row.FirstName,
			LastName:       row.LastName,
			MembershipType: row.MembershipType,
			HomeClub:       row.HomeClub,
			TargetClub:     row.TargetClub,
			TransferType:   "TRANSFER IN",
			TransferDate:   time.Now(),
		}

		transferOut := transferIn
		transferOut.TransferType = "TRANSFER OUT"

		// Add transfers to appropriate clubs
		if _, exists := transfers[row.TargetClub]; !exists {
			transfers[row.TargetClub] = []domain.ClubTransferData{}
		}
		transfers[row.TargetClub] = append(transfers[row.TargetClub], transferIn)

		if _, exists := transfers[row.HomeClub]; !exists {
			transfers[row.HomeClub] = []domain.ClubTransferData{}
		}
		transfers[row.HomeClub] = append(transfers[row.HomeClub], transferOut)
	}

	return transfers, nil
}

// getOutputFileName generates the output file name based on payment type and club name
func getOutputFileName(transferType, clubName string) string {
	if transferType == "DD" {
		return fmt.Sprintf("dd_club_transfer_%s.csv", clubName)
	}
	return fmt.Sprintf("pif_club_transfer_%s.csv", clubName)
}

// writeClubTransferData writes club transfer data to CSV files for each club
func writeClubTransferData(data map[string][]domain.ClubTransferData, transferType string) error {
	for club, transfers := range data {
		clubFileName := getOutputFileName(transferType, club)
		if err := domain.WriteClubTransferCSV(clubFileName, transfers); err != nil {
			return fmt.Errorf("error writing club transfer data for %s: %w", club, err)
		}
	}
	return nil
}

// sendEmailToClub sends emails to clubs with their transfer data
func sendEmailToClub(data map[string][]domain.ClubTransferData, db *db.Pool, cfg Config) error {
	// Get current month and year information for email subject/content
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0).Month().String()
	currentYear := now.Year()

	var subject, bodyContent string
	if cfg.TransferType == "PIF" {
		subject = fmt.Sprintf("Club Transfer for Paid in Full Members (%s %d)", lastMonth, currentYear)
		bodyContent = fmt.Sprintf("Please find attached the Paid in Full club transfer data for your club (%s %d).", lastMonth, currentYear)
	} else {
		lastQuarter := now.AddDate(0, -3, 0).Month().String()
		subject = fmt.Sprintf("Club Transfer for Direct Debit Members (%s - %s %d)", lastQuarter, lastMonth, currentYear)
		bodyContent = fmt.Sprintf("Please find attached the Direct Debit club transfer data for your club (%s - %s %d).", lastQuarter, lastMonth, currentYear)
	}

	body := fmt.Sprintf(`
		<html>
		<head></head>
		<body><p>Hello team,</p>
		<p>%s</p>
		<p>Regards</p>
		</html>
  `, bodyContent)

	// Create location repository
	locationRepo := domain.NewLocationRepository(db)

	clubs := make([]string, 0, len(data))
	for club := range data {
		clubs = append(clubs, club)
	}

	logger.Info("Processing %d clubs for email delivery", len(clubs))

	for _, clubName := range clubs {
		logger.Debug("Processing club: %s", clubName)

		location, err := locationRepo.FindByName(clubName)
		if err != nil {
			logger.Warn("Error finding location for club %s: %v", clubName, err)
			continue
		}

		if location == nil {
			logger.Warn("Location not found for club: %s", clubName)
			continue
		}

		if location.Email == "" {
			logger.Warn("Email not found for club: %s", clubName)
			continue
		}

		email := location.Email
		logger.Debug("Location email for %s: %s", clubName, email)

		clubTransferFile := getOutputFileName(cfg.TransferType, clubName)

		// Determine recipient email
		recipient := email
		if cfg.TestEmail != "" {
			logger.Info("Using test email %s instead of club email %s", cfg.TestEmail, email)
			recipient = cfg.TestEmail
		}

		if err := aws.SendEmailWithAttachment(cfg.Sender, recipient, subject, body, clubTransferFile); err != nil {
			logger.Error("Error sending email for club %s: %v", clubName, err)
			continue
		}

		logger.Info("Email sent successfully to club: %s", clubName)
		time.Sleep(1 * time.Second) // Sleep to avoid overwhelming email service
	}

	return nil
}
