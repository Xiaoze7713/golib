package xhttp

import "errors"

const (
	ServerCodeStatusOK  = 200
	ServerCodeStatusBad = 9000
)

var (
	ServerStatusOK   = errors.New("Server Status OK")
	ErrorServerError = errors.New("服务器异常") // 通用

	ErrorServerTokenError   = errors.New("token异常, 请下线")
	ErrorServerTokenExpire  = errors.New("token过期， 请重新登录")
	ErrorServerTokenInvalid = errors.New("token无效")

	ErrorServerDataExpire      = errors.New("数据已过期， 请重试")
	ErrorServerDataFormatError = errors.New("数据格式异常")
	ErrorServerDataNotExists   = errors.New("数据不存在")
	ErrorServerDataVerifyError = errors.New("数据核验错误")

	ErrorServerThirdPartyError         = errors.New("第三方服务异常")
	ErrorServerThirdPartyInternalError = errors.New("第三方服务异常")

	ErrorServerReqOperatingFrequency = errors.New("请求过于频繁")
	ErrorServerReqFieldError         = errors.New("请求字段异常")

	ErrorServerActionNotAllow = errors.New("操作受限")

	ErrorSensitive    = errors.New("输入包含敏感字段")
	ErrorLongConnLost = errors.New("用户失去连接")
	ErrorLongConnSend = errors.New("消息发送异常")
	InnerErrorKey     = errors.New("未知的key")
)

type Error2Code struct {
	error2code map[error]int32
}

var error2CodeIns = &Error2Code{error2code: map[error]int32{
	ServerStatusOK:   ServerCodeStatusOK,
	ErrorServerError: ServerCodeStatusBad,

	ErrorServerTokenError:   10001,
	ErrorServerTokenExpire:  10002,
	ErrorServerTokenInvalid: 10003,

	ErrorServerDataExpire:      11001,
	ErrorServerDataFormatError: 11002,
	ErrorServerDataNotExists:   11003,
	ErrorServerDataVerifyError: 11004,

	ErrorServerThirdPartyError: 12001,

	ErrorServerReqOperatingFrequency: 13001,
	ErrorServerReqFieldError:         13002,

	ErrorServerActionNotAllow:          14001,
	ErrorSensitive:                     15001,
	ErrorLongConnLost:                  16001,
	ErrorLongConnSend:                  16002,
	InnerErrorKey:                      20001,
	ErrorServerThirdPartyInternalError: 20002,
}}

func RegisterNewErr(err error, code int32) bool {
	if code == ServerCodeStatusOK {
		return false
	}
	if _, ok := error2CodeIns.error2code[err]; ok {
		return false
	}
	error2CodeIns.error2code[err] = code
	return true
}

func ErrorCode(err error) int32 {
	if err == nil {
		return ServerCodeStatusOK
	}
	code, ok := error2CodeIns.error2code[err]
	if !ok {
		return ServerCodeStatusBad
	}
	return code
}

func ErrMsg(err error) string {
	if err == nil {
		return ""
	}
	_, ok := error2CodeIns.error2code[err]
	if !ok {
		return ErrorServerError.Error()
	}
	return err.Error()
}
