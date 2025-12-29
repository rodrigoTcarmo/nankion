package properties

import (
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/rodrigoTcarmo/nankion/pkg/models/ofx"
)

func TransactionDateProperty(date time.Time) notion.DatabasePageProperty {
	notionDate := notion.NewDateTime(date, false)
	tz := "America/Sao_Paulo"
	return notion.DatabasePageProperty{
		ID:   "transaction_date",
		Name: "Transaction Date",
		Type: notion.DBPropTypeDate,
		Date: &notion.Date{
			Start:    notionDate,
			End:      &notionDate,
			TimeZone: &tz,
		},
	}
}

func OperationProperty(operationType ofx.OFXOperationType) notion.DatabasePageProperty {
	return notion.DatabasePageProperty{
		ID:   "operation",
		Name: "Operation",
		Type: notion.DBPropTypeSelect,
		Select: &notion.SelectOptions{
			Name: string(operationType),
		},
	}
}

func DestinationProperty(destination string) notion.DatabasePageProperty {
	return notion.DatabasePageProperty{
		ID:   "destination",
		Name: "Destination",
		Type: notion.DBPropTypeRichText,
		RichText: []notion.RichText{
			{
				Type: "text",
				Text: &notion.Text{Content: destination},
			},
		},
	}
}

func AmountProperty(amount float64) notion.DatabasePageProperty {
	return notion.DatabasePageProperty{
		ID:   "amount",
		Name: "Amount",
		Type: notion.DBPropTypeNumber,
		Number: &amount,
	}
}

func MemoProperty(memo string) notion.DatabasePageProperty {
	return notion.DatabasePageProperty{
		ID:   "memo",
		Name: "Memo",
		Type: notion.DBPropTypeRichText,
		RichText: []notion.RichText{
			{
				Type: "text",
				Text: &notion.Text{Content: memo},
			},
		},
	}
}

func TransactionIDProperty(transactionID string) notion.DatabasePageProperty {
	return notion.DatabasePageProperty{
		ID:   "transaction_id",
		Name: "Transaction ID",
		Type: notion.DBPropTypeRichText,
		RichText: []notion.RichText{
			{
				Type: "text",
				Text: &notion.Text{Content: transactionID},
			},
		},
	}
}
