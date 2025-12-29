package statement

import (
	"fmt"
	"os"
	"sort"

	"github.com/aclindsa/ofxgo"
)

// ReadOFX reads and parses an OFX file, returning the parsed response.
func ReadOFX(filepath string) (*ofxgo.Response, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to open OFX file %s: %w", filepath, err)
	}
	defer f.Close()

	response, err := ofxgo.ParseResponse(f)
	if err != nil {
		return nil, fmt.Errorf("unable to parse OFX file %s: %w", filepath, err)
	}

	return response, nil
}

// CreateReport creates a Report struct from an OFX response.
func CreateReport(response *ofxgo.Response) (*Report, error) {
	var (
		statements   []Statement
		finalBalance float64
	)

	// Process bank statements
	for _, msg := range response.Bank {
		statementResponse, ok := msg.(*ofxgo.StatementResponse)
		if !ok {
			continue
		}
		// Get final balance from LEDGERBAL
		finalBalance, _ = statementResponse.BalAmt.Float64()

		var err error
		statements, err = parseTransactions(statementResponse.BankTranList.Transactions)
		if err != nil {
			return nil, fmt.Errorf("unable to parse transactions: %w", err)
		}
	}

	// Process credit card statements
	for _, msg := range response.CreditCard {
		ccStatementResponse, ok := msg.(*ofxgo.CCStatementResponse)
		if !ok {
			continue
		}
		// Get final balance from LEDGERBAL
		finalBalance, _ = ccStatementResponse.BalAmt.Float64()

		if ccStatementResponse.BankTranList != nil {
			var err error
			statements, err = parseTransactions(ccStatementResponse.BankTranList.Transactions)
			if err != nil {
				return nil, fmt.Errorf("unable to parse credit card transactions: %w", err)
			}
		}
	}

	return &Report{
		Statements:   statements,
		FinalBalance: finalBalance,
	}, nil
}

// sortStatements sorts statements by date (oldest first)
func sortStatements(statements []Statement) []Statement {
	sort.Slice(statements, func(i, j int) bool {
		return statements[i].TransactionDate.Before(statements[j].TransactionDate)
	})
	return statements
}
