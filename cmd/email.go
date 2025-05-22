package cmd

import (
	"fmt"
	"os"

	"coral.daniel-guo.com/internal/clubtransfer"
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
		fmt.Printf("Transfer type: %s, filename: %s, sender: %s, env: %s\n", transferType, input, sender, env)

		if err := clubtransfer.Process(transferType, input, sender, env); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var (
	transferType string
	input        string
	sender       string
	env          string
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
}
