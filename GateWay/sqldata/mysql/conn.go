package mysql

import (
	"Learnos/GateWay/sqldata/model"
	"Learnos/common/config"
	util2 "Learnos/common/util"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var pool *sqlPool

type MyDb struct {
	DB *gorm.DB
}

func init() {
	pool = newPool(newDb, 10)
}

func newDb() *gorm.DB {
	conf := config.GetConf()
	//db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", "docker", "eA2pyJ6BHDy53yz5", "127.0.0.1", 3306, "docker"))
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", conf.GateWay.Mysql.UserName, conf.GateWay.Mysql.PassWord, conf.GateWay.Mysql.Addr, conf.GateWay.Mysql.Port, conf.GateWay.Mysql.DbName))
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	if err != nil {
		log.Printf("db open err :%s", err.Error())
		return nil
	}
	db.SetLogger(util2.Logger{})
	if !db.HasTable(&model.User{}) {
		log.Println("init table")
		_ = db.CreateTable(&model.User{}, &model.Image{}, &model.History{})
	}
	return db
}

func (s *MyDb) Close() {
	pool.put(s.DB)
}

func Get() (*MyDb, error) {
	db, err := pool.get()
	return &MyDb{db}, err
}
