package sms

import "encoding/json"

const (
	SendMsgPreFix   = "/msg/sendCode/Phone/"
	GetResultPreFix = "/msg/Verify/result/"
	IpVerifyPreFix = "/msg/sendCode/Ip/"
)

type SendMsg struct {
	Phone string
	Ip    string
}

type resultMsg struct {
	Captcha string
}

func NewSend(phone, ip string) ([]byte, error) {
	msg := SendMsg{Phone: phone, Ip: ip}
	msgByte, err := json.Marshal(msg)
	return msgByte, err
}

func GetSend(msg []byte) (SendMsg, error) {
	var result SendMsg
	err := json.Unmarshal(msg, &result)
	return result, err
}

func NewSmsMsg(captcha string) ([]byte, error) {
	msg := resultMsg{Captcha: captcha}
	msgByte, err := json.Marshal(msg)
	return msgByte, err
}

func GetSmsMsg(msg []byte) (resultMsg, error) {
	var result resultMsg
	err := json.Unmarshal(msg, &result)
	return result, err
}
