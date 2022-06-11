package service

import (
	"log"
	"net/http"
	rt "reptile/service/router"

	"github.com/gin-gonic/gin"
)
func StartUp(){
	router := gin.Default()
	router.LoadHTMLGlob("templates/*/*")
	router.Static("/static","static")
	handleError(router)
	// new一个对象是吧，笑死
	myrouter := new(rt.MyRouter)
	myrouter.NewMyRouter(router)
	StartData(router)
	router.Run(":4001")
}

// 处理常见错误
func handleError(router *gin.Engine){
	// 加载404错误页面
    router.NoRoute(func(c *gin.Context) {
        // 实现内部重定向
		c.Redirect(http.StatusMovedPermanently, "/static/html/404.html")
    })
	router.Use(errorHttp)
}

// errorHttp 统一500错误处理函数
func errorHttp(c *gin.Context) {
    defer func() {
        if r := recover(); r != nil {
            //打印错误堆栈信息
            log.Printf("panic: %v\n", r)
            //封装通用json返回
            c.Redirect(http.StatusMovedPermanently, "/static/html/500.html")
        }
    }()
    c.Next()
}

