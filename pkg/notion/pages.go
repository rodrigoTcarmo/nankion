package notion

import (
	"context"
	"fmt"

	"github.com/dstotijn/go-notion"
)

type Page interface {
	SearchPages(*notion.Client)
}

type Pages struct {
	PageId string
}

func (p *Pages) SearchPages(client *notion.Client) {
	pagesFound, err := client.FindPageByID(context.Background(), p.PageId)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pagesFound)
}
