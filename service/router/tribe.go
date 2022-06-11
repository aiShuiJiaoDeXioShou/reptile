package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"reptile/middleware/usermiddleware"
)

func (myrouter *MyRouter) NewTribeRouter(router *gin.Engine) {
	rg := router.Group("/tribe")
	rg.Use(usermiddleware.JurMiddleware)
	rg.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "tribe/index.html", gin.H{
			"title": "部落",
		})
	})
}