package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
	"github.com/rodrigoTcarmo/nankion/pkg/notion/client/mocks"
)

func TestGetDatabase(t *testing.T) {
	defaultDatabase := notion.Database{
		ID: "12345678",
		Title: []notion.RichText{
			{PlainText: "my-database"},
		},
	}
	databaseWithProperties := notion.Database{
		ID: "12345678",
		Title: []notion.RichText{
			{PlainText: "my-database"},
		},
		Properties: notion.DatabaseProperties{
			"Name":   notion.DatabaseProperty{Type: notion.DBPropTypeTitle},
			"Status": notion.DatabaseProperty{Type: notion.DBPropTypeSelect},
		},
	}
	archivedDatabase := notion.Database{
		ID: "archived-db",
		Title: []notion.RichText{
			{PlainText: "archived-database"},
		},
		Archived: true,
	}

	tests := []struct {
		name         string
		databaseID   string
		mockDatabase notion.Database
		mockError    error
		wantDatabase *notion.Database
		wantError    error
	}{
		{
			name:         "successfully gets database",
			databaseID:   "12345678",
			mockDatabase: defaultDatabase,
			mockError:    nil,
			wantDatabase: &defaultDatabase,
			wantError:    nil,
		},
		{
			name:         "error trying to get database",
			databaseID:   "12345678",
			mockDatabase: notion.Database{},
			mockError:    fmt.Errorf("invalid auth"),
			wantDatabase: nil,
			wantError:    fmt.Errorf("error trying to return Database by ID: invalid auth"),
		},
		{
			name:         "empty database ID",
			databaseID:   "",
			mockDatabase: notion.Database{},
			mockError:    fmt.Errorf("database not found"),
			wantDatabase: nil,
			wantError:    fmt.Errorf("error trying to return Database by ID: database not found"),
		},
		{
			name:         "successfully gets database with properties",
			databaseID:   "12345678",
			mockDatabase: databaseWithProperties,
			mockError:    nil,
			wantDatabase: &databaseWithProperties,
			wantError:    nil,
		},
		{
			name:         "successfully gets archived database",
			databaseID:   "archived-db",
			mockDatabase: archivedDatabase,
			mockError:    nil,
			wantDatabase: &archivedDatabase,
			wantError:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotDatabaseID string

			mockNotionClient := mocks.MockNotionClient{
				FindDatabaseByIDFunc: func(ctx context.Context, databaseID string) (notion.Database, error) {
					gotDatabaseID = databaseID
					return test.mockDatabase, test.mockError
				},
			}

			databaseClient := database{client: &mockNotionClient}

			gotDatabase, err := databaseClient.GetDatabase(test.databaseID)
			if err != nil {
				if test.wantError == nil {
					t.Fatalf("expected no error, but got: %v", err)
				}
				if diff := cmp.Diff(test.wantError.Error(), err.Error()); diff != "" {
					t.Fatalf("error mismatch (-want +got):\n%s", diff)
				}
				return
			}

			if test.wantError != nil {
				t.Fatalf("expected error %v, but got nil", test.wantError)
			}

			if diff := cmp.Diff(test.databaseID, gotDatabaseID); diff != "" {
				t.Fatalf("databaseID mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(test.wantDatabase, gotDatabase); diff != "" {
				t.Fatalf("database mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestListDatabaseProperties(t *testing.T) {
	defaultProperties := notion.DatabaseProperties{
		"Name":   notion.DatabaseProperty{Type: notion.DBPropTypeTitle},
		"Status": notion.DatabaseProperty{Type: notion.DBPropTypeSelect},
	}
	tests := []struct {
		name           string
		databaseID     string
		mockDatabase   notion.Database
		mockError      error
		wantProperties notion.DatabaseProperties
		wantError      error
	}{
		{
			name:           "successfully lists database properties",
			databaseID:     "12345678",
			mockDatabase:   notion.Database{ID: "12345678", Properties: defaultProperties},
			mockError:      nil,
			wantProperties: defaultProperties,
			wantError:      nil,
		},
		{
			name:           "error trying to get database",
			databaseID:     "12345678",
			mockDatabase:   notion.Database{},
			mockError:      fmt.Errorf("invalid auth"),
			wantProperties: nil,
			wantError:      fmt.Errorf("error trying to return Database by ID: invalid auth"),
		},
		{
			name:           "empty database ID",
			databaseID:     "",
			mockDatabase:   notion.Database{},
			mockError:      fmt.Errorf("database not found"),
			wantProperties: nil,
			wantError:      fmt.Errorf("error trying to return Database by ID: database not found"),
		},
		{
			name:           "database with no properties",
			databaseID:     "12345678",
			mockDatabase:   notion.Database{ID: "12345678", Properties: nil},
			mockError:      nil,
			wantProperties: nil,
			wantError:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotDatabaseID string

			mockNotionClient := mocks.MockNotionClient{
				FindDatabaseByIDFunc: func(ctx context.Context, databaseID string) (notion.Database, error) {
					gotDatabaseID = databaseID
					return test.mockDatabase, test.mockError
				},
			}

			databaseClient := database{client: &mockNotionClient}

			gotProperties, err := databaseClient.ListDatabaseProperties(test.databaseID)
			if err != nil {
				if test.wantError == nil {
					t.Fatalf("expected no error, but got: %v", err)
				}
				if diff := cmp.Diff(test.wantError.Error(), err.Error()); diff != "" {
					t.Fatalf("error mismatch (-want +got):\n%s", diff)
				}
				return
			}

			if test.wantError != nil {
				t.Fatalf("expected error %v, but got nil", test.wantError)
			}

			if diff := cmp.Diff(test.databaseID, gotDatabaseID); diff != "" {
				t.Fatalf("databaseID mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(test.wantProperties, gotProperties); diff != "" {
				t.Fatalf("properties mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpsertDatabaseProperties(t *testing.T) {
	defaultProperties := map[string]*notion.DatabaseProperty{
		"Name":   {Type: notion.DBPropTypeTitle},
		"Status": {Type: notion.DBPropTypeSelect},
	}
	tests := []struct {
		name       string
		databaseID string
		properties map[string]*notion.DatabaseProperty
		wantParams notion.UpdateDatabaseParams
		mockError  error
		wantError  error
	}{
		{
			name:       "successfully upserts database properties",
			databaseID: "12345678",
			properties: defaultProperties,
			wantParams: notion.UpdateDatabaseParams{Properties: defaultProperties},
			mockError:  nil,
			wantError:  nil,
		},
		{
			name:       "error trying to update database",
			databaseID: "12345678",
			properties: defaultProperties,
			wantParams: notion.UpdateDatabaseParams{Properties: defaultProperties},
			mockError:  fmt.Errorf("invalid auth"),
			wantError:  fmt.Errorf("invalid auth"),
		},
		{
			name:       "empty database ID",
			databaseID: "",
			properties: defaultProperties,
			wantParams: notion.UpdateDatabaseParams{Properties: defaultProperties},
			mockError:  fmt.Errorf("database not found"),
			wantError:  fmt.Errorf("database not found"),
		},
		{
			name:       "empty properties",
			databaseID: "12345678",
			properties: map[string]*notion.DatabaseProperty{},
			wantParams: notion.UpdateDatabaseParams{Properties: map[string]*notion.DatabaseProperty{}},
			mockError:  nil,
			wantError:  nil,
		},
		{
			name:       "nil properties",
			databaseID: "12345678",
			properties: nil,
			wantParams: notion.UpdateDatabaseParams{Properties: nil},
			mockError:  nil,
			wantError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotDatabaseID string
			var gotParams notion.UpdateDatabaseParams

			mockNotionClient := mocks.MockNotionClient{
				UpdateDatabaseFunc: func(ctx context.Context, databaseID string, params notion.UpdateDatabaseParams) (notion.Database, error) {
					gotDatabaseID = databaseID
					gotParams = params
					return notion.Database{ID: databaseID}, test.mockError
				},
			}

			databaseClient := database{client: &mockNotionClient}

			err := databaseClient.UpsertDatabaseProperties(test.databaseID, test.properties)
			if err != nil {
				if test.wantError == nil {
					t.Fatalf("expected no error, but got: %v", err)
				}
				if diff := cmp.Diff(test.wantError.Error(), err.Error()); diff != "" {
					t.Fatalf("error mismatch (-want +got):\n%s", diff)
				}
				return
			}

			if test.wantError != nil {
				t.Fatalf("expected error %v, but got nil", test.wantError)
			}

			if diff := cmp.Diff(test.databaseID, gotDatabaseID); diff != "" {
				t.Fatalf("databaseID mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(test.wantParams, gotParams); diff != "" {
				t.Fatalf("params mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
