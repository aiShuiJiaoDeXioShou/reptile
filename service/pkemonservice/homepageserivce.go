package pkemonservice

import (
	"github.com/gin-gonic/gin"
	"log"
	"reptile/comm/res"
	"reptile/dao/hotdao/pkmhot"
	"reptile/dao/pkminfodao"
)

var (
	pki *pkminfodao.PkmDaoManger
)

// NewPokemonWikiHomeService 初始化宝可梦首页服务
func NewPokemonWikiHomeService(router *gin.Engine) {
	pki = pkminfodao.NewPkmDaoManger()
	homepage := router.Group("/homepage")
	{
		homepage.GET("/left_data", ShowHomePageLeftData)
		homepage.GET("/right_data", ShowHomePageRightData)
		homepage.GET("/music", ShowHomePageMusic)
		homepage.GET("/book", ShowHomePageBook)
		homepage.GET("/title", ShowHomePageTitle)
	}
}

// ShowHomePageTitle 就是那个首页的火热滚动条
func ShowHomePageTitle(ctx *gin.Context) {
	// 查询当天的火热滚动条
	nowType, err := pkmhot.QueryByNowType(pkmhot.HotType.TITLE)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, res.Fail("什么鬼啊今天程序员这么懒吗,没有创造一条数据!!!"))
	}
	ctx.JSON(200, res.Ok(nowType))
}

// ShowHomePageBook 就是那个显示书籍的数据
func ShowHomePageBook(ctx *gin.Context) {
	byType, err := pkmhot.QueryByType(pkmhot.HotType.BOOK)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, res.Fail("没有找到一本书,你是不是忘记存库了"))
	}
	ctx.JSON(200, res.Ok(byType))
}

func ShowHomePageMusic(ctx *gin.Context) {
	// TODO
}

// 显示首页右边数据
func ShowHomePageRightData(ctx *gin.Context) {
	nowType, err := pkmhot.QueryByNowType(pkmhot.HotType.POKEMON)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, res.Fail("没有找到一只宝可梦,你是不是忘记存库了"))
	}
	var pkmIds []uint
	// 根据关联对象拿到宝可梦id
	for _, object := range nowType {
		// 拿到object里面所有的关联ID
		pkmIds = append(pkmIds, object.RelationID)
	}
	// 根据宝可梦id查询宝可梦的数据
	pkms := pki.FindPkmByIds(pkmIds)
	ctx.JSON(200, res.Ok(gin.H{
		"pkms":    pkms,
		"now_pkm": nowType,
	}))
}

// 显示首页左边数据
func ShowHomePageLeftData(ctx *gin.Context) {
	// 近日资讯的数据
	news, err := pkmhot.QueryByType(pkmhot.HotType.NEWS)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, res.Error("没有找到近日资讯的数据"))
	}
	// 周边信息
	around, err := pkmhot.QueryByType(pkmhot.HotType.AROUND)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, res.Fail("没有找到周边信息,你是不是忘记存库了"))
	}
	// 没有错误的话,就返回数据
	ctx.JSON(200, res.Ok(gin.H{
		"news":   news,
		"around": around,
	}))
}
