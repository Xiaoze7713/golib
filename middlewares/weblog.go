/**
 * @Author: wenliangzhang
 * @Description:
 * @File: weblog
 * @Version: 1.0.0
 * @Date: 2022/4/21 6:51 PM
 */
package middlewares

import (
	"fmt"
	"git.singularity-ai.com/backend/library/common"
	"git.singularity-ai.com/backend/library/env"
	"git.singularity-ai.com/backend/library/log"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
	"time"
)

func WebLogger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		startTime := start.Format("2006-01-02 15:04:05")

		// traceID
		traceID := GenLogIDFromRequest(c)
		c.Set("traceID", traceID)

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()

		// 执行时间
		cost := end.Sub(start)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 本地IP
		localIP := env.LocalIP()

		//hostname
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}

		// 服务器所处机房
		idc := env.IDC()

		//cookie
		cookies := c.Request.Cookies()
		vs := make([]string, 0, len(cookies))
		for _, c := range cookies {
			vs = append(vs, c.String())
		}
		cookie := strings.Join(vs, ";")

		// get|post 数据
		form := c.Request.Form.Encode()

		//errno
		errno, isExit := c.Get("errno")
		if !isExit {
			errno = 0
		}

		//errmsg
		errmsg, isExit := c.Get("errmsg")
		if !isExit {
			errmsg = ""
		}

		//日志格式
		msg := fmt.Sprintf("traceid[%s] sTime[%s] cost[%s] method[%s] uri[%s] code[%d] clientip[%s] localip[%s] hostname[%s] idc[%s] cookie[%s] form[%s] errno[%d] errmsg[%s]",
			traceID, startTime, cost, reqMethod, reqUri, statusCode, clientIP, localIP, hostname, idc, cookie, form, errno.(int), errmsg)

		if statusCode > 499 {
			log.Error(msg)
		} else if statusCode > 399 {
			log.Warn(msg)
		} else if errno != 0 && errno != 200 {
			log.Warn(msg)
		} else {
			log.Info(msg)
		}

	}
}

// GenLogIDFromRequest 生成日志ID
// 优先级 header中X_TRACEID > header中traceid > 表单中traceid > 自己生成
func GenLogIDFromRequest(c *gin.Context) string {
	request := c.Request
	form := request.URL.Query()
	// 上下游header透传时
	if traceID := strings.TrimSpace(request.Header.Get("X_TRACEID")); traceID != "" {
		return traceID
	}
	if traceID := strings.TrimSpace(request.Header.Get("traceid")); traceID != "" {
		return traceID
	}
	// 上下游querystring透传时
	if traceID := strings.TrimSpace(form.Get("traceid")); traceID != "" {
		return traceID
	}

	return strconv.Itoa(common.GetSFInstance().GetUniqueId())
}
