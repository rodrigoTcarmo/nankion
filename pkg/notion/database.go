package notion

import (
	"context"
	"errors"
	"net/http"

	"github.com/dstotijn/go-notion"
)

type DatabaseClient interface {
	SearchDatabases(string) (*notion.Database, int, error)
}

type Database struct {
	client *NankionClient
}

func (d *Database) SearchDatabases(databaseId string) (*notion.Database, int, error) {
	var notionErr *notion.APIError

	database, err := d.client.Client.FindDatabaseByID(context.Background(), databaseId)
	if err != nil {
		if errors.As(err, &notionErr) {
			return nil, notionErr.Status, err
		}
		return nil, http.StatusInternalServerError, err
	}

	return &database, http.StatusOK, nil
}

func NewDatabase() *Database {
	return &Database{
		client: NewNankionClient(),
	}
}
