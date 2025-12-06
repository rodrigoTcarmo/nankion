package notion

import (
	"context"
	"fmt"

	"github.com/dstotijn/go-notion"
)

type DatabaseClient interface {
	GetDatabase(string) (*notion.Database, error)
	ListDatabaseProperties(string) (*notion.DatabaseProperties, error)
}

type Database struct {
	client *NankionClient
}

func (d *Database) GetDatabase(databaseId string) (*notion.Database, error) {
	database, err := d.client.Client.FindDatabaseByID(context.Background(), databaseId)
	if err != nil {
		return nil, fmt.Errorf("error trying to return Database by ID: %s", err)
	}

	return &database, nil	
}

func (d *Database) ListDatabaseProperties(databaseId string) (*notion.DatabaseProperties, error) {
	database, err := d.GetDatabase(databaseId)
	if err != nil {
		return nil, err
	}

	return &database.Properties, nil

}

func NewDatabase() *Database {
	return &Database{
		client: NewNankionClient(),
	}
}
