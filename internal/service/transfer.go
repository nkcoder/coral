// Package transfer provides club transfer functionality
package service

import (
	"fmt"
	"time"

	"coral.daniel-guo.com/internal/config"
	"coral.daniel-guo.com/internal/csvutil"
	"coral.daniel-guo.com/internal/email"
	"coral.daniel-guo.com/internal/logger"
	"coral.daniel-guo.com/internal/model"
	"coral.daniel-guo.com/internal/repository"
	"coral.daniel-guo.com/internal/secrets"
)

// Service handles club transfer operations
type Service struct {
	config         *config.AppConfig
	secretsManager *secrets.Manager
	emailSender    *email.Sender
}

// NewService creates a new transfer service
func NewService(cfg *config.AppConfig) *Service {
	return &Service{
		config:         cfg,
		secretsManager: secrets.NewManager(cfg.Secrets),
		emailSender:    email.NewSender(cfg.Email),
	}
}

// TransferRequest contains parameters for processing club transfers
type TransferRequest struct {
	TransferType string
	FileName     string
}

// Process handles the club transfer workflow
func (s *Service) Process(req TransferRequest) error {
	// Setup database connection pool
	dbConfig := repository.PoolConfig{
		Environment:    s.config.Environment,
		SecretsManager: s.secretsManager,
	}

	db, err := repository.NewPool(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	logger.Info("Starting club transfer process for type: %s", req.TransferType)

	// Read club transfer data from CSV file
	data, err := s.readClubTransferData(req.FileName)
	if err != nil {
		return fmt.Errorf("failed to read club transfer data: %w", err)
	}
	logger.Info("Successfully read club transfer data from %s", req.FileName)

	// Send emails to clubs
	if err := s.sendEmailToClubs(data, db, req.TransferType); err != nil {
		return fmt.Errorf("failed to send emails to clubs: %w", err)
	}

	logger.Info("Club transfer process completed successfully")
	return nil
}

// readClubTransferData reads the club transfer data from the CSV file
func (s *Service) readClubTransferData(fileName string) (map[string][]model.ClubTransferData, error) {
	// Read CSV and parse data
	clubTransferRows, err := csvutil.ReadClubTransferCSV(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading club transfer data: %w", err)
	}

	transfers := make(map[string][]model.ClubTransferData)
	for _, row := range clubTransferRows {
		transferIn := model.ClubTransferData{
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
			transfers[row.TargetClub] = []model.ClubTransferData{}
		}
		transfers[row.TargetClub] = append(transfers[row.TargetClub], transferIn)

		if _, exists := transfers[row.HomeClub]; !exists {
			transfers[row.HomeClub] = []model.ClubTransferData{}
		}
		transfers[row.HomeClub] = append(transfers[row.HomeClub], transferOut)
	}

	return transfers, nil
}

// getOutputFileName generates the output file name based on payment type and club name
func (s *Service) getOutputFileName(transferType, clubName string) string {
	if transferType == "DD" {
		return fmt.Sprintf("dd_club_transfer_%s.csv", clubName)
	}
	return fmt.Sprintf("pif_club_transfer_%s.csv", clubName)
}

// sendEmailToClubs sends emails to clubs with their transfer data
func (s *Service) sendEmailToClubs(data map[string][]model.ClubTransferData, db *repository.Pool, transferType string) error {
	// Get current month and year information for email subject/content
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0).Month().String()
	currentYear := now.Year()

	var subject, bodyContent string
	if transferType == "PIF" {
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
	locationRepo := repository.NewLocationRepository(db)

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

		recipientEmail := location.Email
		logger.Debug("Location email for %s: %s", clubName, recipientEmail)

		// Generate CSV content in memory
		csvContent, err := csvutil.GenerateCSVContent(data[clubName])
		if err != nil {
			logger.Error("Error generating CSV content for club %s: %v", clubName, err)
			continue
		}

		// Get attachment filename
		attachmentName := s.getOutputFileName(transferType, clubName)

		// Determine recipient email
		if s.config.TestEmail != "" {
			logger.Info("Using test email %s instead of club email %s", s.config.TestEmail, recipientEmail)
			recipientEmail = s.config.TestEmail
		}

		// Send email with in-memory attachment
		if err := s.emailSender.SendWithAttachment(
			s.config.DefaultSender,
			recipientEmail,
			subject,
			body,
			attachmentName,
			csvContent,
		); err != nil {
			logger.Error("Error sending email for club %s: %v", clubName, err)
			continue
		}

		logger.Info("Email sent successfully to club: %s", clubName)
		time.Sleep(1 * time.Second) // Sleep to avoid overwhelming email service
	}

	return nil
}
