package databases

import (
	"encoding/json"
	"fmt"

	"github.com/rodrigoTcarmo/nankion/pkg/notion"
)

func GetDatabase(databaseID string) {
	var db notion.DatabaseClient = notion.NewDatabase()

	database, err := db.GetDatabase(databaseID)
	if err != nil {
		fmt.Printf("error trying to get database: %s", err)
	}

	databaseOutput, err := json.MarshalIndent(&database, "", " ")
	if err != nil {
		fmt.Printf("error trying to encode database info into JSON: %s", err)
	}
	fmt.Println(string(databaseOutput))
}

func ListDatabaseProperties(databaseID string) {
	var db notion.DatabaseClient = notion.NewDatabase()

	database, err := db.ListDatabaseProperties(databaseID)
	if err != nil {
		fmt.Printf("error trying to get database properties: %s", err)
	}

	databasePropertiesOutput, err := json.MarshalIndent(&database, "", " ")
	if err != nil {
		fmt.Printf("error trying to enconde database properties info into JSON: %s", err)
	}
	fmt.Println(string(databasePropertiesOutput))
}
