package usermiddleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reptile/comm"
	"reptile/comm/res"
	"reptile/dao/userdao"
)

// UserLoginJudgeMiddleware 验证用户是否满足登入权限
func UserLoginJudgeMiddleware(c *gin.Context) {
	logJudge(c)
}

// JurMiddleware 权限判断
func JurMiddleware(c *gin.Context, jurs []string) {
	logJudge(c)
	jurJudge(c, jurs)
}

// 角色判断中间件
func RoleMiddleware(c *gin.Context, roles []string) {
	logJudge(c)
	roleJudge(c, roles)
}

// 登入判断
func logJudge(c *gin.Context) {
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

	// 登入之后将值储存到局部域里面
	c.Set("current_user", user)
}

// 角色判断
func roleJudge(c *gin.Context, roles []string) {
	// 判断用户是否具有某个权限,默认为使用这个中间件的路由组件都登入了
	// 如果是管理员路由,那么就要判断是否是管理员
	// 取到当前用户
	get, exists := c.Get("current_user")
	if !exists {
		c.JSON(401, res.Unlogin("用户还未登入"))
		c.Abort()
		return
	}

	user := get.(*userdao.User)

	jurmanager := NewJurJudgeManager(roles, user, 1)

	err := jurmanager.Judge()

	if err != nil {
		c.JSON(403, res.Unauthorized("你没有拥有该角色"))
		c.Abort()
		return
	}

	c.Next()
}

// Jur 权限判断
func jurJudge(c *gin.Context, jurs []string) {
	// 判断用户是否具有某个权限,默认为使用这个中间件的路由组件都登入了
	// 如果是管理员路由,那么就要判断是否是管理员
	// 取到当前用户
	get, exists := c.Get("current_user")
	if !exists {
		c.JSON(401, res.Unlogin("用户还未登入"))
		c.Abort()
		return
	}

	user := get.(*userdao.User)

	jurmanager := NewJurJudgeManager(jurs, user)

	err := jurmanager.Judge()

	if err != nil {
		c.JSON(403, res.Unauthorized("没有权限访问"))
		c.Abort()
		return
	}

	c.Next()
}
