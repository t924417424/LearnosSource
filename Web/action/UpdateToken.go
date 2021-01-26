package action

import (
	"Learnos/Web/callHelper"
	"Learnos/Web/formBind"
	gateway "Learnos/proto/gateway"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateToken(c *gin.Context) {
	var rsp formBind.Rsp
	token, ok := c.Get("token")
	if !ok || token == "" {
		rsp.Code = formBind.TokenErr
		rsp.Msg = "未取得有效Token"
		c.JSON(http.StatusOK, rsp)
		return
	}
	rsp.Code = formBind.FormErr
	rsp.Msg = "参数错误"
	//var gRsp gateway.CallRsp
	//opt := &gateway.Call{
	//	Type:   gateway.CallType_User,
	//	Caller: gateway.CallerType_UpdateToken,
	//	Opt: &gateway.Options{
	//	},
	//}
	opt := gateway.Options{}
	rsp.Code = formBind.GateWayErr
	gRsp, err := callHelper.NewCall(c.ClientIP(), token.(string)).Call(gateway.CallType_User, gateway.CallerType_UpdateToken).Do(opt)
	if err != nil {
		rsp.Data = err.Error()
	} else {
		if gRsp.Status == true {
			rsp.Code = formBind.Success
		}
		rsp.Msg = gRsp.Msg
		rsp.Token = gRsp.Token
	}
	//ctx := metadata.NewContext(context.Background(), map[string]string{"Source-Ip": c.ClientIP(), "Authorization": token.(string)})
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
	c.JSON(http.StatusOK, rsp)
}
