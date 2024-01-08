/**
 * @Author: wenliangzhang
 * @Description:
 * @File: weblog
 * @Version: 1.0.0
 * @Date: 2022/4/21 6:51 PM
 */
package middlewares

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Xiaoze7713/golib/v3/arch/web"
	"github.com/Xiaoze7713/golib/v3/env"
)

func WebLogger() web.WebHandlerFunc {

	return func(c *web.WebContext) {

		// 请求路由
		reqUri := c.Request.RequestURI
		if reqUri == "/metrics" {
			return
		}
		// 开始时间
		start := time.Now()
		startTime := start.Format("2006-01-02 15:04:05")

		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()

		// 执行时间
		cost := end.Sub(start)

		// 请求方式
		reqMethod := c.Request.Method

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 本地IP
		localIP := env.LocalIP()

		// hostname
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}

		// 服务器所处机房
		idc := env.IDC()

		header := c.Request.Header

		// cookie
		cookies := c.Request.Cookies()
		vs := make([]string, 0, len(cookies))
		for _, c := range cookies {
			vs = append(vs, c.String())
		}
		cookie := strings.Join(vs, ";")

		// get|post 数据
		form := c.Request.Form.Encode()

		// errno
		errno, isExit := c.Get("errno")
		if !isExit {
			errno = 0
		}

		// errmsg
		errmsg, isExit := c.Get("errmsg")
		if !isExit {
			errmsg = ""
		}

		// 日志格式
		msg := fmt.Sprintf("sTime[%s] cost[%s] method[%s] uri[%s] code[%d] clientip[%s] localip[%s] hostname[%s] idc[%s] header[%v] cookie[%s] body[%s] form[%s] errno[%d] errmsg[%s]",
			startTime, cost, reqMethod, reqUri, statusCode, clientIP, localIP, hostname, idc, header, cookie, string(body), form, errno.(int), errmsg)

		if statusCode > 499 {
			c.Error(msg)
		} else if statusCode > 399 {
			c.Warn(msg)
		} else if errno != 0 && errno != 200 {
			c.Warn(msg)
		} else {
			c.Info(msg)
		}
		// c.Writer.Header().Set("X_TRACE_ID", c.SerializeSpanContext())
		c.SetTag("idc", env.IDC())
		c.SetTag("header", header)
		c.SetTag("req", string(body))
		c.Fin()
	}
}
