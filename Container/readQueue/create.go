package readQueue

import (
	"Learnos/Container/dockerControl"
	"Learnos/Container/dockerUtil"
	node "Learnos/proto/cnode"
	"errors"
	"github.com/docker/docker/api/types/container"
)

func inspectOpt(req *node.CreateOpt) (*container.Config, *container.HostConfig, error) { //创建容器
	if req.Config == nil {
		return nil, nil, errors.New("请求未初始化")
	}
	if req.Config.HostName == "" {
		req.Config.HostName = "MyServer"
	}
	//if err := dockerUtil.CheckImage(req.Config.Image); err != nil { //查找本地镜像出错则拉取DockerHub远程镜像
	//	if err := dockerUtil.PullImage(req.Config.Image); err != nil { //远程拉取镜像错误
	//		return nil, nil, err
	//	}
	//}
	config := &container.Config{Hostname: req.Config.HostName, AttachStdout: true, AttachStdin: true, AttachStderr: true, OpenStdin: true, Tty: true, StdinOnce: true, Image: req.Config.Image, Cmd: []string{req.Config.Cmd}}
	hostConfig := dockerUtil.GetLimit(req.Resources)
	//err := dockerControl.Create(req.Cid, config, hostConfig)
	//if err != nil {
	//	return err
	//}
	return config, hostConfig, nil
}

func checkImage(Image string) error {
	return dockerUtil.CheckImage(Image)
}

func pullImage(Image string) error {
	return dockerUtil.PullImage(Image)
}

func create(cid string, uid uint64, config *container.Config, hostConfig *container.HostConfig) error {
	return dockerControl.Create(cid, uid, config, hostConfig)
}
