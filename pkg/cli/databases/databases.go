package databases

import (
	"encoding/json"
	"fmt"

	notiondatabase "github.com/rodrigoTcarmo/nankion/pkg/notion/database"
)

type DatabaseCLI struct {
	database notiondatabase.DatabaseClient
}

func NewDatabaseCLI() *DatabaseCLI {
	return &DatabaseCLI{
		database: notiondatabase.NewDatabase(),
	}
}

func (d *DatabaseCLI) GetDatabase(databaseID string) {
	database, err := d.database.GetDatabase(databaseID)
	if err != nil {
		fmt.Printf("error trying to get database: %s", err)
	}

	databaseOutput, err := json.MarshalIndent(&database, "", " ")
	if err != nil {
		fmt.Printf("error trying to encode database info into JSON: %s", err)
	}
	fmt.Println(string(databaseOutput))
}

func (d *DatabaseCLI) ListDatabaseProperties(databaseID string) {
	database, err := d.database.ListDatabaseProperties(databaseID)
	if err != nil {
		fmt.Printf("error trying to get database properties: %s", err)
	}

	databasePropertiesOutput, err := json.MarshalIndent(&database, "", " ")
	if err != nil {
		fmt.Printf("error trying to enconde database properties info into JSON: %s", err)
	}
	fmt.Println(string(databasePropertiesOutput))
}
