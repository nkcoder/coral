package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendEmailCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "should have correct usage",
			args: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic command structure
			assert.Equal(t, "send-email", sendEmailCmd.Use)
			assert.Equal(t, "Send club transfer emails", sendEmailCmd.Short)
			assert.Contains(t, sendEmailCmd.Long, "Send club transfer notification emails")
		})
	}
}

func TestSendEmailCmdFlags(t *testing.T) {
	t.Run("should have all required flags", func(t *testing.T) {
		// Test that all flags are defined
		typeFlag := sendEmailCmd.Flags().Lookup("type")
		require.NotNil(t, typeFlag)
		assert.Equal(t, "t", typeFlag.Shorthand)

		inputFlag := sendEmailCmd.Flags().Lookup("input")
		require.NotNil(t, inputFlag)
		assert.Equal(t, "i", inputFlag.Shorthand)

		senderFlag := sendEmailCmd.Flags().Lookup("sender")
		require.NotNil(t, senderFlag)
		assert.Equal(t, "s", senderFlag.Shorthand)

		envFlag := sendEmailCmd.Flags().Lookup("env")
		require.NotNil(t, envFlag)
		assert.Equal(t, "e", envFlag.Shorthand)

		testEmailFlag := sendEmailCmd.Flags().Lookup("test-email")
		require.NotNil(t, testEmailFlag)

		verboseFlag := sendEmailCmd.Flags().Lookup("verbose")
		require.NotNil(t, verboseFlag)
		assert.Equal(t, "v", verboseFlag.Shorthand)
	})
}

func TestSendEmailCmdValidation(t *testing.T) {
	// Create a temporary CSV file for testing
	tempDir, err := os.MkdirTemp("", "test-csv")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	csvFile := filepath.Join(tempDir, "test.csv")
	csvContent := `Member Id,Fob Number,First Name,Last Name,Membership Type,Home Club,Target Club
12345,FOB001,John,Doe,Premium,CLUB A,CLUB B`

	err = os.WriteFile(csvFile, []byte(csvContent), 0644)
	require.NoError(t, err)

	t.Run("should initialize command properly", func(t *testing.T) {
		// Reset flags
		typeFlag = ""
		inputFlag = ""
		senderFlag = ""
		envFlag = ""
		testEmailFlag = ""
		verboseFlag = false

		// Test init function
		assert.NotPanics(t, func() {
			// The init function is already called, so we just verify the command exists
			assert.NotNil(t, sendEmailCmd.Run)
		})
	})
}
