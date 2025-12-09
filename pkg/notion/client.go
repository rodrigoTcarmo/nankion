package notion

import (
	"os"

	"github.com/dstotijn/go-notion"
)

type NankionClient struct {
	Client *notion.Client
}

func NewNankionClient() *NankionClient{
	secret := os.Getenv("NOTION_CARMO_SECRET")
	if secret == ""{
		panic("notion secret not found!")
	}
	return &NankionClient{
		Client: notion.NewClient(secret),
	}
}
