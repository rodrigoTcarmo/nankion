package notion

import (
	"os"

	"github.com/dstotijn/go-notion"
)

func NotionClient() *notion.Client{
	return notion.NewClient(os.Getenv("NOTION_CARMO_SECRET"))
}
