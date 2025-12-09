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

	queryPageId, ok := ctx.GetQuery(pageId)
	if !ok {
		ctx.AbortWithStatusJSON(422, map[string]string{
			"message": missingPageIdErrorMessage,
		})
		return
	}

	var pg notion.PageClient = notion.NewPage()

	pageFound, code, err := pg.SearchPages(queryPageId)
	if err != nil {
		ctx.AbortWithStatusJSON(code, map[string]any{
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
