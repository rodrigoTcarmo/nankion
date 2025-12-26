package notion

import (
	"fmt"
	"os"

	"github.com/dstotijn/go-notion"
	"github.com/rodrigoTcarmo/nankion/pkg/notion/database"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

var RequiredProperties = map[string]notion.DatabasePropertyType{
	"Transaction Date": notion.DBPropTypeDate,
	"Operation":        notion.DBPropTypeSelect,
	"Destination":      notion.DBPropTypeRichText,
	"Amount":           notion.DBPropTypeNumber,
	"Memo":             notion.DBPropTypeRichText,
	"Transaction ID":   notion.DBPropTypeRichText,
}

type NotionLoader struct {
	notionClient database.DatabaseClient
}

func (n *NotionLoader) UploadReport(report *statement.Report) error {
	// 1 - Check if the database exists
	databaseID := os.Getenv("DATABASE_ID")
	if databaseID == "" {
		return fmt.Errorf("database ID not found in environment variables")
	}

	_, err := n.notionClient.GetDatabase(databaseID)
	if err != nil {
		return fmt.Errorf("error trying to get database by ID: %s", err)
	}

	// 2 - Check if all properties are created and exists in the database
	properties, err := n.notionClient.ListDatabaseProperties(databaseID)
	if err != nil {
		return fmt.Errorf("error trying to list database properties: %s", err)
	}

	var existentProperties = map[string]notion.DatabasePropertyType{}
	for _, property := range properties {
		existentProperties[property.Name] = property.Type
	}

	var missingProperties []string
	for reqPropName, reqPropType := range RequiredProperties {
		if propType, ok := existentProperties[reqPropName]; !ok || propType != reqPropType {
			missingProperties = append(missingProperties, reqPropName)
		}
	}

	if len(missingProperties) > 0 {
		return fmt.Errorf("missing properties: %s", missingProperties)
	}

	// 3 - Check if the page already exists in this database (by id)

	return nil
}
