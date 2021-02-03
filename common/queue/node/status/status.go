package status

import (
	"Learnos/common/queue"
	"Learnos/common/queueMsg/node/status"
	"context"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"strconv"
)

type NodeStatus struct {
	nodeId string
	nodeIp string
	ctx    context.Context
}

func NewNodeStatus(id, ip string, ctx context.Context) *NodeStatus {
	return &NodeStatus{id, ip, ctx}
}

func (m *NodeStatus) Get() {
	log.Println(m)
}

func (m *NodeStatus) Send() error {
	mem := strconv.Itoa(int(GetFree()))
	key := status.ServerStatusPreFix + m.nodeId
	return queue.MClient.Set(key, mem) //上报内存空余，网关选择节点时使用用
}

func (m *NodeStatus) KeepExists() chan error {
	mem := strconv.Itoa(int(GetFree()))
	key := status.ServerStatusPreFix + m.nodeId
	exit := make(chan error, 1)
	ctx, cancel := context.WithCancel(m.ctx)
	go func() {
		queue.MClient.KeepAlive(m.ctx, key, mem, 30, exit)
		select {
		case <-ctx.Done():
			cancel()
			return
		}
	}()
	return exit
}

func (m *NodeStatus) Delete() error {
	key := status.ServerStatusPreFix + m.nodeId
	return queue.MClient.Del(key)
}

func GetFree() uint64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return v.Free / 1024 / 1024
}
