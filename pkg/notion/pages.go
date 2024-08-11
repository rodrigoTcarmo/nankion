package notion

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dstotijn/go-notion"
)

func SearchPages(client *notion.Client, pageId string) (*notion.Page, int, error) {
	notionErr := &notion.APIError{}

	pageFound, err := client.FindPageByID(context.Background(), pageId)
	if err != nil {
		if errors.As(err, &notionErr) {
			return nil, notionErr.Status, fmt.Errorf(notionErr.Message)
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("error trying to find page by ID: %v", err)
	}
	return &pageFound, http.StatusOK, nil
}
