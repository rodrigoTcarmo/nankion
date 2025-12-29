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

func (p *page) validatePageDuplicity(query string) (*notion.SearchResponse, error) {
	searchResponse, err := p.client.Search(context.Background(), &notion.SearchOpts{
		Query: query,
		Filter: &notion.SearchFilter{
			Property: "object",
			Value:    "page",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error trying to search page by transaction ID: %v", err)
	}

	return &searchResponse, nil
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
	searchResponse, err := p.validatePageDuplicity(statement.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("error trying to validate page duplicity: %v", err)
	}

	if len(searchResponse.Results) > 0 {
		fmt.Printf("Warning! Transaction %s may already exist!\n", statement.TransactionID)
	}
	
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
