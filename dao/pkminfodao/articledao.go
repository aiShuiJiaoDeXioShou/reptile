package pkminfodao

import (
	"reptile/dao"

	"gorm.io/gorm"
)

// 更新数据库表的信息
func init(){
	dao.DB.AutoMigrate(&Article{}, &ArticleComment{},&ArticleLike{},&ArticleComment{})
}

func NewArticleManager()(am *ArticleManager){
	am = &ArticleManager{}
	am.Dao = dao.DB
	return
}

type (
	Article struct {
		gorm.Model
		Title string `gorm:"type:text;not null" json:"title" validate:"required" query:"title"`
		Content string `gorm:"type:text;null" json:"content" validate:"required" query:"content"`
		// 文章类型
		Type string `gorm:"type:text;null" json:"type" validate:"required" query:"type"`
		// 文章标签
		Tag string `gorm:"type:text;null" json:"tag" validate:"required" query:"tag"`
		// 文章作者
		Author string `gorm:"type:text;null" json:"author" validate:"required" query:"author"`
		// 文章作者id
		AuthorId uint `gorm:"not null" json:"author_id" validate:"required" query:"author_id"`
	}
	ArticleManager struct {
		Dao *gorm.DB
	}
	// 点赞收藏表
	ArticleLike struct {
		gorm.Model
		ArticleId uint `gorm:"not null" json:"article_id" validate:"required" query:"article_id"`
		Type int `gorm:"not null" json:"type" validate:"required" query:"type"`
		UserId uint `gorm:"not null" json:"user_id" validate:"required" query:"user_id"`
	}
	// 评论表
	ArticleComment struct {
		gorm.Model
		// 这个是关联的评论楼层，如果为零就是顶楼
		RelationCommentId uint `json:"relation_comment_id" validate:"required" query:"relation_comment_id"`
		ArticleId uint `gorm:"not null" json:"article_id" validate:"required" query:"article_id"`
		Content string `gorm:"type:text;null" json:"content" validate:"required" query:"content"`
		UserId uint `gorm:"not null" json:"user_id" validate:"required" query:"user_id"`
	}
)

// 添加一篇文章
func (am *ArticleManager) AddArticle(article *Article) (err error) {
	err = am.Dao.Create(article).Error
	return
}

// 查找自己所有的文章
func (am *ArticleManager) FindAllArticle(userId uint) (articles []Article, err error) {
	err = am.Dao.Where("author_id = ?", userId).Find(&articles).Error
	return
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

// 指定关联楼层的评论
func (am *ArticleManager) FindCommentByRelationCommentId(relationCommentId uint) (comments []ArticleComment, err error) {
	err = am.Dao.Where("relation_comment_id = ?", relationCommentId).Find(&comments).Error
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