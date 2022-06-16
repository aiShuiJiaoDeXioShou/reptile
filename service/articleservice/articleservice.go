package articleservice

import (
	"log"
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

		// 判断当前用户是否收藏过该文章
		rg.GET("/collect/:id", IsCollectArticle)

		// 判断当前用户是否点赞过该文章
		rg.GET("/like/:id", IsLikeArticle)

		// 取消当前用户收藏或者点赞文章
		rg.DELETE("/collect/:id/:type", UnLikeArticle)
		rg.DELETE("/like/:id/:type", UnLikeArticle)
	}

	// 创建一个评论api列表
	comments := rg.Group("/comment")
	{
		// 添加一个评论
		comments.POST("/", AddComment)
		// 删除一条评论
		comments.DELETE("/:id", DeleteComment)
		// 查看这条评论的所有子评论
		comments.GET("/:id", GetChildrenComment)
		// 获取指定文章所有顶楼的评论
		comments.GET("/top/:id", GetTopComment)
	}

}

// 获取指定文章的顶楼评论
func GetTopComment(ctx *gin.Context) {
	id := ctx.Param("id")
	// 将stringid强转为uid类型
	u, err := comm.ParseUint(id)
	if err != nil {
		ctx.JSON(200, res.FormatError("获取评论失败,格式错误!"))
		return
	}

	// 找到所有顶楼的评论
	comments, err := am.FindCommentByRelationCommentId(0, u)
	if err != nil {
		ctx.JSON(200, res.FormatError("获取评论失败,格式错误!"))
		return
	}

	// 根据评论用户id查询用户信息
	var userids []uint

	// 创造返回数据的变量
	var childrenComment []gin.H
	for _, comment := range comments {

		userids = append(userids, comment.UserId)
		user, _ := userdao.FindUserById(comment.UserId)

		childrenComment = append(childrenComment, gin.H{
			"id":      comment.ID,
			"content": comment.Content,
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"image": user.HeadPortrait,
			},
			"create_time": comment.CreatedAt,
		})
	}

	// 查询用户信息
	users, err := userdao.FindUserByIds(userids)

	if err != nil {
		ctx.JSON(200, res.FormatError("获取评论失败,找不到用户!"))
		return
	}

	ctx.JSON(200, res.Ok(gin.H{
		"comments":  comments,
		"users":     users,
		"commentxs": childrenComment,
	}))
}

// 获取子评论
func GetChildrenComment(ctx *gin.Context) {
	id := ctx.Param("id")
	// 将stringid强转为uid类型
	u, err := comm.ParseUint(id)
	if err != nil {
		ctx.JSON(200, res.FormatError("获取评论失败,格式错误!"))
		return
	}
	// 找到所有与此楼相关的评论
	comments, err := am.FindCommentAllChildren(u)

	if err != nil {
		ctx.JSON(200, res.FormatError("获取评论失败,格式错误!"))
		return
	}

	// 根据评论用户id查询用户信息
	var userids []uint
	// 创造返回数据的变量
	var childrenComment []gin.H
	for _, comment := range comments {

		userids = append(userids, comment.UserId)
		user, _ := userdao.FindUserById(comment.UserId)
		touser, _ := userdao.FindUserById(comment.ToUserID)

		childrenComment = append(childrenComment, gin.H{
			"id":      comment.ID,
			"content": comment.Content,
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"image": user.HeadPortrait,
			},
			"touser": gin.H{
				"image": touser.HeadPortrait,
				"name":  touser.Name,
				"id":    touser.ID,
			},
			"create_time": comment.CreatedAt,
		})
	}

	// 查询用户信息
	users, err := userdao.FindUserByIds(userids)

	if err != nil {
		ctx.JSON(200, res.FormatError("获取评论失败,找不到用户!"))
		return
	}

	ctx.JSON(200, res.Ok(gin.H{
		"comments":        comments,
		"users":           users,
		"childrenComment": childrenComment,
	}))
}

func DeleteComment(ctx *gin.Context) {
	id := ctx.Param("id")
	// 将stringid强转为uid类型
	u, err := comm.ParseUint(id)
	if err != nil {
		ctx.JSON(200, res.FormatError("删除评论失败,格式错误!"))
		return
	}
	err = am.DeleteComment(u)
	if err != nil {
		ctx.JSON(200, res.FormatError("删除评论失败,格式错误!"))
		return
	}
	ctx.JSON(200, res.Ok("删除评论成功!"))
}

func AddComment(ctx *gin.Context) {
	// 获得当前的用户对象
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(200, res.FormatError("评论失败,请先登录!"))
		return
	}
	user := value.(*userdao.User)
	// 获取文章的评论
	var comment pkminfodao.ArticleComment
	err := ctx.ShouldBindJSON(&comment)
	comment.UserId = user.ID
	if err != nil {
		ctx.JSON(200, res.FormatError("获取评论失败,格式错误!"))
		return
	}

	err = am.AddComment(&comment)
	if err != nil {
		ctx.JSON(200, res.FormatError("评论失败,格式错误!"))
		return
	}

	ctx.JSON(200, res.Ok("评论成功!"))
}

func LikeArticle(ctx *gin.Context) {
	// 获取当前空间的用户LikeArticle
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(200, res.FormatError("点赞失败,请先登录!"))
		return
	}
	user := value.(*userdao.User)

	// 获取文章id
	id := ctx.Param("id")
	// 将stringid强转为uid类型
	u, err := comm.ParseUint(id)
	if err != nil {
		ctx.JSON(200, res.FormatError("点赞失败,格式错误!"))
		return
	}

	// 获取该条消息是什么类型
	typestr := ctx.Param("type")
	// 强转为int类型
	typeint, err := strconv.Atoi(typestr)
	if err != nil {
		ctx.JSON(200, res.FormatError("点赞失败,格式错误!"))
	}

	// 调用dao层的方法
	err = am.AddLike(&pkminfodao.ArticleLike{
		UserId:    user.ID,
		Type:      typeint,
		ArticleId: u,
	})
	if err != nil {
		ctx.JSON(200, res.FormatError("点赞失败,格式错误!"))
		return
	}
	ctx.JSON(200, res.Ok("点赞成功!"))
}

// 取消该收藏信息
func UnLikeArticle(ctx *gin.Context) {
	// 获取当前空间的用户LikeArticle
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(200, res.FormatError("取消点赞失败,请先登录!"))
		return
	}
	user := value.(*userdao.User)
	// 获取路径中文章的id
	id := ctx.Param("id")
	// 将stringid强转为uid类型
	u, err := comm.ParseUint(id)
	if err != nil {
		ctx.JSON(200, res.FormatError("取消失败,格式错误!"))
		return
	}
	// 获取该条消息是什么类型
	typestr := ctx.Param("type")
	// 强转为int类型
	typeint, err := strconv.Atoi(typestr)
	if err != nil {
		ctx.JSON(200, res.FormatError("取消操作失败,格式错误!"))
	}
	// 调用dao层的方法
	err = am.DeleteLikeOrCollect(&pkminfodao.ArticleLike{
		UserId:    user.ID,
		Type:      typeint,
		ArticleId: u,
	})
	if err != nil {
		log.Println(err)
		ctx.JSON(200, res.FormatError("取消失败,格式错误!"))
		return
	}
	ctx.JSON(200, res.Ok("取消成功!"))
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

// 判断当前用户是否收藏过该文章
func IsCollectArticle(ctx *gin.Context) {
	s := ctx.Param("id")
	// 将stringid强转为uid类型
	articleId, err := comm.ParseUint(s)
	if err != nil {
		ctx.JSON(200, res.FormatError("查询收藏失败,格式错误!"))
		return
	}

	// 获取当前用户
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(200, res.FormatError("查询收藏失败,你还未登录,没有查询收藏的权限!"))
		return
	}
	user := value.(*userdao.User)

	// 查询数据库里面是否有该条数据,0的话是点赞，1是收藏
	isCollect, err2 := am.IsCollect(user.ID, articleId, 1)
	if err2 != nil {
		ctx.JSON(401, res.Ok("未收藏"))
		return
	}
	if isCollect {
		ctx.JSON(200, res.Ok("已收藏"))
	} else {
		ctx.JSON(401, res.Error("未收藏"))
	}

}
// 判断当前用户是否点赞过该文章
func IsLikeArticle(ctx *gin.Context) {
	s := ctx.Param("id")
	// 将stringid强转为uid类型
	articleId, err := comm.ParseUint(s)
	if err != nil {
		ctx.JSON(200, res.FormatError("查询收藏失败,格式错误!"))
		return
	}

	// 获取当前用户
	value, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(200, res.FormatError("查询收藏失败,你还未登录,没有查询收藏的权限!"))
		return
	}
	user := value.(*userdao.User)

	// 查询数据库里面是否有该条数据,0的话是点赞，1是收藏
	isCollect, err2 := am.IsCollect(user.ID, articleId, 0)
	if err2 != nil {
		ctx.JSON(401, res.Ok("未收藏"))
		return
	}
	if isCollect {
		ctx.JSON(200, res.Ok("已收藏"))
	} else {
		ctx.JSON(401, res.Error("未收藏"))
	}
}