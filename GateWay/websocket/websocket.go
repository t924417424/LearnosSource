package websocket

import (
	"Learnos/common/config"
	"Learnos/common/queue"
	"Learnos/common/queueMsg/node/create"
	"Learnos/common/util"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type auth struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

func Run() {
	conf := config.GetConf()
	r := httprouter.New()
	r.GET("/proxy/:cid", proxy)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.WebSocket.GateWay.WsPort), r); err != nil {
		log.Fatal(err.Error())
	}
}

func proxy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer func(w http.ResponseWriter) {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("server err"))
		}
	}(w)
	key := p.ByName("cid")
	wsConn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer wsConn.Close()
	var login auth
	var uid uint
	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			return
		}
		err = json.Unmarshal(msg, &login)
		if err != nil {
			_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
			continue
		}
		u, err := util.ParseToken(login.Token)
		if err != nil {
			_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("身份令牌失效"))
			return
		}
		if err := u.Valid(); err != nil {
			_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("验证不通过"))
			return
		}
		uid = u.UserId
		break
	}
	info, err := queue.MClient.Get(create.ContainerListPreFix + strconv.Itoa(int(uid)) + "/" + key)
	if err != nil {
		_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
		return
	}
	if len(info) == 0 {
		_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("服务不可用，请刷新重试"))
		return
	}
	server := create.GetCreateMessage(info)
	if server == nil {
		_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("服务信息获取失败！"))
		return
	}
	if server.Status != create.OkCreate {
		_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("创建暂未成功，请稍后重试"))
		return
	}
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:8015", server.Addr), Path: fmt.Sprintf("/term/%s", server.Uuid)} //连接容器所在节点
	containerNode, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("节点连接失败！"))
		return
	}
	defer containerNode.Close()
	go func() {
		for {
			_, msg, err := containerNode.ReadMessage()
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
		_ = containerNode.WriteMessage(websocket.BinaryMessage, msg)
	}
}
