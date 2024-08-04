package main

import (
	"fmt"

	"github.com/rodrigoTcarmo/nankion/pkg/api"
	"github.com/rodrigoTcarmo/nankion/pkg/notion"
)

func main() {
	fmt.Println("Nankion app!")
	var pages notion.NotionPages = &notion.NotionObjects{}
	pages.PageFinder(notion.NotionClient(), "ETCD")
	api.StartApi()
}
