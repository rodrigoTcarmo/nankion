package health

import "github.com/gin-gonic/gin"

//	@description	This is an API that communicate with Notion API//	@title			Nankion API
//	@version		0.0.1
//	@produce		json

func GetApiHealth(ctx *gin.Context) {
	ctx.JSON(200, map[string]string{
		"message": "Alive",
	})
}
