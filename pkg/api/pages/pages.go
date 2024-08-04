package pages

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rodrigoTcarmo/nankion/pkg/notion-handler"
)

type PageEndpoit interface{
	GetPages(*gin.Context)
}

type PageQuery string

func (p *PageQuery)GetPages(ctx *gin.Context) {

	query, ok := ctx.GetQuery("pageName")
	if !ok {
		ctx.JSON(422, map[string]string{
			"message": "Missing page name!",
		})
		return
	}
	notionObjects := &notion.NotionObjects{}
	var pages notion.NotionPages = notionObjects
	pages.PageFinder(notion.NotionClient(), query)

	if notionObjects.HttpStatus != http.StatusOK{
		ctx.JSON(404, map[string]string{
			"message": "Notion page or database not found!",
		})
		return
	}

	for _, v := range notionObjects.All {
		err := json.NewEncoder(ctx.Writer).Encode(v)
		if err != nil {
			ctx.JSON(500, map[string]any{
				"error trying to marshal json": err,
			})
		}
	}
}
