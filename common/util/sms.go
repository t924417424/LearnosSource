package util

import (
	"Learnos/common/config"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"log"
	"math/rand"
	"regexp"
	"time"
)

func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func SendSms(phone string) (captcha string, err error) {
	conf := config.GetConf()
	captcha = createCaptcha()
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", conf.GateWay.AliSms.AccessID, conf.GateWay.AliSms.AccessKey)
	request := dysmsapi.CreateSendBatchSmsRequest()
	request.Scheme = "https"
	request.PhoneNumberJson = fmt.Sprintf("[%q]", phone)
	request.SignNameJson = fmt.Sprintf("[%q]", conf.GateWay.AliSms.SignName)
	request.TemplateCode = conf.GateWay.AliSms.Template
	request.TemplateParamJson = fmt.Sprintf("[{\"code\":\"%s\"}]", captcha)
	response, err := client.SendBatchSms(request)
	//if err != nil {
	//	log.Println(err.Error())
	//}
	if response.Code != "OK" {
		err = errors.New("短信服务器错误")
		log.Println(response)
	}
	return
}

func createCaptcha() string {
	return fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
}
