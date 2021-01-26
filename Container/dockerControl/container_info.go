package dockerControl

import (
	"Learnos/common/queue"
	"Learnos/common/queueMsg/node/create"
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"log"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type ContainerInfo struct {
	*sync.RWMutex
	infos map[string]*info
}

type info struct {
	UserID      uint64
	ContainerID string
	deleteSign  int64
	isConnect   bool
	record      bool //记录是否连接过
	AutoRemove  bool
	overflow    bool
	netLimit    uint64 //网络Io限制
	blockLimit  uint64 //blockIo限制
	userDelete  bool   //用户申请删除
	sshStatus   bool   //ssh是否连接
}

var CInfo *ContainerInfo

func init() {
	CInfo = &ContainerInfo{new(sync.RWMutex), make(map[string]*info)}
	go CInfo.autoRm()
}

func (c *ContainerInfo) DelContainer(cid string) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.infos[cid]; ok {
		c.infos[cid].userDelete = true
	}
}

func Exit() { //程序退出时删除所有程序进行创建的容器
	for cid, v := range CInfo.infos {
		d := time.Second
		_ = DockerClient.ContainerStop(context.Background(), v.ContainerID, &d)
		_ = DockerClient.ContainerRemove(context.Background(), v.ContainerID, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true})
		_ = queue.MClient.Del(create.ContainerListPreFix + strconv.Itoa(int(v.UserID)) + "/" + cid) //删除缓存中的创建状态消息，客户端再次请求更新数据库状态为已删除
		create.UpdateStatus(v.ContainerID, v.UserID, create.Deleted)                                //更新缓存状态为已删除
		//_,_ = queue.MClient.Push(status.DeleteContainerPreFix,cid,cid)	//发送至删除队列
	}
}

func (c *ContainerInfo) Get(key string) (*info, bool) {
	//log.Println(c.infos)
	val, ok := c.infos[key]
	return val, ok
}

func (c *ContainerInfo) autoRm() {
	timer := time.NewTicker(time.Second * 12)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			{
				for k, v := range c.infos {
					if !v.isConnect && time.Now().Unix()-v.deleteSign >= 15 || v.userDelete || !v.isConnect && v.record { //创建后15s未连接或断开后自动删除
						c.delete(k, create.Deleted)
						continue
					}
					v.updateRes(k)
					if v.overflow { //检查资源是否超出限制
						c.delete(k, create.Overstep)
					}
				}
			}
		}
	}
}

func (c *ContainerInfo) delete(cid string, status create.CStatus) {
	c.Lock()
	defer c.Unlock()
	var d = 3000 * time.Millisecond
	_ = DockerClient.ContainerStop(context.Background(), c.infos[cid].ContainerID, &d)
	if !c.infos[cid].AutoRemove {
		_ = DockerClient.ContainerRemove(context.Background(), c.infos[cid].ContainerID, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true})
	}
	//_ = queue.MClient.Del(create.ContainerListPreFix + strconv.Itoa(int(c.infos[cid].UserID)) + "/" + cid)
	create.UpdateStatus(cid, c.infos[cid].UserID, status) //更新缓存状态为已删除
	//_,_ = queue.MClient.Push(status.DeleteContainerPreFix,cid,cid)	//发送至删除队列
	delete(c.infos, cid)
}

func (i *info) updateRes(cid string) { //更新容器使用的资源
	rsp, err := DockerClient.ContainerStats(context.Background(), i.ContainerID, false)
	if err != nil {
		log.Println("ContainerID:", i.ContainerID, " Get Res Err:", err.Error())
		return
	}
	var v *types.StatsJSON
	var net uint64
	var block uint64
	err = json.NewDecoder(rsp.Body).Decode(&v)
	if err != nil {
		log.Println("ContainerID:", i.ContainerID, " Decode Res Err:", err.Error())
		return
	}
	if len(v.Networks) == 0 || len(v.BlkioStats.IoServiceBytesRecursive) == 0 {
		return
	}
	for _, v := range v.Networks {
		net += v.TxBytes + v.RxBytes
	}
	block += v.BlkioStats.IoServiceBytesRecursive[4].Value
	create.UpdateRecord(net, block, cid, i.UserID)
	if i.netLimit > 0 { //0表示不限制该资源
		if i.netLimit <= net {
			i.overflow = true
		}
	}
	if i.blockLimit > 0 {
		if i.blockLimit <= block {
			i.overflow = true
		}
	}
}

func (i *info) GetCmd() *exec.Cmd {
	var opt []string
	if !i.isConnect { //未连接过的则进行start
		i.isConnect = true
		opt = append(opt, "start")
		opt = append(opt, "-i")
	} else { //已连接的则使用attach进入
		opt = append(opt, "attach")
	}
	opt = append(opt, i.ContainerID)
	cmd := exec.Command("docker", opt...)
	return cmd
}

func (i *info) GetStatus() bool {
	//DockerClient.ContainerExecResize()
	return i.isConnect
}

func (i *info) Close() error {
	i.isConnect = false
	i.record = true
	i.deleteSign = time.Now().Unix()
	var d = 300 * time.Millisecond
	return dockerClient().ContainerStop(context.Background(), i.ContainerID, &d)
}
