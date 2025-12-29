package page

import (
	"context"
	"strings"

	"github.com/dstotijn/go-notion"
	"github.com/rodrigoTcarmo/nankion/pkg/models/properties"
	notionclient "github.com/rodrigoTcarmo/nankion/pkg/notion/client"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

type Page interface {
	BuildPage(string, statement.Statement) (*PageData, error)
	CreatePage(PageData) error
}

type PageData struct {
	DatabaseId string
	Properties *notion.DatabasePageProperties
}

type page struct {
	client *notion.Client
}

func (p *page) BuildPage(databaseID string, statement statement.Statement) (*PageData, error) {
	properties, err := p.BuildPageProperties(statement)
	if err != nil {
		return nil, err
	}

	pageData := &PageData{
		DatabaseId: databaseID,
		Properties: properties,
	}
	return pageData, nil
}

func (p *page) CreatePageParams(pageData PageData) notion.CreatePageParams {
	return notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               pageData.DatabaseId,
		DatabasePageProperties: pageData.Properties,
	}
}

func (p *page) CreatePage(pageData PageData) error {
	pageParams := p.CreatePageParams(pageData)

	_, err := p.client.CreatePage(context.Background(), pageParams)
	if err != nil {
		return err
	}

	return nil
}

func (p *page) BuildPageProperties(statement statement.Statement) (*notion.DatabasePageProperties, error) {
	title := strings.Join([]string{string(statement.Operation), statement.Destination, statement.TransactionID}, "-")
	return &notion.DatabasePageProperties{
		"Name":             properties.TitleProperty(title),
		"Transaction Date": properties.TransactionDateProperty(statement.TransactionDate),
		"Operation":        properties.OperationProperty(statement.Operation),
		"Destination":      properties.DestinationProperty(statement.Destination),
		"Amount":           properties.AmountProperty(statement.Amount),
		"Memo":             properties.MemoProperty(statement.Memo),
		"Transaction ID":   properties.TransactionIDProperty(statement.TransactionID),
	}, nil
}

func NewPage() Page {
	return &page{
		client: notionclient.NewClient(),
	}
}
