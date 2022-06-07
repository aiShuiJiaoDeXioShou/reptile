package comm

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 验证token的中间件
func TokenHandler(c *gin.Context) {
	// 获取token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "token is empty",
		})
		c.Abort()
		return
	}
	// 验证token
	if !CheckToken(token) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "token is error",
		})
		c.Abort()
		return
	}
	c.Next()
}

// 验证token格式是否正确
func CheckToken(token string)bool{
	s := strings.Split(token, " ")
	if len(s) != 2 {
		return false
	}
	if s[0] != "Bearer" {
		return false
	}
	return true
}

func RestApi(code int, c *gin.Context,data interface{}) {
	c.JSON(code, gin.H{
		"message": "success",
		"data": data,
		"code": code,
	})
}