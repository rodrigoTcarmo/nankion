package properties

import (
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/rodrigoTcarmo/nankion/pkg/models/ofx"
)

func TransactionDateProperty(date time.Time) notion.DatabasePageProperty {
	notionDate := notion.NewDateTime(date, false)
	return notion.DatabasePageProperty{
		Date: &notion.Date{
			Start: notionDate,
		},
	}
}

func OperationProperty(operationType ofx.OFXOperationType) notion.DatabasePageProperty {
	return notion.DatabasePageProperty{
		Select: &notion.SelectOptions{
			Name: string(operationType),
		},
	}
}

func DestinationProperty(destination string) notion.DatabasePageProperty {
	return notion.DatabasePageProperty{
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
		Number: &amount,
	}
}

func MemoProperty(memo string) notion.DatabasePageProperty {
	return notion.DatabasePageProperty{
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
		RichText: []notion.RichText{
			{
				Type: "text",
				Text: &notion.Text{Content: transactionID},
			},
		},
	}
}
