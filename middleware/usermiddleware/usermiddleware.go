package usermiddleware

import (
	"log"
	"net/http"
	"reptile/comm"
	"reptile/dao/userdao"

	"github.com/gin-gonic/gin"
)

// 验证用户是否满足权限
func JurMiddleware(c *gin.Context) {
	cookie, err := c.Cookie("token")
	// 如果连token都没有拿到那肯定没有登入
	if err != nil {
		c.HTML(http.StatusOK, "default/index.html", gin.H{
			"error": "用户还未登入",
		})
		return
	}

	userid, tokenerr := comm.ParseToken("YANGTENG", cookie)
	// 验证失败说明token的值是错的或者过期了需要重新登入
	if tokenerr != nil {
		log.Println(tokenerr)
		c.HTML(http.StatusOK, "default/index.html", gin.H{
			"error": "token的值是错的或者过期了需要重新登入",
		})
		return
	}

	// 将string的值转换为uint类型
	uid := uint(userid["uid"].(float64))
	user, err := userdao.FindUserById(uid)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "default/index.html", gin.H{
			"error": "没有找到该token对应的用户",
		})
		return
	}

	// 判断用户是否具有某个权限,默认为使用这个中间件的路由组件都登入了
	// 获取传递的路由地址
	url := c.Request.URL.Path
	if(url == "/admin"){

	}

	// 登入之后将值储存到局部域里面
	c.Set("current_user", user)
}