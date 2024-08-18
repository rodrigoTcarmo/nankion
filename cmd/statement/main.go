package main

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/rodrigoTcarmo/nankion/cmd/statement/docs"
	"github.com/rodrigoTcarmo/nankion/pkg/api/health"
	"github.com/rodrigoTcarmo/nankion/pkg/api/items"
	"github.com/rodrigoTcarmo/nankion/pkg/api/pages"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           Gin Book Service
// @version         1.0
// @description     A book management service API in Go using Gin framework.
// @termsOfService  https://tos.santoshk.dev

// @contact.name   Santosh Kumar
// @contact.url    https://twitter.com/sntshk
// @contact.email  sntshkmr60@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
func main() {
	fmt.Println("Nankion app!")

	//fmt Print runtime goroot
	fmt.Println(runtime.GOROOT())

	route := gin.Default()

	route.GET("/ping", health.GetApiHealth)
	route.GET("/items", items.SearchItems)
	route.GET("/page", pages.GetPage)
	route.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	route.Run()
}
