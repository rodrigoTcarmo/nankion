package databases

import (
	"testing"

	"github.com/dstotijn/go-notion"
	databasemock "github.com/rodrigoTcarmo/nankion/pkg/notion/database/mocks"
)

func TestGetDatabase(t *testing.T) {
	tests := []struct {
		name   string
		setup func()
		want   *notion.Database
	}{
		{
			name: "successfully gets the database id",
			setup: func ()  {
				
			},
			want: &notion.Database{
				ID:  "super-banana-id",
				URL: "https://notion.so/" + "super-banana-id",
				Title: []notion.RichText{
					{
						PlainText: "my-database",
					},
				},
				Properties: notion.DatabaseProperties{
					"prop1": {
						ID:   "propid1",
						Type: notion.DBPropTypeRichText,
						Name: "prop-description",
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			client := DatabaseCLI{
				database: &databasemock.MockDatabaseClient{
					GetDatabaseFunc: func(databaseID string) (*notion.Database, error) {
						return test.want, nil
					},
					ListDatabasePropertiesFunc: ,
				},
			}
			client.GetDatabase(test.want.ID)
		})
	}

}
