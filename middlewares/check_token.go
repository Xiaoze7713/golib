/**
 * @Author: wenliangzhang
 * @Description:
 * @File: CheckToken
 * @Version: 1.0.0
 * @Date: 2022/5/9 5:33 PM
 */
package middlewares

import (
	"git.singularity-ai.com/backend/library/log"
	"git.singularity-ai.com/backend/library/service/token"
	"github.com/gin-gonic/gin"
)

func CheckToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		//前置校验
		t := c.GetHeader("Token")
		if t == "" {
			log.Info("no token")
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
