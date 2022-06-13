package test

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	log.Println("server start")
	imagUrl := "https://media.52poke.com/wiki/thumb/7/73/004Charmander.png/300px-004Charmander.png"
	get, err2 := http.Get(imagUrl)
	if err2 != nil {
		log.Println(err2)
		return
	}
	imageData, err := ioutil.ReadAll(get.Body)
	if err != nil {
		log.Println(err)
		return
	}
	// 如果文件不存在就创造该文件

	err = ioutil.WriteFile("../static/img/test.png", imageData, 0666)
	if err != nil {
		log.Println(err)
		return
	}
}
