package router

import (
	"log"
	"net/http"
	"reptile/dao/userdao"

	"github.com/gin-gonic/gin"
)

// 这个用于跳转到指定页面,你根据函数名称可以清晰的看到对应的页面
type(
	// 这个是项目路由的管理器，用于管理项目的路由
	MyRouter struct{
		
	}
)

// 启动路由
func (myrouter *MyRouter)NewMyRouter(router *gin.Engine) {
	myrouter.InitForHtml(router)
}

// 注册路由
func (myrouter *MyRouter)InitForHtml(router *gin.Engine) {	
	router.GET("/", myrouter.IndexRouter)
	router.GET("/record", myrouter.RecordRouter)
	router.GET("/body/index", myrouter.BodyRouter)
	router.GET("/body/model", myrouter.ModelRouter)
	myrouter.newUserRouter(router)
	myrouter.newWikiRouter(router)
	myrouter.NewTribeRouter(router)
}

// 模块化的首页
func (myrouter *MyRouter)ModelRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "body/model.html", gin.H{
		"title": "model",
	})
}

// 下面是处理路由的函数
func (myrouter *MyRouter)IndexRouter(c *gin.Context) {
	cookie, err := c.Cookie("user_cookie")
	if err != nil {
		c.HTML(http.StatusOK, "default/index.html", gin.H{
			"title": "用户还未登入",
		})
		return
	}
	log.Println(cookie)
	user,err := userdao.FindUserByPhone(cookie)
	if err != nil {
		log.Println(err)
	}
	log.Println("这个是主页显示的user",user)
	c.HTML(http.StatusOK, "default/index.html", gin.H{
		"title": "用户已经登入了",
		"code":200,
		"data":"用户已经登入",
	})
}

func (myrouter *MyRouter)RecordRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "default/record.html", gin.H{
		"title": "Record",
	})
}

func (myrouter *MyRouter)BodyRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "body/index.html", gin.H{
		"title": "Body",
	})
}