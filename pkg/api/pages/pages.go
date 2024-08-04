package pages

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rodrigoTcarmo/nankion/pkg/notion-handler"
)

func GetPages(ctx *gin.Context) {

	notionObjects := &notion.NotionObjects{}
	var pages notion.NotionPages = notionObjects
	pages.PageFinder(notion.NotionClient(), "roubo")

	for _, v := range notionObjects.All {
		err := json.NewEncoder(ctx.Writer).Encode(v)
		if err != nil {
			ctx.JSON(500, map[string]any{
				"error trying to marshal json": err,
			})
		}
	}
}
