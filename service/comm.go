package service

import (
	"fmt"
	"log"
	"net/http"
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
	upload(router)
}

// 文件上传
func upload(router *gin.Engine) {
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		dst := "./static/upload/" + file.Filename
		// 上传文件至指定的完整文件路径
		c.SaveUploadedFile(file, dst).Error()

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
}
