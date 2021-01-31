package main

import (
	"Learnos/GateWay/handler"
	"Learnos/GateWay/websocket"
	"Learnos/GateWay/wrap"
	"Learnos/common/config"
	"Learnos/common/queue/gateway/sms"
	"Learnos/common/util"
	gateway "Learnos/proto/gateway"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/server"
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
	port, err := util.GetPort()
	if err != nil {
		panic(err.Error())
	}

	conf := config.GetConf()
	//log.SetFlags(log.Lshortfile)
	reg := etcd.NewRegistry(registry.Addrs(conf.Etcd.Addr...), etcd.Auth(conf.Etcd.UserName, conf.Etcd.PassWord))

	service := micro.NewService(
		micro.Name("micro.service.gateway"),
		micro.Registry(reg),
		micro.Address(fmt.Sprintf("0.0.0.0:%d", port)),
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

	//公网部署模式
	if conf.Common.PublicNetWorkMode {
		publicIp := fmt.Sprintf("%s:%d", util.GetPubicIP(), port)
		service.Server().Init(server.Advertise(publicIp)) //将节点公网信息注册到服务中心
	}

	err = gateway.RegisterGateWayHandler(service.Server(), handler.Handler{})
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
