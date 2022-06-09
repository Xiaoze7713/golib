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
	"git.singularity-ai.com/backend/library/log"
	"git.singularity-ai.com/backend/library/service/token"
	"github.com/gin-gonic/gin"
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

func CheckToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		//前置校验
		t := c.GetHeader("Token")
		if t == "" {
			return
		}
		data, err := token.VerifyToken(t)
		if err != nil {
			log.Errorf("read token data failed %s,token=%s", err.Error(), t)
			return
		}
		c.Set("userID", data.UserId)
		c.Set("robotID", data.RobotId)

		c.Next()
	}
}
