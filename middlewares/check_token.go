/**
 * @Author: wenliangzhang
 * @Description:
 * @File: CheckToken
 * @Version: 1.0.0
 * @Date: 2022/5/9 5:33 PM
 */
package middlewares

import (
	"errors"

	"git.singularity-ai.com/backend/library/arch/web"
	"git.singularity-ai.com/backend/library/service/token"
)

var (
	ErrorServerTokenError   = errors.New("token异常, 请下线")
	ErrorServerTokenExpire  = errors.New("token过期， 请重新登录")
	ErrorServerTokenInvalid = errors.New("token无效")
)

var Error2CodeIns = map[error]int{
	ErrorServerTokenError:   10001,
	ErrorServerTokenExpire:  10002,
	ErrorServerTokenInvalid: 10003,
}

func CheckToken() web.WebHandlerFunc {

	return func(ctx *web.WebContext) {
		// 前置校验
		t := ctx.GetHeader("Token")
		if t == "" {
			return
		}
		data, err := token.VerifyToken(ctx, t)
		if err != nil {
			ctx.Errorf("read token data failed %s,token=%s", err.Error(), t)
			return
		}
		ctx.Set("userID", data.UserId)
		ctx.Set("robotID", data.RobotId)

		ctx.Next()
	}
}
