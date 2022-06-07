package userdao

import (
	"log"
	"reptile/comm"

	daox "reptile/dao"

	"gorm.io/gorm"
)

var(
	dao = daox.DB
)

func init(){
	dao.AutoMigrate(&User{})
	dao.AutoMigrate(&Account{})
	dao.AutoMigrate(&Role{})
	dao.AutoMigrate(&UserRole{})
	dao.AutoMigrate(&Jur{})
	dao.AutoMigrate(&JurRole{})
}

type (
	User struct {
		gorm.Model
		Name string `gorm:"not null" json:"name" form:"name" query:"name" validate:"required"`
		Password string `gorm:"not null" json:"password" form:"password" query:"password" validate:"required"`
		Email string `gorm:"null" json:"email" form:"email" query:"email"`
		Phone string `gorm:"null" json:"phone" form:"phone" query:"phone"`
		// 用户头像
		HeadPortrait string `json:"head_portrait" form:"head_portrait" query:"head_portrait"`
		// 用户权限
		RoleId int `json:"role_id" form:"role_id" query:"role_id"`
		// 签名
		Motto string `json:"motto" form:"motto" query:"motto"`
		// 性别int类型默认为零
		Sex int `json:"sex" form:"sex" query:"sex"`
	}
	Account struct{
		gorm.Model
		UserId uint `json:"user_id" form:"user_id" query:"user_id"`
		Account float64 `json:"account" form:"account" query:"account"`
	}
	Role struct{
		gorm.Model
		Name string `json:"name" form:"name" query:"name"`
		Description string `json:"description" form:"description" query:"description"`
	}
	UserRole struct{
		gorm.Model
		UserId uint `json:"user_id" form:"user_id" query:"user_id"`
		RoleId uint `json:"role_id" form:"role_id" query:"role_id"`
	}
	Jur struct{
		gorm.Model
		Name string `json:"name" form:"name" query:"name"`
		Description string `json:"description" form:"description" query:"description"`
	}
	JurRole struct{
		gorm.Model
		JurId uint `json:"jur_id" form:"jur_id" query:"jur_id"`
		RoleId uint `json:"role_id" form:"role_id" query:"role_id"`
	}
)

// 创建一个用户对象
func CreateUser(user *User) error{
	// 将密码进行加密
	password := comm.PasswordEncrypt("YANGTENG",user.Password)
	user.Password = password
	tx := dao.Create(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 根据ID查询一个对象
func FindUserById(id uint) (*User, error){
	var user User
	tx := dao.Where("id = ?", id).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

// 根据手机号查询用户
func FindUserByPhone(phone string) (*User, error){
	var user User
	tx := dao.Debug().Where("phone = ?", phone).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

// 根据手机号拿到用户密码
func FindUserPasswordByPhone(phone string) (string, error){
	var user User
	tx := dao.Where("phone = ?", phone).First(&user)
	if tx.Error != nil {
		return "", tx.Error
	}
	return user.Password, nil
}

func UpdateUser(user *User) error{
	log.Println("user_sex:",user.Sex)
	tx := dao.Debug().Model(user).Updates(user)
	if tx.Error != nil {
		log.Println(tx.Error)
		return tx.Error
	}
	return nil
}

// 查询用户的权限信息
func SelectUserJur(userid uint){
	
}