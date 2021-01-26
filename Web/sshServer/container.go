package sshServer

import (
	"Learnos/Web/callHelper"
	"Learnos/common/config"
	gateway "Learnos/proto/gateway"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type client struct {
	ip    string
	token string
	cid   string
	addr  string
}

const (
	loading      uint32 = 0
	pullImage    uint32 = 1
	startCreate  uint32 = 2
	okCreate     uint32 = 3
	errorCreate  uint32 = 4
	pullImageErr uint32 = 5
	deleted      uint32 = 6
	overstep     uint32 = 7
)

var containerMsg = map[uint32]string{
	loading:      "等待创建",
	pullImage:    "正在拉取镜像",
	pullImageErr: "拉取镜像失败",
	startCreate:  "开始创建",
	okCreate:     "创建成功",
	errorCreate:  "创建失败",
	deleted:      "已删除",
	overstep:     "触发资源限制",
}

func newClient(ip, token string) client {
	return client{ip, token, "", ""}
}

func (c client) getWebSocket() (wsUrl string) {
	conf := config.GetConf()	//获取配置，连接container的websocket服务（因为在终端操作，会话信息保存在服务端，此处省去一个gateway节点鉴权的步骤，直连container节点）
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", c.addr, conf.WebSocket.Container.WsPort), Path: fmt.Sprintf("/term/%s", c.cid)}
	return u.String()
}

func (c client) getImage() (list map[string]uint32, err error) {
	list = make(map[string]uint32)
	err = errors.New("获取可用镜像出错！")
	rsp, err := callHelper.NewCall(c.ip, c.token).Call(gateway.CallType_Container, gateway.CallerType_GetImageList).Do(gateway.Options{})
	if err != nil {
		//log.Println(err.Error())
		return
	}
	if rsp.Status == true {
		for _, v := range rsp.ImageList {
			list[v.ImageName] = v.Id
		}
		err = nil
	} else {
		err = errors.New(rsp.Msg)
	}
	return
}

func (c *client) createContainer(imageId uint32) (err error) {
	err = errors.New("创建失败！")
	c.cid = ""
	c.addr = ""
	opt := gateway.Options{
		Create: &gateway.ImageInfo{
			ID: imageId,
		},
	}
	gRsp, err := callHelper.NewCall(c.ip, c.token).Call(gateway.CallType_Container, gateway.CallerType_CreateContainer).Do(opt)
	if err != nil {
		return
	} else {
		if gRsp.Status == true {
			c.cid = gRsp.Cid
			err = nil
		} else {
			err = errors.New(gRsp.Msg)
		}
	}
	return
}

func (c *client) getContainerStatus() (status uint32, limit string, err error) {
	err = errors.New("暂无可用容器！")
	opt := gateway.Options{
		Cid: c.cid,
	}
	gRsp, err := callHelper.NewCall(c.ip, c.token).Call(gateway.CallType_Container, gateway.CallerType_GetContainerStatus).Do(opt)
	if err != nil {
		return
	} else {
		if gRsp.Status == true {
			if gRsp.Data.NetWorkLimit == 0 {
				limit = "网络流量：无限制"
			} else {
				limit = "网络流量：" + strconv.Itoa(int(gRsp.Data.NetWorkLimit/1000/1024)) + "M"
			}
			status = gRsp.Data.Status
			c.addr = gRsp.Data.Addr
			err = nil
		} else {
			err = errors.New(gRsp.Msg)
		}
	}
	return
}

func (c client) deleteContainer() {
	if c.cid == "" {
		return
	}
	opt := gateway.Options{
		Cid: c.cid,
	}
	_, err := callHelper.NewCall(c.ip, c.token).Call(gateway.CallType_Container, gateway.CallerType_DeleteContainer).Do(opt)
	if err != nil {
		//rsp.Data = err.Error()
	}
}
