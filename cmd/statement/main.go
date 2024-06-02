package main

import (
	"fmt"

	"github.com/rodrigoTcarmo/nankion/pkg/notion"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

func main() {
	fmt.Println("Nankion app!")
	fmt.Println("Readind CSV file...")

	allRecords := statement.ReadCsvFile("extrato.csv")

	for _, v := range allRecords {
		fmt.Println(v)
		for _, s := range v {
			fmt.Println(s)
		}
	}

	notion.PagePrinter()
}
