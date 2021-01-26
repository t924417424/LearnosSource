package sms

import (
	"Learnos/common/queue"
	"Learnos/common/queueMsg/gateway/sms"
	"Learnos/common/util"
	"log"
	"time"
)

func Recv() {
	for {
		msg, err := queue.MClient.Pop(sms.SendMsgPreFix)
		if err != nil {
			log.Println(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		if len(msg) == 0 {
			time.Sleep(10 * time.Second)
			continue
		}
		info,_ := sms.GetSend(msg)
		if util.VerifyMobileFormat(info.Phone) {
			key := sms.GetResultPreFix + info.Phone
			ipTmp := sms.IpVerifyPreFix + info.Ip
			if ok, _ := queue.MClient.Exists(key); !ok {
				captcha, err := util.SendSms(info.Phone)
				if err == nil{
					msg, err := sms.NewSmsMsg(captcha)
					if err == nil {
						queue.MClient.SetEx(ipTmp,info.Ip,120)     //缓存IP——Key，避免频繁请求接口
						queue.MClient.SetEx(key, string(msg), 120) //将结果发送至消息队列
					}
				}
			}
		}
	}
}
