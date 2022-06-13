package test

import (
	"io/ioutil"
	"log"
	"net/http"
	"reptile/comm/reptilewrite"
	"reptile/dao/pkminfodao"
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

func Test123(t *testing.T) {
	err := reptilewrite.ReptileWrite("https://media.52poke.com/wiki/thumb/7/73/004Charmander.png/300px-004Charmander.png", "小火龙官方图")
	if err != nil {
		return
	}
}

func Test124(t *testing.T) {
	// 获取所有pkminfodao.Pkemon的实例,拿到他们的图片,并且修改他们的image路径
	pkm := pkminfodao.NewPkmDaoManger()
	allPkm := pkm.FindAllPkm()
	for _, p := range *allPkm {
		err := reptilewrite.ReptileWrite("https:"+p.Image, p.Name+"官方图")
		if err != nil {
			continue
		}
		p.Image = "/static/img/" + p.Name + "官方图.png"
	}
}
