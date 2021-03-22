package router

import (
	"YozComment/controller/comment"
	"YozComment/controller/manage"
	"YozComment/middleware"
	"YozComment/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter is setup router setting
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()
	engine.LoadHTMLFiles("templates/index.html", "templates/manage.html", "templates/login.html")
	engine.Use(middleware.LoggerMiddleware())

	if util.Config.CROSEnabled == true {
		engine.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}

			c.Next()
		})
	}

	engine.Static("/static", "./templates/static")
	engine.StaticFile("client.js", "./templates/client.js")

	engine.GET("/index.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	api := engine.Group("/api")
	api.GET("/page", comment.GetComment)
	api.POST("/comment", comment.Save)

	// manageApi.POST("/comment", comment.Save)
	manageApi := engine.Group(util.Config.ManageRouter, middleware.AuthCheck())
	manageApi.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "manage.html", gin.H{})
	})
	manageApi.GET("/comment/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "manage.html", gin.H{})
	})
	manageApi.GET("/comment/page", manage.GetPage)
	manageApi.POST("/comment/delete", manage.Delete)

	return engine
}
