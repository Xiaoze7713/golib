package xlog_base

var xLogger LoggerIF

type LoggerIF interface {
	Debug(specialKV map[string]interface{}, args ...interface{})
	Info(specialKV map[string]interface{}, args ...interface{})
	Warn(specialKV map[string]interface{}, args ...interface{})
	Error(specialKV map[string]interface{}, args ...interface{})
	Fatal(specialKV map[string]interface{}, args ...interface{})
	Debugf(specialKV map[string]interface{}, fmt string, args ...interface{})
	Infof(specialKV map[string]interface{}, fmt string, args ...interface{})
	Warnf(specialKV map[string]interface{}, fmt string, args ...interface{})
	Errorf(specialKV map[string]interface{}, fmt string, args ...interface{})
	Fatalf(specialKV map[string]interface{}, fmt string, args ...interface{})
	DebugKV(specialKV map[string]interface{}, KV map[string]interface{})
	InfoKV(specialKV map[string]interface{}, KV map[string]interface{})
	WarnKV(specialKV map[string]interface{}, KV map[string]interface{})
	ErrorKV(specialKV map[string]interface{}, KV map[string]interface{})
	FatalKV(specialKV map[string]interface{}, KV map[string]interface{})
}
