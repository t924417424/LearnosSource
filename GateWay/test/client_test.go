package main

import (
	"Learnos/proto/gateway"
	"context"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-plugins/registry/etcd"
	"testing"
	"time"
)

func TestGateWay(t *testing.T) {
	reg := etcd.NewRegistry(registry.Addrs("127.0.0.1:2379"), etcd.Auth("root", "123456"))
	c := client.NewClient(
		client.Registry(reg),
		client.Retries(3),
		client.Transport(
			transport.NewTransport(transport.Secure(true)),
		),
	)
	//var req proto_gateway_service.Call
	var rsp proto_gateway_service.CallRsp
	req := &proto_gateway_service.Call{
		Type:   proto_gateway_service.CallType_User,
		Caller: proto_gateway_service.CallerType_SendCode,
		Opt: &proto_gateway_service.Options{
			User: &proto_gateway_service.UserOpt{
				Iphone: "15552321035",
			},
		},
	}
	for {
		r := c.NewRequest("micro.service.gateway", "GateWay.Service", req)
		ctx := metadata.NewContext(context.Background(),map[string]string{"Source-Ip":"127.0.0.1"})
		err := c.Call(ctx, r, &rsp)
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(rsp)
		time.Sleep(time.Second)
	}
}
