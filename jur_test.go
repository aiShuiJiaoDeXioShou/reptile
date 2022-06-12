package main

import (
	"reptile/dao/userdao"
	"testing"
)

// 创建一条角色
func Test110(t *testing.T) {
	err := userdao.CreateRole(&userdao.Role{
		Name:        "admin",
		Description: "是这个系统的管理员，拥有所有的权限",
	})
	if err != nil {
		t.Error(err)
	}
}

// 创建一条权限
func Test120(t *testing.T) {
	err := userdao.CreateJur(&userdao.Jur{
		Name:        "root",
		Description: "超级权限",
	})
	if err != nil {
		t.Error(err)
	}
}

// 为admin赋予root权限
func Test130(t *testing.T) {
	err := userdao.CreateJurRole(&userdao.JurRole{
		JurId:  1,
		RoleId: 1,
	})
	if err != nil {
		t.Error(err)
	}
}

// 为用户id为1赋予admin角色
func Test140(t *testing.T) {
	err := userdao.CreateUserRole(&userdao.UserRole{
		UserId: 1,
		RoleId: 1,
	})
	if err != nil {
		t.Error(err)
	}
}

// 查询用户为i的是否拥有admin角色
func Test150(t *testing.T) {
	role := userdao.SelectUserRole(1, 1)
	if !role {
		t.Error("该用户没有admin角色")
		return
	}
	t.Log("该用户拥有admin角色")
}

// 查询该用户是否拥有该权限
func Test160(t *testing.T) {
	// 查询权限name为root的id
	id, err := userdao.SelectJurIdByName("root")
	if err != nil {
		t.Error(err)
		return
	}
	jur := userdao.SelectUserJurByJurId(1, id)
	if !jur {
		t.Error("该用户没有root权限")
		return
	}
	t.Log("该用户拥有root权限")
}
