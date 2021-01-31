package create

import (
	"Learnos/Container/nodeStatus"
	"Learnos/common/queue"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	ContainerCreatePreFix = "/container/create/"
	ContainerDeletePreFix = "/container/delete/"
	ContainerListPreFix   = "/container/list/"
)

type CStatus uint32

const (
	Loading      CStatus = 0
	PullImage    CStatus = 1
	StartCreate  CStatus = 2
	OkCreate     CStatus = 3
	ErrorCreate  CStatus = 4
	PullImageErr CStatus = 5
	Deleted      CStatus = 6
	Overstep     CStatus = 7
)

//type RecordType uint
//
//const (
//	NetWork RecordType = 0
//	Block   RecordType = 1
//)

var msg = map[CStatus]string{
	Loading:      "等待创建",
	PullImage:    "正在拉取镜像",
	PullImageErr: "拉取镜像失败",
	StartCreate:  "开始创建",
	OkCreate:     "创建成功",
	ErrorCreate:  "创建失败",
	Deleted:      "已删除",
	Overstep:     "触发资源限制",
}

type limit struct {
	Network networkIo
	BlockIo blockIo
}

type networkIo struct {
	Record uint64
	Limit  uint64
}

type blockIo struct {
	Record uint64
	Limit  uint64
}

type createStatus struct {
	Status CStatus //0等待创建，1：下载镜像，2：正在创建，3：创建成功，4：创建失败
	Msg    string
	Time   int64
	Uid    uint
	NodeId string
	Addr   string
	Uuid   string
	Image  string
	limit
}

func NewCreateMessage(s CStatus, uid uint, nodeId, addr string, uuid string, image string, networkIoLimit uint64, blockIoLimit uint64) ([]byte, error) {
	addr, _, _ = net.SplitHostPort(addr)
	return json.Marshal(createStatus{s, msg[s], time.Now().Unix(), uid, nodeId, addr, uuid, image, limit{Network: networkIo{Limit: networkIoLimit}, BlockIo: blockIo{Limit: blockIoLimit}}})
}

func GetCreateMessage(data []byte) *createStatus {
	var msg createStatus
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
	}
	return &msg
}

func NodeGatCreate(addr, nodeId string, cid string, uid uint64) {
	statusKey := ContainerListPreFix + strconv.Itoa(int(uid)) + "/" + cid
	msg, err := queue.MClient.Get(statusKey)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if len(msg) == 0 { //消息未找到（或者已过期的情况），自动删除
		_ = queue.MClient.Del(statusKey)
		return
	}
	host,_,_ := net.SplitHostPort(addr)
	msg = GetCreateMessage(msg).updateNodeInfo(host,nodeId).toJson()
	_ = queue.MClient.SetEx(statusKey, string(msg), 60*10)
}

func UpdateStatus(cid string, uid uint64, status CStatus) {
	statusKey := ContainerListPreFix + strconv.Itoa(int(uid)) + "/" + cid
	msg, err := queue.MClient.Get(statusKey)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if len(msg) == 0 { //消息未找到（或者已过期的情况），自动删除
		_ = queue.MClient.Del(statusKey)
		return
	}
	//log.Println(statusKey)
	msg = GetCreateMessage(msg).updateStatus(status).toJson()
	if status == OkCreate {
		_ = queue.MClient.Set(statusKey, string(msg)) //防止容器服务器意外退出导致的无法重新创建
	} else if status == PullImage || status == StartCreate {
		_ = queue.MClient.SetEx(statusKey, string(msg), 60*10)
	} else {
		_ = queue.MClient.SetEx(statusKey, string(msg), 15) //防止频繁创建
	}
	_ = nodeStatus.NodeStatus.Send()
}

func (c *createStatus) updateStatus(s CStatus) *createStatus {
	c.Status = s
	c.Msg = msg[s]
	c.Time = time.Now().Unix()
	return c
}

func (c *createStatus) updateNodeInfo(addr,nodeId string) *createStatus {
	c.Addr = addr
	c.NodeId = nodeId
	c.Time = time.Now().Unix()
	return c
}

func GetLimit(cid string, uid uint64) (netLimit, blockLimit uint64, err error) {
	statusKey := ContainerListPreFix + strconv.Itoa(int(uid)) + "/" + cid
	msg, err := queue.MClient.Get(statusKey)
	if err != nil {
		return 0, 0, err
	}
	info := GetCreateMessage(msg)
	return info.Network.Limit, info.BlockIo.Limit, nil
}

func UpdateRecord(net, block uint64, cid string, uid uint64) { //更新容器网络和磁盘Io使用信息
	statusKey := ContainerListPreFix + strconv.Itoa(int(uid)) + "/" + cid
	msg, err := queue.MClient.Get(statusKey)
	if err != nil {
		log.Println(err.Error())
		return
	}
	msg = GetCreateMessage(msg).updateNetWorkRecord(net, block).toJson()
	_ = queue.MClient.Set(statusKey, string(msg))
	_ = nodeStatus.NodeStatus.Send()
}

func (c *createStatus) updateNetWorkRecord(net, block uint64) *createStatus {
	c.Network.Record = net
	c.BlockIo.Record = block
	return c
}

//func (c *createStatus) updateBlockRecord(val uint64) *createStatus {
//	c.blockIo.record = val
//	return c
//}

func (c *createStatus) toJson() []byte {
	msg, err := json.Marshal(c)
	if err != nil {
		log.Println(err.Error())
	}
	return msg
}