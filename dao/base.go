package dao

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var(
	DB *gorm.DB
	err error
	newLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
		  SlowThreshold: time.Second,   // 慢 SQL 阈值
		  LogLevel:      logger.Silent, // 日志级别
		  IgnoreRecordNotFoundError: false,   // 忽略ErrRecordNotFound（记录未找到）错误
		  Colorful:      true,         // 禁用彩色打印
		},
	  )
)


func init() {
	DB, err = gorm.Open(sqlite.Open("terror.db"), &gorm.Config{
		// 是否启用Logger
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		}, // 使用单数表名
	})
	if err != nil {
		panic("failed to connect database")
	}
	UpdateTable()
}

type EffectType struct{
	TypeId int
	Name string
	Description string // 描述
}

var (
	EFFECT_TYPE_NONE = EffectType{
		TypeId: 1,
		Name: "无",
		Description: "没有任何属性值的变化",
	}
	EFFECT_TYPE_BRIEF = EffectType{
		TypeId: 2,
		Name: "一回合类型",
		Description: "一回合类型",
	}
	EFFECT_TYPE_SUSTAINED = EffectType{
		TypeId: 3,
		Name: "持续",
		Description:"持续数个回合,具体回合请参考kp",
	}
	EFFECT_TYPE = []EffectType{
		EFFECT_TYPE_NONE,
		EFFECT_TYPE_BRIEF,
		EFFECT_TYPE_SUSTAINED,
	}
)

type RoadEntity struct {
	gorm.Model
	Name         string `gorm:"type:varchar(255);not null" json:"name" form:"name" query:"name" validate:"required"`
	// 这个是对RoadEntity的描述
	Introduction string `gorm:"type:varchar(255)" json:"introduction" form:"introduction" query:"introduction"`
	// 没有的话就是零咯
	EffectType int `gorm:"type:int(5)" json:"effect_type" form:"effect_type" query:"effect_type"`
}

type FightingEntity struct {
	RoadEntity
	// 冷却时间,回合制
	Cooling int `gorm:"type:int(100)" json:"cooling" form:"cooling" query:"cooling"`
	Damage int `gorm:"type:int" json:"damage" form:"damage" query:"damage"`
}

// 更新数据库列表
func UpdateTable(){
	// 同步数据库的表
	DB.AutoMigrate(&RoadEntity{})
	DB.AutoMigrate(&FightingEntity{})
}

// 分页显示数据
func SelectRecord(page int) []RoadEntity{
	var roads []RoadEntity
	DB.Limit(10).Offset(page).Find(&roads)
	log.Println(roads)
	return roads
}

// 插入数据
func InsertRecord(record RoadEntity){
	DB.Create(&record)
}

// 更新数据
func UpdateRecord(record RoadEntity){
	DB.Save(&record)
}

// 删除数据
func DeleteRecord(record RoadEntity){
	DB.Delete(&record)
}