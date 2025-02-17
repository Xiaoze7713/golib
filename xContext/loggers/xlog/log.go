package xlog

import (
	"fmt"
	"github.com/golib/v3/utils"
	jsoniter "github.com/json-iterator/go"
	"log"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	LevelFlags = [...]string{"TRACE", "DEBUG", "INFO", "PUBLIC", "WARN", "ERROR", "FATAL"}
	recordPool *sync.Pool
)

const (
	TRACE = iota
	DEBUG
	INFO
	PUBLIC
	WARNING
	ERROR
	FATAL
)

const tunnelSizeDefault = 1024

const (
	FmtTypeRaw FmtType = iota
	FmtTypeJson
	FmtTypeColorRaw
)

type Record struct {
	time      string
	code      string
	traceID   interface{}
	spanID    interface{}
	specialKV map[string]interface{}
	kvMap     map[string]interface{}
	argFields []interface{}
	content   string
	level     int
	FmtType   FmtType
}

func (m *Record) Flush() {
	m.time = ""
	m.code = ""
	m.traceID = ""
	m.spanID = ""
	m.specialKV = nil
	m.kvMap = nil
	m.argFields = nil
	m.content = ""
	m.level = 0
	m.FmtType = 0
}

type JsonFmt struct {
	Time      string                 `json:"time"`
	Code      string                 `json:"code"`
	TraceID   interface{}            `json:"trace_id"`
	SpanID    interface{}            `json:"span_id"`
	Level     string                 `json:"level"`
	SpecialKV map[string]interface{} `json:"special_kv"`
	KvMap     map[string]interface{} `json:"kv_map"`
}

func (r *Record) String() string {
	// todo
	switch r.FmtType {
	case FmtTypeRaw:
		return r.RawString()
	case FmtTypeJson:
		return r.JsonString()
	case FmtTypeColorRaw:
		return r.ColorString()
	default:
		return r.JsonString()
	}
}

func (r *Record) RawString() string {
	content := r.content
	if r.kvMap != nil {
		content = fmt.Sprintf("%v", utils.MustJson(r.kvMap))
	}
	if r.spanID != nil && r.spanID != "" {
		content = fmt.Sprintf("[span_id:%v]", r.spanID) + content
	}
	if r.traceID != nil && r.traceID != "" {
		content = fmt.Sprintf("[trace_id:%v]", r.traceID) + content
	}
	return fmt.Sprintf("[%s][%s][%s] %s\n", r.time, LevelFlags[r.level], r.code, content)
}

func (r *Record) JsonString() string {
	// todo
	jf := JsonFmt{}
	jf.Time = r.time
	jf.Level = strings.ToLower(LevelFlags[r.level])
	jf.Code = r.code
	jf.TraceID = r.traceID
	jf.SpanID = r.spanID
	jf.SpecialKV = r.specialKV
	if r.kvMap != nil { // kv模式
		jf.KvMap = r.kvMap
	} else if r.content != "" {
		jf.KvMap = map[string]interface{}{"content": r.content}
	} else if len(r.argFields) != 0 {
		jf.KvMap = map[string]interface{}{}
		argFields := r.argFields
		if len(r.argFields) >= 1000 {
			argFields = argFields[:1000]
		}
		for i, arg := range argFields {
			jf.KvMap[fmt.Sprintf("f%04d", i)] = arg
		}
	}
	msg, _ := jsoniter.MarshalToString(jf)
	return fmt.Sprintf("%s\n", msg)
}

func (r *Record) ColorString() string {
	content := r.content
	if r.kvMap != nil {
		content = fmt.Sprintf("[kv:%+v]", r.kvMap) + content
	}
	if r.spanID != nil {
		content = fmt.Sprintf("[span_id:%v]", r.spanID) + content
	}
	if r.traceID != nil {
		content = fmt.Sprintf("[trace_id:%v]", r.traceID) + content
	}
	//return fmt.Sprintf("[%s][%s][%s] %s\n", LevelFlags[r.level], r.time, r.code, r.content)
	// time
	prefix := fmt.Sprintf("\033[0m[\033[96m%s\033[0m]", r.time)
	// level
	switch r.level {
	case TRACE:
		prefix += fmt.Sprintf("\033[0m[\033[96m%s\033[0m]", LevelFlags[r.level])
	case DEBUG:
		prefix += fmt.Sprintf("\033[0m[\033[94m%s\033[0m]", LevelFlags[r.level])
	case INFO:
		prefix += fmt.Sprintf("\033[0m[\033[92m%s\033[0m]", LevelFlags[r.level])
	case WARNING:
		prefix += fmt.Sprintf("\033[0m[\033[103m%s\033[0m]", LevelFlags[r.level])
	case ERROR:
		prefix += fmt.Sprintf("\033[0m[\033[101m%s\033[0m]", LevelFlags[r.level])
	case FATAL:
		prefix += fmt.Sprintf("\033[0m[\033[105m%s\033[0m]", LevelFlags[r.level])
	case PUBLIC:
		prefix += fmt.Sprintf("\033[0m[\033[96m%s\033[0m]", LevelFlags[r.level])
	}
	// code
	prefix += fmt.Sprintf("\033[0m[\033[47;30m%s\033[0m] ", r.code)
	content = prefix + content + "\n"
	return content
}

type FmtType int

type Writer interface {
	Init() error
	Write(*Record) error
	Fmt() FmtType
}

type Rotated interface {
	Rotate() error
	SetPathPattern(string) error
}

type Flusher interface {
	Flush() error
}

type Logger struct {
	writers     []Writer
	tunnel      chan *Record
	level       int
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
	skipStr     string
	skipStr2    string
	FormatType  FmtType
}

func NewLogger() *Logger {
	if loggerDefault != nil && takeUp == false {
		takeUp = true
		return loggerDefault
	}
	l := new(Logger)
	l.writers = make([]Writer, 0, 2)
	l.tunnel = make(chan *Record, tunnelSizeDefault)
	l.c = make(chan bool, 1)
	l.level = DEBUG
	l.layout = "2006-01-02 15:04:05.999999"
	l.skipStr = "xContextLog.go"

	go boostrapLogWriter(l)

	return l
}

func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

func (l *Logger) SetLevel(lvl int) {
	l.level = lvl
}

func (l *Logger) SetFmtJson() {
	l.FormatType = FmtTypeJson
}

func (l *Logger) SetFmtRaw() {
	l.FormatType = FmtTypeRaw
}

func (l *Logger) SetFmtColor() {
	l.FormatType = FmtTypeColorRaw
}

func (l *Logger) SetFmt(fmtType FmtType) {
	l.FormatType = fmtType
}

func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

func (l *Logger) Tracef(specialKV map[string]interface{}, fmt string, args ...interface{}) {
	l.deliverRecordToWriter(TRACE, specialKV, fmt, args...)
}

func (l *Logger) Debugf(specialKV map[string]interface{}, fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, specialKV, fmt, args...)
}

func (l *Logger) Warnf(specialKV map[string]interface{}, fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, specialKV, fmt, args...)
}

func (l *Logger) Infof(specialKV map[string]interface{}, fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, specialKV, fmt, args...)
}

func (l *Logger) Errorf(specialKV map[string]interface{}, fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, specialKV, fmt, args...)
}

func (l *Logger) Fatalf(specialKV map[string]interface{}, fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, specialKV, fmt, args...)
}

func (l *Logger) Trace(specialKV map[string]interface{}, args ...interface{}) {
	l.deliverRecordToWriter(TRACE, specialKV, "", args...)
}

func (l *Logger) Debug(specialKV map[string]interface{}, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, specialKV, "", args...)
}

func (l *Logger) Warn(specialKV map[string]interface{}, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, specialKV, "", args...)
}

func (l *Logger) Info(specialKV map[string]interface{}, args ...interface{}) {
	l.deliverRecordToWriter(INFO, specialKV, "", args...)
}

func (l *Logger) Error(specialKV map[string]interface{}, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, specialKV, "", args...)
}

func (l *Logger) Fatal(specialKV map[string]interface{}, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, specialKV, "", args...)
}

func (l *Logger) Public(specialKV map[string]interface{}, args ...interface{}) {
	l.deliverRecordToWriter(PUBLIC, specialKV, "", args...)
}

func (l *Logger) TraceKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	l.deliverKVRecordToWriter(TRACE, specialKV, kvFields)
}

func (l *Logger) DebugKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	l.deliverKVRecordToWriter(DEBUG, specialKV, kvFields)

}

func (l *Logger) WarnKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	l.deliverKVRecordToWriter(WARNING, specialKV, kvFields)

}

func (l *Logger) InfoKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	l.deliverKVRecordToWriter(INFO, specialKV, kvFields)

}

func (l *Logger) ErrorKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	l.deliverKVRecordToWriter(ERROR, specialKV, kvFields)

}

func (l *Logger) FatalKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	l.deliverKVRecordToWriter(FATAL, specialKV, kvFields)

}

func (l *Logger) PublicKV(specialKV map[string]interface{}, kvFields map[string]interface{}) {
	l.deliverKVRecordToWriter(PUBLIC, specialKV, kvFields)

}

func (l *Logger) Close() {
	close(l.tunnel)
	<-l.c

	for _, w := range l.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				log.Println(err)
			}
		}
	}
}

func (l *Logger) CodeLine(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	//fmt.Println(file, l.skipStr, l.skipStr2)
	if ok {
		if l.skipStr != "" && l.skipStr == path.Base(file) {
			_, file, line, ok = runtime.Caller(skip + 1)
			if ok {
				if l.skipStr2 != "" && l.skipStr2 == path.Base(file) {
					_, file, line, ok = runtime.Caller(skip + 2)
				}
			}
		}
	}
	if ok {
		return path.Base(file) + ":" + strconv.Itoa(line)
	}
	return ""
}

func (l *Logger) deliverRecordToWriter(level int, specialKV map[string]interface{}, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	// source code, file and line num
	// format time
	now := time.Now()
	if now.UnixNano() != l.lastTime {
		l.lastTime = now.UnixNano()
		l.lastTimeStr = now.Format(l.layout)
	}
	r := recordPool.Get().(*Record)
	r.Flush()
	r.FmtType = l.FormatType
	if spanID, ok := specialKV["span_id"]; ok {
		r.spanID = spanID
	}
	if traceID, ok := specialKV["trace_id"]; ok {
		r.traceID = traceID
	}
	if format != "" {
		r.content = fmt.Sprintf(format, args...)
	} else {
		r.argFields = args
		r.content = fmt.Sprint(args...)
		//inf = fmt.Sprint(args...)
	}

	r.code = l.CodeLine(3)
	//for i := 0; i < 6; i++ {
	//	fmt.Println(l.CodeLine(i))
	//}
	r.time = l.lastTimeStr
	r.level = level
	l.tunnel <- r
}

func (l *Logger) deliverKVRecordToWriter(level int, specialKV map[string]interface{}, kvFields map[string]interface{}) {
	if level < l.level {
		return
	}

	// source code, file and line num
	// format time
	now := time.Now()
	if now.UnixNano() != l.lastTime {
		l.lastTime = now.UnixNano()
		l.lastTimeStr = now.Format(l.layout)
	}
	r := recordPool.Get().(*Record)
	r.Flush()
	r.FmtType = l.FormatType
	r.code = l.CodeLine(3)
	traceID, ok := specialKV["trace_id"]
	if ok {
		r.traceID = traceID
	} else {
		traceID, ok = kvFields["trace_id"]
		if ok && traceID != nil {
			r.traceID = traceID
		}
	}
	spanID, ok := specialKV["span_id"]
	if ok {
		r.spanID = spanID
	} else {
		spanID, ok = kvFields["span_id"]
		if ok && spanID != nil {
			r.spanID = spanID
		}
	}
	//for i := 0; i < 6; i++ {
	//	fmt.Println(l.CodeLine(i))
	//}
	r.time = l.lastTimeStr
	r.level = level
	//fmt.Println(kvFields)
	if len(kvFields) > 0 {
		r.kvMap = kvFields
	}
	l.tunnel <- r
}

func boostrapLogWriter(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	var (
		r  *Record
		ok bool
	)

	if r, ok = <-logger.tunnel; !ok {
		logger.c <- true
		return
	}

	for _, w := range logger.writers {
		if err := w.Write(r); err != nil {
			log.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 200)
	rotateTimer := time.NewTimer(time.Second * 10)

	for {
		select {
		case r, ok = <-logger.tunnel:
			if !ok {
				logger.c <- true
				return
			}

			for _, w := range logger.writers {
				if err := w.Write(r); err != nil {
					log.Println(err)
				}
			}

			recordPool.Put(r)

		case <-flushTimer.C:
			for _, w := range logger.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 200)

		case <-rotateTimer.C:
			for _, w := range logger.writers {
				if r, ok := w.(Rotated); ok {
					if err := r.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
		}
	}
}
