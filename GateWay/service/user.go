package service

import (
	"Learnos/GateWay/sqldata/model"
	"Learnos/GateWay/sqldata/mysql"
	"Learnos/common/queue"
	"Learnos/common/queueMsg/gateway/sms"
	"Learnos/common/queueMsg/gateway/user"
	util2 "Learnos/common/util"
	gateway "Learnos/proto/gateway"
	"context"
	"errors"
	"time"
)

type UserImpl interface {
	Registry() error
	UserLogin() (string, error)
	UpdateToken() (string, error)
	CheckPhone() error
	Logout() error
	//SendCode()
}

type Info struct {
	UserName string
	Password string
	Phone    string
	Token    string
	Verify   string
	ClientIP string
}

func newInfo(opt *gateway.Options, ClientIP string, token string) UserImpl {
	opts := &gateway.UserOpt{}
	if opt.User != nil{
		opts = opt.User
	}
	return &Info{UserName: opts.Username, Password: opts.Password, Phone: opts.Phone, Token: token, Verify: opts.Verify, ClientIP: ClientIP}
}

func (u *Info) Registry() error {
	if u.UserName == "" || u.Password == "" || u.Phone == "" || u.Verify == "" {
		return errors.New("参数错误")
	}
	var tmpKey = sms.GetResultPreFix + u.Phone
	ok, val := queue.MClient.Exists(tmpKey)
	if ok == false {
		return errors.New("验证失败")
	}
	result, err := sms.GetSmsMsg(val)
	if err != nil {
		return err
	}
	//if result.Error != "" {
	//	return errors.New(result.Error)
	//}
	if result.Captcha != u.Verify {
		return errors.New("验证码错误！")
	}
	db, err := mysql.Get()
	if err != nil {
		return errors.New("数据库连接失败！")
	}
	defer db.Close()
	var count int
	check := db.DB.Model(model.User{}).Where(model.User{Username: u.UserName}).Count(&count)
	if check.Error != nil {
		return check.Error
	}
	if count > 0 {
		return errors.New("用户名已被使用！")
	}
	create := db.DB.Model(model.User{}).Create(&model.User{Username: u.UserName, Password: util2.Md5(u.Password), Phone: u.Phone})
	if create.Error != nil {
		return create.Error
	}
	if create.RowsAffected < 1 {
		return errors.New("创建用户失败！")
	}
	return nil
}

//func (u *Info) SendCode() {
//	//log.Println(u.Iphone,"发送验证码成功！")
//}

func (u *Info) UserLogin() (token string, err error) {
	if u.UserName == "" || u.Password == "" {
		return token, errors.New("参数错误")
	}
	db, err := mysql.Get()
	if err != nil {
		return token, errors.New("数据库连接失败！")
	}
	var userModel model.User
	result := db.DB.Model(model.User{}).Where(model.User{Username: u.UserName, Password: util2.Md5(u.Password)}).First(&userModel)
	if result.RowsAffected != 1 {
		return token, errors.New("用户名或密码错误")
	}
	token, err = util2.ReleaseToken(userModel.ID, u.ClientIP)
	if err != nil {
		return token, errors.New("token签发错误！")
	}
	return token, nil
}

func (u *Info) UpdateToken() (newToken string, err error) {
	if u.Token == "" {
		return newToken, errors.New("token不存在")
	}
	token, err := util2.ParseToken(u.Token)
	if err != nil {
		return newToken, errors.New("token错误")
	}
	if err := token.Valid(); err != nil {
		return newToken, errors.New("token已失效")
	}
	newToken, err = util2.ReleaseToken(token.UserId, token.Ip)
	if err != nil {
		return newToken, errors.New("token签发失败")
	}
	u.Logout()	//使当前Token失效
	return newToken, nil
}

func (u *Info) Logout() error {
	if u.Token == "" {
		return errors.New("token不存在")
	}
	token, err := util2.ParseToken(u.Token)
	if err != nil {
		return errors.New("token错误")
	}
	if err := token.Valid(); err != nil {
		return errors.New("token已失效")
	}
	logoutTime := token.ExpiresAt - time.Now().Unix()
	return queue.MClient.SetEx(user.LogOutPreFix+util2.Md5(u.Token), "", logoutTime)
}

func (u *Info) CheckPhone() error {
	if u.Phone == "" || !util2.VerifyMobileFormat(u.Phone) {
		return errors.New("手机号格式错误！")
	}
	if ok, _ := queue.MClient.Exists(sms.SendMsgPreFix + u.Phone); ok { //手机号发送限制验证
		return errors.New("请勿频繁发送验证码！")
	}
	if ok, _ := queue.MClient.Exists(sms.IpVerifyPreFix + u.ClientIP); ok { //客户端IP发送限制验证
		return errors.New("请勿频繁发送验证码！")
	}
	db, err := mysql.Get()
	if err != nil {
		return errors.New("数据库连接失败！")
	}
	defer db.Close()
	var count int
	result := db.DB.Model(&model.User{}).Where(model.User{Phone: u.Phone}).Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return errors.New("该手机号已被注册！")
	}
	send, err := sms.NewSend(u.Phone, u.ClientIP)
	if err != nil {
		return errors.New("验证码生成失败！")
	}
	_, err = queue.MClient.Push(sms.SendMsgPreFix, u.Phone, string(send))
	if err != nil {
		return errors.New("验证码发送失败！")
	}
	return nil
}

func UserHandler(opt *gateway.Options, caller gateway.CallerType, rsp *gateway.CallRsp, ctx context.Context) {
	var clientIP string
	var token string
	if ctx.Value("clientIp") != nil {
		clientIP = ctx.Value("clientIp").(string)
	}
	if ctx.Value("token") != nil {
		token = ctx.Value("token").(string)
	}
	userInfo := newInfo(opt, clientIP, token)
	if caller == gateway.CallerType_Register {
		err := userInfo.Registry()
		if err != nil {
			rsp.Msg = err.Error()
			//rsp.Data = err.Error()
			return
		}
		rsp.Msg = "注册成功"
	} else if caller == gateway.CallerType_UserLogin {
		token, err := userInfo.UserLogin()
		if err != nil {
			rsp.Msg = err.Error()
			//rsp.Data = err.Error()
			return
		}
		rsp.Msg = "登录成功"
		rsp.Token = token
	} else if caller == gateway.CallerType_UpdateToken {
		newToken,err := userInfo.UpdateToken()
		if err != nil {
			rsp.Msg = err.Error()
			//rsp.Data = err.Error()
			return
		}
		rsp.Token = newToken
		rsp.Msg = "签发成功"
	} else if caller == gateway.CallerType_UserLogout {
		err := userInfo.Logout()
		if err != nil {
			rsp.Msg = err.Error()
			//rsp.Data = err.Error()
			return
		}
		rsp.Msg = "退出成功"
	} else if caller == gateway.CallerType_SendCode {
		if err := userInfo.CheckPhone(); err != nil {
			rsp.Msg = err.Error()
			//rsp.Data = err.Error()
			return
		}
		rsp.Msg = "发送成功"
	}
	rsp.Code = 1
	rsp.Status = true
}
