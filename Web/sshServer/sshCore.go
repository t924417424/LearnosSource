package sshServer

import (
	"Learnos/common/microClient"
	"Learnos/common/microClient/ServiceCall"
	gateway "Learnos/proto/gateway"
	"context"
	"github.com/gliderlabs/ssh"
	"github.com/micro/go-micro/metadata"
	"net"
	"strings"
)

func Run(addr string) error {
	ssh.Handle(console)
	return ssh.ListenAndServe(addr, nil, ssh.PasswordAuth(passwordVerify))
}

func passwordVerify(ctx ssh.Context, password string) bool {
	ctx.SetValue("ip",getIp(ctx.RemoteAddr().String()))
	var gRsp gateway.CallRsp
	opt := &gateway.Call{
		Type:   gateway.CallType_User,
		Caller: gateway.CallerType_UserLogin,
		Opt: &gateway.Options{
			User: &gateway.UserOpt{
				Username: ctx.User(),
				Password: password,
			},
		},
	}
	mCtx := metadata.NewContext(context.Background(), map[string]string{"Source-Ip": ctx.Value("ip").(string)})
	client := microClient.Get()
	defer client.Close()
	if err := client.Call(mCtx, client.NewRequest(ServiceCall.GateWayServer, ServiceCall.GateWayService, opt), &gRsp); err != nil {
		return false
	} else {
		if gRsp.Status == true {
			ctx.SetValue("token","Bearer " + gRsp.Token)	//登陆成功，保存token到上下文
			return true
		}
	}
	return false
}

func getIp(addr string) (ip string) {
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(addr)); err == nil {
		return ip
	}
	return
}
