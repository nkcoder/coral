package cmd

import (
	"fmt"

	"coral.daniel-guo.com/internal/clubtransfer"
	"github.com/spf13/cobra"
)

var sendEmailCmd = &cobra.Command{
	Use:   "send-email",
	Short: "Send an email",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sending email...")
		fmt.Printf("Transfer type: %s, filename: %s, sender: %s, env: %s\n", transferType, input, sender, env)
		clubtransfer.Process(transferType, input, sender, env)
	},
}

var (
	transferType string
	input        string
	sender       string
	env          string
)

func init() {
	sendEmailCmd.Flags().StringVarP(&transferType, "type", "t", "", "Club transfer type: PIF or DD")
	sendEmailCmd.MarkFlagRequired("type")
	sendEmailCmd.Flags().StringVarP(&input, "input", "i", "", "The input file name")
	sendEmailCmd.MarkFlagRequired("input")
	sendEmailCmd.Flags().StringVarP(&sender, "sender", "s", "", "The email address to send the email from")
	sendEmailCmd.MarkFlagRequired("sender")
	sendEmailCmd.Flags().StringVarP(&env, "env", "e", "", "The environment to send the email to")
	sendEmailCmd.MarkFlagRequired("env")
}
