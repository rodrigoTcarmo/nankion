package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rodrigoTcarmo/nankion/pkg/api/health"
	"github.com/rodrigoTcarmo/nankion/pkg/api/items"
	"github.com/rodrigoTcarmo/nankion/pkg/api/pages"
)

func StartApi() {
	route := gin.Default()
	route.HEAD("/", func(ctx *gin.Context) {
		ctx.Set("Content-Type", "application/json")
	})

	route.GET("/ping", health.GetApiHealth)
	route.GET("/items", items.SearchItems)
	route.GET("/page", pages.GetPage)
	route.Run()

}
