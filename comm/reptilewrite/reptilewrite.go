package reptilewrite

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func ReptileWrite(url, filename string) error {
	log.Println("start image")

	// 获取图片转化为图片链接形式
	image, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}

	imageData, errRead := ioutil.ReadAll(image.Body)
	if errRead != nil {
		log.Println(errRead)
		return errRead
	}

	// 如果文件不存在就创造该文件
	err = ioutil.WriteFile(fmt.Sprintf("D:/有天道/img/%s.png", filename), imageData, 0666)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("end image")
	return nil
}
