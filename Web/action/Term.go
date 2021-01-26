package action

import (
	"Learnos/common/config"
	"Learnos/common/microClient"
	"Learnos/common/microClient/ServiceCall"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net"
	"net/http"
	"net/url"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Term(c *gin.Context){
	defer func() {
		if err := recover(); err != nil{
			c.JSON(http.StatusBadGateway,err)
		}
	}()
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil{
		c.JSON(http.StatusBadGateway,err.Error())
	}
	defer wsConn.Close()
	client := microClient.Get()
	defer client.Close()
	serList, err := client.Options().Registry.GetService(ServiceCall.GateWayServer) //获取所有服务节点
	if err != nil {
		_ = wsConn.WriteMessage(websocket.BinaryMessage,[]byte(err.Error()))
		return
	}
	if len(serList) < 1 {
		_ = wsConn.WriteMessage(websocket.BinaryMessage,[]byte("暂无可用网关节点"))
		return
	}
	cid := c.Param("cid")
	ip := serList[rand.Intn(len(serList))].Nodes[0].Address
	host,_,_ := net.SplitHostPort(ip)
	//获取配置，连接gateway的websocket服务
	conf := config.GetConf()
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", host,conf.WebSocket.GateWay.WsPort), Path: fmt.Sprintf("/proxy/%s", cid)}
	gatewayWs,_,err := websocket.DefaultDialer.Dial(u.String(),nil)
	if err != nil {
		_ = wsConn.WriteMessage(websocket.BinaryMessage,[]byte(err.Error()))
	}
	defer gatewayWs.Close()
	go func() {
		for {
			_, msg, err := gatewayWs.ReadMessage()
			if err != nil {
				_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("Failed to receive data!"))
				return
			}
			_ = wsConn.WriteMessage(websocket.BinaryMessage, msg)
		}
	}()
	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("Failed to send data!"))
			return
		}
		_ = gatewayWs.WriteMessage(websocket.BinaryMessage, msg)
	}
}