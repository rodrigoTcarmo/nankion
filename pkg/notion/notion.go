package notion

import (
	"fmt"

	"github.com/dstotijn/go-notion"
	notionclient "github.com/rodrigoTcarmo/nankion/pkg/notion/client"
	"github.com/rodrigoTcarmo/nankion/pkg/notion/database"
)

type Notion struct {
	client   *notion.Client
	database database.Database
}

func NewNotionLoader() *Notion {
	return &Notion{
		client:   notionclient.NewClient(),
		database: database.NewDatabase(),
	}
}

// UploadReport uploads a report to the Notion database
func (n *Notion) UploadReport(databaseID string) error {
	_, err := n.database.GetDatabase(databaseID)
	if err != nil {
		return fmt.Errorf("error trying to get database by ID: %s", err)
	}

	// 2 - Check if all properties are created and exists in the database
	properties, err := n.database.ListDatabaseProperties(databaseID)
	if err != nil {
		return fmt.Errorf("error trying to list database properties: %s", err)
	}

	var existentProperties = map[string]notion.DatabasePropertyType{}
	for _, property := range properties {
		existentProperties[property.Name] = property.Type
	}

	// Create a copy of RequiredProperties to track missing ones
	missingProperties := make(map[string]*notion.DatabaseProperty)
	for reqPropName, reqProp := range RequiredProperties {
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

	return nil
}
