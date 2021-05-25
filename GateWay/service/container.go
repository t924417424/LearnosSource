package service

import (
	"Learnos/GateWay/sqldata/model"
	"Learnos/GateWay/sqldata/mysql"
	"Learnos/common/microClient"
	"Learnos/common/microClient/ServiceCall"
	"Learnos/common/queue"
	"Learnos/common/queueMsg/node/create"
	"Learnos/common/queueMsg/node/status"
	node "Learnos/proto/cnode"
	gateway "Learnos/proto/gateway"
	"context"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

type createInfo struct {
	Token   string
	Uid     uint
	ImageID uint32
}

func newCreateInfo(opt *gateway.Options, token string, user model.User) createInfo {
	imageId := uint32(0)
	if opt.Create != nil {
		imageId = opt.Create.ID
	}
	return createInfo{token, user.ID, imageId}
}

func (c *createInfo) getImageList() (list []*gateway.ImageList, err error) {
	db, err := mysql.Get()
	if err != nil {
		return list, errors.New("数据库连接失败！")
	}
	//defer db.DB.Close()
	var listInfo []model.Image
	result := db.DB.Model(model.Image{}).Find(&listInfo)
	if result.Error != nil || result.RowsAffected == 0 {
		return list, errors.New("获取可用镜像失败")
	}
	for _, v := range listInfo {
		list = append(list, &gateway.ImageList{Id: uint32(v.ID), ImageName: v.Name, Logo: v.Logo})
	}
	return list, nil
}

func (c *createInfo) createContainer() (string, error) {
	var cid string
	if c.ImageID == 0 {
		return "", errors.New("系统参数错误")
	}
	count, _, err := queue.MClient.GetByPreFix(create.ContainerListPreFix + strconv.Itoa(int(c.Uid)) + "/") //检查已创建的容器
	if err != nil {
		return cid, errors.New("连接缓存服务器失败！")
	}
	if count > 0 {
		return cid, errors.New(fmt.Sprintf("您当前有 %d 个容器还在使用或待删除，请稍候重试！", count))
	}
	client := microClient.Get()
	defer client.Close()
	serList, err := client.Options().Registry.GetService(ServiceCall.ContainerServer) //获取所有服务节点
	if err != nil {
		return "", err
	}
	if len(serList) < 1 {
		return "", errors.New("暂无可用节点")
	}
	var mem string
	//var ip string
	//var nodeId string
	for _, v := range serList {
		for _, node := range v.Nodes {
			nid := strings.SplitN(node.Id, "-", 2)[1]
			memTmp, err := queue.MClient.Get(status.ServerStatusPreFix + nid)
			if err != nil {
				continue
			}
			if string(memTmp) > mem {
				mem = string(memTmp) //选取最多空闲内存的节点信息
				//ip = node.Address
				//nodeId = nid
			}
		}
	}
	if mem < "100" { //最多空闲节点的内容小于100M
		return "", errors.New("所有节点繁忙，请错峰重试！")
	}
	db, err := mysql.Get()
	if err != nil {
		return "", errors.New("数据库连接失败！")
	}
	//defer db.DB.Close()
	var imageInfo model.Image
	imageInfo.ID = uint(c.ImageID)
	result := db.DB.Where(imageInfo).First(&imageInfo)
	if result.Error != nil || result.RowsAffected != 1 {
		return "", errors.New("系统信息获取失败！")
	}
	if imageInfo.Name == "" || imageInfo.Cmd == "" {
		return "", errors.New("基础创建参数无效！")
	}
	createOpt := &node.CreateOpt{
		Config: &node.CreateConfig{
			Image:           imageInfo.Name,
			Cmd:             imageInfo.Cmd,
			NetworkDisabled: imageInfo.Network,
		},
		Resources: &node.Resources{
			Memory: imageInfo.Memory,
		},
	}
	cid = uuid.New().String()
	createOpt.Cid = cid
	createOpt.Uid = uint64(c.Uid)
	CreateOpt, err := proto.Marshal(createOpt)
	if err != nil {
		return "", errors.New("配置文件创建失败！")
	}

	CreateMsg, err := create.NewCreateMessage(create.Loading, c.Uid, "", "", cid, createOpt.Config.Image, imageInfo.NetWorkIoLimit, imageInfo.BlockIoLimit)
	if err != nil {
		return "", errors.New("创建消息索引失败！")
	}

	//c.CreateOpt.Status = cNode.CreateStatus_Loading

	//db, err := mysql.Get()
	//if err != nil {
	//	return "", errors.New("数据库连接失败！")
	//}
	//defer db.Close()
	createHistory := db.DB.Begin() //开启事务
	result = db.DB.Create(&model.History{Cid: cid, Status: 0, ImagesName: createOpt.Config.Image, BindUser: c.Uid})
	if result.Error != nil || result.RowsAffected < 1 {
		createHistory.Rollback()
		return "", errors.New("创建使用记录失败")
	}
	err = queue.MClient.SetEx(create.ContainerCreatePreFix+cid, string(CreateOpt), 60*9) //发送至节点队列，9分钟未处理则丢掉
	if err != nil {
		createHistory.Rollback()
		return cid, err
	}
	err = queue.MClient.SetEx(create.ContainerListPreFix+strconv.Itoa(int(c.Uid))+"/"+cid, string(CreateMsg), 60*10) //10分钟自动创建失败
	if err != nil {
		createHistory.Rollback()
		return cid, err
	}
	createHistory.Commit()
	return cid, nil
}

func (c *createInfo) getContainerStatus(cid string) (containers *gateway.ContainerStatus, status string, err error) {
	var data []byte
	if cid != "" {
		data, err = queue.MClient.Get(create.ContainerListPreFix + strconv.Itoa(int(c.Uid)) + "/" + cid)
	} else {
		n, infos, err := queue.MClient.GetByPreFix(create.ContainerListPreFix + strconv.Itoa(int(c.Uid)) + "/")
		if n != 0 && err == nil {
			for _, v := range infos {
				tmp := create.GetCreateMessage(v).Status
				if tmp == create.OkCreate {
					data = v
					break
				}
			}
		}
	}
	if err != nil || len(data) == 0 {
		return containers, status, errors.New("无可用容器")
	}
	//for _, v := range data {
	//	tmp := create.GetCreateMessage(data)
	//	if tmp.Status == create.Overstep || tmp.Status == create.Deleted || tmp.Status == create.PullImageErr || tmp.Status == create.ErrorCreate {
	//		//存入数据库并删除缓存
	//	}
	//	containers = &gateway.ContainerStatus{Cid: tmp.Uuid, NodeId: tmp.NodeId, Status: tmp.Msg, Image: tmp.Image}
	//}
	tmp := create.GetCreateMessage(data)
	//if tmp.Status == create.Overstep || tmp.Status == create.Deleted || tmp.Status == create.PullImageErr || tmp.Status == create.ErrorCreate {
	//	//存入数据库并删除缓存
	//}
	if tmp.Uid == c.Uid {
		containers = &gateway.ContainerStatus{Cid: tmp.Uuid, NodeId: tmp.NodeId, Status: uint32(tmp.Status), Image: tmp.Image, NetWorkRecord: tmp.Network.Record, NetWorkLimit: tmp.Network.Limit, Addr: tmp.Addr}
	} else {
		return nil, status, errors.New("权限认证失败！")
	}
	return containers, status, nil
}

func (c *createInfo) deleteContainer(cid string) (err error) {
	if cid == "" {
		return errors.New("实例参数错误")
	}
	data, err := queue.MClient.Get(create.ContainerListPreFix + strconv.Itoa(int(c.Uid)) + "/" + cid)
	if err != nil {
		return errors.New("获取实例信息失败")
	}
	if len(data) == 0 {
		return errors.New("实例不存在")
	}
	info := create.GetCreateMessage(data)
	if info == nil {
		return errors.New("实例所在节点信息失败")
	}
	opt := &node.CreateOpt{
		Cid:    cid,
		Status: node.CreateStatus_Delete,
	}
	delOpt, err := proto.Marshal(opt)
	if err != nil {
		return errors.New("消息创建失败！")
	}

	ok, err := queue.MClient.Push(create.ContainerDeletePreFix, info.NodeId+"/"+cid, string(delOpt))
	if err != nil || !ok {
		return errors.New("消息发送失败！")
	}
	return nil
}

func ContainerHandler(opt *gateway.Options, caller gateway.CallerType, rsp *gateway.CallRsp, ctx context.Context) {
	var token string
	var userInfo model.User
	//log.Println(ctx.Value("token"))
	//log.Println(ctx.Value("userInfo"))
	if ctx.Value("token") == nil || ctx.Value("userInfo") == nil {
		rsp.Msg = "参数错误"
		return
	}
	token = ctx.Value("token").(string)
	userInfo = ctx.Value("userInfo").(model.User)
	info := newCreateInfo(opt, token, userInfo)
	if caller == gateway.CallerType_CreateContainer {
		rsp.Msg = "请求成功"
		cid, err := info.createContainer()
		if err != nil {
			rsp.Msg = err.Error()
			//rsp.Data = err.Error()
			return
		}
		rsp.Cid = cid
		//rsp.Data = cid
	} else if caller == gateway.CallerType_GetContainerStatus {
		containerStatus, status, err := info.getContainerStatus(opt.Cid)
		if err != nil {
			rsp.Msg = err.Error()
			return
		}
		rsp.Msg = status
		rsp.Data = containerStatus
	} else if caller == gateway.CallerType_DeleteContainer {
		rsp.Msg = "请求成功"
		err := info.deleteContainer(opt.Cid)
		if err != nil {
			rsp.Msg = err.Error()
			//rsp.Data = err.Error()
			return
		}
	} else if caller == gateway.CallerType_GetImageList {
		rsp.Msg = "请求成功"
		list, err := info.getImageList()
		if err != nil {
			rsp.Msg = err.Error()
			//rsp.Data = err.Error()
			return
		}
		rsp.ImageList = list
	}
	rsp.Code = 1
	rsp.Status = true
}
