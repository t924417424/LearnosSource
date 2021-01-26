package handler

import (
	"Learnos/Container/service"
	"Learnos/common/microClient/Dlog"
	node "Learnos/proto/cnode"
	"errors"
)
import "context"

type Handler struct {
}

func (c Handler) Service(ctx context.Context, opt *node.CallOpt, rsp *node.CallRsp) error {
	var err error
	if opt.Type == node.CallType_CreateContainer && opt.Create != nil {
		err = service.InspectOpt(opt.Create)
		if err != nil {
			rsp.Status = false
			rsp.Msg = err.Error()
		} else {
			rsp.Status = true
			rsp.Data = opt.Create.Cid
		}
	} else {
		err = errors.New("请求参数错误")
	}
	return err
}

func (c Handler) CreateMsg(ctx context.Context,opt *node.CreateNotice,rsp *node.CreateRsp) error {
	Dlog.Println(opt.Cid)
	return nil
}
