package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sendEmailCmd = &cobra.Command{
	Use:   "send-email",
	Short: "Send an email",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sending email...")
		fmt.Println("Transfer type:", transferType)
		fmt.Println("Filename:", filename)
	},
}

var (
	transferType string
	filename     string
)

func init() {
	sendEmailCmd.Flags().StringVarP(&transferType, "type", "t", "", "Club transfer type: PIF or DD")
	sendEmailCmd.Flags().StringVarP(&filename, "filename", "f", "", "The input file name")
}
