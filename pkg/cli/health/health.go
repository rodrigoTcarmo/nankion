package health

import "github.com/gin-gonic/gin"

func GetApiHealth(ctx *gin.Context) {
	ctx.JSON(200, map[string]string{
		"message": "Alive",
	})
}
