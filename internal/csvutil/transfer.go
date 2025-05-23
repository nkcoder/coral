// Package csvutil provides utilities for working with CSV data
package csvutil

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"coral.daniel-guo.com/internal/model"
)

// ReadClubTransferCSV reads a CSV file with club transfer data
func ReadClubTransferCSV(fileName string) ([]model.ClubTransferRow, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			// Log the error, but don't override the main error
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", cerr)
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no records found")
	}

	headers := records[0]
	colMap := make(map[string]int)
	for i, colName := range headers {
		colMap[colName] = i
	}

	requiredCols := []string{
		"Member Id",
		"Fob Number",
		"First Name",
		"Last Name",
		"Membership Type",
		"Home Club",
		"Target Club",
	}

	for _, col := range requiredCols {
		if _, ok := colMap[col]; !ok {
			return nil, fmt.Errorf("column %s not found", col)
		}
	}

	var result []model.ClubTransferRow
	for i, record := range records {
		if i == 0 {
			continue
		}

		row := model.ClubTransferRow{
			MemberID:       record[colMap["Member Id"]],
			FobNumber:      record[colMap["Fob Number"]],
			FirstName:      record[colMap["First Name"]],
			LastName:       record[colMap["Last Name"]],
			MembershipType: record[colMap["Membership Type"]],
			HomeClub:       strings.ToUpper(record[colMap["Home Club"]]),
			TargetClub:     strings.ToUpper(record[colMap["Target Club"]]),
		}

		result = append(result, row)
	}

	return result, nil
}

// GenerateCSVContent generates CSV content in memory as []byte
func GenerateCSVContent(data []model.ClubTransferData) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{
		"Member Id",
		"Fob Number",
		"First Name",
		"Last Name",
		"Membership Type",
		"Home Club",
		"Target Club",
		"Transfer Type",
		"Transfer Date",
	}

	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write headers: %w", err)
	}

	for _, transfer := range data {
		record := []string{
			transfer.MemberID,
			transfer.FobNumber,
			transfer.FirstName,
			transfer.LastName,
			transfer.MembershipType,
			transfer.HomeClub,
			transfer.TargetClub,
			transfer.TransferType,
			transfer.TransferDate.Format("2006-01-02"),
		}

		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing csv writer: %w", err)
	}

	return buf.Bytes(), nil
}
