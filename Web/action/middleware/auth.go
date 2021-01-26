package middleware

import (
	"Learnos/Web/formBind"
	"Learnos/common/util"
	"github.com/gin-gonic/gin"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp formBind.Rsp
		jwtToken := c.GetHeader("Authorization")
		if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
			resp.Code = formBind.TokenErr
			resp.Msg = "Token不正确"
			c.JSON(200, resp)
			c.Abort()
			return
		}
		token := jwtToken[7:]
		claims, err := util.ParseToken(token)
		if err != nil {
			resp.Code = formBind.TokenErr
			resp.Msg = "Token错误，请重新登录"
			c.JSON(200, resp)
			c.Abort()
			return
		}
		valid := claims.Valid()
		if valid != nil {
			resp.Code = formBind.TokenExpErr
			resp.Msg = "用户登录超时，请重新登录"
			c.JSON(200, resp)
			c.Abort()
			return
		}
		c.Set("token",jwtToken)
		c.Next()
	}
}

