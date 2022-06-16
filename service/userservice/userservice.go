package userservice

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reptile/comm"
	"reptile/comm/res"
	"reptile/dao/pkminfodao"
	"reptile/dao/userdao"
	"reptile/middleware/usermiddleware"

	"time"

	"strings"

	"github.com/gin-gonic/gin"
)

var am = pkminfodao.NewArticleManager()

// 路由
func UserSerice(router *gin.Engine) {

	r := router.Group("/user")

	{
		r.GET("/", getUser)

		r.GET("/:phone", getUserByPhone)

		// 创建一个用户
		r.POST("/", userRegister)

		// 更新用户信息
		r.PUT("/", userUpdate)

		// 用户登入验证
		r.POST("/login", userLogin)

		// 判断用户是否为管理员
		r.GET("/admin", func(context *gin.Context) {
			usermiddleware.RoleMiddleware(context, []string{"admin"})
			context.JSON(200, res.Ok("该用户拥有管理员权限"))
		})

		// 退出登入UserSerice
		r.GET("/logout", userLogout)

		// 获取当前用户收藏的文章
		r.GET("/collect/article", usermiddleware.UserLoginJudgeMiddleware, UserCollectArticle)

		// 获取当前用户喜欢的文章
		r.GET("/like/article", usermiddleware.UserLoginJudgeMiddleware, UserLikeArticle)

		// 上传或者修改用户头像
		uploadUserHeader(router)
	}

}



// 文件上传
func uploadUserHeader(router *gin.Engine) {
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/userheader/upload", usermiddleware.UserLoginJudgeMiddleware,func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		dst := "static/upload/" + file.Filename
		// 判断dst路径是否有文件存在，没有就创造一个文件
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			fmt.Println("文件不存在，创建文件")
			_, err := os.Create(dst)
			if err != nil {
				fmt.Println("创造文件失败:", err)
				res.Fail("创造文件失败")
				return
			}
		}

		// 上传文件至指定的完整文件路径
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			fmt.Println("保存文件失败:", err)
			return
		}

		// 上传成功修改用户头像信息
		u := res.NowUser(c)
		if u == nil {
			res.Unlogin("请先登录")
			return
		}
		err2 := userdao.UpdateUserHeadPortrait(u.ID, "/"+dst)
		if err2 != nil {
			fmt.Println("修改用户头像失败:", err2)
			res.Fail("修改用户头像失败")
			return
		}

		// 上传文件成功,返回文件的路径
		c.JSON(http.StatusOK, gin.H{
			"success": 1,
			"message": "上传成功",
			"url":     "http://localhost:4001/" + dst,
		})
	})
}

// 获取当前用户喜欢的文章
func UserLikeArticle(ctx *gin.Context) {
	// 获取当前用户的id
	user := res.NowUser(ctx)
	if user == nil {
		res.Unlogin("请先登录")
		return
	}
	id := user.ID
	// 获取当前用户喜欢的文章
	articleIDs, err := am.FindCollectOrLikeByUserId(id, "0")
	if err != nil {
		log.Println(err)
		res.Error("获取文章失败")
		return
	}
	// 根据article的文章id获取文章的详细信息
	var articles []pkminfodao.Article
	for _, v := range articleIDs {
		article, err := am.FindArticleById(v.ArticleId)
		articles = append(articles, article)
		if err != nil {
			log.Println(err)
		}
	}
	comm.RestApi(200, ctx, articles)
}

// 获取当前用户收藏的文章
func UserCollectArticle(ctx *gin.Context) {
	// 获取当前用户的id
	user := res.NowUser(ctx)
	if user == nil {
		res.Unlogin("请先登录")
		return
	}
	id := user.ID
	// 获取当前用户喜欢的文章
	articleIDs, err := am.FindCollectOrLikeByUserId(id, "1")
	if err != nil {
		log.Println(err)
		res.Error("获取文章失败")
		return
	}
	// 根据article的文章id获取文章的详细信息
	var articles []pkminfodao.Article
	for _, v := range articleIDs {
		article, err := am.FindArticleById(v.ArticleId)
		articles = append(articles, article)
		if err != nil {
			log.Println(err)
		}
	}
	comm.RestApi(200, ctx, articles)
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
	ctx.BindJSON(&user)
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
	token, err := comm.GetToken("YANGTENG", timestamp, 24*60*60*7, daouser.ID)
	if err != nil {
		log.Println(err)
		comm.RestApi(500, ctx, "生成token失败")
		return
	}

	// 将user存入Cookie,设置过期时间为7天
	ctx.SetCookie("token", token, 24*60*60*7, "/", "localhost", false, true)
	ctx.JSON(200, gin.H{
		"code":  200,
		"msg":   "登入成功",
		"token": token,
	})
}

// 更新用户
func userUpdate(ctx *gin.Context) {
	var user userdao.User
	ctx.BindJSON(&user)
	err := userdao.UpdateUser(&user)
	if err != nil {
		comm.RestApi(201, ctx, "用户更新操作失败！")
		return
	}
	comm.RestApi(200, ctx, user)
}

// 退出登入
func userLogout(ctx *gin.Context) {
	// 清除该tonken
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "退出登录成功",
	})
}

// 查询当前用户是否关注了指定id用户
func userIsFollow(ctx *gin.Context) {
	// 获取当前用户的id
	user := res.NowUser(ctx)
	if user == nil {
		res.Unlogin("请先登录")
		return
	}
	id := user.ID
	// 获取指定用户id
	userId := ctx.Param("id")
	// 查询是否关注了指定用户
	isFollow, err := userdao.IsFollow(id, userId)
	if err != nil {
		log.Println(err)
		res.Error("查询失败")
		return
	}
	comm.RestApi(200, ctx, isFollow)
}

// 取消关注一个人
func userCancelFollow(ctx *gin.Context) {
	// 获取当前用户的id
	user := res.NowUser(ctx)
	if user == nil {
		res.Unlogin("请先登录")
		return
	}
	id := user.ID
	// 获取指定用户id
	userId := ctx.Param("id")
	// 取消关注指定用户
	err := userdao.Unfollow(id, userId)
	if err != nil {
		log.Println(err)
		res.Error("取消关注失败")
		return
	}
	res.Ok("取消关注成功")
}