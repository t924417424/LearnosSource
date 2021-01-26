package util

import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"regexp"
	"strings"
)

func GetPulicIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
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
