package dockerUtil

import (
	"Learnos/Container/dockerControl"
	node "Learnos/proto/cnode"
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"io"
	"log"
)

func CheckImage(Image string) error {
	if Image == "" {
		return errors.New("镜像名不能为空")
	}
	list, err := dockerControl.DockerClient.ImageList(context.Background(), types.ImageListOptions{All: false, Filters: filters.NewArgs(filters.Arg("reference", Image))})
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return errors.New("image no find")
	}
	return err
}

func PullImage(Image string) error {
	if Image == "" {
		return errors.New("镜像名不能为空")
	}
	rsp, err := dockerControl.DockerClient.ImageSearch(context.Background(), Image, types.ImageSearchOptions{Filters: filters.NewArgs(filters.Arg("is-official", "true")), Limit: 1})
	if err != nil || len(rsp) < 1 {
		return errors.New("查找镜像文件出错")
	}
	pio, err := dockerControl.DockerClient.ImagePull(context.Background(), Image, types.ImagePullOptions{})
	defer pio.Close()
	if err != nil {
		return errors.New("下载镜像文件出错")
	}
	var buff = make([]byte, 1024)
	for {
		_, err := pio.Read(buff)
		if err == io.EOF {
			log.Println("Pull Images Over")
			break
		}
	}
	return nil
}

func GetLimit(resources *node.Resources) *container.HostConfig {
	if resources == nil {
		return nil
	}
	limit := &container.HostConfig{}
	//limit.ConsoleSize = [2]uint{32,130}
	limit.AutoRemove = resources.AutoRemove
	if resources.Memory > 0 {
		limit.Memory = resources.Memory
	}
	//if resources.BlkioDeviceReadBps.Path != "" && resources.BlkioDeviceReadBps.Rate > 0 {
	//	limit.BlkioDeviceReadBps = append(limit.BlkioDeviceReadBps, &blkiodev.ThrottleDevice{Path: resources.BlkioDeviceReadBps.Path, Rate: resources.BlkioDeviceReadBps.Rate})
	//}
	//if resources.BlkioDeviceWriteBps.Path != "" && resources.BlkioDeviceWriteBps.Rate > 0 {
	//	limit.BlkioDeviceReadBps = append(limit.BlkioDeviceWriteBps, &blkiodev.ThrottleDevice{Path: resources.BlkioDeviceWriteBps.Path, Rate: resources.BlkioDeviceWriteBps.Rate})
	//}
	//if resources.KernelMemoryTCP > 0 {
	//	limit.KernelMemoryTCP = resources.KernelMemoryTCP
	//}
	return limit
}
