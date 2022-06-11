package test

import (
	"fmt"
	"log"
	"strings"

	"testing"

	"github.com/PuerkitoBio/goquery"
)

func Test1(t *testing.T) {
	baseurl := "https://wiki.52poke.com"
	url := baseurl + "/wiki/%E5%AE%9D%E5%8F%AF%E6%A2%A6%E5%88%97%E8%A1%A8%EF%BC%88%E6%8C%89%E5%85%A8%E5%9B%BD%E5%9B%BE%E9%89%B4%E7%BC%96%E5%8F%B7%EF%BC%89/%E7%AE%80%E5%8D%95%E7%89%88#.E7.AC.AC.E4.B8.80.E4.B8.96.E4.BB.A3"
	d, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	d.Find(".eplist").Each(func(i int, s *goquery.Selection) {
		node := s.Find("td")
		node.Find("a").Each(func(i int, s *goquery.Selection) {
			if s.Text() == "烈咬陆鲨" {
				// 序号
				name := s.Parent().Parent().Children().Eq(0).Text()
				log.Print("序号:",name)
				// 名称
				log.Print("昵称:",s.Text())
				if val, exists := s.Attr("href"); exists {
					拿到指定宝可梦的详细数据(baseurl + val)
				}
				return
			}
		})
	})
}

func 拿到指定宝可梦的详细数据(url string) {
	d, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	// 拿到数据显示的table表单
	table := d.Find(".mw-parser-output").Find(".roundy.a-r.at-c")

	// 通过表单拿到立绘
	image := table.Find("table.roundy.bgwhite.fulltable").Find(".image").Find("img")
	val, exists := image.Attr("data-url")
	if exists {
		log.Println("官方立绘为：", val)
	}

	// 拿到宝可梦的属性值
	s := table.Find("tbody>tr").Eq(4)
	s.Find("a").Each(func(i int, s *goquery.Selection) {
		if str, err := s.Attr("title"); err && strings.Contains(str, "（属性）") {
			{
				log.Println(str)
			}
		}

		if str, err := s.Attr("title"); err && strings.Contains(str, "分类") {
			{
				s2 := s.Parent().Next().Find("td")
				log.Println(s2.Text())
			}
		}
	})
	
	fater := table.Find("table.roundy.bgwhite.fulltable>tbody>tr")

	// 特性
	texin := fater.Eq(3)
	texin.Find("a").Each(func(i int, s *goquery.Selection) {
		if str, err := s.Attr("title"); err && strings.Contains(str, "（特性）") {
			{
				if len(str) > 0 {
					s2 := s.Next().Next()
					if len(s2.Text()) > 0 {
						log.Println(str, s2.Text())
					}else {
						log.Println(str)
					}
				}
			}
		}
	})

	// 经验值
	ex := fater.Eq(4)
	log.Println("经验值:", ex.Find("td").Text())

	// 地区图鉴编号
	// 关都地区
	basepath := "#mw-content-text > div > table.roundy.a-r.at-c > tbody >"
	guandu := d.Find("td.bd-关都.bw-1.roundyright").Eq(0)
	log.Println("关都地区 ",guandu.Text())

	// 成都地区
	chendu := d.Find(".bd-城都.bw-1").Eq(0)
	log.Println("成都地区 ",chendu.Text())

	// 卡洛斯地区
	kls := d.Find(".bd-卡洛斯.bw-1.roundyright").Eq(0)
	log.Println("卡洛斯地区 ",kls.Text())

	// 伽勒尔地区
	jls := d.Find("td.bd-伽勒尔.bw-1.roundyright").Eq(0)
	log.Println("伽勒尔地区 ",jls.Text())

	// 身高
	var height *goquery.Selection
	height = d.Find(fmt.Sprintf("%s tr:nth-child(8) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td",basepath))
	if len(height.Text()) == 0 {
		height = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(8) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td")
	}

	log.Println("身高 ",height.Text())

	// 体重
	var weight *goquery.Selection
	weight = d.Find(fmt.Sprintf("%s tr:nth-child(8) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td",basepath))
	if len(weight.Text()) == 0 {
		weight = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(8) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td")
	}
	log.Println("体重",weight.Text())

	// 体形
	var tixing *goquery.Selection
	tixing = d.Find(fmt.Sprintf("%s tr:nth-child(9) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td > a > img",basepath))
	tixingval, tixingexists := tixing.Attr("data-url")
	if tixingexists {
		log.Println("体形 ",url+tixingval)
	}else {
		tixing = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(9) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td > a > img")
		tixingval, tixingexists := tixing.Attr("data-url")
		if tixingexists {
			log.Println("体形 ",url+tixingval)
		}
	}

	// 脚印
	jiaoyin := d.Find(fmt.Sprintf("%s tr:nth-child(9) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td > a",basepath))
	jiaoyinval, jiaoyinexists := jiaoyin.Attr("data-url")
	if jiaoyinexists {
		log.Println("体形 ",url+jiaoyinval)
	}

	// 图鉴颜色
	var color *goquery.Selection
	color = d.Find(fmt.Sprintf("%s tr:nth-child(10) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td > a > span",basepath))
	if len(color.Text()) == 0 {
		color = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(10) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td > a > span")
	}
	log.Println("图鉴颜色 ",color.Text())

	// 捕获概率
	var buchang *goquery.Selection
	buchang = d.Find(fmt.Sprintf("%s tr:nth-child(10) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td",basepath))
	if len(buchang.Text()) == 0 {
		buchang = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(10) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td")
	}
	log.Println("捕获概率 ",buchang.Text())
	// 性别比例
	xinbie := d.Find(fmt.Sprintf("%s tr:nth-child(11) > td > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr:nth-child(2) > td > span:nth-child",basepath))
	if len(xinbie.Text()) == 0 {
		xinbie = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(11) > td > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr:nth-child(2) > td")
	}
	xinbie.Each(func(i int, s *goquery.Selection) {
		log.Println("性别比例 ",s.Text())
	})

	// 培育
	peiyu := d.Find(fmt.Sprintf("%s tr:nth-child(12) > td > table > tbody > tr > td > table > tbody > tr > td:nth-child(2) > small",basepath))
	if len(peiyu.Text()) == 0 {
		peiyu = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(12) > td > table > tbody > tr > td > table > tbody > tr > td:nth-child(2) > small")
	}
	log.Println("培育 ",peiyu.Text())

	// 基础点数
	hp := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-HP.bd-HP",basepath))
	log.Println("hp ",hp.Text())

	// 攻击力
	attack := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-攻击.bd-攻击",basepath))
	log.Println("attack",attack.Text())

	// 防御力
	defense := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-防御.bd-防御",basepath))
	log.Println("defense",defense.Text())

	// 特攻
	tattack := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-特攻.bd-特攻",basepath))
	log.Println("tattack",tattack.Text())
	// 特防
	tdefense := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-特防.bd-特防",basepath))
	log.Println("tdefense",tdefense.Text())

	// 速度
	speed := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-速度.bd-速度",basepath))
	log.Println("speed",speed.Text())

	// 基础经验值
	exp := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(2) > td:nth-child(1) > span:nth-child(3)",basepath))
	log.Println("exp",exp.Text())
	// 对战经验值
	attackexp := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(2) > td:nth-child(2) > span:nth-child(3)",basepath))
	log.Println("attackexp",attackexp.Text())

}
