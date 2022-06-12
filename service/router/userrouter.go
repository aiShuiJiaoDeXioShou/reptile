package router

import (
	"log"
	"net/http"
	"reptile/dao/userdao"
	"reptile/middleware/usermiddleware"

	"github.com/gin-gonic/gin"
)

func (myrouter *MyRouter) newUserRouter(router *gin.Engine) {
	userrouter := router.Group("/user")
	// 注册
	userrouter.GET("/register", myrouter.RegisterRouter)
	// 登录
	userrouter.GET("/login", myrouter.LoginRouter)
	// 个人中心
	userrouter.GET("/setting", myrouter.SeetingRouter)
	// 个人中心
	userrouter.GET("/personal", usermiddleware.UserLoginJudgeMiddleware, myrouter.PersonalRouter)
	// 写文章
	userrouter.GET("/write", usermiddleware.UserLoginJudgeMiddleware, myrouter.WriteRouter)
	// 后台管理页面
	userrouter.GET("/adminview", func(context *gin.Context) {
		// 只有管理员才能访问
		usermiddleware.JurMiddleware(context, []string{"root"})
	}, func(context *gin.Context) {
		context.HTML(http.StatusOK, "admin/home_admin.html", gin.H{
			"title": "后台管理页面",
		})
	})
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
		log.Println("解析中间件失败")
	}
	log.Println("这个是当前用户的user", value.(*userdao.User))
	c.HTML(http.StatusOK, "user/personal.html", gin.H{
		"title": "用户个人中心",
	})
}

func (myrouter *MyRouter) WriteRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "user/write.html", gin.H{
		"title": "写文章",
	})
}
