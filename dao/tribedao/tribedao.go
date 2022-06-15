package tribedao

import (
	"gorm.io/gorm"
	"reptile/dao"
	"reptile/dao/hotdao/pkmhot"
)

type TribeDao struct {
	DB *gorm.DB
}

func NewTribeDao() *TribeDao {
	return &TribeDao{
		DB: dao.DB,
	}
}

// 根据type值获取榜单信息
func (t *TribeDao) GetTribeByTypeAndPage(typeId string, page int) (tribeList []pkmhot.HotObject, err error) {
	err = t.DB.Debug().Where("type = ?", typeId).Offset(page * 10).Limit(10).Find(&tribeList).Error
	return
}

// 根据type值获取榜单信息
func (t *TribeDao) GetTribeByType(typeId string) (tribeList []pkmhot.HotObject, err error) {
	err = t.DB.Where("type = ?", "hot").Find(&tribeList).Error
	return
}
