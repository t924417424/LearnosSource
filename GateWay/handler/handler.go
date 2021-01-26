package handler

import (
	"Learnos/GateWay/service"
	gateway "Learnos/proto/gateway"
	"context"
	"errors"
	"log"
)

type Handler struct {
}

func (h Handler) Service(ctx context.Context, call *gateway.Call, rsp *gateway.CallRsp) error {
	defer func(rsp *gateway.CallRsp) {
		if err := recover(); err != nil {
			log.Println(err)
			rsp.Msg = "网关错误"
		}
	}(rsp)
	if call.Opt == nil {
		return errors.New("参数错误")
	}
	if call.Type == gateway.CallType_User && call.Opt != nil {
		service.UserHandler(call.Opt, call.Caller, rsp, ctx)
	} else if call.Type == gateway.CallType_Container {
		service.ContainerHandler(call.Opt, call.Caller, rsp, ctx)
	} else {
		rsp.Code = 0
		rsp.Status = false
		rsp.Msg = "参数错误"
		return errors.New("参数错误")
	}
	return nil
}
