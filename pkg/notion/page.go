package notion

import (
	"context"
	"fmt"
	"log"
	"os"
)

func PagePrinter() {
	client := NotionClient()

	page, err := client.FindPageByID(context.Background(), os.Getenv("NOTION_NOTAS_RAPIDAS_PAGE_ID"))
	if err != nil {
		log.Panicf("Unable to get database page: ", err)
	}
	fmt.Println(page.Properties)
}
