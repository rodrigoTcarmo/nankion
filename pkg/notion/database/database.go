package database

import (
	"context"
	"fmt"

	"github.com/dstotijn/go-notion"
	notionclient "github.com/rodrigoTcarmo/nankion/pkg/notion/client"
)

type Database interface {
	GetDatabase(string) (*notion.Database, error)
	ListDatabaseProperties(string) (notion.DatabaseProperties, error)
	UpsertDatabaseProperties(string, map[string]*notion.DatabaseProperty) error
}

type database struct {
	client *notion.Client
}

func (d *database) GetDatabase(databaseId string) (*notion.Database, error) {
	database, err := d.client.FindDatabaseByID(context.Background(), databaseId)
	if err != nil {
		return nil, fmt.Errorf("error trying to return Database by ID: %s", err)
	}

	return &database, nil
}

func (d *database) ListDatabaseProperties(databaseId string) (notion.DatabaseProperties, error) {
	database, err := d.GetDatabase(databaseId)
	if err != nil {
		return nil, err
	}

	return database.Properties, nil

}

func (d *database) UpsertDatabaseProperties(databaseId string, properties map[string]*notion.DatabaseProperty) error {
	databaseParams := notion.UpdateDatabaseParams{
		Properties: properties,
	}

	_, err := d.client.UpdateDatabase(context.Background(), databaseId, databaseParams)
	if err != nil {
		return err
	}

	return nil
}

func NewDatabase() *database {
	return &database{
		client: notionclient.NewClient(),
	}
}
