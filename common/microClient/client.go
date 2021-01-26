package microClient

import (
	"Learnos/common/config"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/transport"
	"log"
	"sync"
	"time"
)

type microClient struct {
	client.Client
}

var pool *sync.Pool

func init() {
	err := config.ReadConf("config.toml")
	if err != nil {
		log.Fatal(err.Error())
	}
	conf := config.GetConf()
	pool = &sync.Pool{
		New: func() interface{} {
			reg := etcd.NewRegistry(registry.Addrs(conf.Etcd.Addr...), etcd.Auth(conf.Etcd.UserName, conf.Etcd.PassWord))
			c := client.NewClient(
				client.Registry(reg),
				client.Retries(3),
				client.DialTimeout(10*time.Second),
				client.Transport(
					transport.NewTransport(transport.Secure(true), transport.Timeout(10*time.Second)),
				),
			)
			return &microClient{c}
		},
	}
}

func Get() *microClient {
	c := pool.Get().(*microClient)
	return c
}

func (m *microClient) Close() {
	pool.Put(m)
}
