package notion

import (
	"fmt"
	"os"

	"github.com/dstotijn/go-notion"
)

type NankionClient struct {
	Client *notion.Client
}

func NewNankionClient() *NankionClient{
	secret := os.Getenv("NOTION_CARMO_SECRET")
	if secret == ""{
		fmt.Errorf("notion secret not found!")
		return nil
	}
	return &NankionClient{
		Client: notion.NewClient(secret),
	}
}
