package pkminfodao

import (
	"log"
	"reptile/dao"

	"gorm.io/gorm"
)

// 这里存放关于宝可梦相关的信息
type (
	Pokemon struct {
		gorm.Model
		// 官方编号
		Number string `gorm:"column:number"`
		// 姓名
		Name string `json:"name" validate:"required" query:"name"`
		// 宝可梦类型
		Type string `json:"type" validate:"required" query:"type"`
		// 宝可梦属性,这个属性默认为json数据
		Property string `json:"property" validate:"required" query:"property"`
		// 宝可梦分类Pokemon
		Category string `json:"category" validate:"required" query:"category"`
		// 特性,先为JSON数据吧,方便后期的管理
		Ability string `json:"ability" validate:"required" query:"ability"`
		// 地区编号,这个也先显示为json数据
		Area string `json:"area" validate:"required" query:"area"`
		// 宝可梦的官方绘图
		Image string `json:"image" validate:"required" query:"image"`
		// 宝可梦的身高
		Height string `json:"height" validate:"required" query:"height"`
		// 宝可梦的体重
		Weight string `json:"weight" validate:"required" query:"weight"`
		// 体型路径
		Body string `json:"body" validate:"required" query:"body"`
		// 脚印路径
		Foot string `json:"foot" validate:"required" query:"foot"`
		// 图鉴颜色
		Color string `json:"color" validate:"required" query:"color"`
		// 捕获率
		CaptureRate string `json:"capture_rate" validate:"required" query:"capture_rate"`
		// 培育周期
		EvolutionCycle string `json:"evolution_cycle" validate:"required" query:"evolution_cycle"`
		// 基础点数
		// Hp
		Hp int `json:"hp" validate:"required" query:"hp"`
		// 攻击力
		Attack int `json:"attack" validate:"required" query:"attack"`
		// 防御力
		Defense int `json:"defense" validate:"required" query:"defense"`
		// 特攻
		SpAttack int `json:"sp_attack" validate:"required" query:"sp_attack"`
		// 特防
		SpDefense int `json:"sp_defense" validate:"required" query:"sp_defense"`
		// 速度
		Speed int `json:"speed" validate:"required" query:"speed"`
		// 基础经验值
		BaseExp string `json:"base_exp" validate:"required" query:"base_exp"`
		// 对战经验值
		BattleExp int `json:"battle_exp" validate:"required" query:"battle_exp"`
	}

	PkmDaoManger struct {
		db *gorm.DB
	}
)

func NewPkmDaoManger() *PkmDaoManger {
	return &PkmDaoManger{
		db: dao.DB,
	}
}

// 插入一个数组的数据
func (p *PkmDaoManger) InsertPkm(pkm *[]Pokemon) {
	for _, v := range *pkm {
		if len(v.Name) > 0 {
			p.db.Create(&v)
		}
	}
}

// 根据id删除一条数据
func (p *PkmDaoManger) DeletePkm(id int) {
	p.db.Where("id = ?", id).Delete(&Pokemon{})
}

// 根据name查询一条数据
func (p *PkmDaoManger) FindPkmByName(name string) *Pokemon {
	var pkm Pokemon
	p.db.Where("name = ?", name).First(&pkm)
	return &pkm
}

// 根据id查询一条数据
func (p *PkmDaoManger) FindPkmById(id int) *Pokemon {
	var pkm Pokemon
	p.db.Where("id = ?", id).First(&pkm)
	return &pkm
}

// 分页查询数据
func (p *PkmDaoManger) FindPkmByPage(page int, limit int) *[]Pokemon {
	var pkm []Pokemon
	p.db.Limit(limit).Offset(page).Find(&pkm)
	return &pkm
}

// 查询一组ids的数据
func (p *PkmDaoManger) FindPkmByIds(ids []uint) *[]Pokemon {
	var pkm []Pokemon
	p.db.Where("id in (?)", ids).Find(&pkm)
	return &pkm
}

// 根据名称查询一条数据
func (p *PkmDaoManger) FindPkmByNameLike(name string) *[]Pokemon {
	var pkm []Pokemon
	p.db.Where("name like ?", "%"+name+"%").Find(&pkm)
	return &pkm
}

// 拿到这个表所有的数据
func (p *PkmDaoManger) FindAllPkm() *[]Pokemon {
	var pkm []Pokemon
	err := p.db.Find(&pkm).Error
	if err != nil {
		log.Println(err)
	}
	return &pkm
}

// 更新数据库里面的image数据
func (p *PkmDaoManger) UpdatePkmImage(id uint, image string) {
	p.db.Model(&Pokemon{}).Where("id = ?", id).Update("image", image)
}
