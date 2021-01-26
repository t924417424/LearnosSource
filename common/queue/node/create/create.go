package create

import (
	"Learnos/common/queue"
	"Learnos/common/queueMsg/node/create"
	cnode "Learnos/proto/cnode"
	"github.com/gogo/protobuf/proto"
	"log"
	"time"
)

func QueueCreate(nodeId string) {
	for {
		msg, err := queue.MClient.Pop(create.ContainerCreatePreFix + nodeId)
		if err != nil {
			log.Println(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		if len(msg) == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
		var config cnode.CreateOpt
		err = proto.UnmarshalMerge(msg, &config)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println(config)
		log.Println(config.Config.Image)
	}
}
