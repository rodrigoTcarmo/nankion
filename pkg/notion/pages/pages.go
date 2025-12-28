package notion

import (
	"context"
	"fmt"

	"github.com/dstotijn/go-notion"
	notionclient "github.com/rodrigoTcarmo/nankion/pkg/notion/client"
)

type PageClient interface {
	SearchPages(string) (*notion.Page, error)
	CreatePage(string, string, *notion.DatabasePageProperties) error
}

type PageData struct {
	Title      string
	databaseId string
	properties *notion.DatabasePageProperties
}

type Page struct {
	client *notion.Client
}

func (p *Page) SearchPages(pageId string) (*notion.Page, error) {
	pageFound, err := p.client.FindPageByID(context.Background(), pageId)
	if err != nil {
		return nil, fmt.Errorf("error trying to find page by ID: %v", err)
	}

	return &pageFound, nil
}

func (p *Page) CreatePage(pageData PageData) error {
	pageParams := notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               pageData.databaseId,
		DatabasePageProperties: pageData.properties,
		Title:                  []notion.RichText{{PlainText: pageData.Title}},
	}

	_, err := p.client.CreatePage(context.Background(), pageParams)
	if err != nil {
		return err
	}

	return nil
}

func NewPage() *Page {
	return &Page{
		client: notionclient.NewClient(),
	}
}
