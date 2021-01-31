package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

func GetPubicIP() string {
	//conn, _ := net.Dial("udp", "8.8.8.8:80")
	//defer conn.Close()
	//localAddr := conn.LocalAddr().String()
	//idx := strings.LastIndex(localAddr, ":")
	//return localAddr[0:idx]
	response, err := http.Get("http://ip.cip.cc")
	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()
	res := ""
	// 类似的API应当返回一个纯净的IP地址
	for {
		tmp := make([]byte, 32)
		n, err := response.Body.Read(tmp)
		if err != nil {
			if err != io.EOF {
				panic(err.Error())
			}
			res += string(tmp[:n])
			break
		}
		res += string(tmp[:n])
	}
	return strings.TrimSpace(res)
}

func CheckIp(ip string) bool {
	addr := strings.Trim(ip, " ")
	regStr := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	if match, _ := regexp.MatchString(regStr, addr); match {
		return true
	}
	return false
}

func Md5(str string) string {
	m := md5.Sum([]byte(str))
	return hex.EncodeToString(m[:])
}

func GetPort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
