package items

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rodrigoTcarmo/nankion/pkg/notion"
)

const (
	itemsQueryKey            = "queryValue"
	missingQueryErrorMessage = "Missing page name!"
)

func SearchItems(ctx *gin.Context) {

	query, ok := ctx.GetQuery(itemsQueryKey)
	if !ok {
		ctx.AbortWithStatusJSON(422, map[string]string{
			"message": missingQueryErrorMessage,
		})
		return
	}

	notionObjects := &notion.NotionObjects{}
	var items notion.NotionItems = notionObjects
	items.SearchItems(notion.NotionClient(), query)

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
