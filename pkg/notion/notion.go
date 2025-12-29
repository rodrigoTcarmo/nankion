package notion

import (
	"errors"
	"fmt"
	"log/slog"

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

// UploadStatement uploads a report to the Notion database
func (n *Notion) UploadStatement(databaseID, filePath string) error {
	if err := n.validateDatabase(databaseID); err != nil {
		return fmt.Errorf("error trying to validate database: %s", err)
	}

	if err := n.validateDatabaseProperties(databaseID); err != nil {
		return fmt.Errorf("error trying to validate database properties: %s", err)
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
		newPage, err := n.page.BuildPage(databaseID, statement)
		if err != nil {
			errs = append(errs, fmt.Errorf("error trying to build page for %s: %w", statement.Destination, err))
			continue
		}
		if err := n.page.CreatePage(*newPage); err != nil {
			errs = append(errs, fmt.Errorf("error trying to create page for %s: %w", statement.Destination, err))
		} else {
			slog.Info("Page successfully created", "memo", statement.Memo, "destination", statement.Destination)
		}
	}

	return errors.Join(errs...)
}

func (n *Notion) validateDatabase(databaseID string) error {
	_, err := n.database.GetDatabase(databaseID)
	if err != nil {
		return err
	}
	return nil
}

func (n *Notion) validateDatabaseProperties(databaseID string) error {
	existentProperties, err := n.getExistentProperties(databaseID)
	if err != nil {
		return err
	}

	missingProperties := n.mapMissingProperties(existentProperties)

	if len(missingProperties) == 0 {
		return nil
	}

	return n.upsertDatabaseProperties(databaseID, missingProperties)
}

func (n *Notion) getExistentProperties(databaseID string) (map[string]notion.DatabasePropertyType, error) {
	databaseProperties, err := n.database.ListDatabaseProperties(databaseID)
	if err != nil {
		return nil, fmt.Errorf("error trying to list database properties: %s", err)
	}

	var existentProperties = map[string]notion.DatabasePropertyType{}
	for _, property := range databaseProperties {
		existentProperties[property.Name] = property.Type
	}
	return existentProperties, nil
}

func (n *Notion) mapMissingProperties(existentProperties map[string]notion.DatabasePropertyType) map[string]*notion.DatabaseProperty {
	missingProperties := make(map[string]*notion.DatabaseProperty)
	for reqPropName, reqProp := range properties.RequiredProperties {
		if propType, ok := existentProperties[reqPropName]; !ok || propType != reqProp.Type {
			missingProperties[reqPropName] = reqProp
		}
	}

	return missingProperties
}

func (n *Notion) upsertDatabaseProperties(databaseID string, missingProperties map[string]*notion.DatabaseProperty) error {
	if err := n.database.UpsertDatabaseProperties(databaseID, missingProperties); err != nil {
		return fmt.Errorf("error trying to upsert new database properties: %s", err)
	}

	return nil
}
