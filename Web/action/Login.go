package action

import (
	"Learnos/Web/callHelper"
	"Learnos/Web/formBind"
	gateway "Learnos/proto/gateway"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var rsp formBind.Rsp
	rsp.Code = formBind.FormErr
	rsp.Msg = "参数错误"
	var args formBind.Login
	if c.Bind(&args) == nil {
		//var gRsp gateway.CallRsp
		//opt := &gateway.Call{
		//	Type:   gateway.CallType_User,
		//	Caller: gateway.CallerType_UserLogin,
		//	Opt: &gateway.Options{
		//		User: &gateway.UserOpt{
		//			Username: args.Username,
		//			Password: args.Password,
		//		},
		//	},
		//}
		opt := gateway.Options{
			User: &gateway.UserOpt{
				Username: args.Username,
				Password: args.Password,
			},
		}
		rsp.Code = formBind.GateWayErr
		gRsp, err := callHelper.NewCall(c.ClientIP(), "").Call(gateway.CallType_User, gateway.CallerType_UserLogin).Do(opt)
		if err != nil {
			rsp.Data = err.Error()
		} else {
			if gRsp.Status == true {
				rsp.Code = formBind.Success
			}
			rsp.Msg = gRsp.Msg
			rsp.Token = gRsp.Token
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
		//	rsp.Token = gRsp.Token
		//}
	}
	c.JSON(http.StatusOK, rsp)
}
