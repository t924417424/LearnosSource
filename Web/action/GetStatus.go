package action

import (
	"Learnos/Web/callHelper"
	"Learnos/Web/formBind"
	gateway "Learnos/proto/gateway"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

func GetStatus(c *gin.Context) {
	var rsp formBind.Rsp
	rsp.Code = formBind.FormErr
	rsp.Msg = "参数错误"
	var args formBind.ContainerApi
	token, ok := c.Get("token")
	if !ok || token == "" {
		rsp.Code = formBind.TokenErr
		rsp.Msg = "未取得有效Token"
		c.JSON(http.StatusOK, rsp)
		return
	}
	if err := c.Bind(&args); err == nil {
		//var gRsp gateway.CallRsp
		//opt := &gateway.Call{
		//	Type:   gateway.CallType_Container,
		//	Caller: gateway.CallerType_GetContainerStatus,
		//	Opt: &gateway.Options{
		//		Cid: args.Cid,
		//	},
		//}
		opt := gateway.Options{
			Cid: args.Cid,
		}
		rsp.Code = formBind.GateWayErr
		gRsp, err := callHelper.NewCall(c.ClientIP(), token.(string)).Call(gateway.CallType_Container, gateway.CallerType_GetContainerStatus).Do(opt)
		if err != nil {
			rsp.Data = err.Error()
		} else {
			if gRsp.Status == true {
				rsp.Code = formBind.Success
				rsp.Data = gRsp.Data.Status
				if args.Limit {
					if gRsp.Data.NetWorkLimit > 0{
						rsp.Data = math.Ceil(float64(gRsp.Data.NetWorkRecord) / float64(gRsp.Data.NetWorkLimit) * 100)	//取流量百分比
					}else{
						rsp.Data = 0
					}
				}
			}
			rsp.Msg = gRsp.Msg
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
		//		rsp.Data = gRsp.Data.Status
		//		if args.Limit {
		//			if gRsp.Data.NetWorkLimit > 0{
		//				rsp.Data = math.Ceil(float64(gRsp.Data.NetWorkRecord) / float64(gRsp.Data.NetWorkLimit) * 100)	//取流量百分比
		//			}else{
		//				rsp.Data = 0
		//			}
		//		}
		//	}
		//	rsp.Msg = gRsp.Msg
		//}
	}
	c.JSON(http.StatusOK, rsp)
}