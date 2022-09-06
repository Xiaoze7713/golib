package http_base

import (
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type request struct {
	Data interface{} `json:"data"`
}

func GinMustBind(gc *gin.Context, binding binding.Binding, data interface{}) (err error) {
	req := request{data}
	err = gc.MustBindWith(req, binding)
	return
}

type response struct {
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
	resp := response{
		Code:    ErrorCode(err),
		Message: errorMsg,
		Data:    data,
	}
	g.Set("resp", resp)
}

func SetErrorResponse(g *gin.Context, code int32, msg string) {
	resp := response{
		Code:    code,
		Message: msg,
	}
	g.Set("resp", resp)
}
