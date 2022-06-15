package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reptile/comm/res"
)

func (mr *MyRouter) NewVipRouter(r *gin.Engine) {
	vip := r.Group("/vip")
	vip.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "vip/index.html", res.Ok("Vip"))
	})
}
