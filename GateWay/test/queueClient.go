package main

import (
	"Learnos/common/cache/etcd"
	"Learnos/common/queue/gateway/sms"
	"log"
	"strconv"
	"time"
)

func main() {
	go sms.Recv()
	const smsPrefix = "/msg/sendCode/"
	c, err := etcd.NewCache(etcd.Addr([]string{"127.0.0.1:2379"}), etcd.TimeOut(5*time.Second), etcd.Auth("root", "123456"))
	if err != nil {
		log.Fatal(err.Error())
	}
	for i := 0; i < 1000; i++ {
		_, err := c.Push(smsPrefix, strconv.Itoa(i), strconv.Itoa(i))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	time.Sleep(100 * time.Second)
}
