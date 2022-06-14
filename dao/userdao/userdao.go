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
	err := dao.AutoMigrate(&User{})
	if err != nil {
		return
	}
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
		Email    string `gorm:"not null" json:"email" form:"email" query:"email"`
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

// CreateUser 创建一个用户对象
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

// FindUserById 根据ID查询一个对象
func FindUserById(id uint) (*User, error) {
	var user User
	tx := dao.Where("id = ?", id).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

// FindUserByPhone 根据手机号查询用户
func FindUserByPhone(phone string) (*User, error) {
	var user User
	tx := dao.Debug().Where("phone = ?", phone).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

// FindUserPasswordByPhone 根据手机号拿到用户密码
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

// 管理员创建一种角色
func CreateRole(role *Role) (err error) {
	err = dao.Create(role).Error
	return
}

// 创建一种权限
func CreateJur(jur *Jur) (err error) {
	err = dao.Create(jur).Error
	return
}

// 为某个角色赋予权限
func CreateJurRole(jurRole *JurRole) (err error) {
	err = dao.Create(jurRole).Error
	return
}

// 赋予某个用户一种角色
func CreateUserRole(userRole *UserRole) (err error) {
	err = dao.Create(userRole).Error
	return
}

// 更改权限的使用范围
func UpdateJur(jur *Jur) (err error) {
	err = dao.Model(jur).Updates(jur).Error
	return
}

// 删除用户指定角色
func DeleteUserRole(role_id, userId uint) (err error) {
	err = dao.Where("role_id = ? and user_id = ?", role_id, userId).Delete(&UserRole{}).Error
	return
}

// 删除角色指定权限
func DeleteJurRole(jur_id, roleId uint) (err error) {
	err = dao.Where(" jur_id = ? and role_id = ?", jur_id, roleId).Delete(&JurRole{}).Error
	return
}

// id删除指定权限
func DeleteJur(id uint) (err error) {
	err = dao.Where("id = ?", id).Delete(&Jur{}).Error
	return
}

// id删除指定角色
func DeleteRole(id uint) (err error) {
	err = dao.Where("id = ?", id).Delete(&Role{}).Error
	return
}

// 根据一组id查询用户的头像和姓名
func FindUserByIds(ids []uint) (users []User, err error) {
	err = dao.Select("id,head_portrait,name").Where("id in (?)", ids).Find(&users).Error
	return
}

// 修改角色
func UpdateRole(role *Role) (err error) {
	err = dao.Model(role).Updates(role).Error
	return
}

// SelectUserJur 查询用户拥有的所有角色
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

// SelectUserRole 查询用户是否包含该角色
func SelectUserRole(userid uint, roleid uint) bool {
	// 通过用户id去与权限关联的表中取到相应的数据
	var userRole UserRole
	// 拿到用户所有的角色信息
	d := dao.Model(&UserRole{}).Where("user_id = ? and role_id = ?", userid, roleid).First(&userRole)
	if d.Error != nil {
		log.Println(d.Error)
		return false
	}
	if userRole.ID == 0 {
		return false
	}
	return true
}

// SelectRoleJur 查询该角色是否包含该权限
func SelectRoleJur(roleid uint, jurid uint) bool {
	// 通过用户id去与权限关联的表中取到相应的数据
	var jurRole JurRole
	// 拿到用户所有的角色信息
	d := dao.Model(&JurRole{}).Where("role_id = ? and jur_id = ?", roleid, jurid).First(&jurRole)
	if d.Error != nil {
		log.Println(d.Error)
		return false
	}
	if jurRole.ID == 0 {
		return false
	}
	return true
}

// SelectRoleJurByRoleId 根据角色找到对应的权限
func SelectRoleJurByRoleId(roleid uint) []JurRole {
	// 通过用户id去与权限关联的表中取到相应的数据
	var jurRole []JurRole
	// 拿到用户所有的角色信息
	d := dao.Model(&JurRole{}).Where("role_id = ?", roleid).Find(&jurRole)
	if d.Error != nil {
		log.Println(d.Error)
		return nil
	}
	return jurRole
}

// 根据权限名称查询该权限id
func SelectJurIdByName(name string) (id uint, err error) {
	var jur Jur
	d := dao.Model(&Jur{}).Where("name = ?", name).First(&jur)
	if d.Error != nil {
		return 0, d.Error
	}
	return jur.ID, nil
}

// 根据角色名称查询该角色id
func SelectRoleIdByName(name string) (id uint, err error) {
	var role Role
	d := dao.Model(&Role{}).Where("name = ?", name).First(&role)
	if d.Error != nil {
		return 0, d.Error
	}
	return role.ID, nil
}

// SelectUserJurByJurId 查询该用户是否拥有指定权限信息
func SelectUserJurByJurId(userid uint, jurid uint) bool {
	//先通过userid拿到用户全部的角色
	var userRole []UserRole
	err := dao.Model(&UserRole{}).Where("user_id = ?", userid).Find(&userRole).Error
	if err != nil {
		log.Println(err)
		return false
	}

	// 通过权限 id拿到拥有这个权限的所有角色
	var jurRole []JurRole
	err = dao.Model(&JurRole{}).Where("jur_id = ?", jurid).Find(&jurRole).Error
	if err != nil {
		log.Println(err)
		return false
	}

	// 判断用户是否拥有该权限
	for _, v := range userRole {
		for _, vv := range jurRole {
			if v.RoleId == vv.RoleId {
				return true
			}
		}
	}
	return false
}
