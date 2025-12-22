package databases

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/dstotijn/go-notion"
	databasemock "github.com/rodrigoTcarmo/nankion/pkg/notion/database/mocks"
)

func TestGetDatabase(t *testing.T) {
	tests := []struct {
		name       string
		databaseID string
		database   *notion.Database
		wantOutput string
		wantError  error
	}{
		{
			name:       "successfully gets the database id",
			databaseID: "super-banana-id",
			database: &notion.Database{
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
			wantOutput: `{
 "id": "super-banana-id",
 "created_time": "0001-01-01T00:00:00Z",
 "created_by": {
  "id": ""
 },
 "last_edited_time": "0001-01-01T00:00:00Z",
 "last_edited_by": {
  "id": ""
 },
 "url": "https://notion.so/super-banana-id",
 "title": [
  {
   "plain_text": "my-database"
  }
 ],
 "description": null,
 "properties": {
  "prop1": {
   "id": "propid1",
   "type": "rich_text",
   "name": "prop-description"
  }
 },
 "parent": {},
 "archived": false,
 "is_inline": false
}
`,
			wantError: nil,
		},
		{
			name:       "returns error when database not found",
			databaseID: "non-existent-id",
			database:   nil,
			wantOutput: "error trying to get database: database not foundnull\n",
			wantError:  errors.New("database not found"),
		},
		{
			name:       "successfully gets database with multiple properties",
			databaseID: "multi-prop-db",
			database: &notion.Database{
				ID:  "multi-prop-db",
				URL: "https://notion.so/multi-prop-db",
				Title: []notion.RichText{
					{
						PlainText: "Task Tracker",
					},
				},
				Properties: notion.DatabaseProperties{
					"Name": {
						ID:   "title-id",
						Type: notion.DBPropTypeTitle,
						Name: "Name",
					},
					"Status": {
						ID:   "status-id",
						Type: notion.DBPropTypeSelect,
						Name: "Status",
					},
					"Due Date": {
						ID:   "date-id",
						Type: notion.DBPropTypeDate,
						Name: "Due Date",
					},
					"Assignee": {
						ID:   "people-id",
						Type: notion.DBPropTypePeople,
						Name: "Assignee",
					},
				},
			},
			wantOutput: `{
 "id": "multi-prop-db",
 "created_time": "0001-01-01T00:00:00Z",
 "created_by": {
  "id": ""
 },
 "last_edited_time": "0001-01-01T00:00:00Z",
 "last_edited_by": {
  "id": ""
 },
 "url": "https://notion.so/multi-prop-db",
 "title": [
  {
   "plain_text": "Task Tracker"
  }
 ],
 "description": null,
 "properties": {
  "Assignee": {
   "id": "people-id",
   "type": "people",
   "name": "Assignee"
  },
  "Due Date": {
   "id": "date-id",
   "type": "date",
   "name": "Due Date"
  },
  "Name": {
   "id": "title-id",
   "type": "title",
   "name": "Name"
  },
  "Status": {
   "id": "status-id",
   "type": "select",
   "name": "Status"
  }
 },
 "parent": {},
 "archived": false,
 "is_inline": false
}
`,
			wantError: nil,
		},
		{
			name:       "successfully gets database with empty properties",
			databaseID: "empty-props-db",
			database: &notion.Database{
				ID:  "empty-props-db",
				URL: "https://notion.so/empty-props-db",
				Title: []notion.RichText{
					{
						PlainText: "Empty Database",
					},
				},
				Properties: notion.DatabaseProperties{},
			},
			wantOutput: `{
 "id": "empty-props-db",
 "created_time": "0001-01-01T00:00:00Z",
 "created_by": {
  "id": ""
 },
 "last_edited_time": "0001-01-01T00:00:00Z",
 "last_edited_by": {
  "id": ""
 },
 "url": "https://notion.so/empty-props-db",
 "title": [
  {
   "plain_text": "Empty Database"
  }
 ],
 "description": null,
 "properties": {},
 "parent": {},
 "archived": false,
 "is_inline": false
}
`,
			wantError: nil,
		},
		{
			name:       "successfully gets archived database",
			databaseID: "archived-db",
			database: &notion.Database{
				ID:       "archived-db",
				URL:      "https://notion.so/archived-db",
				Archived: true,
				Title: []notion.RichText{
					{
						PlainText: "Archived Database",
					},
				},
				Properties: notion.DatabaseProperties{},
			},
			wantOutput: `{
 "id": "archived-db",
 "created_time": "0001-01-01T00:00:00Z",
 "created_by": {
  "id": ""
 },
 "last_edited_time": "0001-01-01T00:00:00Z",
 "last_edited_by": {
  "id": ""
 },
 "url": "https://notion.so/archived-db",
 "title": [
  {
   "plain_text": "Archived Database"
  }
 ],
 "description": null,
 "properties": {},
 "parent": {},
 "archived": true,
 "is_inline": false
}
`,
			wantError: nil,
		},
		{
			name:       "successfully gets inline database",
			databaseID: "inline-db",
			database: &notion.Database{
				ID:       "inline-db",
				URL:      "https://notion.so/inline-db",
				IsInline: true,
				Title: []notion.RichText{
					{
						PlainText: "Inline Database",
					},
				},
				Properties: notion.DatabaseProperties{},
			},
			wantOutput: `{
 "id": "inline-db",
 "created_time": "0001-01-01T00:00:00Z",
 "created_by": {
  "id": ""
 },
 "last_edited_time": "0001-01-01T00:00:00Z",
 "last_edited_by": {
  "id": ""
 },
 "url": "https://notion.so/inline-db",
 "title": [
  {
   "plain_text": "Inline Database"
  }
 ],
 "description": null,
 "properties": {},
 "parent": {},
 "archived": false,
 "is_inline": true
}
`,
			wantError: nil,
		},
		{
			name:       "successfully gets database with empty title",
			databaseID: "no-title-db",
			database: &notion.Database{
				ID:         "no-title-db",
				URL:        "https://notion.so/no-title-db",
				Title:      []notion.RichText{},
				Properties: notion.DatabaseProperties{},
			},
			wantOutput: `{
 "id": "no-title-db",
 "created_time": "0001-01-01T00:00:00Z",
 "created_by": {
  "id": ""
 },
 "last_edited_time": "0001-01-01T00:00:00Z",
 "last_edited_by": {
  "id": ""
 },
 "url": "https://notion.so/no-title-db",
 "title": [],
 "description": null,
 "properties": {},
 "parent": {},
 "archived": false,
 "is_inline": false
}
`,
			wantError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			databaseMock := &databasemock.MockDatabaseClient{
				GetDatabaseFunc: func(databaseID string) (*notion.Database, error) {
					return test.database, test.wantError
				},
			}
			client := DatabaseCLI{
				database: databaseMock,
			}
			output := captureStdout(func() {
				client.GetDatabase(test.databaseID)
			})

			if output != test.wantOutput {
				t.Errorf("expected output to be %q\n\n, got %q", test.wantOutput, output)
			}
		})
	}

}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
