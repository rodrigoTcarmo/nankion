package notion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dstotijn/go-notion"
)

type Objects []interface{}

type NotionPages interface {
	SearchItems(*notion.Client, string)
}

type NotionObjects struct {
	All        Objects
	Page       []*notion.Page
	Database   []*notion.Database
	HttpStatus int
}

func (n *NotionObjects) SearchItems(client *notion.Client, query string) {
	gotResult, err := client.Search(context.Background(), &notion.SearchOpts{Query: query})
	if err != nil {
		fmt.Println(err)
	}

	if results := len(gotResult.Results); results != 0 {
		asserted := make(Objects, results)

		for i, v := range gotResult.Results {
			switch found := v.(type) {
			case notion.Database:
				n.Database = append(n.Database, &found)
				asserted[i] = found
			case notion.Page:
				n.Page = append(n.Page, &found)
				asserted[i] = found
			default:
				fmt.Println("unsupported result type gotten")
			}
		}
		n.All = asserted
		n.HttpStatus = http.StatusOK
	} else {
		n.HttpStatus = http.StatusInternalServerError
	}
}

func (n *NotionObjects) ObjectFoundIdentifier() {

}
