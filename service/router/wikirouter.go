package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (myrouter *MyRouter) newWikiRouter(router *gin.Engine) {
	wiki := router.Group("/wiki")
	wiki.GET("/index", myrouter.WikiheadRouter)
	wiki.GET("/detail", myrouter.WikiDetailRouter)
	wiki.GET("/atlas", myrouter.WikiAtlasRouter)
	wiki.GET("/atlaslist", myrouter.WikiAtlasListRouter)
}

// 跳转到维基页面 ps：这个是首页来着
func (myrouter *MyRouter) WikiheadRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "wikifirst/index.html", gin.H{
		"title": "Wikihead",
	})
}

// 跳转到资讯详情
func (myrouter *MyRouter) WikiDetailRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "wikifirst/detail.html", gin.H{
		"title": "Wikihead",
	})
}

// 全国图鉴
func (myrouter *MyRouter) WikiAtlasRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "wikifirst/atlas.html", gin.H{
		"title": "Wikihead",
	})
}

func (myrouter *MyRouter) WikiAtlasListRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "wikifirst/atlaslist.html", gin.H{
		"title": "Wikihead",
	})
}

// 跳转到最近资讯

// 跳转到部落

// 跳转到对战室

// 跳转到游戏部落
