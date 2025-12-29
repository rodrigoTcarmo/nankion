package notion

import (
	"errors"
	"fmt"

	"github.com/dstotijn/go-notion"
	"github.com/rodrigoTcarmo/nankion/pkg/models/properties"
	notionclient "github.com/rodrigoTcarmo/nankion/pkg/notion/client"
	"github.com/rodrigoTcarmo/nankion/pkg/notion/database"
	"github.com/rodrigoTcarmo/nankion/pkg/notion/page"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

type Notion struct {
	client   *notion.Client
	database database.Database
	page     page.Page
}

func NewNotionLoader() *Notion {
	return &Notion{
		client:   notionclient.NewClient(),
		database: database.NewDatabase(),
		page:     page.NewPage(),
	}
}

// UploadReport uploads a report to the Notion database
func (n *Notion) UploadReport(databaseID, filePath string) error {
	_, err := n.database.GetDatabase(databaseID)
	if err != nil {
		return fmt.Errorf("error trying to get database by ID: %s", err)
	}

	// 2 - Check if all properties are created and exists in the database
	databaseProperties, err := n.database.ListDatabaseProperties(databaseID)
	if err != nil {
		return fmt.Errorf("error trying to list database properties: %s", err)
	}

	var existentProperties = map[string]notion.DatabasePropertyType{}
	for _, property := range databaseProperties {
		existentProperties[property.Name] = property.Type
	}

	// Create a copy of RequiredProperties to track missing ones
	missingProperties := make(map[string]*notion.DatabaseProperty)
	for reqPropName, reqProp := range properties.RequiredProperties {
		if propType, ok := existentProperties[reqPropName]; !ok || propType != reqProp.Type {
			missingProperties[reqPropName] = reqProp
		}
	}

	if len(missingProperties) > 0 {
		err := n.database.UpsertDatabaseProperties(databaseID, missingProperties)
		if err != nil {
			return fmt.Errorf("error trying to upsert new database properties: %s", err)
		}
	}

	report, err := statement.LoadReport(filePath)
	if err != nil {
		return fmt.Errorf("error trying to load statement report: %s", err)
	}

	if err := n.uploadPages(databaseID, report); err != nil {
		return fmt.Errorf("error uploading pages: %w", err)
	}

	return nil
}

func (n *Notion) uploadPages(databaseID string, report *statement.Report) error {
	var errs []error
	for _, statement := range report.Statements {
		newPage := n.page.BuildPage(databaseID, statement)
		err := n.page.CreatePage(newPage)
		if err != nil {
			errs = append(errs, fmt.Errorf("error trying to create page for %s: %w", statement.Destination, err))
		}
	}

	return errors.Join(errs...)
}
