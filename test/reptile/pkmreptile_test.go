package test

import (
	"encoding/json"
	io "io/ioutil"
	"log"
	"regexp"
	"reptile/comm"
	"reptile/pkmreptile"
	"testing"
	"time"
)

func Test2(t *testing.T) {
	log.Println("test2")
	pkr := pkmreptile.NewPkmReptile()
	m := pkr.GrabBasicsPkm()
	b, _ := json.Marshal(m)
	io.WriteFile("D:\\4.goCode\\space-working\\reptile\\static\\json\\pkm.json", b, 0666)
}

func Test3(t *testing.T){
	reg1 := regexp.MustCompile(`^第(.*)世代`)
	if reg1.MatchString("海第鬼"){
		log.Println("match")
	}else {
		log.Println("not match")
	}
}

func Test4(t *testing.T){
	go Grab()
	select{
	case <-time.After(time.Second * 10):
		log.Println("timeout 这是十秒之后的事情了！")
	}
}

func Grab(){
	pkr := pkmreptile.NewPkmReptile()
	m := pkr.GrabPkomAllMap()
	b, _ := json.Marshal(m)
	io.WriteFile("D:\\4.goCode\\space-working\\reptile\\static\\json\\pkomallmap.json", b, 0666)
}

func Test6(t *testing.T){
	u, err := comm.ParseUint("1")
	if err != nil {
		log.Println(err)
	}
	log.Println(u)
}