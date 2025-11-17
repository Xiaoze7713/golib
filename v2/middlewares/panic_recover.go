/**
* Author: xumengnan
* Description:
* File: panic_recover
* Date: 2021/11/12 5:40 下午
 */
package middlewares

import (
	"fmt"
	"os"
	"runtime"

	"github.com/golib/v2/arch/web"
	"github.com/golib/v2/common"
	"github.com/golib/v2/env"
	"github.com/golib/v2/service/robot"
)

const (
	RobotUrl = "https://open.feishu.cn/open-apis/bot/v2/hook/e08c2085-c0b2-4aa3-9df6-5e6068ac4cfa"
)

func PanicRecover() web.WebHandlerFunc {
	return func(ctx *web.WebContext) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				stack := make([]byte, 2<<10)
				length := runtime.Stack(stack, false)
				hostname, _ := os.Hostname()
				homeDir, _ := os.UserHomeDir()
				curDir, _ := os.Getwd()
				ip := env.LocalIP()
				realErr := fmt.Sprintf("err:%s \nstack:%s \nhostname:%s \nhomeDir:%s \ncurDir:%s \n", err.Error(), stack[:length], hostname, homeDir, curDir)
				ctx.Fatal(realErr)
				// debug模式 下打印错误内容到终端，线上模式 下问题发送飞书
				if env.IsDebug() {
					fmt.Println(realErr)
				} else {
					realErr = fmt.Sprintf("#### 错误信息:\n%s \n#### traceID:\n%s \n#### 堆栈信息:\n%s \n#### Hostname:\n%s \n#### IP:\n%s \nHomeDir:%s \nCurDir:%s \n", err.Error(), ctx.TraceID().String(), stack[:length/2], hostname, ip, homeDir, curDir)
					request := [][]map[string]string{
						{
							{
								"tag":  "text",
								"text": realErr,
							},
							{
								"tag":     "at",
								"user_id": "all",
							},
						},
					}
					err = robot.SendMsg(ctx, RobotUrl, "Panic报警", request)
					if err != nil {
						ctx.Fatal(err)
					}
				}
				common.Failed(ctx, common.ErrUnknown, "内部服务异常！", nil)
			}
		}()
		ctx.Next()
	}
}
