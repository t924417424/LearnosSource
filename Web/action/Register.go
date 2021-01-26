package action

import (
	"Learnos/Web/callHelper"
	"Learnos/Web/formBind"
	gateway "Learnos/proto/gateway"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	var rsp formBind.Rsp
	rsp.Code = formBind.FormErr
	rsp.Msg = "参数错误"
	var args formBind.Register
	if c.Bind(&args) == nil {
		//var gRsp gateway.CallRsp
		//opt := &gateway.Call{
		//	Type:   gateway.CallType_User,
		//	Caller: gateway.CallerType_Register,
		//	Opt: &gateway.Options{
		//		User: &gateway.UserOpt{
		//			Username: args.Username,
		//			Password: args.Password,
		//			Phone:    args.Phone,
		//			Verify:   args.Verify,
		//		},
		//	},
		//}
		opt := gateway.Options{
			User: &gateway.UserOpt{
				Username: args.Username,
				Password: args.Password,
				Phone:    args.Phone,
				Verify:   args.Verify,
			},
		}
		rsp.Code = formBind.GateWayErr
		gRsp, err := callHelper.NewCall(c.ClientIP(), "").Call(gateway.CallType_User, gateway.CallerType_Register).Do(opt)
		if err != nil {
			rsp.Data = err.Error()
		} else {
			if gRsp.Status == true {
				rsp.Code = formBind.Success
			}
			rsp.Msg = gRsp.Msg
		}
		//ctx := metadata.NewContext(context.Background(), map[string]string{"Source-Ip": c.ClientIP()})
		//client := microClient.Get()
		//defer client.Close()
		//rsp.Code = formBind.GateWayErr
		//if err := client.Call(ctx, client.NewRequest(ServiceCall.GateWayServer, ServiceCall.GateWayService, opt), &gRsp); err != nil {
		//	rsp.Data = err.Error()
		//} else {
		//	if gRsp.Status == true {
		//		rsp.Code = formBind.Success
		//	}
		//	rsp.Msg = gRsp.Msg
		//}
	}
	c.JSON(http.StatusOK, rsp)
}
