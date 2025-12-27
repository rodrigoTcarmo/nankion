package notion

import "github.com/dstotijn/go-notion"

var RequiredProperties = map[string]*notion.DatabaseProperty{
	"Transaction Date": {
		ID:   "transaction_date",
		Name: "Transaction Date",
		Type: notion.DBPropTypeDate,
		Date: &notion.EmptyMetadata{},
	},
	"Operation": {
		ID:   "operation",
		Name: "Operation",
		Type: notion.DBPropTypeSelect,
		Select: &notion.SelectMetadata{Options: []notion.SelectOptions{
			{Name: "CREDIT"},
			{Name: "DEBIT"},
			{Name: "INT"},
			{Name: "DIV"},
			{Name: "FEE"},
			{Name: "SRVCHG"},
			{Name: "DEP"},
			{Name: "ATM"},
			{Name: "POS"},
			{Name: "XFER"},
			{Name: "CHECK"},
			{Name: "PAYMENT"},
			{Name: "CASH"},
			{Name: "DIRECTDEP"},
			{Name: "DIRECTDEBIT"},
			{Name: "REPEATPMT"},
			{Name: "OTHER"},
		}},
	},
	"Destination": {
		ID:       "destination",
		Name:     "Destination",
		Type:     notion.DBPropTypeRichText,
		RichText: &notion.EmptyMetadata{},
	},
	"Amount": {
		ID:     "amount",
		Name:   "Amount",
		Type:   notion.DBPropTypeNumber,
		Number: &notion.NumberMetadata{Format: notion.NumberFormatReal},
	},
	"Memo": {
		ID:       "memo",
		Name:     "Memo",
		Type:     notion.DBPropTypeRichText,
		RichText: &notion.EmptyMetadata{},
	},
	"Transaction ID": {
		ID:       "transaction_id",
		Name:     "Transaction ID",
		Type:     notion.DBPropTypeRichText,
		RichText: &notion.EmptyMetadata{},
	},
}
