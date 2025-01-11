package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/qwerqy/ynab-csv-converter/internal/bank"
)

func main() {
	// Input and output file paths
	inputFile := "TransactionHistory.csv"
	outputFile := "ynab_formatted.csv"

	b := bank.NewBank()

	// Open the input file
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Read the CSV content
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable-length rows
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	// Create the output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer outFile.Close()

	// Write cleaned data to the new file
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Add YNAB headers
	ynabHeaders := []string{"Date", "Payee", "Memo", "Amount"}
	if err := writer.Write(ynabHeaders); err != nil {
		fmt.Printf("Error writing headers: %v\n", err)
		return
	}

	// Process each row and write to the new file
	b.HSBC.Credit.Process(writer, rows)

	fmt.Println("Conversion complete. Output saved to:", outputFile)
}
