package statement

import (
	"fmt"
	"time"

	"github.com/rodrigoTcarmo/nankion/pkg/models/ofx"
)

type Report struct {
	Statements   []Statement
	FinalBalance float64
}

// Statement represents a bank transaction from an OFX statement.
type Statement struct {
	TransactionDate time.Time
	Operation       ofx.OFXOperationType
	Destination     string
	Amount          float64
	Memo            string
	TransactionID   string
}

// LoadReport reads an OFX file and returns a Report struct.
func LoadReport(filepath string) (*Report, error) {
	response, err := ReadOFX(filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to read OFX file %s: %w", filepath, err)
	}
	report, err := CreateReport(response)
	if err != nil {
		return nil, fmt.Errorf("unable to parse OFX file %s: %w", filepath, err)
	}
	return report, nil
}
