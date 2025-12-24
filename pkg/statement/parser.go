package statement

import (
	"fmt"
	"strings"
	"time"

	"github.com/aclindsa/ofxgo"
)

// ParseOFXDate parses an OFX date string in "YYYYMMDD" or "YYYYMMDDHHMMSS" format.
func ParseOFXDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	if len(s) >= 14 {
		return time.Parse("20060102150405", s[:14])
	}
	if len(s) >= 8 {
		return time.Parse("20060102", s[:8])
	}

	return time.Time{}, &time.ParseError{Layout: "20060102", Value: s}
}

func parseTransactions(transactions []ofxgo.Transaction) ([]Statement, error) {
	var statements []Statement
	// Parse each transaction
	for _, transaction := range transactions {
		date, err := ParseOFXDate(transaction.DtPosted.String())
		if err != nil {
			return nil, fmt.Errorf("invalid date %q: %w", transaction.DtPosted.String(), err)
		}

		amount, _ := transaction.TrnAmt.Float64()

		statements = append(statements, Statement{
			TransactionDate: date,
			Operation:       transaction.TrnType.String(),
			Destination:     transaction.Name.String(),
			Amount:          amount,
			Memo:            transaction.Memo.String(),
			TransactionID:   transaction.FiTID.String(),
		})
	}
	return sortStatements(statements), nil
}
