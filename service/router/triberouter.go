package router

import (
	"log"
	"net/http"
	"reptile/comm"
	"reptile/dao/pkminfodao"
	"reptile/dao/userdao"
	"reptile/middleware/usermiddleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	am = pkminfodao.NewArticleManager()
)

func (myrouter *MyRouter) NewTribeRouter(router *gin.Engine) {
	rg := router.Group("/tribe")
	rg.Use(usermiddleware.UserLoginJudgeMiddleware)

	rg.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "tribe/index.html", gin.H{
			"title": "部落",
		})
	})

	// 修改文章跳转的路径
	rg.GET("/article/put/:id", func(ctx *gin.Context) {
		s := ctx.Param("id")
		id, _ := comm.ParseUint(s)
		a, err := am.FindArticleById(id)
		if err != nil {
			ctx.JSON(500, gin.H{
				"code": 0,
				"msg":  "查询文章失败,请确定是否有该文章的信息！",
			})
			return
		}
		ctx.HTML(200, "user/updatearticle.html", gin.H{
			"code": 200,
			"data": a,
			"msg":  "查询文章成功！",
		})
	})

	rg.GET("/article/:id", func(ctx *gin.Context) {
		// 获取路径参数
		id := ctx.Param("id")
		log.Println("获取到的id是", id)
		var iduid uint

		// 绑定iduint
		ctx.ShouldBindUri(&iduid)
		log.Println("获取到的iduid是", iduid)

		// 将路径id强转为uint类型
		idUint, err := strconv.Atoi(id)
		if err != nil {
			log.Println("路径参数转换失败", err)
			ctx.HTML(http.StatusOK, "tribe/index.html", gin.H{
				"title": "显示文章失败",
			})
			return
		}

		// 根据id从数据库获取文章
		comments, err := am.FindArticleById(uint(idUint))
		if err != nil {
			log.Println("根据id从数据库获取文章失败", err)
			ctx.HTML(201, "tribe/index.html", gin.H{
				"title": "显示文章失败",
			})
			return
		}

		ctx.HTML(http.StatusOK, "combat/article.html", gin.H{
			"title":   "文章",
			"id":      id,
			"article": comments,
		})
	})

	rg.GET("/articleadmin/:pagesize", func(ctx *gin.Context) {
		var pagesize int
		ctx.BindUri(&pagesize)
		value, exists := ctx.Get("current_user")
		if !exists {
			ctx.HTML(http.StatusOK, "tribe/index.html", gin.H{
				"code":  http.StatusUnauthorized,
				"title": "你还未登入，请先登入账号",
			})
			return
		}
		user := value.(*userdao.User)
		articles, err := am.FindArticleByPageAndUserId(pagesize, user.ID)
		if err != nil {
			log.Println("根据id从数据库获取文章失败", err)
			ctx.HTML(201, "tribe/index.html", gin.H{
				"title": "显示文章失败",
			})
			return
		}
		ctx.HTML(200, "user/articleadmin.html", gin.H{
			"code": 200,
			"msg":  "获取文章列表成功",
			"data": articles,
		})
	})
}
