package reptilewrite

import (
	"io/ioutil"
	"log"
	"net/http"
)

func ReptileWrite(url string) error {
	log.Println("start image")

	get, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}

	imageData, errRead := ioutil.ReadAll(get.Body)
	if errRead != nil {
		log.Println(errRead)
		return errRead
	}

	// 如果文件不存在就创造该文件
	err = ioutil.WriteFile("../static/img/test.png", imageData, 0666)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("end image")
	return nil
}
