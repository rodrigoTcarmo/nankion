package notion

import (
	"os"

	"github.com/dstotijn/go-notion"
)

func NewClient() *notion.Client{
	secret := os.Getenv("NOTION_SECRET")
	if secret == ""{
		panic("notion secret not found!")
	}
	return notion.NewClient(secret)
}
