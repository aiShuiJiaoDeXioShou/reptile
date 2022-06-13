package articleservice

import (
	"reptile/comm"
	"reptile/comm/res"
	"reptile/dao/pkminfodao"
	"reptile/dao/userdao"
	"reptile/middleware/usermiddleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	am *pkminfodao.ArticleManager
)

func init() {
	am = pkminfodao.NewArticleManager()
}

// 文章服务类
func StartArticleService(router *gin.Engine) {
	rg := router.Group("/article", usermiddleware.UserLoginJudgeMiddleware)

	{
		// 添加一个文章
		rg.POST("/", AddArticle)

		// 删除一个文章
		rg.DELETE("/:id", DeleteArticle)

		// 修改该文章
		rg.PUT("/:id", UpdateArticle)

		// 文章点赞
		rg.POST("/like/:id/:type", LikeArticle)

		// 文章收藏
		rg.POST("collect/:id/:type", LikeArticle)
	}

}

func LikeArticle(ctx *gin.Context) {
	// 获取当前空间的用户LikeArticle
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(500, gin.H{
			"code": 0,
			"msg":  "点赞失败,你还未登录,没有点赞的权限!",
		})
		return
	}
	user := value.(*userdao.User)

	// 获取文章id
	id := ctx.Param("id")
	// 将stringid强转为uid类型
	u, err := comm.ParseUint(id)
	if err != nil {
		ctx.JSON(500, res.FormatError("点赞失败,格式错误!"))
		return
	}

	// 获取该条消息是什么类型
	typestr := ctx.Param("type")
	// 强转为int类型
	typeint, err := strconv.Atoi(typestr)
	if err != nil {
		ctx.JSON(500, res.FormatError("点赞失败,格式错误!"))
	}

	// 调用dao层的方法
	err = am.AddLike(&pkminfodao.ArticleLike{
		UserId:    user.ID,
		Type:      typeint,
		ArticleId: u,
	})
}

func AddArticle(ctx *gin.Context) {
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(500, gin.H{
			"code": 0,
			"msg":  "添加文章失败,你还未登录,没有添加文章的权限!",
		})
		return
	}
	user := value.(*userdao.User)

	var article pkminfodao.Article
	ctx.ShouldBind(&article)
	article.AuthorId = user.ID
	article.Author = user.Name

	err := am.AddArticle(&article)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 0,
			"msg":  "添加文章失败",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "添加文章成功",
	})
}

func DeleteArticle(ctx *gin.Context) {
	id := ctx.Param("id")
	// 将stringid强转为uid类型
	i, qerr := strconv.Atoi(id)

	if qerr != nil {
		ctx.JSON(500, gin.H{
			"code": 0,
			"msg":  "删除文章失败",
		})
		return
	}

	err := am.DeleteArticle(uint(i))
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 0,
			"msg":  "删除文章失败",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 1,
		"msg":  "删除文章成功",
	})
}

// 修改文章
func UpdateArticle(ctx *gin.Context) {
	// 从上下文获取当前空间的用户
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(500, gin.H{
			"code": 0,
			"msg":  "修改文章失败,你还未登录,没有修改文章的权限!",
		})
		return
	}
	user := value.(*userdao.User)
	id := ctx.Param("id")
	// 将表单数据绑定到一个结构体上UpdateArticle
	var article pkminfodao.Article
	ctx.ShouldBind(&article)

	// 将stringid强转为uid类型
	i, qerr := strconv.Atoi(id)
	if qerr != nil {
		ctx.JSON(500, gin.H{
			"code": 0,
			"msg":  "修改文章失败",
		})
		return
	}
	article.ID = uint(i)
	article.AuthorId = user.ID
	article.Author = user.Name
	err := am.UpdateArticle(&article)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code": 0,
			"msg":  "修改文章失败",
		})
		return
	}
	// 返回修改成功的消息
	ctx.JSON(200, gin.H{
		"code": 1,
		"msg":  "修改文章成功",
	})
}
