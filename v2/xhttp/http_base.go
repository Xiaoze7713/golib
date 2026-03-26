package xhttp

import (
	"bytes"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golib/v2/xContext/loggers/xlog"
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
	Code     int32       `json:"code"`
	Message  string      `json:"code_msg"`
	Reason   string      `json:"code_reason,omitempty"`
	RespData interface{} `json:"resp_data,omitempty"`
}

type Error interface {
	error
	Reason() string
}

type ReasonError struct {
	ReasonStr string `json:"reason"`
	error
}

func (e ReasonError) Reason() string {
	return e.ReasonStr
}

func NewError(reason string, err error) ReasonError {
	return ReasonError{reason, err}
}

func ConvertError(err error) ReasonError {
	newErr, ok := err.(ReasonError)
	if ok {
		return newErr
	}
	return ReasonError{err.Error(), err}
}

func NewResponse(err error) Response {
	errorMsg := ""
	if err != nil {
		xlog.Error(err)
		errorMsg = ErrMsg(err)
	}
	resp := Response{
		Code:    ErrorCode(err),
		Message: errorMsg,
	}
	return resp
}

func NewReasonResponse(err Error) Response {
	resp := NewResponse(err)
	if err != nil {
		var rErr ReasonError
		ok := errors.As(err, &rErr)
		if ok {
			resp.Reason = err.Reason()
		} else {
			resp.Reason = err.Error()
		}
	}
	return resp
}

func SetResponseWithReason(g *gin.Context, err Error, data interface{}) {
	resp := NewReasonResponse(err)
	resp.RespData = data
	g.Set(RespJson, resp)
}

func SetHttpCodeResponseWithReason(g *gin.Context, httpCode int, err Error, data interface{}) {
	SetResponseWithReason(g, err, data)
	g.Set(HttpResponseCode, httpCode)
}

func SetResponse(g *gin.Context, err error, data interface{}) {
	resp := NewResponse(err)
	if err != nil {
		resp.Reason = err.Error()
	}
	if data != nil {
		resp.RespData = data
	} else {
		resp.RespData = make(map[string]interface{})
	}
	g.Set(RespJson, resp)
}

func SetHttpBodyResponse(g *gin.Context, httpCode int, data interface{}) {
	g.Set(RespJson, data)
	g.Set(HttpResponseCode, httpCode)
}

func SetBytesResponse(g *gin.Context, data []byte) {
	g.Set(RespBytes, data)
}

func SetHttpBytesResponse(g *gin.Context, httpCode int, data []byte) {
	g.Set(RespBytes, data)
	g.Set(HttpResponseCode, httpCode)
}

func SetErrorResponse(g *gin.Context, code int32, msg string) {
	resp := Response{
		Code:    code,
		Message: msg,
	}
	g.Set(RespJson, resp)
}

func WriteStreamBytes(g *gin.Context, data []byte, addLine bool) (err error) {
	_, err = g.Writer.Write(data)
	if err != nil {
		return err
	}
	if addLine {
		_, err2 := g.Writer.WriteString("\n\n")
		if err2 != nil {
			return err2
		}
	}
	g.Writer.Flush()
	stream, ok := g.Get(RespStream)
	if ok {
		streamRespList, ok1 := stream.([]string)
		if ok1 {
			streamRespList = append(streamRespList, string(data))
		} else {
			streamRespList = append([]string{}, string(data))
		}
		g.Set(RespStream, streamRespList)
	}
	if !ok {
		streamRespList := append([]string{}, string(data))
		g.Set(RespStream, streamRespList)
	}
	return nil
}

const DataPrefix = "data: "

func WriteStreamDataBlock(g *gin.Context, data []byte, addLine bool) (err error) {
	mergeData := bytes.NewBuffer([]byte(DataPrefix))
	mergeData.Write(data)
	err = WriteStreamBytes(g, mergeData.Bytes(), addLine)
	if err != nil {
		return err
	}
	return nil
}
