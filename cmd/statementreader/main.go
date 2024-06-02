package main

import (
	"cmd/statement-reader/pkg/statementreader"
	"fmt"
)

func main() {
	fmt.Println("Nankion app!")
	fmt.Println("Readind CSV file...")

	allRecords := statementreader.ReadCsvFile("extrato.csv")

	for _, v := range allRecords {
		fmt.Println(v)
		for _, s := range v {
			fmt.Println(s)
		}
	}
}
