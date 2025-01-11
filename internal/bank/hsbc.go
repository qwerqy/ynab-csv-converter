package bank

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

type HSBCBank struct {
	Credit interface {
		Process(*csv.Writer, [][]string) error
	}
	Debit interface {
		Process(*csv.Writer, [][]string) error
	}
}

type HSBCCredit struct{}
type HSBCDebit struct{}

func (b *HSBCCredit) Process(w *csv.Writer, rows [][]string) error {
	for _, row := range rows {
		// Map to YNAB format
		date := row[0]
		memo := row[1]
		payee := parsePayeeFromMemo(memo)
		amount := row[2]

		// Convert amount to a standardized format (removing commas, ensuring proper sign)
		cleanedAmount, err := strconv.ParseFloat(strings.ReplaceAll(amount, ",", ""), 64)
		if err != nil {
			fmt.Printf("Skipping row with invalid amount: %v\n", row)
			continue
		}

		// Write the YNAB row
		ynabRow := []string{date, payee, memo, fmt.Sprintf("%.2f", cleanedAmount)}
		if err := w.Write(ynabRow); err != nil {
			fmt.Printf("Error writing row: %v\n", err)
			continue
		}
	}

	return nil
}

func (b *HSBCDebit) Process(w *csv.Writer, rows [][]string) error {
	for _, row := range rows {
		// Map to YNAB format
		date := row[0]
		payee := row[1]
		memo := ""
		amount := row[2]

		// Convert amount to a standardized format (removing commas, ensuring proper sign)
		cleanedAmount, err := strconv.ParseFloat(strings.ReplaceAll(amount, ",", ""), 64)
		if err != nil {
			fmt.Printf("Skipping row with invalid amount: %v\n", row)
			continue
		}

		// Write the YNAB row
		ynabRow := []string{date, payee, memo, fmt.Sprintf("%.2f", cleanedAmount)}
		if err := w.Write(ynabRow); err != nil {
			fmt.Printf("Error writing row: %v\n", err)
			continue
		}
	}

	return nil
}

func parsePayeeFromMemo(memo string) string {
	parts := strings.Fields(memo)

	if parts[0] == "ATM" {
		memo = strings.Join(parts[7:len(parts)-2], " ")
	} else if parts[0] == "TRANSFER" && parts[1] == "DEBIT" {
		indexOfHIB := 0

		for index, item := range parts {
			if item == "HIB-" {
				indexOfHIB = index
				break
			}
		}

		parts = append(parts[:indexOfHIB], parts[indexOfHIB+2:]...)

		memo = strings.Join(parts[2:], " ")
	} else if parts[0] == "TRANSFER" && parts[1] == "CREDIT" {
		memo = strings.Join(parts[2:], " ")
	} else if parts[0] == "CDM" {
		memo = parts[0]
	}

	fmt.Println(memo)

	return memo
}
