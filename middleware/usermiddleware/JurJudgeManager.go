package usermiddleware

import (
	"errors"
	"log"
	"reptile/dao/userdao"
)

type (
	JurJudgeManager struct {
		// 必须要满足的权限
		JurList []string
		// 判断是基于什么做判断的目前可选项为：基于角色1，基于基础权限0，默认为0
		JurType int
		// 这个是要满足权限的用户
		User *userdao.User
	}
)

// 创建一个新JurJudgeManager对象
func NewJurJudgeManager(jurList []string, user *userdao.User, jurtype ...int) *JurJudgeManager {
	var judgeManager = &JurJudgeManager{
		JurList: jurList,
		JurType: 0,
		User:    user,
	}
	if len(jurtype) > 0 {
		judgeManager.JurType = jurtype[0]
	}
	return judgeManager
}

func (jmg *JurJudgeManager) Judge() error {

	// 判断一下他是不是拥有root权限，如果有直接放行
	if userdao.SelectUserJurByJurId(jmg.User.ID, 1) || userdao.SelectUserRole(jmg.User.ID, 1) {
		return nil
	}

	for _, jurstr := range jmg.JurList {

		// 根据权限名称拿到权限id如果type值为1则拿到角色id
		var juedgeid uint
		var err error

		if jmg.JurType == 0 {
			// 拿到权限id
			juedgeid, err = userdao.SelectJurIdByName(jurstr)
		} else if jmg.JurType == 1 {
			// 拿到角色id
			juedgeid, err = userdao.SelectRoleIdByName(jurstr)
		}

		if err != nil {
			log.Println(err)
			return err
		}

		var isJur = false
		if jmg.JurType == 0 {
			// 根据拿到的权限id去查询该用户是否拥有该权限
			isJur = userdao.SelectUserJurByJurId(jmg.User.ID, juedgeid)
		} else if jmg.JurType == 1 {
			// 根据拿到的角色id去查询该用户是否拥有该角色
			isJur = userdao.SelectUserRole(jmg.User.ID, juedgeid)
		}

		if !isJur {
			return errors.New("没有权限")
		}

	}
	// 如果都没问题就返回true
	return nil
}
