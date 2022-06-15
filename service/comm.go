package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	pkemonservice.StartPkmonService(router)
	upload(router)
}

// 文件上传
func upload(router *gin.Engine) {
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("editormd-image-file")
		log.Println(file.Filename)

		dst := "static/upload/" + file.Filename
		// 判断dst路径是否有文件存在，没有就创造一个文件
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			fmt.Println("文件不存在，创建文件")
			_, err := os.Create(dst)
			if err != nil {
				fmt.Println("创造文件失败:", err)
				return
			}
		}

		// 上传文件至指定的完整文件路径
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			fmt.Println("保存文件失败:", err)
			return
		}

		// 上传文件成功,返回文件的路径
		c.JSON(http.StatusOK, gin.H{
			"success": 1,
			"message": "上传成功",
			"url":     "http://localhost:4001/" + dst,
		})
	})
}
