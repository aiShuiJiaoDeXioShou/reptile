package router

import (
	"log"
	"net/http"
	"reptile/dao/userdao"
	"reptile/middleware/usermiddleware"

	"github.com/gin-gonic/gin"
)

func (myrouter *MyRouter)newUserRouter(router *gin.Engine) {
	userrouter := router.Group("/user")
	userrouter.GET("/register", myrouter.RegisterRouter)
	userrouter.GET("/login", myrouter.LoginRouter)
	userrouter.GET("/setting", myrouter.SeetingRouter)
	userrouter.GET("/personal",usermiddleware.JurMiddleware, myrouter.PersonalRouter)
}

func (myrouter *MyRouter) RegisterRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "user/register.html", gin.H{
		"title": "Register",
	})
}

func (myrouter *MyRouter) LoginRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "user/login.html", gin.H{
		"title": "Login",
	})
}

func (myrouter *MyRouter) SeetingRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "user/setting.html", gin.H{
		"title": "Setting",
	})
}

func (myrouter *MyRouter) PersonalRouter(c *gin.Context) {
	value, exists := c.Get("current_user")
	if !exists {
		log.Println("解析中间件失败");
	}
	log.Println("这个是当前用户的user", value.(*userdao.User))
	c.HTML(http.StatusOK, "user/personal.html", gin.H{
		"title": "用户个人中心",
	})
}

