package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rodrigoTcarmo/nankion/pkg/api/health"
	"github.com/rodrigoTcarmo/nankion/pkg/api/pages"
)

func StartApi() {
	route := gin.Default()
	route.GET("/ping", health.GetApiHealth)
	route.GET("/hi", pages.GetHi)
	route.Run()

}
