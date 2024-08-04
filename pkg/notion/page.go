package notion

import (
	"context"
	"fmt"

	"github.com/dstotijn/go-notion"
)

type NotionPages interface {
	PageFinder(*notion.Client, string)
}

type NotionObjects struct {
	Page     []*notion.Page
	Database []*notion.Database
}

func (n *NotionObjects) PageFinder(client *notion.Client, query string) {
	gotResult, err := client.Search(context.Background(), &notion.SearchOpts{Query: query})
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range gotResult.Results {
		switch found := v.(type) {
		case notion.Database:
			n.Database = append(n.Database, &found)
		case notion.Page:
			n.Page = append(n.Page, &found)
		default:
			fmt.Println("unsupported result type gotten")
		}
	}
}
