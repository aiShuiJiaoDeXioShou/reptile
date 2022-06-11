package articleservice

import (
	"reptile/dao/pkminfodao"
	"reptile/dao/userdao"
	"reptile/middleware/usermiddleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

var(
	am *pkminfodao.ArticleManager
)

func init(){
	am = pkminfodao.NewArticleManager()
}

// 文章服务类
func StartArticleService(router *gin.Engine) {
	 rg := router.Group("/article",usermiddleware.JurMiddleware)
	 
	 // 添加一个文章
	 rg.POST("/",AddArticle)

	 // 删除一个文章
	 rg.DELETE("/:id",DeleteArticle)

	 // 修改该文章
	 rg.PUT("/:id",UpdateArticle)
}

func AddArticle(ctx *gin.Context){
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