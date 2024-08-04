package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rodrigoTcarmo/nankion/pkg/api/health"
	"github.com/rodrigoTcarmo/nankion/pkg/api/pages"
)

func StartApi() {
	route := gin.Default()
	route.HEAD("/", func(ctx *gin.Context) {
		ctx.Set("Content-Type", "application/json")
	})

	var pageQuery pages.PageQuery = "ETCD"
	var searchQuery pages.PageEndpoit = &pageQuery
	route.GET("/ping", health.GetApiHealth)
	route.GET("/pages", searchQuery.GetPages)
	
	route.Run()

}
