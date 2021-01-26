package main

import (
	"Learnos/GateWay/handler"
	"Learnos/GateWay/websocket"
	"Learnos/GateWay/wrap"
	"Learnos/common/config"
	"Learnos/common/queue/gateway/sms"
	gateway "Learnos/proto/gateway"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/transport"
	"time"
)

func init() {
	err := config.ReadConf("config.toml")
	if err != nil {
		panic(err)
	}
}

func main() {
	conf := config.GetConf()
	//log.SetFlags(log.Lshortfile)
	reg := etcd.NewRegistry(registry.Addrs(conf.Etcd.Addr...), etcd.Auth(conf.Etcd.UserName, conf.Etcd.PassWord))

	service := micro.NewService(
		micro.Name("micro.service.gateway"),
		micro.Registry(reg),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15), //超时重新注册
		micro.WrapHandler(wrap.Verify),
		micro.BeforeStart(beforeStart),
		micro.Transport(
			transport.NewTransport(
				transport.Secure(true),
			),
		),
	)
	err := gateway.RegisterGateWayHandler(service.Server(), handler.Handler{})
	if err != nil {
		panic(err.Error())
	}
	if err := service.Run(); err != nil {
		panic(err.Error())
	}
}

func beforeStart() error {
	go websocket.Run()
	go sms.Recv()
	return nil
}