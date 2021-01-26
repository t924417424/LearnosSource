package websocket

import (
	"Learnos/Container/dockerControl"
	"Learnos/common/config"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Run() {
	conf := config.GetConf()
	//conf := config.GetConf()
	r := httprouter.New()
	r.GET("/term/:key", xtermM)
	if err := http.ListenAndServe(fmt.Sprintf(":%v",conf.WebSocket.Container.WsPort), r); err != nil {
		log.Fatal(err.Error())
	}
}

func xtermM(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //防止程序panic后导致的主程序宕机
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("server err"))
		}
	}()
	xterm(w, r, p)
}

func xterm(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	key := p.ByName("key")
	wsConn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer wsConn.Close()
	if key == "" {
		_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("key invalid"))
		return
	}
	info, ok := dockerControl.CInfo.Get(key)
	if !ok {
		_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("key non-existent"))
		return
	}
	cmd := info.GetCmd()
	term, err := newTerm(wsConn, cmd, info.ContainerID)
	if err != nil {
		_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
		return
	}
	exit := make(chan struct{}, 3)
	go term.start(exit)
	go term.outPutWs(exit)
	go term.inputTerm(exit)

	<-exit
	if err := info.Close(); err != nil {
		log.Println(err.Error())
	}
	log.Println("exit")
}
