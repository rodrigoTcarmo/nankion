package notion

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	databasemock "github.com/rodrigoTcarmo/nankion/pkg/notion/database/mocks"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

func TestUploadReport(t *testing.T) {
	databaseId := "1234567890"
	allProperties := notion.DatabaseProperties{
		"Transaction Date": notion.DatabaseProperty{Name: "Transaction Date", Type: notion.DBPropTypeDate},
		"Operation":        notion.DatabaseProperty{Name: "Operation", Type: notion.DBPropTypeSelect},
		"Destination":      notion.DatabaseProperty{Name: "Destination", Type: notion.DBPropTypeRichText},
		"Amount":           notion.DatabaseProperty{Name: "Amount", Type: notion.DBPropTypeNumber},
		"Memo":             notion.DatabaseProperty{Name: "Memo", Type: notion.DBPropTypeRichText},
		"Transaction ID":   notion.DatabaseProperty{Name: "Transaction ID", Type: notion.DBPropTypeRichText},
	}
	tests := []struct {
		name                   string
		envDatabaseId          string
		mockDatabase           func(string) (*notion.Database, error)
		mockDatabaseProperties func(string) (notion.DatabaseProperties, error)
		wantError              error
	}{
		{
			name:          "get all expected properties from notion",
			envDatabaseId: databaseId,
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return &notion.Database{
					ID:         databaseId,
					Properties: allProperties,
				}, nil
			},
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return allProperties, nil
			},
			wantError: nil,
		},
		{
			name:          "succeed when database has extra properties beyond required ones",
			envDatabaseId: databaseId,
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return &notion.Database{ID: databaseId}, nil
			},
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				propsWithExtra := notion.DatabaseProperties{
					"Transaction Date": notion.DatabaseProperty{Name: "Transaction Date", Type: notion.DBPropTypeDate},
					"Operation":        notion.DatabaseProperty{Name: "Operation", Type: notion.DBPropTypeSelect},
					"Destination":      notion.DatabaseProperty{Name: "Destination", Type: notion.DBPropTypeRichText},
					"Amount":           notion.DatabaseProperty{Name: "Amount", Type: notion.DBPropTypeNumber},
					"Memo":             notion.DatabaseProperty{Name: "Memo", Type: notion.DBPropTypeRichText},
					"Transaction ID":   notion.DatabaseProperty{Name: "Transaction ID", Type: notion.DBPropTypeRichText},
					"Extra Property":   notion.DatabaseProperty{Name: "Extra Property", Type: notion.DBPropTypeCheckbox},
					"Another Extra":    notion.DatabaseProperty{Name: "Another Extra", Type: notion.DBPropTypeURL},
				}
				return propsWithExtra, nil
			},
			wantError: nil,
		},
		{
			name:          "return error if DATABASE_ID env var is empty",
			envDatabaseId: "",
			wantError:     errors.New("database ID not found in environment variables"),
		},
		{
			name:          "return error if database does not exist",
			envDatabaseId: databaseId,
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return nil, errors.New("database not found")
			},
			wantError: errors.New("error trying to get database by ID: database not found"),
		},
		{
			name:          "return error if listing database properties fails",
			envDatabaseId: databaseId,
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return &notion.Database{ID: databaseId}, nil
			},
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return nil, errors.New("unauthorized access")
			},
			wantError: errors.New("error trying to list database properties: unauthorized access"),
		},
		{
			name:          "return error if all properties are missing",
			envDatabaseId: databaseId,
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return &notion.Database{
					ID: databaseId,
				}, nil
			},
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return notion.DatabaseProperties{}, nil
			},
			wantError: errors.New("missing properties: [Transaction Date Operation Destination Amount Memo Transaction ID]"),
		},
		{
			name:          "return error if Transaction ID property is missing",
			envDatabaseId: databaseId,
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return &notion.Database{
					ID: databaseId,
				}, nil
			},
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return notion.DatabaseProperties{
					"Transaction Date": notion.DatabaseProperty{Name: "Transaction Date", Type: notion.DBPropTypeDate},
					"Operation":        notion.DatabaseProperty{Name: "Operation", Type: notion.DBPropTypeSelect},
					"Destination":      notion.DatabaseProperty{Name: "Destination", Type: notion.DBPropTypeRichText},
					"Amount":           notion.DatabaseProperty{Name: "Amount", Type: notion.DBPropTypeNumber},
					"Memo":             notion.DatabaseProperty{Name: "Memo", Type: notion.DBPropTypeRichText},
				}, nil
			},
			wantError: errors.New("missing properties: [Transaction ID]"),
		},
		{
			name:          "return error if multiple properties are missing",
			envDatabaseId: databaseId,
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return &notion.Database{ID: databaseId}, nil
			},
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return notion.DatabaseProperties{
					"Transaction Date": notion.DatabaseProperty{Name: "Transaction Date", Type: notion.DBPropTypeDate},
					"Operation":        notion.DatabaseProperty{Name: "Operation", Type: notion.DBPropTypeSelect},
					"Destination":      notion.DatabaseProperty{Name: "Destination", Type: notion.DBPropTypeRichText},
				}, nil
			},
			wantError: errors.New("missing properties: [Amount Memo Transaction ID]"),
		},
		{
			name:          "return error if property exists but with wrong type",
			envDatabaseId: databaseId,
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return &notion.Database{ID: databaseId}, nil
			},
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return notion.DatabaseProperties{
					"Transaction Date": notion.DatabaseProperty{Name: "Transaction Date", Type: notion.DBPropTypeDate},
					"Operation":        notion.DatabaseProperty{Name: "Operation", Type: notion.DBPropTypeSelect},
					"Destination":      notion.DatabaseProperty{Name: "Destination", Type: notion.DBPropTypeRichText},
					"Amount":           notion.DatabaseProperty{Name: "Amount", Type: notion.DBPropTypeRichText}, // Should be Number
					"Memo":             notion.DatabaseProperty{Name: "Memo", Type: notion.DBPropTypeRichText},
					"Transaction ID":   notion.DatabaseProperty{Name: "Transaction ID", Type: notion.DBPropTypeRichText},
				}, nil
			},
			wantError: errors.New("missing properties: [Amount]"),
		},
		{
			name:          "return error if multiple properties have wrong types",
			envDatabaseId: databaseId,
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return &notion.Database{ID: databaseId}, nil
			},
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return notion.DatabaseProperties{
					"Transaction Date": notion.DatabaseProperty{Name: "Transaction Date", Type: notion.DBPropTypeRichText}, // Should be Date
					"Operation":        notion.DatabaseProperty{Name: "Operation", Type: notion.DBPropTypeRichText},        // Should be Select
					"Destination":      notion.DatabaseProperty{Name: "Destination", Type: notion.DBPropTypeRichText},
					"Amount":           notion.DatabaseProperty{Name: "Amount", Type: notion.DBPropTypeNumber},
					"Memo":             notion.DatabaseProperty{Name: "Memo", Type: notion.DBPropTypeRichText},
					"Transaction ID":   notion.DatabaseProperty{Name: "Transaction ID", Type: notion.DBPropTypeRichText},
				}, nil
			},
			wantError: errors.New("missing properties: [Transaction Date Operation]"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			loader := &NotionLoader{
				notionClient: &databasemock.MockDatabaseClient{
					GetDatabaseFunc:            test.mockDatabase,
					ListDatabasePropertiesFunc: test.mockDatabaseProperties,
				},
			}

			os.Setenv("DATABASE_ID", test.envDatabaseId)

			err := loader.UploadReport(&statement.Report{})
			if err != nil {
				if test.wantError != nil {
					
					// check if the missing properties obtained are the expected ones
					if strings.Contains(err.Error(), "missing properties") {
						for propName := range allProperties {

							// Unexpected missing property case
							if strings.Contains(err.Error(), propName) && !strings.Contains(test.wantError.Error(), propName) {
								t.Errorf("%s property was not expected to be missing", propName)
							}

							// Unexpected existing property case
							if !strings.Contains(err.Error(), propName) && strings.Contains(test.wantError.Error(), propName) {
								t.Errorf("%s property was expected to be missing", propName)
							}
						}
					} else {
						if diff := cmp.Diff(test.wantError.Error(), err.Error(),
							cmpopts.EquateErrors(),
						); diff != "" {
							t.Fatalf("Mistach (-want +got):\n%s", diff)
						}
					}
				} else {
					t.Fatalf("expected no error, got: %v", err)
				}
			} else if test.wantError != nil {
				t.Errorf("expected %v error, got nil", test.wantError)
			}
		})
	}
}
