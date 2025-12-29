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
		wantUpsertProperties   []string
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
			name:          "create all properties when all are missing",
			envDatabaseId: databaseId,
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return &notion.Database{
					ID: databaseId,
				}, nil
			},
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return notion.DatabaseProperties{}, nil
			},
			wantError:            nil,
			wantUpsertProperties: []string{"Transaction Date", "Operation", "Destination", "Amount", "Memo", "Transaction ID"},
		},
		{
			name:          "create Transaction ID property when it is missing",
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
			wantError:            nil,
			wantUpsertProperties: []string{"Transaction ID"},
		},
		{
			name:          "create multiple properties when they are missing",
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
			wantError:            nil,
			wantUpsertProperties: []string{"Amount", "Memo", "Transaction ID"},
		},
		{
			name:          "create property when it exists but with wrong type",
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
			wantError:            nil,
			wantUpsertProperties: []string{"Amount"},
		},
		{
			name:          "create multiple properties when they have wrong types",
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
			wantError:            nil,
			wantUpsertProperties: []string{"Transaction Date", "Operation"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			upsertCalled := false
			var upsertedProperties map[string]*notion.DatabaseProperty

			mockUpsert := func(databaseId string, properties map[string]*notion.DatabaseProperty) error {
				upsertCalled = true
				upsertedProperties = properties
				return nil
			}

			loader := &Notion{
				database: &databasemock.MockDatabaseClient{
					GetDatabaseFunc:              test.mockDatabase,
					ListDatabasePropertiesFunc:   test.mockDatabaseProperties,
					UpsertDatabasePropertiesFunc: mockUpsert,
				},
			}

			os.Setenv("DATABASE_ID", test.envDatabaseId)

			err := loader.UploadReport(test.envDatabaseId)

			// Check error expectations
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

			// Verify the correct properties were passed to upsert
			if test.wantUpsertProperties != nil && upsertCalled {
				if len(upsertedProperties) != len(test.wantUpsertProperties) {
					t.Errorf("expected %d properties to be upserted, got %d", len(test.wantUpsertProperties), len(upsertedProperties))
				}
				for _, propName := range test.wantUpsertProperties {
					if _, ok := upsertedProperties[propName]; !ok {
						t.Errorf("expected property %s to be upserted, but it wasn't", propName)
					}
				}
			} else if test.wantUpsertProperties == nil && upsertCalled {
				t.Error("unexpected UpsertDatabaseProperties func call")
			} else if test.wantUpsertProperties != nil && upsertCalled == false {
				t.Error("expected UpsertDatabaseProperties func to be called, but it was not.")
			}
		})
	}
}
