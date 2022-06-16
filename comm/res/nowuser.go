package comm

import (
	"github.com/gin-gonic/gin"
	"reptile/dao/userdao"
)

func NowUser(ctx *gin.Context) *userdao.User {
	// 获取当前的用户对象
	value, exists := ctx.Get("current_user")
	if !exists {
		return nil
	}
	user := value.(*userdao.User)
	return user
}
