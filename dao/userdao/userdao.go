package userdao

import (
	"log"
	"reptile/comm"

	daox "reptile/dao"

	"gorm.io/gorm"
)

var (
	dao = daox.DB
)

func init() {
	dao.AutoMigrate(&User{})
	dao.AutoMigrate(&Account{})
	dao.AutoMigrate(&Role{})
	dao.AutoMigrate(&UserRole{})
	dao.AutoMigrate(&Jur{})
	dao.AutoMigrate(&JurRole{})
}

type (
	// 用户详细信息
	User struct {
		gorm.Model
		Name     string `gorm:"not null" json:"name" form:"name" query:"name" validate:"required"`
		Password string `gorm:"not null" json:"password" form:"password" query:"password" validate:"required"`
		Email    string `gorm:"null" json:"email" form:"email" query:"email"`
		Phone    string `gorm:"null" json:"phone" form:"phone" query:"phone"`
		// 用户头像
		HeadPortrait string `json:"head_portrait" form:"head_portrait" query:"head_portrait"`
		// 用户权限
		RoleId int `json:"role_id" form:"role_id" query:"role_id"`
		// 签名
		Motto string `json:"motto" form:"motto" query:"motto"`
		// 性别int类型默认为零
		Sex int `json:"sex" form:"sex" query:"sex"`
	}
	// 用户账户信息
	Account struct {
		gorm.Model
		UserId  uint    `json:"user_id" form:"user_id" query:"user_id"`
		Account float64 `json:"account" form:"account" query:"account"`
	}
	// 角色表
	Role struct {
		gorm.Model
		Name        string `json:"name" form:"name" query:"name"`
		Description string `json:"description" form:"description" query:"description"`
	}
	// 用户角色表
	UserRole struct {
		gorm.Model
		UserId uint `json:"user_id" form:"user_id" query:"user_id"`
		RoleId uint `json:"role_id" form:"role_id" query:"role_id"`
	}
	// 权限表
	Jur struct {
		gorm.Model
		Name        string `json:"name" form:"name" query:"name"`
		Description string `json:"description" form:"description" query:"description"`
	}
	// 权限角色表
	JurRole struct {
		gorm.Model
		JurId  uint `json:"jur_id" form:"jur_id" query:"jur_id"`
		RoleId uint `json:"role_id" form:"role_id" query:"role_id"`
	}
)

// 创建一个用户对象
func CreateUser(user *User) error {
	// 将密码进行加密
	password := comm.PasswordEncrypt("YANGTENG", user.Password)
	user.Password = password
	tx := dao.Create(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 根据ID查询一个对象
func FindUserById(id uint) (*User, error) {
	var user User
	tx := dao.Where("id = ?", id).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

// 根据手机号查询用户
func FindUserByPhone(phone string) (*User, error) {
	var user User
	tx := dao.Debug().Where("phone = ?", phone).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

// 根据手机号拿到用户密码
func FindUserPasswordByPhone(phone string) (string, error) {
	var user User
	tx := dao.Where("phone = ?", phone).First(&user)
	if tx.Error != nil {
		return "", tx.Error
	}
	return user.Password, nil
}

func UpdateUser(user *User) error {
	log.Println("user_sex:", user.Sex)
	tx := dao.Debug().Model(user).Updates(user)
	if tx.Error != nil {
		log.Println(tx.Error)
		return tx.Error
	}
	return nil
}

// 查询用户的权限信息
func SelectUserJur(userid uint) ([]UserRole, error) {
	// 通过用户id去与权限关联的表中取到相应的数据
	var userRole []UserRole
	// 拿到用户所有的角色信息
	d := dao.Model(&UserRole{}).Where("user_id = ?", userid).Find(&userRole)
	if d.Error != nil {
		return nil, d.Error
	}
	return userRole, nil
}

// 查询用户是否包含该角色
func SelectUserRole(userid uint, roleid uint) (bool, error) {
	// 通过用户id去与权限关联的表中取到相应的数据
	var userRole UserRole
	// 拿到用户所有的角色信息
	d := dao.Model(&UserRole{}).Where("user_id = ? and role_id = ?", userid, roleid).First(&userRole)
	if d.Error != nil {
		return false, d.Error
	}
	if userRole.ID == 0 {
		return false, nil
	}
	return true, nil
}

// 查询该角色是否包含该权限
func SelectRoleJur(roleid uint, jurid uint) (bool, error) {
	// 通过用户id去与权限关联的表中取到相应的数据
	var jurRole JurRole
	// 拿到用户所有的角色信息
	d := dao.Model(&JurRole{}).Where("role_id = ? and jur_id = ?", roleid, jurid).First(&jurRole)
	if d.Error != nil {
		return false, d.Error
	}
	if jurRole.ID == 0 {
		return false, nil
	}
	return true, nil
}

// 根据角色找到对应的权限
func SelectRoleJurByRoleId(roleid uint) ([]JurRole, error) {
	// 通过用户id去与权限关联的表中取到相应的数据
	var jurRole []JurRole
	// 拿到用户所有的角色信息
	d := dao.Model(&JurRole{}).Where("role_id = ?", roleid).Find(&jurRole)
	if d.Error != nil {
		return nil, d.Error
	}
	return jurRole, nil
}
