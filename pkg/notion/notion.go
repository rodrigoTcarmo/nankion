package notion

import (
	"fmt"
	"os"

	"github.com/rodrigoTcarmo/nankion/pkg/notion/database"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

func UploadReport(report *statement.Report) error {
	databaseClient := database.NewDatabase()

	// 1 - Check if the database exists
	databaseID := os.Getenv("DATABASE_ID")
	if databaseID == "" {
		return fmt.Errorf("database ID not found in environment variables")
	}

	_, err := databaseClient.GetDatabase(databaseID)
	if err != nil {
		return fmt.Errorf("error trying to get database by ID: %s", err)
	}

	// 2 - Check if all properties are created and exists in the database
	properties, err := databaseClient.ListDatabaseProperties(databaseID)
	if err != nil {
		return fmt.Errorf("error trying to list database properties: %s", err)
	}

	for _, property := range *properties {
		fmt.Println(property.Name)
	}

	// 3 - Check if the page already exists in this database (by id)

	return nil
}
