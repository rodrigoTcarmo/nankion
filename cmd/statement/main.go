package main

import (
	"fmt"
	"log"

	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

func main() {
	report, err := statement.LoadReport("")
	if err != nil {
		log.Fatalf("Error loading statements: %v", err)
	}
	fmt.Printf("Loaded %d statements\n\n", len(report.Statements))
	for _, s := range report.Statements {
		fmt.Printf("Date: %s | Type: %s | Amount: %.2f\n",
			s.TransactionDate.Format("02/01/2006"),
			s.Operation,
			s.Amount,
		)
		if s.Destination != "" {
			fmt.Printf("  To: %s\n", s.Destination)
		}
		if s.Memo != "" {
			fmt.Printf("  Memo: %s\n", s.Memo)
		}
		fmt.Println()
	}
}
