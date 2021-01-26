package callHelper

import (
	"Learnos/common/microClient"
	"Learnos/common/microClient/ServiceCall"
	gateway "Learnos/proto/gateway"
	"context"
	"github.com/micro/go-micro/metadata"
)

type info struct {
	ip    string
	token string
}

type callOpt struct {
	info       info
	callType   gateway.CallType
	callerType gateway.CallerType
}

func NewCall(ip, token string) info {
	return info{ip, token}
}

func (i info) Call(callType gateway.CallType, callerType gateway.CallerType) callOpt {
	opt := callOpt{
		info:       i,
		callType:   callType,
		callerType: callerType,
	}
	return opt
}

func (c callOpt) Do(option gateway.Options) (gRsp gateway.CallRsp, err error) {
	opt := &gateway.Call{
		Type:   c.callType,
		Caller: c.callerType,
		Opt:    &option,
	}
	metaData := make(map[string]string)
	//map[string]string{"Source-Ip": c.info.ip, "Authorization": c.info.token}
	metaData["Source-Ip"] = c.info.ip
	if c.info.token != "" {
		metaData["Authorization"] = c.info.token
	}
	ctx := metadata.NewContext(context.Background(), metaData)
	client := microClient.Get()
	defer client.Close()
	err = client.Call(ctx, client.NewRequest(ServiceCall.GateWayServer, ServiceCall.GateWayService, opt), &gRsp)
	//log.Println(gRsp)
	return
}
