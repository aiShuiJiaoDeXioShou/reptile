package pkmhot

import (
	"gorm.io/gorm"
	"reptile/dao"
	"time"
)

type (
	// HotObject 这个是一个关联表,用于控制其他的表显示数据
	HotObject struct {
		gorm.Model
		// 热门标题
		Title string `json:"title" form:"title" query:"title"`
		// 关联ID
		RelationID uint `json:"article_id" form:"article_id" query:"article_id"`
		// 热门类型,这里面只能在pokemon,article,user,title,activity,book,music,news,around中选填
		Type string `json:"type" form:"type"`
		// image
		ImageUrl string `json:"image_url" form:"image_url" query:"image_url"`
		// 介绍
		Description string `json:"description" form:"description" query:"description"`
	}

	// HotObjectType 这个是HotObjectType的枚举
	hotObjectType struct {
		POKEMON  string `json:"pokemon" yaml:"pokemon" form:"pokemon" query:"pokemon"`
		ARTICLE  string `json:"article" yaml:"article" form:"article" query:"article"`
		USER     string `json:"user" yaml:"user" form:"user" query:"user"`
		TITLE    string `json:"title" yaml:"title" form:"title" query:"title"`
		ACTIVITY string `json:"activity" yaml:"activity" form:"activity" query:"activity"`
		BOOK     string `json:"book" yaml:"book" form:"book" query:"book"`
		MUSIC    string `json:"music" yaml:"music" form:"music" query:"music"`
		NEWS     string `json:"news" yaml:"news" form:"news" query:"news"`
		AROUND   string `json:"around" yaml:"around" form:"around" query:"around"`
	}
)

const (
	// 初始化HotObjectType
	local_POKEMON  = "pokemon"
	local_ARTICLE  = "article"
	local_USER     = "user"
	local_TITLE    = "title"
	local_ACTIVITY = "activity"
	local_BOOK     = "book"
	local_MUSIC    = "music"
	local_NEWS     = "news"
	local_AROUND   = "around"
)

var (
	HotType = &hotObjectType{
		local_POKEMON,
		local_ARTICLE,
		local_USER,
		local_TITLE,
		local_ACTIVITY,
		local_BOOK,
		local_MUSIC,
		local_NEWS,
		local_AROUND,
	}
)

// Insert 插入一条热门数据
func Insert(hot *HotObject) error {
	return dao.DB.Create(hot).Error
}

// Delete 删除一条热门数据，这个是逻辑删除哦
func Delete(hot *HotObject) error {
	return dao.DB.Delete(hot).Error
}

// Update 更新一条热门数据
func Update(hot *HotObject) error {
	return dao.DB.Save(hot).Error
}

// QueryByType QueryByTitle 这个查就有点花式多样
// 查询type是xx的全部热门数据
func QueryByType(typeName string) ([]HotObject, error) {
	var hot []HotObject
	err := dao.DB.Where("type = ?", typeName).Find(&hot).Error
	return hot, err
}

// QueryByNowType QueryByType 根据type查询当天的热门数据
func QueryByNowType(typeName string) ([]HotObject, error) {
	var hots []HotObject
	err := dao.DB.Where("type = ?", typeName).Find(&hots).Error
	// 获取现在年月日的时间
	now := time.Now()
	// 获取现在年月日的时间戳
	nowTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	// 过滤掉当天之前的数据
	for index, hot := range hots {
		if hot.CreatedAt.Before(nowTime) {
			hots = append(hots[:index], hots[index+1:]...)
		}
	}
	return hots, err
}
