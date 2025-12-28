package page

import (
	"context"
	"fmt"
	"strings"

	"github.com/dstotijn/go-notion"
	"github.com/rodrigoTcarmo/nankion/pkg/models/properties"
	notionclient "github.com/rodrigoTcarmo/nankion/pkg/notion/client"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

type PageClient interface {
	SearchPages(string) (*notion.Page, error)
	CreatePage(string, string, *notion.DatabasePageProperties) error
}

type PageData struct {
	Title      string
	DatabaseId string
	Properties *notion.DatabasePageProperties
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

func (p *Page) BuildPage(databaseID string, statement statement.Statement) PageData {
	return PageData{
		Title:      strings.Join([]string{string(statement.Operation), statement.Destination}, "-"),
		DatabaseId: databaseID,
		Properties: BuildPageProperties(statement),
	}
}

func (p *Page) CreatePageParams(pageData PageData) notion.CreatePageParams {
	return notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               pageData.DatabaseId,
		DatabasePageProperties: pageData.Properties,
		Title:                  []notion.RichText{{PlainText: pageData.Title}},
	}
}

func (p *Page) CreatePage(pageData PageData) error {
	pageParams := p.CreatePageParams(pageData)

	_, err := p.client.CreatePage(context.Background(), pageParams)
	if err != nil {
		return err
	}

	return nil
}

func BuildPageProperties(statement statement.Statement) *notion.DatabasePageProperties {
	return &notion.DatabasePageProperties{
		"Transaction Date": properties.TransactionDateProperty(statement.TransactionDate),
		"Operation":        properties.OperationProperty(statement.Operation),
		"Destination":      properties.DestinationProperty(statement.Destination),
		"Amount":           properties.AmountProperty(statement.Amount),
		"Memo":             properties.MemoProperty(statement.Memo),
		"Transaction ID":   properties.TransactionIDProperty(statement.TransactionID),
	}
}

func NewPage() Page {
	return Page{
		client: notionclient.NewClient(),
	}
}
