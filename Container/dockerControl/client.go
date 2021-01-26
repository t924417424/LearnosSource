package dockerControl

import (
	"Learnos/common/queueMsg/node/create"
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"log"
	"time"
)

var DockerClient *client.Client

func init() {
	DockerClient = dockerClient()
}

func dockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	_, err = cli.Ping(context.Background())
	if err != nil {
		log.Fatal("Docker Client 创建失败，请检查本机是否安装Docker！")
	}
	return cli
}

func ContainerList() error {
	list, err := DockerClient.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err == nil {
		for _, v := range list {
			log.Println(v.ID)
		}
	}
	return err
}

func Create(cid string, uid uint64, config *container.Config, hostConfig *container.HostConfig) (err error) {
	if cid == "" {
		return errors.New("CID不能为空")
	}
	cInfo, err := DockerClient.ContainerCreate(context.Background(), config, hostConfig, &network.NetworkingConfig{}, nil,"")
	if err != nil {
		return err
	}
	var autoRemove bool
	if hostConfig == nil {
		autoRemove = false
	} else {
		autoRemove = hostConfig.AutoRemove
	}
	net, block, err := create.GetLimit(cid, uid)
	if err != nil {
		return err
	}
	//key = uuid.Must(uuid.NewV4(), nil).String()
	CInfo.infos[cid] = &info{uid, cInfo.ID, time.Now().Unix(), false, false, autoRemove, false, net, block,false,false}
	return
}
