package main

import (
	GateWay "Learnos/proto/gateway"
	"context"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/transport"
	"log"
	"time"
)

func main(){
	reg := etcd.NewRegistry(registry.Addrs("127.0.0.1:2379"), etcd.Auth("root", "123456"))
	c := client.NewClient(
		client.Registry(reg),
		client.Retries(3),
		client.Transport(
			transport.NewTransport(transport.Secure(true)),
		),
	)
	//var req proto_gateway_service.Call
	var rsp GateWay.CallRsp
	req := &GateWay.Call{
		Type:   GateWay.CallType_User,
		Caller: GateWay.CallerType_UserLogin,
		Opt: &GateWay.Options{
			User: &GateWay.UserOpt{
				Username: "123456",
				Password: "123456",
			},
		},
	}
	for {
		ctx := metadata.NewContext(context.Background(),map[string]string{"Source-Ip":"127.0.0.1"})
		r := c.NewRequest("micro.service.gateway", "GateWay.Service", req)
		err := c.Call(ctx, r, &rsp)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println(rsp)
		time.Sleep(time.Second)
	}
}
