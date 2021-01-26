package main

import (
	"Learnos/GateWay/sqldata/model"
	"Learnos/GateWay/sqldata/mysql"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)
	db, err := mysql.Get()
	if err != nil{
		log.Println(err.Error())
		return
	}
	defer func() {
		log.Println("回收")
		db.Close()
	}()
	var count int
	result := db.DB.Create(&model.User{Username: "123456",Password: "123456",Phone: ""})
	log.Println(result.Error)
	log.Println(result.RowsAffected)
	var t model.User
	t.Phone = ""
	result2 := db.DB.Where(t).First(&t)
	log.Println(result2.RowsAffected)
	log.Println(t.CreatedAt)
	result3 := db.DB.Model(&model.User{}).Where(model.User{Phone: ""}).Count(&count)
	log.Println(result3.RowsAffected)
	log.Println(count)
	time.Sleep(time.Second * 5)
}
