package notion

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	databasemock "github.com/rodrigoTcarmo/nankion/pkg/notion/database/mocks"
	"github.com/rodrigoTcarmo/nankion/pkg/notion/page"
	pagemock "github.com/rodrigoTcarmo/nankion/pkg/notion/page/mocks"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

// getTestDataPath returns the absolute path to the testdata directory
func getTestDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "testdata", filename)
}

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
		filePath               string
		mockDatabase           func(string) (*notion.Database, error)
		mockDatabaseProperties func(string) (notion.DatabaseProperties, error)
		wantError              error
		wantUpsertProperties   []string
		wantCreatePageCalls    int
	}{
		{
			name:          "get all expected properties from notion and create pages",
			envDatabaseId: databaseId,
			filePath:      getTestDataPath("test.ofx"),
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return allProperties, nil
			},
			wantError:           nil,
			wantCreatePageCalls: 1,
		},
		{
			name:          "succeed when database has extra properties beyond required ones",
			envDatabaseId: databaseId,
			filePath:      getTestDataPath("test.ofx"),
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
			wantError:           nil,
			wantCreatePageCalls: 1,
		},
		{
			name:          "return error if database does not exist",
			envDatabaseId: databaseId,
			filePath:      "nonexistent.ofx",
			mockDatabase: func(databaseId string) (*notion.Database, error) {
				return nil, errors.New("database not found")
			},
			wantError:           errors.New("error trying to validate database: database not found"),
			wantCreatePageCalls: 0,
		},
		{
			name:          "return error if listing database properties fails",
			envDatabaseId: databaseId,
			filePath:      "nonexistent.ofx",
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return nil, errors.New("unauthorized access")
			},
			wantError:           errors.New("error trying to validate database properties: error trying to list database properties: unauthorized access"),
			wantCreatePageCalls: 0,
		},
		{
			name:          "create all properties when all are missing",
			envDatabaseId: databaseId,
			filePath:      getTestDataPath("test.ofx"),
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return notion.DatabaseProperties{}, nil
			},
			wantError:            nil,
			wantUpsertProperties: []string{"Transaction Date", "Operation", "Destination", "Amount", "Memo", "Transaction ID"},
			wantCreatePageCalls:  1,
		},
		{
			name:          "create Transaction ID property when it is missing",
			envDatabaseId: databaseId,
			filePath:      getTestDataPath("test.ofx"),
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
			wantCreatePageCalls:  1,
		},
		{
			name:          "create multiple properties when they are missing",
			envDatabaseId: databaseId,
			filePath:      getTestDataPath("test.ofx"),
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return notion.DatabaseProperties{
					"Transaction Date": notion.DatabaseProperty{Name: "Transaction Date", Type: notion.DBPropTypeDate},
					"Operation":        notion.DatabaseProperty{Name: "Operation", Type: notion.DBPropTypeSelect},
					"Destination":      notion.DatabaseProperty{Name: "Destination", Type: notion.DBPropTypeRichText},
				}, nil
			},
			wantError:            nil,
			wantUpsertProperties: []string{"Amount", "Memo", "Transaction ID"},
			wantCreatePageCalls:  1,
		},
		{
			name:          "create property when it exists but with wrong type",
			envDatabaseId: databaseId,
			filePath:      getTestDataPath("test.ofx"),
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
			wantCreatePageCalls:  1,
		},
		{
			name:          "create multiple properties when they have wrong types",
			envDatabaseId: databaseId,
			filePath:      getTestDataPath("test.ofx"),
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
			wantCreatePageCalls:  1,
		},
		{
			name:          "return error when loading report fails",
			envDatabaseId: databaseId,
			filePath:      "nonexistent.ofx",
			mockDatabaseProperties: func(databaseId string) (notion.DatabaseProperties, error) {
				return allProperties, nil
			},
			wantError:           errors.New("error trying to load statement report"),
			wantCreatePageCalls: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			upsertCalled := false
			var upsertedProperties map[string]*notion.DatabaseProperty
			createPageCalls := 0

			mockUpsert := func(databaseId string, properties map[string]*notion.DatabaseProperty) error {
				upsertCalled = true
				upsertedProperties = properties
				return nil
			}

			// Setup page mock with call tracking
			var mockPage *pagemock.MockPageClient
			mockPage = &pagemock.MockPageClient{
				BuildPageFunc: func(dbID string, stmt statement.Statement) page.PageData {
					return page.PageData{DatabaseId: dbID}
				},
				CreatePageFunc: func(pageData page.PageData) error {
					createPageCalls++
					return nil
				},
			}

			loader := &Notion{
				database: &databasemock.MockDatabaseClient{
					GetDatabaseFunc: func(databaseId string) (*notion.Database, error) {
						if test.mockDatabase != nil {
							return test.mockDatabase(databaseId)
						}
						return &notion.Database{ID: databaseId, Properties: allProperties}, nil
					},
					ListDatabasePropertiesFunc:   test.mockDatabaseProperties,
					UpsertDatabasePropertiesFunc: mockUpsert,
				},
				page: mockPage,
			}

			os.Setenv("DATABASE_ID", test.envDatabaseId)

			err := loader.UploadReport(test.envDatabaseId, test.filePath)

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
						// Use strings.Contains for partial error matching
						if !strings.Contains(err.Error(), test.wantError.Error()) {
							if diff := cmp.Diff(test.wantError.Error(), err.Error(),
								cmpopts.EquateErrors(),
							); diff != "" {
								t.Fatalf("Mistach (-want +got):\n%s", diff)
							}
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

			// Verify page creation calls
			if createPageCalls != test.wantCreatePageCalls {
				t.Errorf("expected %d CreatePage calls, got %d", test.wantCreatePageCalls, createPageCalls)
			}
		})
	}
}
