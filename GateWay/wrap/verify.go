package wrap

import (
	"Learnos/GateWay/sqldata/model"
	"Learnos/GateWay/sqldata/mysql"
	"Learnos/common/queue"
	"Learnos/common/queueMsg/gateway/user"
	"Learnos/common/util"
	"context"
	"errors"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"log"
	"strings"
)

func Verify(next server.HandlerFunc) server.HandlerFunc {
	log.SetFlags(log.Lshortfile)
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		data, _ := metadata.FromContext(ctx)
		clientIp := data["Source-Ip"]
		Authorization := data["Authorization"]
		if clientIp == "" || !util.CheckIp(clientIp) {
			return errors.New("客户端IP错误！")
		}
		ctx = context.WithValue(ctx, "clientIp", clientIp)
		if Authorization != "" { //如果Token存在，则开始校验Token
			if !strings.HasPrefix(Authorization, "Bearer ") {
				return errors.New("token格式错误")
			}
			token := Authorization[7:]
			c, err := util.ParseToken(token)
			if err != nil {
				return err
			}
			if err := c.Valid(); err != nil {
				return errors.New("token失效")
			}
			if exists, _ := queue.MClient.Exists(user.LogOutPreFix + util.Md5(token)); exists {
				return errors.New("客户端已退出")
			}
			if c.Ip != clientIp {
				return errors.New("权限验证失败！")
			}
			db,err := mysql.Get()
			if err != nil{
				return err
			}
			defer db.Close()
			var userInfo model.User
			userInfo.ID = c.UserId
			db.DB.Where(userInfo).First(&userInfo)
			ctx = context.WithValue(ctx, "token", token)
			ctx = context.WithValue(ctx, "userInfo", userInfo)
		}
		//Dlog.Println(req.Service(), "调用成功")
		err := next(ctx, req, rsp)
		return err
	}
}
