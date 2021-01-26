package sshServer

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/gorilla/websocket"
	"io"
	"strings"
	"time"
)

func console(session ssh.Session) {
	defer func() {
		if err := recover(); err != nil {
			println("console：", err)
			//session.Write([]byte{230,156,141, 229, 138, 161, 233, 148 ,153 ,232, 175, 175 ,239 ,188,129})	//服务错误
		}
	}()
	defer session.Close()
	ip := session.Context().Value("ip").(string)
	token := session.Context().Value("token").(string)
	client := newClient(ip, token)
	//println(ip)
	//println(token)
	//检查是否创建container
	//_, _ = io.WriteString(session, "\u001B[2J")
	_, _ = io.WriteString(session, fmt.Sprintf("\x1b[2J\x1b[31;47m欢迎使用LearnosTerminal 当前登陆用户名: %s\x1b[0m\n", session.User()))
	images, err := client.getImage()
	if err != nil {
		_, _ = io.WriteString(session, fmt.Sprintf("\u001B[0;31m%s\t\u001B[0m\r\n", err.Error()))
		return
	}
	_, _ = io.WriteString(session, "当前支持的镜像列表\n\n")
	for k, _ := range images {
		_, err = io.WriteString(session, fmt.Sprintf("\u001B[0;32m%s\r\n\u001B[0m", k))
		if err != nil {
			return
		}
	}
	_, _ = io.WriteString(session, fmt.Sprintf("\n\n"))
	reader := bufio.NewReader(session)
	var imageId uint32
create:
	for {
		imageName := bytes.NewBufferString("")
		_, err = io.WriteString(session, fmt.Sprintf("\u001B[s请输入你要使用的系统（支持tab补全）："))
		if err != nil {
			return
		}
		for {
			//var tmp []byte
			key, err := reader.ReadByte()
			if err != nil {
				return
			}
			//println(key)
			if key == 13 { //回车
				break
			} else if key == 8 { //删除前一个输入
				//println(key)
				//println(imageName.String())
				if imageName.Len() >= 1 {
					imageName.Truncate(imageName.Len() - 1)
					//_, _ = session.Write([]byte{127})
					_, _ = io.WriteString(session, fmt.Sprintf("\u001B[1D\u001B[K"))
				} else {
					imageName.Reset()
				}
			} else if key == 9 { //tab补全
				for k, _ := range images {
					if k[:imageName.Len()] == imageName.String() {
						tabStr := k[imageName.Len():]
						_, _ = io.WriteString(session, k[imageName.Len():])
						imageName.WriteString(tabStr)
						break
					}
				}
			} else if key >= 65 && key <= 90 || key >= 97 && key <= 122 || key >= 47 && key <= 57 { //大小写字母/123456789
				imageName.WriteByte(key)
				_, _ = session.Write([]byte{key})
			}
		}
		if id, ok := images[imageName.String()]; ok {
			imageName.Reset()
			imageId = id
			break
		} else {
			imageName.Reset()
			_, _ = io.WriteString(session, fmt.Sprintf("\u001B[2K\u001B[u"))
		}
	}
	_, _ = io.WriteString(session, "\r\n")
	err = client.createContainer(imageId)
	if err != nil {
		_, _ = io.WriteString(session, fmt.Sprintf("\u001B[0;31m%s\t\u001B[0m\r\n", err.Error()))
		goto create
		//return
	}
	defer client.deleteContainer()
	for {
		status, netLimit, err := client.getContainerStatus()
		if err != nil {
			_, _ = io.WriteString(session, fmt.Sprintf("\u001B[0;31m%s\t\u001B[0m\r\n", err.Error()))
			client.deleteContainer()
			goto create
		}
		_, _ = io.WriteString(session, fmt.Sprintf("\u001B[s\u001B[2K%s\u001B[0;33m\u001B[5m%s\u001B[0m\u001B[0m\u001B[u\u001B[?25l", "实例状态：", containerMsg[status]))
		if status == errorCreate || status == pullImageErr || status == deleted || status == overstep {
			_, _ = io.WriteString(session, fmt.Sprintf("\u001B[s\u001B[2K%s\u001B[0;31m%s\u001B[0m\u001B[u\u001B[?25l", "实例状态：", containerMsg[status]))
			goto create
		} else if status == okCreate {
			_, _ = io.WriteString(session, fmt.Sprintf("\u001B[s\u001B[2K%s\u001B[0;32m%s\t%s\u001B[0m\u001B[u\u001B[?25l", "实例状态：", containerMsg[status], netLimit))
			_, _ = io.WriteString(session, fmt.Sprintf("\r\n\u001B[?25h\u001B[0;32m%s\u001B[0m\r\n", "创建成功，开始连接..."))
			break
		}
		time.Sleep(3 * time.Second)
	}
	containerNode, _, err := websocket.DefaultDialer.Dial(client.getWebSocket(), nil)
	//log.Println(client.getWebSocket())
	if err != nil {
		_, _ = session.Write([]byte(err.Error()))
		return
	}
	//_, _ = io.WriteString(session, "\u001B[2J")	//清屏
	defer containerNode.Close()
	go func() { //远程返回数据写入ssh
		for {
			_, msg, err := containerNode.ReadMessage()
			if err != nil {
				_, _ = session.Write([]byte("\r\nFailed to receive data!"))
				return
			}
			_, _ = session.Write(msg)
		}
	}()
	go func() { //监听窗口大小
		_, winCh, isPty := session.Pty()
		if isPty {
			for v := range winCh {
				//log.Println(fmt.Sprintf("{type: \"resize\", data: \"%d,%d\"}",v.Width,v.Height)) //这里用于发送设置窗口大小的命令
				_ = containerNode.WriteMessage(websocket.BinaryMessage, []byte(fmt.Sprintf("{\"type\": \"resize\", \"data\": \"%d,%d\"}", v.Width, v.Height)))
			}
		}
	}()
	//reader = bufio.NewReader(session)
	for { //ssh数据写入远程
		data, err := reader.ReadByte()
		if err != nil {
			//log.Println(err.Error())
			return
		}
		//println()
		//println(fmt.Sprintf("{\"type\": \"cmd\", \"data\": %q}", string(data)))
		//log.Println(fmt.Sprintf("{type: \"cmd\", data: \"%s\"}",string(data)))
		err = containerNode.WriteMessage(websocket.BinaryMessage, []byte(fmt.Sprintf("{\"type\": \"cmd\", \"data\": %s}", xTou(fmt.Sprintf("%q", string(data)))))) //%q安全转义，引号包裹
		if err != nil {
			_, _ = session.Write([]byte("Failed to send data!"))
			return
		}
	}
	//println(cid)
	//connToContainer(session)
	//log.Println(imageId)
}

func xTou(x string) (u string) {
	u = x
	if strings.HasPrefix(x, "\"\\x") {
		u = strings.Replace(x, "\"\\x", "\"\\u00", 1)
	}
	return u
}

//func renderbar(count, total int) {
//	barwidth := wscol - len("Progress: 100% []")
//	done := int(float64(barwidth) * float64(count) / float64(total))
//
//	fmt.Printf("Progress: \x1b[33m%3d%%\x1b[0m ", count*100/total)
//	fmt.Printf("[%s%s]",
//		strings.Repeat("=", done),
//		strings.Repeat("-", barwidth-done))
//}
