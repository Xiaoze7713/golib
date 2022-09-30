/**
 * @Author: wenliangzhang
 * @Description:
 * @File: response
 * @Version: 1.0.0
 * @Date: 2022/4/12 5:42 PM
 */
package common

import (
	"encoding/json"
	"net/http"

	"git.singularity-ai.com/backend/library/arch/web"
)

type Resp struct {
	TraceID  string      `json:"trace_id"`
	Code     int         `json:"code"`
	CodeMsg  string      `json:"code_msg"`
	RespData interface{} `json:"resp_data"`
}

// Success 通用格式化成功返回
func Success(ctx *web.WebContext, data interface{}) {
	Common(ctx, http.StatusOK, 200, "success", data)
}

// Failed 通用格式化错误返回
func Failed(ctx *web.WebContext, errno int, errmsg string, data interface{}) {
	Common(ctx, http.StatusOK, errno, errmsg, data)
}

// Common 通用格式化错误返回
func Common(ctx *web.WebContext, httpCode int, errno int, errmsg string, data interface{}) {
	traceId := ctx.SpanID().String()
	resp := Resp{
		traceId,
		errno,
		errmsg,
		data,
	}
	if errno != 0 {
		ctx.Set("errno", errno)
		ctx.Set("errmsg", errmsg)
	}
	respBody, _ := json.Marshal(resp)
	ctx.Infof("response :%s", string(respBody))
	ctx.JSON(httpCode, resp)
}
