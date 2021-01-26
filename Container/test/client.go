package main

import (
	"Learnos/proto/cnode"
	"context"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-plugins/registry/nats"
	"log"
)

func main(){
	reg := nats.NewRegistry(registry.Addrs("192.168.0.106:4222"))
	node := micro.NewService(
		micro.Name("micro.service.container.node.client"),
		micro.Registry(reg),
		micro.Transport(
			transport.NewTransport(
				transport.Secure(true),
			),
		),
	)

	client := proto_cnode_service.NewNodeService("micro.service.container.node",node.Client())
	var opt proto_cnode_service.CallOpt
	var create proto_cnode_service.CreateOpt
	var config proto_cnode_service.CreateConfig
	opt.Type = proto_cnode_service.CallType_CreateContainer
	config.Image = "centos"
	config.Cmd = "/bin/bash"
	opt.Create = &create
	opt.Create.Config = &config
	var rsp *proto_cnode_service.CallRsp
	rsp,err := client.Service(context.Background(),&opt)
	if err != nil{
		log.Fatal(err.Error())
	}
	log.Println(rsp)
}