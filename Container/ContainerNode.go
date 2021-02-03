package main

import (
	"Learnos/Container/dockerControl"
	_ "Learnos/Container/dockerControl"
	"Learnos/Container/handler"
	"Learnos/Container/nodeStatus"
	"Learnos/Container/readQueue"
	"Learnos/Container/websocket"
	"Learnos/common/config"
	"Learnos/common/queue/node/status"
	"Learnos/common/util"
	node "Learnos/proto/cnode"
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/transport"
	"log"
	"time"
)

var service micro.Service
var statusCtx context.Context
var statusCancel context.CancelFunc
var queueCtx context.Context
var queueCancel context.CancelFunc
var queue *readQueue.CreateQueue
var publicIp string

func init() {
	err := config.ReadConf("config.toml")
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	port, err := util.GetPort()
	if err != nil {
		panic(err.Error())
	}
	conf := config.GetConf()
	reg := etcd.NewRegistry(registry.Addrs(conf.Etcd.Addr...), etcd.Auth(conf.Etcd.UserName, conf.Etcd.PassWord))
	service = micro.NewService(
		micro.Name("micro.service.container.node"),
		micro.Registry(reg),
		micro.Address(fmt.Sprintf("0.0.0.0:%d", port)),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15), //超时重新注册
		micro.AfterStart(afterStart),
		micro.AfterStop(stopContainer),
		micro.Transport(
			transport.NewTransport(
				transport.Secure(true),
			),
		),
	)
	//公网部署模式
	if conf.Common.PublicNetWorkMode {
		publicIp = fmt.Sprintf("%s:%d", util.GetPubicIP(), port)
		service.Server().Init(server.Advertise(publicIp)) //将节点公网信息注册到服务中心
	}

	err = node.RegisterNodeHandler(service.Server(), handler.Handler{})
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := service.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func afterStart() error {
	conf := config.GetConf()
	go websocket.Run()
	if !conf.Common.PublicNetWorkMode {
		publicIp = service.Server().Options().Advertise
	}
	statusCtx, statusCancel = context.WithCancel(context.TODO())
	nodeStatus.NodeStatus = status.NewNodeStatus(service.Server().Options().Id, publicIp, statusCtx)
	go func(ctx context.Context) {
		select {
		case <-nodeStatus.NodeStatus.KeepExists():
			log.Println("自动上报内存功能出错，请检查缓存（Etcd）配置")
			return
		case <-ctx.Done():
			log.Println("结束自动上报内存信息")
			nodeStatus.NodeStatus.Delete()
			return
		}
	}(statusCtx)
	/**
	 * 消息投递方式1：网关发送通知，节点进行应答，开始创建
	 * 消息投递方式2：节点主动从自身维护的消息队列读取
	 */
	queueCtx, queueCancel = context.WithCancel(context.Background())
	queue = readQueue.NewCreateQueue(service.Server().Options().Id, publicIp)
	go func(ctx context.Context) {
		for {
			select {
			case <-queue.StartQueueRecv(ctx):
				log.Println("容器创建已停止")
				return
			case <-ctx.Done():
				log.Println("结束容器创建功能") //清空队列以及serverList
				return
			}
		}
	}(queueCtx)
	//err := nodeStatus.NodeStatus.Send() //发送当前节点内存状态，每次创建容器重复操作，用于网关选择创建docker镜像的节点
	//if err != nil {
	//	log.Println(err.Error())
	//}
	return nil
}

func stopContainer() error {
	statusCancel()                     //停止程序，取消自动续约
	queueCancel()                      //结束自动创建
	_ = nodeStatus.NodeStatus.Delete() //停止程序，自动删除内存上报
	queue.ClearQueue()                 //停止程序自动删除该节点维护的队列
	dockerControl.Exit()               //停止所有程序创建的容器，并删除状态缓存消息
	return nil
}
