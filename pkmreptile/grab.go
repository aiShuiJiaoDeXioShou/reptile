package pkmreptile

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"reptile/comm"
	"reptile/dao/pkminfodao"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type (
	// 宝可梦指定爬虫管理器
	PkmReptile struct {
		Url         string              // 需要抓取的链接
		pkm         *pkminfodao.PkmDaoManger // 宝可梦数据库数据源
		d 			*goquery.Document // 全局的抓取对象,方便链式调用
		SaveFileUrl string              // 指定读取的文件路径
		Pkmmap      map[string]string   // 用于存放的map
	}
	areaID struct{
		GuanDu string
		ChenDu string
		Kaluosi string
		LaJiaEr string
	}
)
// 被抓取值的首页路径
var basepath string

// 抓取指定宝可梦的数据,返回被抓取的宝可梦对象
func (pkr *PkmReptile) Grab() (pkm *pkminfodao.Pokemon) {

	return
}

// 只有先调用这个函数才能使用该结构体的所有函数
func NewPkmReptile() *PkmReptile {
	// 设置抓取的宝可梦路径
	baseurl := "https://wiki.52poke.com"
	url := baseurl + "/wiki/%E5%AE%9D%E5%8F%AF%E6%A2%A6%E5%88%97%E8%A1%A8%EF%BC%88%E6%8C%89%E5%85%A8%E5%9B%BD%E5%9B%BE%E9%89%B4%E7%BC%96%E5%8F%B7%EF%BC%89/%E7%AE%80%E5%8D%95%E7%89%88#.E7.AC.AC.E4.B8.80.E4.B8.96.E4.BB.A3"

	// 创建抓取的宝可梦对象
	d, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	return &PkmReptile{
		Url: "https://wiki.52poke.com",
		d: d,
		pkm: pkminfodao.NewPkmDaoManger(),
		SaveFileUrl: "./pkm.txt",
		Pkmmap: make(map[string]string),
	}
}

// 抓取宝可梦基本信息
func (pkr *PkmReptile) GrabInfo(pkm *pkminfodao.Pokemon,name string) (*pkminfodao.Pokemon) {
	// 抓取该目标值
	pkr.d.Find(".eplist").Each(func(i int, s *goquery.Selection) {
		node := s.Find("td")
		node.Find("a").Each(func(i int, s *goquery.Selection) {
			if s.Text() == name {
				// 序号
				number := s.Parent().Parent().Children().Eq(0).Text()
				pkm.Number = number
				// 名称
				pkm.Name = s.Text()
				if val, exists := s.Attr("href"); exists {
					pkr.grabDetail(pkm, pkr.Url+val)
				}
				return
			}
		})
	})

	return pkm
}

// 抓取所有宝可梦的链接,以及名称的map
func (pkr *PkmReptile) GrabBasicsPkm() map[string]gin.H {
	hash := make(map[string]gin.H)
	// 抓取该目标值
	pkr.d.Find(".eplist").Each(func(i int, s *goquery.Selection) {
		node := s.Find("td")
		node.Find("a").Each(func(i int, s *goquery.Selection) {
			// comm.IsChinese 是判断抓取目标是不是中文
			number := s.Parent().Parent().Children().Eq(0).Text()
			if val, exists := s.Attr("href");exists&&comm.IsChinese(s.Text()) {
				hash[s.Text()] = gin.H{
					"url": pkr.Url+val,
					"name": s.Text(),
					"number": number,
				}
			}

		})
	})
	return hash
}

// 抓取所有的宝可梦信息
func (pkr *PkmReptile) GrabPkomAll()(*[]pkminfodao.Pokemon) {
	pkmmap := pkr.GrabBasicsPkm()
	pkms := make([]pkminfodao.Pokemon,len(pkmmap))
	// 遍历pkmmap
	for _, v := range pkmmap {
		var pkm = pkminfodao.Pokemon{
			Name: v["name"].(string),
			Number: v["number"].(string),
		}
		// 获取宝可梦的详细信息
		pkr.grabDetail(&pkm, v["url"].(string))
		pkms = append(pkms, pkm)
	}
	return &pkms
}

// 抓取所有的宝可梦信息返回值为Map
func (pkr *PkmReptile) GrabPkomAllMap()(*map[string]*pkminfodao.Pokemon) {
	pkmmap := pkr.GrabBasicsPkm()
	pkms := make(map[string]*pkminfodao.Pokemon)
	// 遍历pkmmap
	for _, v := range pkmmap {
		var pkm = &pkminfodao.Pokemon{
			Name: v["name"].(string),
			Number: v["number"].(string),
		}
		// 获取宝可梦的详细信息
		pkr.grabDetail(pkm, v["url"].(string))
		pkms[pkm.Name] = pkm
	}
	return &pkms
}


// 抓取详细信息

func (pkr *PkmReptile) grabDetail(pkm *pkminfodao.Pokemon, url string) *pkminfodao.Pokemon {
	d, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	s := d.Find("#firstHeading")
	if regexp.MustCompile(`^第(.*)世代`).MatchString(s.Text()) {
		return pkm
	}
	// 拿到数据显示的table表单
	table := d.Find(".mw-parser-output").Find(".roundy.a-r.at-c")

	// 这个是宝可梦详细表单的父类
	fater := table.Find("table.roundy.bgwhite.fulltable>tbody>tr")
	basepath = "#mw-content-text > div > table.roundy.a-r.at-c > tbody >"

	// 抓取立绘
	grabimage(table, pkm)

	// 拿到宝可梦的属性值
	grabattr(table, pkm)

	// 特性
	grabability(fater, pkm)

	// 经验值
	grabexp(fater, pkm)

	// 地区图鉴编号
	grabareaid(d, pkm)

	// 身高
	grabheight(d, pkm)

	// 体重
	grabweight(d, pkm)

	// 体形
	grabshape(d, pkm, url)

	// 脚印
	grabfoot(d, pkm, url)

	// 图鉴颜色
	grabcolor(d, pkm)

	// 捕获概率
	grabcatchrate(d, pkm)

	// 性别比例
	grabsexrate(d, pkm)

	// 培育
	grabevo(d, pkm)

	// 基本属性
	grabbaseattr(d, pkm)

	// 种族值
	grabrace(d, pkm)
	return pkm
}

// 把该条数据储存到数据库
func (pkr *PkmReptile) Save(pkm *pkminfodao.Pokemon) {

}

// 把该条数据写入指定文件
func (pkr *PkmReptile) Write(pkm *pkminfodao.Pokemon) {

}

// 把该条数据JSON化保存
func (pkr *PkmReptile) SaveJson(pkm *pkminfodao.Pokemon) {

}

// 获取pkm立绘
func grabimage(table *goquery.Selection, pkm *pkminfodao.Pokemon) {
	// 通过表单拿到立绘
	image := table.Find("table.roundy.bgwhite.fulltable").Find(".image").Find("img")
	val, exists := image.Attr("data-url")
	if exists {
		log.Println("官方立绘为：", val)
		pkm.Image = val
	}
}

// 拿到宝可梦的属性值
func grabattr(table *goquery.Selection, pkm *pkminfodao.Pokemon) {
	s := table.Find("tbody>tr").Eq(4)
	var prop string
	var category string
	s.Find("a").Each(func(i int, s *goquery.Selection) {
		if str, err := s.Attr("title"); err && strings.Contains(str, "（属性）") {
			prop += str
		}

		if str, err := s.Attr("title"); err && strings.Contains(str, "分类") {
			s2 := s.Parent().Next().Find("td")
			category += s2.Text()
		}
	})
	pkm.Property = prop
	pkm.Category = category
}

// 拿到宝可梦的特性
func grabability(fater *goquery.Selection, pkm *pkminfodao.Pokemon) {
	var ability string
	texin := fater.Eq(3)
	texin.Find("a").Each(func(i int, s *goquery.Selection) {
		if str, err := s.Attr("title"); err && strings.Contains(str, "（特性）") {

			if len(str) > 0 {
				s2 := s.Next().Next()
				if len(s2.Text()) > 0 {
					ability += str + ":" + s2.Text() + ";"
				} else {
					ability += str
				}

			}

		}
	})
	pkm.Ability = ability
}

// 宝可梦的基本属性
func grabbaseattr(d *goquery.Document, pkm *pkminfodao.Pokemon) {
	// 基础点数
	hp := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-HP.bd-HP", basepath))
	log.Println("hp ", hp.Text())

	// 攻击力
	attack := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-攻击.bd-攻击", basepath))
	log.Println("attack", attack.Text())

	// 防御力
	defense := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-防御.bd-防御", basepath))
	log.Println("defense", defense.Text())

	// 特攻
	tattack := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-特攻.bd-特攻", basepath))
	log.Println("tattack", tattack.Text())
	// 特防
	tdefense := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-特防.bd-特防", basepath))
	log.Println("tdefense", tdefense.Text())

	// 速度
	speed := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td.roundy.bw-1.bgl-速度.bd-速度", basepath))
	log.Println("speed", speed.Text())

	// 基础经验值
	exp := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(2) > td:nth-child(1) > span:nth-child(3)", basepath))
	log.Println("exp", exp.Text())
	// 对战经验值
	attackexp := d.Find(fmt.Sprintf("%s tr:nth-child(13) > td > table > tbody > tr > td > table > tbody > tr:nth-child(2) > td:nth-child(2) > span:nth-child(3)", basepath))
	log.Println("attackexp", attackexp.Text())
}

// 下面就是一些比较奇怪的东西
// 脚印
func grabfoot(d *goquery.Document, pkm *pkminfodao.Pokemon, url string) {
	jiaoyin := d.Find(fmt.Sprintf("%s tr:nth-child(9) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td > a", basepath))
	jiaoyinval, jiaoyinexists := jiaoyin.Attr("data-url")
	if jiaoyinexists {
		pkm.Foot = url + jiaoyinval
	}

}

// 图鉴颜色
func grabcolor(d *goquery.Document, pkm *pkminfodao.Pokemon) {
	var color *goquery.Selection
	color = d.Find(fmt.Sprintf("%s tr:nth-child(10) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td > a > span", basepath))
	if len(color.Text()) == 0 {
		color = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(10) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td > a > span")
	}
	pkm.Color = color.Text()
}

// 捕获概率
func grabcatchrate(d *goquery.Document, pkm *pkminfodao.Pokemon) {
	var buchang *goquery.Selection
	buchang = d.Find(fmt.Sprintf("%s tr:nth-child(10) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td", basepath))
	if len(buchang.Text()) == 0 {
		buchang = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(10) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td")
	}
	pkm.CaptureRate = buchang.Text()
}

// 性别比例
func grabsexrate(d *goquery.Document, pkm *pkminfodao.Pokemon) {
	sexrate := d.Find(fmt.Sprintf("%s tr:nth-child(11) > td > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr:nth-child(2) > td > span:nth-child", basepath))
	if len(sexrate.Text()) == 0 {
		sexrate = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(11) > td > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr:nth-child(2) > td")
	}
	sexrate.Each(func(i int, s *goquery.Selection) {
		log.Println("性别比例 ", s.Text())
	})
}

// 培育指数
func grabevo(d *goquery.Document, pkm *pkminfodao.Pokemon) {
	peiyu := d.Find(fmt.Sprintf("%s tr:nth-child(12) > td > table > tbody > tr > td > table > tbody > tr > td:nth-child(2) > small", basepath))
	if len(peiyu.Text()) == 0 {
		peiyu = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(12) > td > table > tbody > tr > td > table > tbody > tr > td:nth-child(2) > small")
	}
	pkm.EvolutionCycle = peiyu.Text()
}

// 身高
func grabheight(d *goquery.Document, pkm *pkminfodao.Pokemon) {
	var height *goquery.Selection
	height = d.Find(fmt.Sprintf("%s tr:nth-child(8) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td", basepath))
	if len(height.Text()) == 0 {
		height = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(8) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td")
	}
	pkm.Height = height.Text()
}

// 体重
func grabweight(d *goquery.Document, pkm *pkminfodao.Pokemon) {
	var weight *goquery.Selection
	weight = d.Find(fmt.Sprintf("%s tr:nth-child(8) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td", basepath))
	if len(weight.Text()) == 0 {
		weight = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(8) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td")
	}
	pkm.Weight = weight.Text()
}

// 体型
func grabshape(d *goquery.Document, pkm *pkminfodao.Pokemon, url string) {
	var tixing *goquery.Selection
	tixing = d.Find(fmt.Sprintf("%s tr:nth-child(9) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td > a > img", basepath))
	tixingval, tixingexists := tixing.Attr("data-url")
	if tixingexists {
		pkm.Body = url + tixingval
	} else {
		tixing = d.Find("#mw-content-text > div > table:nth-child(2) > tbody > tr._toggle.form1 > td > table > tbody > tr:nth-child(9) > td:nth-child(1) > table > tbody > tr > td > table > tbody > tr > td > a > img")
		tixingval, tixingexists := tixing.Attr("data-url")
		if tixingexists {
			pkm.Body = url + tixingval
		}
	}
}



// 地区图鉴编号
func grabareaid(d *goquery.Document, pkm *pkminfodao.Pokemon) {
	var areaid areaID
	// 关都地区
	guandu := d.Find("td.bd-关都.bw-1.roundyright").Eq(0)
	areaid.GuanDu = "关都地区->"+guandu.Text()

	// 成都地区
	chendu := d.Find(".bd-城都.bw-1").Eq(0)
	areaid.ChenDu = "成都地区->"+chendu.Text()

	// 卡洛斯地区
	kls := d.Find(".bd-卡洛斯.bw-1.roundyright").Eq(0)
	areaid.Kaluosi = "卡洛斯地区->"+ kls.Text()

	// 伽勒尔地区
	jls := d.Find("td.bd-伽勒尔.bw-1.roundyright").Eq(0)
	areaid.LaJiaEr = "伽勒尔地区 ->"+jls.Text()
	b, err := json.Marshal(areaid)
	if err != nil {
		panic(err)
	}
	pkm.Area = string(b)
}

// 经验值
func grabexp(fater *goquery.Selection, pkm *pkminfodao.Pokemon) {
	ex := fater.Eq(4)
	log.Println("经验值:", ex.Find("td").Text())
	pkm.BaseExp = ex.Find("td").Text()
}

func grabrace(d *goquery.Document, pkm *pkminfodao.Pokemon){
	// TODO
	
}
