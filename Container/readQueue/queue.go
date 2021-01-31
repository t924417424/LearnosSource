package readQueue

import (
	"Learnos/Container/dockerControl"
	"Learnos/common/queue"
	status2 "Learnos/common/queue/node/status"
	create2 "Learnos/common/queueMsg/node/create"
	"Learnos/common/queueMsg/node/status"
	node "Learnos/proto/cnode"
	"context"
	"errors"
	"github.com/gogo/protobuf/proto"
	"log"
	"time"
)

type CreateQueue struct {
	nodeId string
	addr   string
}

//type containerStatus struct {
//	cid string
//	err error
//}

func NewCreateQueue(nodeId, addr string) *CreateQueue {
	return &CreateQueue{nodeId: nodeId, addr: addr}
}

func (c *CreateQueue) StartQueueRecv(ctx context.Context) chan error {
	exit := make(chan error, 1)
	go func(exit chan<- error, ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				exit <- errors.New("CreateQueue process exit")
				return
			default:
				//先判断一下内存
				if status2.GetFree() < 100 { //可用内存小于100M
					time.Sleep(10 * time.Second)
					continue
				}
				//if queue.MClient.Lock() != nil { //加分布式锁
				//	time.Sleep(3 * time.Second)
				//	continue
				//}
				msg, err := queue.MClient.Pop(create2.ContainerCreatePreFix)
				//queue.MClient.UnLock() ////解开分布式锁
				if err != nil {
					log.Println(err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
				if len(msg) == 0 {
					time.Sleep(10 * time.Second)
					continue
				}
				var config node.CreateOpt
				err = proto.UnmarshalMerge(msg, &config)
				if err != nil {
					log.Println(err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
				if config.Status == node.CreateStatus_Delete {
					dockerControl.CInfo.DelContainer(config.Cid)
				} else {
					c.createContainer(&config)
				}
				//err := inspectOpt(config)
				//log.Println(config)
				//log.Println(config.Config.Image)
				/*开始创建容器*/
				//log.Println("开始创建容器：",config.Cid,"-",config.Config.Image)
				//c.createContainer(&config)
				//这里会阻塞，暂时不使用协程进行创建操作，避免内存占用过量
			}
		}
	}(exit, ctx)
	//自动删除
	go func(exit chan<- error, ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				exit <- errors.New("DeleteQueue process exit")
				return
			default:
				msg, err := queue.MClient.Pop(create2.ContainerDeletePreFix + c.nodeId)
				if err != nil {
					log.Println(err.Error())
					time.Sleep(15 * time.Second)
					continue
				}
				if len(msg) == 0 {
					time.Sleep(10 * time.Second)
					continue
				}
				var config node.CreateOpt
				err = proto.UnmarshalMerge(msg, &config)
				if err != nil {
					log.Println(err.Error())
					time.Sleep(10 * time.Second)
					continue
				}
				dockerControl.CInfo.DelContainer(config.Cid)
			}
		}
	}(exit, ctx)
	return exit
}

func updateStatus(cid string, uid uint64, status create2.CStatus) {
	//statusKey := create2.ContainerListPreFix + strconv.Itoa(int(uid)) + "/" + cid
	//msg, err := queue.MClient.Get(statusKey)
	//if err != nil {
	//	log.Println(err.Error())
	//}
	//msg = create2.GetCreateMessage(msg).UpdateStatus(status).ToJson()
	//_ = queue.MClient.SetEx(statusKey, string(msg), 60*59)
	//_ = nodeStatus.NodeStatus.Send()
	create2.UpdateStatus(cid, uid, status)
}

func (c *CreateQueue) createContainer(opt *node.CreateOpt) {
	c.startWork()
	defer c.stopWork()
	create2.NodeGatCreate(c.addr, c.nodeId, opt.Cid, opt.Uid)	//节点那倒创建信息后，更改etcd中的信息
	updateStatus(opt.Cid, opt.Uid, create2.StartCreate)
	config, hostConfig, err := inspectOpt(opt)
	if err != nil {
		//log.Println(err.Error())
		updateStatus(opt.Cid, opt.Uid, create2.ErrorCreate)
		return
	}
	if config.Image == "" {
		updateStatus(opt.Cid, opt.Uid, create2.PullImageErr)
		return
	}
	err = checkImage(config.Image)
	if err != nil {
		//log.Println(err.Error())
		updateStatus(opt.Cid, opt.Uid, create2.PullImage)
		err = pullImage(config.Image)
		if err != nil {
			log.Println("Pull Image", err.Error())
			updateStatus(opt.Cid, opt.Uid, create2.PullImageErr)
			return
		}
	}
	updateStatus(opt.Cid, opt.Uid, create2.StartCreate)

	err = create(opt.Cid, opt.Uid, config, hostConfig)
	if err != nil {
		//log.Println(err.Error())
		//_, _ = queue.MClient.Push(status.DeleteContainerPreFix, opt.Cid, opt.Cid) //创建失败
		updateStatus(opt.Cid, opt.Uid, create2.ErrorCreate)
		return
	}
	updateStatus(opt.Cid, opt.Uid, create2.OkCreate)
}

func (c *CreateQueue) startWork() {
	_ = queue.MClient.Set(status.CreatingPreFix+c.nodeId, "true")
}

func (c *CreateQueue) stopWork() {
	_ = queue.MClient.Set(status.CreatingPreFix+c.nodeId, "false")
}

func (c *CreateQueue) ClearQueue() {
	_ = queue.MClient.Del(status.CreatingPreFix + c.nodeId)
	//for {
	//	msg, err := queue.MClient.Pop(create2.ContainerCreatePreFix + c.nodeId)
	//	if err != nil {
	//		return
	//	}
	//	if len(msg) == 0 {
	//		return
	//	}
	//}
}

