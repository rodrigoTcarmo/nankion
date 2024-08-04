package items

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rodrigoTcarmo/nankion/pkg/notion-handler"
)

const (
	PageNameQueryKey         = "pageName"
	MissingQueryErrorMessage = "Missing page name!"
)

func SearchItems(ctx *gin.Context) {

	query := getUriQuery(ctx)
	if query == "" {
		return
	}
	notionObjects := &notion.NotionObjects{}
	var pages notion.NotionPages = notionObjects
	pages.SearchItems(notion.NotionClient(), query)

	if notionObjects.HttpStatus != http.StatusOK {
		ctx.AbortWithStatusJSON(404, map[string]string{
			"message": "Notion page or database not found!",
		})
		return
	}

	for _, v := range notionObjects.All {
		err := json.NewEncoder(ctx.Writer).Encode(v)
		if err != nil {
			ctx.AbortWithStatusJSON(500, map[string]any{
				"error trying to marshal json": err,
			})
			return
		}
	}
}

func getUriQuery(ctx *gin.Context) string {
	query, ok := ctx.GetQuery(PageNameQueryKey)
	if !ok {
		ctx.AbortWithStatusJSON(422, map[string]string{
			"message": MissingQueryErrorMessage,
		})
		return ""
	}
	return query
}
