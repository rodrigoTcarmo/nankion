package pages

import (
	"github.com/gin-gonic/gin"
)

func GetHi(ctx *gin.Context) {
	ctx.JSON(200, map[string]any{
		"message": "hi hi hi",
	})
}
