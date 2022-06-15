package userdao

import (
	"reptile/dao/hotdao/pkmhot"
	"reptile/dao/pkminfodao"
	"time"
)

// 关注一个人
func Follow(friend *Friend) error {
	err := dao.Create(&friend).Error
	return err
}

// 取消关注一个人
func Unfollow(selfid uint, id string) error {
	err := dao.Where("user_id = ? and friend_id = ?", selfid, id).Delete(&Friend{}).Error
	return err
}

// 查询是否关注了某人
func IsFollow(selfID, friendID uint) (bool, error) {
	var friend Friend
	err := dao.Where("user_id = ? and friend_id = ?", selfID, friendID).First(&friend).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取自己的所有好友
func GetFriends(userID uint) (friends []Friend, err error) {
	err = dao.Where("user_id = ?", userID).Find(&friends).Error
	return
}

// 获取某人最近的动态
func GetRecentDynamic(userID uint) (hotObject []pkmhot.HotObject, err error) {
	// 根据传入的userID去查询该用户的最近动态
	// 查询创作时间是二十四小时之内的文章
	err = dao.Where("user_id = ? and created_at > ?", userID, time.Now().Add(-24*time.Hour)).Find(&pkminfodao.Article{}).Error
	if err != nil {
		return nil, err
	}
	err = dao.Where("user_id = ? and created_at > ?", userID, time.Now().Add(-24*time.Hour)).Find(&pkminfodao.ArticleLike{}).Error
	if err != nil {
		return nil, err
	}
	return
}
