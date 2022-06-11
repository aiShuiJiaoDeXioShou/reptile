package main

import (
	"encoding/json"
	io "io/ioutil"
	"reptile/dao/pkminfodao"
	"reptile/pkmreptile"
)

func main() {
	// make一个bool类型的通道
	/* c := make(chan *[]pkminfodao.Pokemon)
	go 爬取数据(c) */
	爬取数据储存为map文件()
	
	/* b := <-c
	if len(*b) > 0  {
		// 这里执行数据库操作
		pdm := pkminfodao.NewPkmDaoManger()
		log.Println("开始插入数据库",len(*b))
		pdm.InsertPkm(b)
		log.Println("插入数据成功")
	} */

	// service.StartUp()
}

func 爬取数据(c chan *[]pkminfodao.Pokemon) {
	// service.StartUp()
	pkr := pkmreptile.NewPkmReptile()
	m := pkr.GrabPkomAll()
	b, _ := json.Marshal(m)
	io.WriteFile("D:\\4.goCode\\space-working\\reptile\\static\\json\\pkomallarr.json", b, 0666)
	c <- m
}

func 爬取数据储存为map文件() {
	// service.StartUp()
	pkr := pkmreptile.NewPkmReptile()
	m := pkr.GrabPkomAllMap()
	b, _ := json.Marshal(m)
	io.WriteFile("D:\\4.goCode\\space-working\\reptile\\static\\json\\pkomallmap.json", b, 0666)
}
