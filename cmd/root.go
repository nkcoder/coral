/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"coral.daniel-guo.com/internal/config"
	"coral.daniel-guo.com/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "club-transfer",
	Short: "Club transfer email notification tool",
	Long:  "A CLI application for processing club transfer data and sending notification emails to clubs.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize configuration
		if err := config.Init(); err != nil {
			logger.Error("Failed to initialize configuration: %v", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.AddCommand(sendEmailCmd)
}
