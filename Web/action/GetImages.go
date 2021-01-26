package action

import (
	"Learnos/Web/callHelper"
	"Learnos/Web/formBind"
	gateway "Learnos/proto/gateway"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetImages(c *gin.Context) {
	var rsp formBind.Rsp
	var list []map[string]interface{}
	token, ok := c.Get("token")
	if !ok || token == "" {
		rsp.Code = formBind.TokenErr
		rsp.Msg = "未取得有效Token"
		c.JSON(http.StatusOK, rsp)
		return
	}
	//var gRsp gateway.CallRsp
	//opt := &gateway.Call{
	//	Type:   gateway.CallType_Container,
	//	Caller: gateway.CallerType_GetImageList,
	//	Opt: &gateway.Options{
	//
	//	},
	//}
	opt := gateway.Options{}
	rsp.Code = formBind.GateWayErr
	gRsp, err := callHelper.NewCall(c.ClientIP(), token.(string)).Call(gateway.CallType_Container, gateway.CallerType_GetImageList).Do(opt)
	if err != nil {
		rsp.Data = err.Error()
	} else {
		if gRsp.Status == true {
			rsp.Code = formBind.Success
			for _, v := range gRsp.ImageList {
				//log.Println(gRsp.ImageList)
				list = append(list, map[string]interface{}{
					"Id":   v.Id,
					"Name": v.ImageName,
					"Logo": v.Logo,
				})
			}
		}
		rsp.Msg = gRsp.Msg
	}
	rsp.Data = list
	//ctx := metadata.NewContext(context.Background(), map[string]string{"Source-Ip": c.ClientIP(), "Authorization": token.(string)})
	//client := microClient.Get()
	//defer client.Close()
	//rsp.Code = formBind.GateWayErr
	//if err := client.Call(ctx, client.NewRequest(ServiceCall.GateWayServer, ServiceCall.GateWayService, opt), &gRsp); err != nil {
	//	rsp.Data = err.Error()
	//} else {
	//	if gRsp.Status == true {
	//		rsp.Code = formBind.Success
	//		for _, v := range gRsp.ImageList {
	//			log.Println(gRsp.ImageList)
	//			list = append(list, map[string]interface{}{
	//				"Id":   v.Id,
	//				"Name": v.ImageName,
	//				"Logo": v.Logo,
	//			})
	//		}
	//	}
	//	rsp.Msg = gRsp.Msg
	//	rsp.Data = list
	//}
	c.JSON(http.StatusOK, rsp)
}
