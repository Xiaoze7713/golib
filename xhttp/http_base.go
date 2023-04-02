package xhttp

import (
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Request struct {
	Data interface{} `json:"data,omitempty"`
}

func GinMustBind(gc *gin.Context, binding binding.Binding, data interface{}) (err error) {
	req := &Request{data}
	err = gc.ShouldBindWith(req, binding)
	return
}

type Response struct {
	Code    int32       `json:"code"`
	Message string      `json:"code_msg"`
	Data    interface{} `json:"resp_data,omitempty"`
}

func SetResponse(g *gin.Context, err error, data interface{}) {
	errorMsg := ""
	if err != nil {
		xlog.Error(err)
		errorMsg = ErrMsg(err)
	}
	resp := Response{
		Code:    ErrorCode(err),
		Message: errorMsg,
		Data:    data,
	}
	g.Set("resp", resp)
}

func SetErrorResponse(g *gin.Context, code int32, msg string) {
	resp := Response{
		Code:    code,
		Message: msg,
	}
	g.Set("resp", resp)
}
