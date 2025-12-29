package page

import (
	"testing"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/rodrigoTcarmo/nankion/pkg/models/ofx"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

func TestBuildPageProperties(t *testing.T) {
	// expectedTitle computes the expected title from a statement
	expectedTitle := func(stmt statement.Statement) string {
		return string(stmt.Operation) + "-" + stmt.Destination + "-" + stmt.TransactionID
	}

	tests := []struct {
		name      string
		statement statement.Statement
	}{
		{
			name: "build properties from complete statement",
			statement: statement.Statement{
				TransactionDate: time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
				Operation:       ofx.OFXOperationTypeDebit,
				Destination:     "Amazon Store",
				Amount:          -150.99,
				Memo:            "Online purchase",
				TransactionID:   "TXN123456",
			},
		},
		{
			name: "build properties for credit transaction",
			statement: statement.Statement{
				TransactionDate: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
				Operation:       ofx.OFXOperationTypeCredit,
				Destination:     "Salary Deposit",
				Amount:          5000.00,
				Memo:            "Monthly salary",
				TransactionID:   "SAL202511",
			},
		},
		{
			name: "build properties with empty memo",
			statement: statement.Statement{
				TransactionDate: time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),
				Operation:       ofx.OFXOperationTypePOS,
				Destination:     "Coffee Shop",
				Amount:          -5.50,
				Memo:            "",
				TransactionID:   "POS789",
			},
		},
		{
			name: "build properties for ATM withdrawal",
			statement: statement.Statement{
				TransactionDate: time.Date(2025, 9, 5, 0, 0, 0, 0, time.UTC),
				Operation:       ofx.OFXOperationTypeATM,
				Destination:     "ATM Withdrawal",
				Amount:          -200.00,
				Memo:            "Cash withdrawal",
				TransactionID:   "ATM001",
			},
		},
		{
			name: "build properties with zero amount",
			statement: statement.Statement{
				TransactionDate: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
				Operation:       ofx.OFXOperationTypeOther,
				Destination:     "Fee Reversal",
				Amount:          0.00,
				Memo:            "Reversed fee",
				TransactionID:   "REV999",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &page{}
			stmt := test.statement

			result, err := p.BuildPageProperties(stmt)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("expected result, got nil")
			}

			// Verify Name (Title) property - computed from statement fields
			nameProperty := (*result)["Name"]
			if len(nameProperty.Title) == 0 {
				t.Error("expected Name property to have title, got empty")
			} else if got, want := nameProperty.Title[0].Text.Content, expectedTitle(stmt); got != want {
				t.Errorf("Name: expected %q, got %q", want, got)
			}

			// Verify Transaction Date property
			dateProperty := (*result)["Transaction Date"]
			if dateProperty.Date == nil {
				t.Error("expected Transaction Date property to have date, got nil")
			} else if got := dateProperty.Date.Start.Time; !got.Equal(stmt.TransactionDate) {
				t.Errorf("Transaction Date: expected %v, got %v", stmt.TransactionDate, got)
			}

			// Verify Operation property
			opProperty := (*result)["Operation"]
			if opProperty.Select == nil {
				t.Error("expected Operation property to have select, got nil")
			} else if got, want := opProperty.Select.Name, string(stmt.Operation); got != want {
				t.Errorf("Operation: expected %q, got %q", want, got)
			}

			// Verify Destination property
			destProperty := (*result)["Destination"]
			if len(destProperty.RichText) == 0 {
				t.Error("expected Destination property to have rich text, got empty")
			} else if got := destProperty.RichText[0].Text.Content; got != stmt.Destination {
				t.Errorf("Destination: expected %q, got %q", stmt.Destination, got)
			}

			// Verify Amount property
			amountProperty := (*result)["Amount"]
			if amountProperty.Number == nil {
				t.Error("expected Amount property to have number, got nil")
			} else if got := *amountProperty.Number; got != stmt.Amount {
				t.Errorf("Amount: expected %f, got %f", stmt.Amount, got)
			}

			// Verify Memo property
			memoProperty := (*result)["Memo"]
			if len(memoProperty.RichText) == 0 {
				t.Error("expected Memo property to have rich text, got empty")
			} else if got := memoProperty.RichText[0].Text.Content; got != stmt.Memo {
				t.Errorf("Memo: expected %q, got %q", stmt.Memo, got)
			}

			// Verify Transaction ID property
			txIDProperty := (*result)["Transaction ID"]
			if len(txIDProperty.RichText) == 0 {
				t.Error("expected Transaction ID property to have rich text, got empty")
			} else if got := txIDProperty.RichText[0].Text.Content; got != stmt.TransactionID {
				t.Errorf("Transaction ID: expected %q, got %q", stmt.TransactionID, got)
			}
		})
	}
}

func TestCreatePageParams(t *testing.T) {
	tests := []struct {
		name       string
		pageData   PageData
		wantDbID   string
		wantParent notion.ParentType
	}{
		{
			name: "create params with valid page data",
			pageData: PageData{
				DatabaseId: "db-123456",
				Properties: &notion.DatabasePageProperties{
					"Name": notion.DatabasePageProperty{
						Title: []notion.RichText{{Text: &notion.Text{Content: "Test"}}},
					},
				},
			},
			wantDbID:   "db-123456",
			wantParent: notion.ParentTypeDatabase,
		},
		{
			name: "create params with different database id",
			pageData: PageData{
				DatabaseId: "another-db-789",
				Properties: &notion.DatabasePageProperties{},
			},
			wantDbID:   "another-db-789",
			wantParent: notion.ParentTypeDatabase,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &page{}

			result := p.CreatePageParams(test.pageData)

			if result.ParentID != test.wantDbID {
				t.Errorf("ParentID: expected %q, got %q", test.wantDbID, result.ParentID)
			}

			if result.ParentType != test.wantParent {
				t.Errorf("ParentType: expected %v, got %v", test.wantParent, result.ParentType)
			}

			if result.DatabasePageProperties != test.pageData.Properties {
				t.Error("DatabasePageProperties: expected properties to be passed through")
			}
		})
	}
}
