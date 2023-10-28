package xlog

import "sync"

// default
var (
	loggerDefault *Logger
	takeUp        = false
)

func GetLogger() *Logger {
	return loggerDefault
}

func SetLevel(lvl int) {
	loggerDefault.level = lvl
}

func SetLayout(layout string) {
	loggerDefault.layout = layout
}

func TraceKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	loggerDefault.deliverKVRecordToWriter(TRACE, specialKV, kvFields)
}

func DebugKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	loggerDefault.deliverKVRecordToWriter(DEBUG, specialKV, kvFields)
}

func WarnKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	loggerDefault.deliverKVRecordToWriter(WARNING, specialKV, kvFields)

}

func InfoKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	loggerDefault.deliverKVRecordToWriter(INFO, specialKV, kvFields)

}

func ErrorKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	loggerDefault.deliverKVRecordToWriter(ERROR, specialKV, kvFields)

}

func FatalKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	loggerDefault.deliverKVRecordToWriter(FATAL, specialKV, kvFields)

}

func PublicKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	loggerDefault.deliverKVRecordToWriter(PUBLIC, specialKV, kvFields)

}

func Tracef(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(TRACE, nil, fmt, args...)
}

func Debugf(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(DEBUG, nil, fmt, args...)
}

func Warnf(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(WARNING, nil, fmt, args...)
}

func Infof(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(INFO, nil, fmt, args...)
}

func Errorf(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(ERROR, nil, fmt, args...)
}

func Fatalf(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(FATAL, nil, fmt, args...)
}

func Publicf(fmt string, args ...interface{}) {
	loggerDefault.deliverRecordToWriter(PUBLIC, nil, fmt, args...)
}

func Trace(args ...interface{}) {
	loggerDefault.deliverRecordToWriter(TRACE, nil, "", args...)
}

func Debug(args ...interface{}) {
	loggerDefault.deliverRecordToWriter(DEBUG, nil, "", args...)
}

func Warn(args ...interface{}) {
	loggerDefault.deliverRecordToWriter(WARNING, nil, "", args...)
}

func Info(args ...interface{}) {
	loggerDefault.deliverRecordToWriter(INFO, nil, "", args...)
}

func Error(args ...interface{}) {
	loggerDefault.deliverRecordToWriter(ERROR, nil, "", args...)
}

func Fatal(args ...interface{}) {
	loggerDefault.deliverRecordToWriter(FATAL, nil, "", args...)
}

func Public(args ...interface{}) {
	loggerDefault.deliverRecordToWriter(PUBLIC, nil, "", args...)
}

func Register(w Writer) {
	loggerDefault.Register(w)
}

func Close() {
	loggerDefault.Close()
}

func init() {
	loggerDefault = NewLogger()
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
}
