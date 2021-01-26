package service

import (
	"Learnos/Container/dockerControl"
	"Learnos/Container/dockerUtil"
	node "Learnos/proto/cnode"
	"errors"
	"github.com/docker/docker/api/types/container"
)

func InspectOpt(req *node.CreateOpt) error { //创建容器
	//if err := checkImage(req.Config.Image); err != nil {              //查找镜像文件出错
	//	return "", err
	//}
	if req.Config == nil {
		return errors.New("请求未初始化")
	}
	if req.Config.HostName == "" {
		req.Config.HostName = "MyServer"
	}
	config := &container.Config{Hostname: req.Config.HostName, AttachStdout: true, AttachStdin: true, AttachStderr: true, OpenStdin: true, Tty: true, StdinOnce: true, Image: req.Config.Image, Cmd: []string{req.Config.Cmd}}
	hostConfig := dockerUtil.GetLimit(req.Resources)
	err := dockerControl.Create("", req.Uid, config, hostConfig)
	if err != nil {
		return err
	}
	return nil
}
