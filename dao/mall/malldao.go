package malldao

import "gorm.io/gorm"

// 商品对象
type Commodity struct {
	gorm.Model
	Name string `json:"name" query:"name" form:"name"`
	// 商品类型
	Type        string `json:"type" query:"type" form:"type"`
	Description string `json:"description" query:"description" form:"description"`
	Price       int64  `json:"price" query:"price" form:"price"`
	Stock       int64  `json:"stock" query:"stock" form:"stock"`
	Status      int64  `json:"status" query:"status" form:"status"`
	Image       string `json:"image" query:"image" form:"image"`
	// 是否为热门商品, 0:否, 1:是, 默认为0
	IsHot int64 `json:"is_hot" query:"is_hot" form:"is_hot" gorm:"default:0"`
	// 是否为新品, 0:否, 1:是, 默认为0
	IsNew int64 `json:"is_new" query:"is_new" form:"is_new" gorm:"default:0"`
	// 是否为精品, 0:否, 1:是, 默认为0
	IsBest int64 `json:"is_best" query:"is_best" form:"is_best" gorm:"default:0"`
}

// 地址
type Address struct {
	gorm.Model
	UserID    uint   `json:"user_id" query:"user_id" form:"user_id"`
	Name      string `json:"name" query:"name" form:"name"`
	Phone     string `json:"phone" query:"phone" form:"phone"`
	Address   string `json:"address" query:"address" form:"address"`
	IsDefault int64  `json:"is_default" query:"is_default" form:"is_default"`
}

// 订单
type Order struct {
	gorm.Model
	UserID     uint  `json:"user_id" query:"user_id" form:"user_id"`
	CartID     uint  `json:"cart_id" query:"cart_id" form:"cart_id"`
	AddressID  uint  `json:"address_id" query:"address_id" form:"address_id"`
	TotalPrice int64 `json:"total_price" query:"total_price" form:"total_price"`
	Status     int64 `json:"status" query:"status" form:"status"`
}

// 购物车
type Cart struct {
	gorm.Model
	UserID      uint  `json:"user_id" query:"user_id" form:"user_id"`
	CommodityID uint  `json:"commodity_id" query:"commodity_id" form:"commodity_id"`
	Count       int64 `json:"count" query:"count" form:"count"`
	// 购物车总金额
	TotalPrice int64 `json:"total_price" query:"total_price" form:"total_price"`
	// 是否有折扣，默认为0
	IsDiscount int64 `json:"is_discount" query:"is_discount" form:"is_discount" gorm:"default:0"`
}

// 对商品的增删改查
type CommodityDao interface {
	// 查询商品列表
	QueryCommodityList(page, pageSize uint, name string) ([]Commodity, uint, error)
	// 查询商品详情
	QueryCommodityDetail(id uint) (Commodity, error)
	// 添加商品
	AddCommodity(commodity Commodity) error
	// 修改商品
	UpdateCommodity(commodity Commodity) error
	// 删除商品
	DeleteCommodity(id uint) error
}

// 实现该接口
type commodityDaoImpl struct {
	db *gorm.DB
}

// 查询一个商品
func (c *commodityDaoImpl) QueryCommodityDetail(id uint) (Commodity, error) {
	var commodity Commodity
	if err := c.db.Where("id = ?", id).First(&commodity).Error; err != nil {
		return commodity, err
	}
	return commodity, nil
}

// 添加一个商品
func (c *commodityDaoImpl) AddCommodity(commodity Commodity) error {
	if err := c.db.Create(&commodity).Error; err != nil {
		return err
	}
	return nil
}

// 删除一个商品
func (c *commodityDaoImpl) DeleteCommodity(id uint) error {
	if err := c.db.Where("id = ?", id).Delete(&Commodity{}).Error; err != nil {
		return err
	}
	return nil
}

// 修改一个商品
func (c *commodityDaoImpl) UpdateCommodity(commodity Commodity) error {
	if err := c.db.Save(&commodity).Error; err != nil {
		return err
	}
	return nil
}
