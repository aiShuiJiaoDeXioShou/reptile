package main

import (
	"encoding/json"
	io "io/ioutil"
	"reptile/pkmreptile"
)

// import "reptile/service"

func main() {
	// make一个管道
	chanel := make(chan bool)
	var ok bool
	go func() {
		// service.StartUp()
		pkr := pkmreptile.NewPkmReptile()
		m := pkr.GrabPkomAllMap()
		b, _ := json.Marshal(m)
		io.WriteFile("D:\\4.goCode\\space-working\\reptile\\static\\json\\pkomallmap.json", b, 0666)
		chanel <- true
	}()
	if ok {
		print("程序执行完毕！")
	}
	
}
