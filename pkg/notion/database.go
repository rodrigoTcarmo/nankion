package notion

import (
	"context"
	"errors"
	"net/http"

	"github.com/dstotijn/go-notion"
)

func SearchDatabases(client *notion.Client, databaseId string) (*notion.Database, int, error) {
	var notionErr *notion.APIError

	database, err := client.FindDatabaseByID(context.Background(), databaseId)
	if err != nil {
		if errors.As(err, &notionErr) {
			return nil, notionErr.Status, err
		}
		return nil, http.StatusInternalServerError, err
	}

	return &database, http.StatusOK, nil
}
