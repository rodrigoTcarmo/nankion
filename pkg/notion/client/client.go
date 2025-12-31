package notion

import (
	"context"
	"os"

	"github.com/dstotijn/go-notion"
)

type NotionClient interface {
	FindDatabaseByID(context.Context, string) (notion.Database, error)
	UpdateDatabase(context.Context, string, notion.UpdateDatabaseParams) (notion.Database, error)
	CreatePage(context.Context, notion.CreatePageParams) (notion.Page, error)
}

func NewClient() NotionClient {
	secret := os.Getenv("NOTION_SECRET")
	if secret == "" {
		panic("notion secret not found!")
	}
	return notion.NewClient(secret)
}
