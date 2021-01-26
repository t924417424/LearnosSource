package websocket

import (
	"Learnos/Container/dockerControl"
	"bytes"
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type wsBufferWriter struct {
	dataSize    int
	consoleSize [2]int
	buffer      bytes.Buffer
	mu          sync.Mutex
}

func (w *wsBufferWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

func flushComboOutput(w *wsBufferWriter, wsConn *websocket.Conn) error {
	if w.buffer.Len() != 0 {
		//log.Println(w.buffer.Bytes())
		//超出consoleSize的宽度会从当前行重新写入修复,因为是用户输入，所以不需要递归处理,影响正常归位符，暂时注释
		//end := []byte{13}
		//endOk := []byte{13, 10}
		//index := bytes.Index(w.buffer.Bytes(), end)
		//index2 := bytes.Index(w.buffer.Bytes(), endOk)
		//if index >= 0 && index2 == -1 { //cmd.exec不同于ssh，row超出并不会有换行符，会导致前端xterm.js从头覆盖写入
		//	var before = make([]byte, len(w.buffer.Bytes()[:index+1]))
		//	copy(before, w.buffer.Bytes()[:index+1])
		//	newData := bytes.NewBuffer(before)
		//	after := w.buffer.Bytes()[index+len(end):]
		//	newData.Write([]byte{10})
		//	newData.Write(after)
		//	w.buffer.Reset() //释放原来的Buffer
		//	w.buffer = *newData
		//	newData.Reset() //因为上一行进行的是值拷贝，所以要手动释放一下newData
		//}
		//超出consoleSize的宽度会从当前行重新写入修复

		//DockerSdk未提供创建容器和Reset ConsoleSize的方法，此处自行实现，为适应不同分辨率的前端xterm，需要进行递归处理
		//w.dataSize += w.buffer.Len()
		//newLine := bytes.LastIndex(w.buffer.Bytes(), []byte{13, 10})
		//if newLine >= 0 {
		//	w.dataSize = len(w.buffer.Bytes()[newLine+2:])
		//}
		//DockerSdk未提供创建容器和Reset ConsoleSize的方法，此处自行实现，为适应不同分辨率的前端xterm

		err := wsConn.WriteMessage(websocket.BinaryMessage, w.buffer.Bytes())
		if err != nil {
			return err
		}
		w.buffer.Reset()
	}
	return nil
}

type wsTermConn struct {
	containerID string
	wsConn      *websocket.Conn
	stdinPipe   io.WriteCloser
	stdoutPut   *wsBufferWriter
	cmd         *exec.Cmd
}

func newTerm(ws *websocket.Conn, cmd *exec.Cmd, containerID string) (term *wsTermConn, err error) {
	output := new(wsBufferWriter)
	term = &wsTermConn{}
	term.stdoutPut = output
	cmd.Stdout = output
	cmd.Stderr = output
	term.containerID = containerID
	term.wsConn = ws
	term.cmd = cmd
	term.stdinPipe, err = cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	return term, nil
}

func (term *wsTermConn) start(exit chan struct{}) {
	defer setQuit(exit)
	if err := term.cmd.Run(); err != nil {
		log.Println("err:", err.Error())
		term.wsConn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
		return
	}
}

func (term *wsTermConn) getConsoleSize() (height, width int) {
	return term.stdoutPut.consoleSize[0], term.stdoutPut.consoleSize[1]
}

func (term *wsTermConn) setConsoleSize(width, height int) {
	term.stdoutPut.mu.Lock()
	defer term.stdoutPut.mu.Unlock()
	term.stdoutPut.consoleSize[0], term.stdoutPut.consoleSize[1] = height, width
	//dockerControl.DockerClient.ContainerResize()
	err := dockerControl.DockerClient.ContainerResize(context.Background(), term.containerID, types.ResizeOptions{Height: uint(height), Width: uint(width)})
	if err != nil {
		log.Println(err.Error())
	}
}

func (term *wsTermConn) outPutWs(exit chan struct{}) {
	defer setQuit(exit)
	tick := time.NewTicker(time.Millisecond * time.Duration(100))
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if err := flushComboOutput(term.stdoutPut, term.wsConn); err != nil {
				//log.Println(err.Error())
				return
			}
		case <-exit:
			return
		}
	}
}

type recvData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (term *wsTermConn) inputTerm(exit chan struct{}) {
	defer setQuit(exit)
	var recv recvData
	for {
		select {
		case <-exit:
			return
		default:
			_, msg, err := term.wsConn.ReadMessage()
			if err != nil {
				//log.Println(err.Error())
				return
			}
			//log.Println(bytes2str(msg))
			err = json.Unmarshal(msg, &recv)
			if err != nil {
				//log.Println(err.Error())
				continue
			}
			if recv.Type == "cmd" {
				term.stdinPipe.Write(str2bytes(recv.Data))
			} else if recv.Type == "resize" {
				size := strings.Split(recv.Data, ",")
				if len(size) == 2 {
					height, err := strconv.Atoi(size[0])
					if err != nil {
						//log.Println(err.Error())
						continue
					}
					width, err := strconv.Atoi(size[1])
					if err != nil {
						//log.Println(err.Error())
						continue
					}
					//dockerControl.DockerClient.ContainerExecResize(context.Background(),term.)
					//DockerClient.ContainerExecResize()
					//DockerClient.ContainerExecResize
					term.setConsoleSize(height, width)
				}
			}
		}
	}
}

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func setQuit(ch chan struct{}) {
	ch <- struct{}{}
}
