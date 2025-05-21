package clubtransfer

import (
	"fmt"
	"log"
	"time"

	"coral.daniel-guo.com/internal/aws"
	"coral.daniel-guo.com/internal/config"
	"coral.daniel-guo.com/internal/db"
	"coral.daniel-guo.com/internal/domain"
)

func Process(transferType string, fileName string, sender string, env string) {
	dbConfig, err := config.LoadDBConfig(env)
	if err != nil {
		log.Fatalf("Failed to load database configuration: %v", err)
	}

	// Setup database connection pool
	db, err := db.NewPool(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Read club transfer data from CSV file
	data, err := readClubTransferData(fileName)
	if err != nil {
		log.Fatalf("Failed to read club transfer data: %v", err)
	}

	// Write club transfer data to CSV files for each club
	if err := writeClubTransferData(data, transferType); err != nil {
		log.Fatalf("Failed to write club transfer data: %v", err)
	}

	// Send emails to clubs
	if err := sendEmailToClub(data, db, transferType, sender); err != nil {
		log.Fatalf("Failed to send emails to clubs: %v", err)
	}

	fmt.Println("Club transfer process completed successfully")
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
func sendEmailToClub(data map[string][]domain.ClubTransferData, db *db.Pool, transferType string, sender string) error {
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
	locationRepo := domain.NewLocationRepository(db)

	clubs := make([]string, 0, len(data))
	for club := range data {
		clubs = append(clubs, club)
	}

	fmt.Printf("Total: %d clubs\n", len(clubs))

	for _, clubName := range clubs {
		fmt.Printf("Processing club: %s\n", clubName)

		location, err := locationRepo.FindByName(clubName)
		if err != nil {
			fmt.Printf("Error finding location for club %s: %v\n", clubName, err)
			continue
		}

		if location == nil {
			fmt.Printf("--- Location not found for club: %s ---\n", clubName)
			continue
		}

		if location.Email == "" {
			fmt.Printf("--- Email not found for club: %s ---\n", clubName)
			continue
		}

		email := location.Email
		fmt.Printf("Location email: %s\n", email)

		clubTransferFile := getOutputFileName(transferType, clubName)

		// For testing purpose, send to a specific email
		toDaniel := "daniel.guo@vivalabs.com.au"
		if err := aws.SendEmailWithAttachment(sender, toDaniel, subject, body, clubTransferFile); err != nil {
			fmt.Printf("Error sending email for club %s: %v\n", clubName, err)
			continue
		}

		fmt.Printf("Process club: %s completed\n", clubName)
		time.Sleep(1 * time.Second) // Sleep to avoid overwhelming email service
	}

	return nil
}
