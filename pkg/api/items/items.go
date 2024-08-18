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

//	 ListItems godoc
//		@Summary		Search Items
//		@Description	Search for items
//		@Tags			item
//		@Produce		json
//
//		@Param			id	path		int	true	"Account ID"
//
//		@Success		200	{array}		model.Account
//		@Failure		422	{object}	httputil.HTTPError
//		@Failure		404	{object}	httputil.HTTPError
//		@Failure		500	{object}	httputil.HTTPError
//		@Router			/items/{id} [get]
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
