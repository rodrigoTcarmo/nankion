package statement

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/rodrigoTcarmo/nankion/pkg/models/ofx"
)

// getTestDataPath returns the absolute path to the testdata directory
func getTestDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "testdata", filename)
}

func TestLoadReport(t *testing.T) {
	tests := []struct {
		name            string
		filePath        string
		wantErr         bool
		wantErrContains string
		wantStatements  int
		wantBalance     float64
		validateReport  func(t *testing.T, report *Report)
	}{
		{
			name:           "successfully load valid OFX file",
			filePath:       getTestDataPath("valid.ofx"),
			wantErr:        false,
			wantStatements: 2,
			wantBalance:    1400.00,
			validateReport: func(t *testing.T, report *Report) {
				// Statements should be sorted by date (oldest first)
				if len(report.Statements) != 2 {
					t.Fatalf("expected 2 statements, got %d", len(report.Statements))
				}

				// First statement (oldest - 2025-12-10)
				stmt1 := report.Statements[0]
				expectedDate1 := time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC)
				if !stmt1.TransactionDate.Equal(expectedDate1) {
					t.Errorf("statement 1 date = %v, want %v", stmt1.TransactionDate, expectedDate1)
				}
				if stmt1.Operation != ofx.OFXOperationType("PAYMENT") {
					t.Errorf("statement 1 operation = %v, want PAYMENT", stmt1.Operation)
				}
				if stmt1.Destination != "Test Destination" {
					t.Errorf("statement 1 destination = %v, want 'Test Destination'", stmt1.Destination)
				}
				if stmt1.Amount != -100.00 {
					t.Errorf("statement 1 amount = %v, want -100.00", stmt1.Amount)
				}
				if stmt1.Memo != "Test Payment" {
					t.Errorf("statement 1 memo = %v, want 'Test Payment'", stmt1.Memo)
				}
				if stmt1.TransactionID != "202512100001" {
					t.Errorf("statement 1 transaction ID = %v, want '202512100001'", stmt1.TransactionID)
				}

				// Second statement (newer - 2025-12-15)
				stmt2 := report.Statements[1]
				expectedDate2 := time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)
				if !stmt2.TransactionDate.Equal(expectedDate2) {
					t.Errorf("statement 2 date = %v, want %v", stmt2.TransactionDate, expectedDate2)
				}
				if stmt2.Operation != ofx.OFXOperationType("CREDIT") {
					t.Errorf("statement 2 operation = %v, want CREDIT", stmt2.Operation)
				}
				if stmt2.Destination != "Employer Inc" {
					t.Errorf("statement 2 destination = %v, want 'Employer Inc'", stmt2.Destination)
				}
				if stmt2.Amount != 500.00 {
					t.Errorf("statement 2 amount = %v, want 500.00", stmt2.Amount)
				}
				if stmt2.Memo != "Salary Deposit" {
					t.Errorf("statement 2 memo = %v, want 'Salary Deposit'", stmt2.Memo)
				}
				if stmt2.TransactionID != "202512150002" {
					t.Errorf("statement 2 transaction ID = %v, want '202512150002'", stmt2.TransactionID)
				}
			},
		},
		{
			name:            "file does not exist",
			filePath:        "/nonexistent/path/to/file.ofx",
			wantErr:         true,
			wantErrContains: "unable to read OFX file",
		},
		{
			name:            "invalid OFX file content",
			filePath:        getTestDataPath("invalid.ofx"),
			wantErr:         true,
			wantErrContains: "unable to read OFX file",
		},
		{
			name:            "empty file path",
			filePath:        "",
			wantErr:         true,
			wantErrContains: "unable to read OFX file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := LoadReport(tt.filePath)

			// Check error expectations
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.wantErrContains != "" && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("expected error to contain %q, got: %v", tt.wantErrContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Validate report structure
			if report == nil {
				t.Fatal("expected report, got nil")
			}

			if len(report.Statements) != tt.wantStatements {
				t.Errorf("expected %d statements, got %d", tt.wantStatements, len(report.Statements))
			}

			if report.FinalBalance != tt.wantBalance {
				t.Errorf("expected balance %v, got %v", tt.wantBalance, report.FinalBalance)
			}

			// Run custom validation if provided
			if tt.validateReport != nil {
				tt.validateReport(t, report)
			}
		})
	}
}
