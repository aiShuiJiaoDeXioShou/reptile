package userservice

import (
	"log"
	"net/http"
	"reptile/comm"
	"reptile/dao/userdao"

	"time"

	"strings"

	"github.com/gin-gonic/gin"
)

// 路由
func UserSerice(router *gin.Engine) {

	r := router.Group("/user")

	r.GET("/", getUser)

	r.GET("/:phone", getUserByPhone)

	// 创建一个用户
	r.POST("/", userRegister)

	// 更新用户信息
	r.PUT("/",userUpdate)

	// 用户登入验证
	r.POST("/login",userLogin)
}

// 根据token值获取user
func getUser(ctx *gin.Context) {
	comm.TokenHandler(ctx)
	token := ctx.GetHeader("Authorization")
	s := strings.Split(token, " ")
	uid, err := comm.ParseToken("YANGTENG", s[1])
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{
			"code": 0,
			"msg":  "token错误",
		})
		return
	}
	iduid := uint(uid["uid"].(float64))
	user, err := userdao.FindUserById(iduid)
	if err != nil {
		log.Println(err)
	}
	comm.RestApi(200, ctx, user)
}

// 根据电话号码获取User
func getUserByPhone(ctx *gin.Context) {
	phone := ctx.Param("phone")
	user, err := userdao.FindUserByPhone(phone)
	if err != nil {
		log.Println(err)
		return
	}
	comm.RestApi(200, ctx, user)
}

// 用户注册
func userRegister(ctx *gin.Context) {
	var user userdao.User
	ctx.BindJSON(&user)
	userdao.CreateUser(&user)
	comm.RestApi(200, ctx, user)
}

// 用户登录验证
func userLogin(ctx *gin.Context) {
	var user userdao.User
	ctx.ShouldBind(&user)
	daouser, err := userdao.FindUserByPhone(user.Phone)
	if err != nil {
		log.Println(err)
		comm.RestApi(500, ctx, "该用户账号不存在")
		return
	}
	if comm.PasswordEncrypt("YANGTENG", user.Password) != daouser.Password {
		comm.RestApi(500, ctx, "密码错误")
		return
	}

	// 登入成功后返回token值和把用户储存在Coke里面
	//获取当前时间
	now := time.Now()
	//将当前时间转化为int64
	timestamp := now.Unix()
	token, err := comm.GetToken("YANGTENG", timestamp, 604800, daouser.ID)
	if err != nil {
		log.Println(err)
		comm.RestApi(500, ctx, "生成token失败")
		return
	}

	// 将user存入Cookie,设置过期时间为7天
	ctx.SetCookie("token", token, 604800, "/", "localhost", false, true)
	ctx.HTML(http.StatusOK, "default/index.html", gin.H{
		"token": token,
		"code":  200,
		"msg":   "success",
	})
}

// 更新用户
func userUpdate(ctx *gin.Context) {
	var user userdao.User
	ctx.BindJSON(&user)
	err := userdao.UpdateUser(&user)
	if err != nil {
		comm.RestApi(201,ctx,"用户更新操作失败！")
		return
	}
	comm.RestApi(200, ctx, user)
}