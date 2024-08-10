package notion

import (
	"context"

	"github.com/dstotijn/go-notion"
)

type NotionPages interface {
	SearchPages(*notion.Client, string) (notion.Page, error)
}

type NotionPageObjects struct {
	PageId string
	Code   string
}

func (p *NotionPageObjects) SearchPages(client *notion.Client, pageId string) (notion.Page, error) {
	pageFound, err := client.FindPageByID(context.Background(), pageId)
	if err != nil {
		var notionErr notion.APIError
		p.Code = notionErr.Code
		return notion.Page{}, err
	}
	return pageFound, nil
}
