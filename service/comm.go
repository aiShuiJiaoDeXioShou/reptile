package service

import (
	"reptile/service/userservice"
	"reptile/service/team"

	"github.com/gin-gonic/gin"
)

func StartData(router *gin.Engine) {
	team.CharactersAdd(router)
	team.RecordDataOpring(router)
	userservice.UserSerice(router)
}