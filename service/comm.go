package service

import (
	"reptile/service/articleservice"
	"reptile/service/pkemonservice"
	"reptile/service/team"
	"reptile/service/userservice"

	"github.com/gin-gonic/gin"
)

func StartData(router *gin.Engine) {
	team.CharactersAdd(router)
	team.RecordDataOpring(router)
	userservice.UserSerice(router)
	articleservice.StartArticleService(router)
	pkemonservice.NewPokemonWikiHomeService(router)
}
