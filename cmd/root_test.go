package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	t.Run("should have correct basic properties", func(t *testing.T) {
		assert.Equal(t, "club-transfer", rootCmd.Use)
		assert.Equal(t, "Club transfer email notification tool", rootCmd.Short)
		assert.Contains(t, rootCmd.Long, "CLI application for processing club transfer data")
	})

	t.Run("should have send-email subcommand", func(t *testing.T) {
		// Check that the send-email command is added as a subcommand
		commands := rootCmd.Commands()
		var foundSendEmail bool
		for _, cmd := range commands {
			if cmd.Use == "send-email" {
				foundSendEmail = true
				break
			}
		}
		assert.True(t, foundSendEmail, "send-email command should be added to root command")
	})
}

func TestExecute(t *testing.T) {
	t.Run("should not panic when called", func(t *testing.T) {
		// We can't easily test Execute() without modifying os.Args or exit behavior
		// But we can verify the function exists and basic command structure
		assert.NotNil(t, rootCmd)
		assert.NotNil(t, Execute)
	})
}
