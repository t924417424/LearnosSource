package cache

import (
	"Learnos/common/cache/etcd"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	c, err := etcd.NewCache(etcd.Addr([]string{"127.0.0.1:2379"}), etcd.TimeOut(5*time.Second),etcd.Auth("root","123456"))
	if err != nil {
		t.Fatal(err.Error())
	}
	err = c.Set("test", "interface test")
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err := c.Get("test")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(data))
	err = c.Del("test")
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err = c.Get("test")
	if err != nil {
		t.Fatal(err.Error())
	}
	exists,_ := c.Exists("test")
	t.Log(exists)

	go func() {
		i := 0
		for {
			data, err := c.Pop("/msg/queue/")
			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			}
			if len(data) == 0 {
				time.Sleep(3 * time.Second)
				continue
			}
			i++
			log.Println("收到消息:", string(data))
			t.Log("Count:", i)
		}
	}()
	go func() {
		for i := 1000; i < 2000; i++ {
			_, err = c.Push("/msg/queue/", strconv.Itoa(i), fmt.Sprintf("发布消息%d", i))
			if err != nil {
				t.Fatal(err.Error())
			}
		}
	}()
	for i := 0; i < 1000; i++ {
		_, err = c.Push("/msg/queue/", strconv.Itoa(i), fmt.Sprintf("发布消息%d", i))
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	time.Sleep(time.Second * 100)
}
