package pages

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rodrigoTcarmo/nankion/pkg/notion"
)

const (
	pageId                    = "pageId"
	missingPageIdErrorMessage = "Missing page ID!"
)

func GetPage(ctx *gin.Context) {
	notionPageObjects := &notion.NotionPageObjects{}
	var pages notion.NotionPages = notionPageObjects

	pageId, ok := ctx.GetQuery(pageId)
	if !ok {
		ctx.AbortWithStatusJSON(422, map[string]string{
			"message": missingPageIdErrorMessage,
		})
		return
	}
	pageFound, err := pages.SearchPages(notion.NotionClient(), pageId)
	if err != nil {
		ctx.AbortWithStatusJSON(400, map[string]any{
			"message":    err.Error(),
			"Page Found": pageFound,
		})
		return
	}

	err = json.NewEncoder(ctx.Writer).Encode(pageFound)
	if err != nil {
		ctx.AbortWithStatusJSON(500, map[string]any{
			"error trying to marshal json": err,
		})
		return
	}
}
