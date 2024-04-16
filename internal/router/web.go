package router

import (
	"voidcloud-server/controllers"
	"voidcloud-server/middlewares"

	gcors "github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type WebRouter struct {
	*gin.Engine
}

func NewWebRouter() *WebRouter {
	router := &WebRouter{gin.Default()}
	router.Init()
	return router
}

func (engine *WebRouter) Init() {
	// Use Gzip
	engine.Engine.Use(gzip.Gzip(gzip.DefaultCompression))

	// Cors
	engine.Engine.Use(gcors.Default())

	// engine.Engine.Any("/drive", func(c *gin.Context) {
	// 	fmt.Println("URL:", c.Request.URL)
	// 	c.Redirect(http.StatusPermanentRedirect, "http://voidcloud.net:8080/drive")
	// })

	// API 路由
	apiGroup := engine.Engine.Group("api")
	{
		// User
		userGroup := apiGroup.Group("user")
		{
			userGroup.POST("/register", controllers.Register)
			userGroup.POST("/login", controllers.Login)
			userGroup.POST("/validate", middlewares.RequireAuth, controllers.Validate)
		}

		fileGroup := apiGroup.Group("file")
		{
			// fileGroup.GET("/*folder", controllers.GetFile)
			fileGroup.Use(middlewares.RequireAuth)
			// fileGroup.DELETE("/", controllers.DeleteFile)
			// fileGroup.GET("/share", controllers.GetFileShareInfo)
			fileGroup.POST("/share", controllers.ShareFile)
		}
	}
}
