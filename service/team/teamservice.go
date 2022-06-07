package team

import (
	"log"
	"reptile/dao"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Note struct {
	Title string `form:"title" json:"title"`
	Command string `form:"command" json:"command"`
	Content string `form:"content" json:"content"`
}

func CharactersAdd(router *gin.Engine) {
	var note Note
	router.POST("/add/note",func (c *gin.Context) {
		c.ShouldBind(&note)
		log.Println(note)
		c.JSON(200, gin.H{
			"message": "Note added successfully!",
		})
	})
}

func RecordDataOpring(router *gin.Engine){
	// 添加数据
	var recordRouter = router.Group("/record")
	recordRouter.POST("/add",func(c *gin.Context){
		var road dao.RoadEntity
		c.ShouldBind(&road)
		log.Println("从前端传来的note值",road)
		dao.InsertRecord(road)
		record := dao.SelectRecord(0)
		c.JSON(200,gin.H{
			"code":200,
			"data":record,
		})
	})

	// 删除一条数据
	recordRouter.POST("/delete",func(c *gin.Context){
		var road dao.RoadEntity
		c.ShouldBind(&road)
		log.Println("从前端传来的note值",road)
		dao.DeleteRecord(road)
		record := dao.SelectRecord(0)
		c.JSON(200,gin.H{
			"code":200,
			"data":record,
		})
	})

	// 查询指定页码数据
	recordRouter.GET("/select",func(c *gin.Context){
		page := c.Query("page")
		// 将page转化为int类型
		pageInt,err := strconv.Atoi(page)
		if err != nil{
			log.Println("record/select--->","page转化为int类型失败")
			c.JSON(500,gin.H{
				"code":500,
				"data":"page转化为int类型失败",
			})
			return
		}
		record := dao.SelectRecord(pageInt)
		c.JSON(200,gin.H{
			"code":200,
			"data":record,
		})
	})
}