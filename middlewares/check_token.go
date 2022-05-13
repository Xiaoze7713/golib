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
	"git.singularity-ai.com/backend/library/service"
	"github.com/gin-gonic/gin"
)

func CheckToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		//前置校验
		token := c.GetHeader("Token")
		if token == "" {
			log.Info("no token")
			return
		}
		data, err := service.VerifyToken(token)
		if err != nil {
			log.Errorf("read token data failed %s,token=%s", err.Error(), token)
			return
		}
		c.Set("userID", data.UserId)
		c.Set("robotID", data.RobotId)

		c.Next()
	}
}
