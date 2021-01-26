package queue

import (
	"Learnos/common/cache"
	"Learnos/common/cache/etcd"
	"Learnos/common/config"
	"log"
	"time"
)

var MClient cache.Cache

func init() {
	MClient = newQueue()
}

func newQueue() cache.Cache {
	err := config.ReadConf("config.toml")
	if err != nil {
		log.Fatal(err.Error())
	}
	conf := config.GetConf()
	c, err := etcd.NewCache(etcd.Addr(conf.Etcd.Addr), etcd.TimeOut(5*time.Second), etcd.Auth(conf.Etcd.UserName, conf.Etcd.PassWord))
	if err != nil {
		log.Fatal(err.Error())
	}
	return c
}