package pkemonservice

import (
	"github.com/gin-gonic/gin"
	"reptile/comm/res"
	"strconv"
)

// 这里写的是宝可梦图鉴信息的service
func StartPkmonService(router *gin.Engine) {
	rg := router.Group("/pkm")
	{
		// 获取宝可梦图鉴信息
		rg.GET("/:name", GetPkmonInfo)
		// 拿到宝可梦世代的信息
		rg.GET("/g/:generation", GetPkmonInfoByGeneration)
		// 拿到所有宝可梦的信息
		rg.GET("/", GetAllPkmonInfo)
		// 模糊搜索宝可梦的信息
		rg.GET("/search/:name", GetPkmonInfoByNameLike)
	}
}

// 根据名字查询宝可梦信息
func GetPkmonInfo(ctx *gin.Context) {
	pkm, err := pki.FindPkmByNameLike(ctx.Param("name"))
	if err != nil {
		ctx.JSON(200, res.Error(err.Error()))
		return
	}
	ctx.JSON(200, res.Ok(pkm))
}

// 根据世代拿到宝可梦的信息
func GetPkmonInfoByGeneration(ctx *gin.Context) {
	generation := ctx.Param("generation")
	// generation 强转为int类型
	generationInt, err := strconv.Atoi(generation)
	if err != nil {
		ctx.JSON(200, "参数格式错误")
		return
	}
	pkm, err := pki.FindPkmByGeneration(generationInt)
	if err != nil {
		ctx.JSON(200, res.Error(err.Error()))
		return
	}
	ctx.JSON(200, res.Ok(pkm))
}

// 拿到所有宝可梦的信息
func GetAllPkmonInfo(ctx *gin.Context) {
	pkm, err := pki.FindPkmByGeneration(0)
	if err != nil {
		ctx.JSON(200, res.Error(err.Error()))
		return
	}
	ctx.JSON(200, res.Ok(pkm))
}

// 模糊搜索宝可梦的信息
func GetPkmonInfoByNameLike(ctx *gin.Context) {
	pkm, err := pki.FindPkmByNameLike(ctx.Param("name"))
	if err != nil {
		ctx.JSON(200, res.Error(err.Error()))
		return
	}
	ctx.JSON(200, res.Ok(pkm))
}
