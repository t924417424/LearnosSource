package Dlog

import (
	protoLog "Learnos/proto/log"
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"log"
)

var publisher micro.Publisher

func NewLog(c micro.Service) {
	publisher = micro.NewPublisher("micro.serve.log", c.Client())
}

func Println(v ...interface{}) {
	if publisher != nil {
		var logMsg protoLog.LogMsg
		logMsg.Level = protoLog.LogType_Info
		logMsg.Content = fmt.Sprint(v...)
		_ = publisher.Publish(context.Background(), &logMsg)
	}
	log.Println(v)
}

func WarnInfo(v ...interface{}) {
	if publisher != nil {
		var logMsg protoLog.LogMsg
		logMsg.Level = protoLog.LogType_Warn
		logMsg.Content = fmt.Sprint(v...)
		_ = publisher.Publish(context.Background(), &logMsg)
	}
}

func Debug(v ...interface{}) {
	if publisher != nil {
		var logMsg protoLog.LogMsg
		logMsg.Level = protoLog.LogType_Debug
		logMsg.Content = fmt.Sprint(v...)
		_ = publisher.Publish(context.Background(), &logMsg)
	}
}

func Danger(v ...interface{}) {
	if publisher != nil {
		var logMsg protoLog.LogMsg
		logMsg.Level = protoLog.LogType_Danger
		logMsg.Content = fmt.Sprint(v...)
		_ = publisher.Publish(context.Background(), &logMsg)
	}
}
