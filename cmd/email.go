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
		if verbose {
			logger.SetLevel(logger.DebugLevel)
			logger.Debug("Debug logging enabled")
		}

		logger.Info("Transfer type: %s, filename: %s, sender: %s, env: %s", transferType, input, sender, env)

		// Create app configuration
		appConfig := config.NewAppConfig(env).WithSender(sender).WithTestEmail(testEmail)

		// Create transfer service
		transferService := service.NewService(appConfig)

		// Create transfer request
		req := service.TransferRequest{
			TransferType: transferType,
			FileName:     input,
		}

		// Process the request
		if err := transferService.Process(req); err != nil {
			logger.Error("Failed to process club transfers: %v", err)
			os.Exit(1)
		}
	},
}

var (
	transferType string
	input        string
	sender       string
	env          string
	testEmail    string
	verbose      bool
)

func init() {
	sendEmailCmd.Flags().StringVarP(&transferType, "type", "t", "", "Club transfer type: PIF (Paid in Full) or DD (Direct Debit)")
	sendEmailCmd.MarkFlagRequired("type")

	sendEmailCmd.Flags().StringVarP(&input, "input", "i", "", "CSV input file with transfer data")
	sendEmailCmd.MarkFlagRequired("input")

	sendEmailCmd.Flags().StringVarP(&sender, "sender", "s", "", "Sender email address")
	sendEmailCmd.MarkFlagRequired("sender")

	sendEmailCmd.Flags().StringVarP(&env, "env", "e", "", "Environment (dev, staging, prod)")
	sendEmailCmd.MarkFlagRequired("env")

	sendEmailCmd.Flags().StringVarP(&testEmail, "test-email", "", "", "Test email address (if set, all emails go here instead of to clubs)")

	sendEmailCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose debugging output")
}
