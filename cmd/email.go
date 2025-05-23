package cmd

import (
	"os"

	"coral.daniel-guo.com/internal/config"
	"coral.daniel-guo.com/internal/logger"
	"coral.daniel-guo.com/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		// Get configuration values from command line flags
		transferType := viper.GetString("transfer.type")
		input := viper.GetString("transfer.input")

		logger.Info("Transfer type: %s, filename: %s, env: %s",
			transferType, input, viper.GetString("env"))

		// Load application configuration
		appConfig := config.LoadConfig()

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

var verbose bool

func init() {
	// Define command-line flags with viper bindings
	sendEmailCmd.Flags().StringP("type", "t", "", "Club transfer type: PIF (Paid in Full) or DD (Direct Debit)")
	if err := viper.BindPFlag("transfer.type", sendEmailCmd.Flags().Lookup("type")); err != nil {
		logger.Error("Failed to bind flag 'type': %v", err)
		os.Exit(1)
	}
	if err := sendEmailCmd.MarkFlagRequired("type"); err != nil {
		logger.Error("Failed to mark flag 'type' as required: %v", err)
		os.Exit(1)
	}

	sendEmailCmd.Flags().StringP("input", "i", "", "CSV input file with transfer data")
	if err := viper.BindPFlag("transfer.input", sendEmailCmd.Flags().Lookup("input")); err != nil {
		logger.Error("Failed to bind flag 'input': %v", err)
		os.Exit(1)
	}
	if err := sendEmailCmd.MarkFlagRequired("input"); err != nil {
		logger.Error("Failed to mark flag 'input' as required: %v", err)
		os.Exit(1)
	}

	sendEmailCmd.Flags().StringP("sender", "s", "", "Sender email address")
	viper.BindPFlag("email.sender", sendEmailCmd.Flags().Lookup("sender"))

	sendEmailCmd.Flags().StringP("env", "e", "", "Environment (dev, staging, prod)")
	viper.BindPFlag("env", sendEmailCmd.Flags().Lookup("env"))

	sendEmailCmd.Flags().String("test-email", "", "Test email address (if set, all emails go here instead of to clubs)")
	viper.BindPFlag("email.test", sendEmailCmd.Flags().Lookup("test-email"))

	sendEmailCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose debugging output")
}
