package cmd

import (
	"os"

	"coral.daniel-guo.com/internal/config"
	"coral.daniel-guo.com/internal/logger"
	"coral.daniel-guo.com/internal/service"
	"github.com/spf13/cobra"
)

// sendEmailCmd represents the send-email command for sending club transfer emails
var sendEmailCmd = &cobra.Command{
	Use:   "send-email",
	Short: "Send club transfer emails",
	Long: `Send club transfer notification emails to clubs.
This command processes club transfer data from a CSV file and sends 
personalized emails to each club with their relevant transfer information.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set logging level based on verbose flag
		if verboseFlag {
			logger.SetLevel(logger.DebugLevel)
			logger.Debug("Debug logging enabled")
		}

		logger.Info("Transfer type: %s, filename: %s, env: %s",
			typeFlag, inputFlag, envFlag)

		// Load application configuration
		appConfig := config.NewAppConfig(envFlag, testEmailFlag, senderFlag)

		// Create transfer service
		transferService := service.NewService(appConfig)

		// Create transfer request
		req := service.TransferRequest{
			TransferType: typeFlag,
			FileName:     inputFlag,
		}

		// Process the request
		if err := transferService.Process(req); err != nil {
			logger.Error("Failed to process club transfers: %v", err)
			os.Exit(1)
		}
	},
}

var (
	typeFlag      string
	inputFlag     string
	senderFlag    string
	envFlag       string
	testEmailFlag string
	verboseFlag   bool
)

func init() {
	sendEmailCmd.Flags().StringVarP(&typeFlag, "type", "t", "", "Club transfer type: PIF (Paid in Full) or DD (Direct Debit)")

	sendEmailCmd.Flags().StringVarP(&inputFlag, "input", "i", "", "CSV input file with transfer data")

	sendEmailCmd.Flags().StringVarP(&senderFlag, "sender", "s", "", "Sender email address")
	sendEmailCmd.Flags().StringVarP(&envFlag, "env", "e", "", "Environment (dev, staging, prod)")

	sendEmailCmd.Flags().StringVarP(&testEmailFlag, "test-email", "", "", "Test email address (if set, all emails go here instead of to clubs)")

	sendEmailCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose debugging output")
}
