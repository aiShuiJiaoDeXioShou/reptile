package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (mr *MyRouter) NewMallRouter(r *gin.Engine) {
	mall := r.Group("/bookmall")
	mall.GET("/", mr.MallRouter)
	mall.GET("/cart", mr.MallCartRouter)
}

func (myrouter *MyRouter) MallRouter(ctx *gin.Context) {
	log.Println("MallRouter")
	ctx.HTML(http.StatusOK, "mall/index.html", gin.H{
		"title": "BookMall",
	})
}

func (myrouter *MyRouter) MallCartRouter(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "vip/index.html", gin.H{
		"title": "Vip",
	})
}
