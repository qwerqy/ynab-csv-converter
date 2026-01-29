package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/qwerqy/ynab-csv-converter/internal/bank"
)

func main() {
	// Prompt for input file path (loop until file exists)
	var inputFile string
	for {
		inputPrompt := &survey.Input{
			Message: "Enter the input file path:",
		}
		if err := survey.AskOne(inputPrompt, &inputFile); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}

		if _, err := os.Stat(inputFile); err == nil {
			break
		}

		fmt.Println("File not found. Please try again.")
	}

	// Prompt for output file path (optional, defaults to ynab_formatted.csv)
	var outputFile string
	outputPrompt := &survey.Input{
		Message: "Enter the output file path:",
		Default: "ynab_formatted.csv",
	}
	if err := survey.AskOne(outputPrompt, &outputFile); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	// Use default if output file is empty
	if outputFile == "" {
		outputFile = "ynab_formatted.csv"
	}

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
