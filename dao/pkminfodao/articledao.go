package pkminfodao

import (
	"log"
	"reptile/dao"

	"gorm.io/gorm"
)

// 更新数据库表的信息
func init() {
	dao.DB.AutoMigrate(&Article{}, &ArticleComment{}, &ArticleLike{}, &ArticleComment{})
}

func NewArticleManager() *ArticleManager {
	return &ArticleManager{
		Dao: dao.DB,
	}
}

type (
	Article struct {
		gorm.Model
		Title   string `gorm:"type:text;not null" json:"title" validate:"required" query:"title" form:"title"`
		Content string `form:"content" gorm:"type:text;null" json:"content" validate:"required" query:"content"`
		// 文章类型
		Type string `gorm:"type:text;null" json:"type" validate:"required" query:"type"`
		// 文章标签
		Tag string `form:"tag" gorm:"type:text;null" json:"tag" validate:"required" query:"tag"`
		// 文章作者
		Author string `form:"author" gorm:"type:text;null" json:"author" validate:"required" query:"author"`
		// 文章作者id
		AuthorId uint `form:"author_id" gorm:"not null" json:"author_id" validate:"required" query:"author_id"`
	}
	ArticleManager struct {
		Dao *gorm.DB
	}
	// 点赞收藏表
	ArticleLike struct {
		gorm.Model
		ArticleId uint `gorm:"not null" json:"article_id" validate:"required" query:"article_id"`
		Type      int  `gorm:"not null" json:"type" validate:"required" query:"type"`
		UserId    uint `gorm:"not null" json:"user_id" validate:"required" query:"user_id"`
	}
	// 评论表
	ArticleComment struct {
		gorm.Model
		// 这个是关联的评论楼层，如果为零就是顶楼
		RelationCommentId uint   `json:"relation_comment_id" validate:"required" query:"relation_comment_id"`
		ArticleId         uint   `gorm:"not null" json:"article_id" validate:"required" query:"article_id"`
		Content           string `gorm:"type:text;null" json:"content" validate:"required" query:"content"`
		// 评论用户
		UserId uint `gorm:"not null" json:"user_id" validate:"required" query:"user_id"`
		// 被评论用户
		ToUserID uint `gorm:"not null" json:"to_user_id" validate:"required" query:"to_user_id"`
	}
)

// 添加一篇文章
func (am *ArticleManager) AddArticle(article *Article) error {
	err := am.Dao.Create(article).Error
	return err
}

// 根据id查找文章
func (am *ArticleManager) FindArticleById(articleId uint) (article Article, err error) {
	err = am.Dao.Where("id = ?", articleId).Find(&article).Error
	return
}

// 查找自己所有的文章
func (am *ArticleManager) FindAllArticle(userId uint) ([]Article, error) {
	var articles []Article
	err := am.Dao.Where("author_id = ?", userId).Find(&articles).Error
	return articles, err
}

// 删除自己的指定文章
func (am *ArticleManager) DeleteArticle(articleId uint) (err error) {
	err = am.Dao.Where("id = ?", articleId).Delete(&Article{}).Error
	return
}

// 修改自己的文章
func (am *ArticleManager) UpdateArticle(article *Article) (err error) {
	err = am.Dao.Save(article).Error
	return
}

// 根据tag值推荐文章
func (am *ArticleManager) FindArticleByTag(tag string) (articles []Article, err error) {
	articles = []Article{}
	err = am.Dao.Where("tag = ?", tag).Find(&articles).Error
	return
}

// 根据分页查找文章
func (am *ArticleManager) FindArticleByPage(page int) (articles []Article, err error) {
	err = am.Dao.Limit(10).Offset(page).Find(&articles).Error
	return
}

// 根据分页查找自己的文章
func (am *ArticleManager) FindArticleByPageAndUserId(page int, userId uint) (articles []Article, err error) {
	err = am.Dao.Where("author_id = ?", userId).Limit(10).Offset(page).Find(&articles).Error
	return
}

// 添加一条评论
func (am *ArticleManager) AddComment(comment *ArticleComment) (err error) {
	err = am.Dao.Create(comment).Error
	return
}

// 删除一条评论
func (am *ArticleManager) DeleteComment(commentId uint) (err error) {
	err = am.Dao.Where("id = ?", commentId).Delete(&ArticleComment{}).Error
	return
}

// 评论无法修改，该文章下面的所有评论
func (am *ArticleManager) FindCommentByArticleId(articleId uint) (comments []ArticleComment, err error) {
	err = am.Dao.Where("article_id = ?", articleId).Find(&comments).Error
	return
}

// 评论无法修改，该文章下面的所有评论
func (am *ArticleManager) FindCommentByRelationCommentId(relationCommentId uint, article_id uint) (comments []ArticleComment, err error) {
	err = am.Dao.Where("relation_comment_id = ? and article_id = ?", relationCommentId, article_id).Find(&comments).Error
	return
}

// 获取指定楼层的评论
func (am *ArticleManager) FindCommentChildren(comment uint) (comments []ArticleComment, err error) {
	err = am.Dao.Where("relation_comment_id = ?", comment).Find(&comments).Error
	return
}

// 获取指定楼层所有的评论
func (am *ArticleManager) FindCommentAllChildren(comment uint) (comments []ArticleComment, err error) {
	var commentsChildrenTop []ArticleComment
	err = am.Dao.Where("relation_comment_id = ?", comment).Find(&commentsChildrenTop).Error

	for _, comment := range commentsChildrenTop {

		commentsChildren, commentsChildrenErr := am.FindCommentAllChildren(comment.ID)

		if commentsChildrenErr != nil {
			log.Print(commentsChildrenErr)
			return
		}

		comments = append(comments, comment)

		comments = append(comments, commentsChildren...)
	}

	return
}

// 判断指定关联楼层是否有评论
func (am *ArticleManager) IsCommentByRelationCommentId(relationCommentId uint) (isComment bool, err error) {
	var comment ArticleComment
	err = am.Dao.Where("relation_comment_id = ?", relationCommentId).First(&comment).Error
	if err != nil {
		isComment = true
	}
	return
}

// 添加一条点赞
func (am *ArticleManager) AddLike(like *ArticleLike) (err error) {
	err = am.Dao.Create(like).Error
	return
}

// 添加一条收藏
func (am *ArticleManager) AddCollect(collect *ArticleLike) (err error) {
	err = am.Dao.Create(collect).Error
	return
}

// 删除一条收藏或者点赞
func (am *ArticleManager) DeleteLikeOrCollect(like *ArticleLike) (err error) {
	err = am.Dao.Debug().Where("article_id = ? and user_id = ? and type = ?",like.ArticleId,like.UserId,like.Type).Delete(like).Error
	return
}

// 判断当前用户是否收藏过该文章
func (am *ArticleManager) IsCollect(userId, articleId uint, typestr int) (isCollect bool, err error) {
	var collect ArticleLike
	err = am.Dao.Debug().Where("user_id = ? and article_id = ? and type = ?", userId, articleId, typestr).First(&collect).Error
	if err != nil {
		isCollect = false
	}
	isCollect = true
	return
}

// 获取指定用户所有收藏Or喜欢的文章
func (am *ArticleManager) FindCollectOrLikeByUserId(userId uint, typestr string) (collects []ArticleLike, err error) {
	err = am.Dao.Where("user_id = ? and type = ?", userId, typestr).Find(&collects).Error
	return
}
