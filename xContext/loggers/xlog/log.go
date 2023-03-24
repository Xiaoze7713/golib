package xlog

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	LEVEL_FLAGS = [...]string{"TRACE", "DEBUG", "INFO", "PUBLIC", "WARN", "ERROR", "FATAL"}
	recordPool  *sync.Pool
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

const tunnel_size_default = 1024

type Record struct {
	time  string
	code  string
	info  string
	level int
}

func (r *Record) String() string {
	return fmt.Sprintf("[%s][%s][%s] %s\n", LEVEL_FLAGS[r.level], r.time, r.code, r.info)
}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Rotater interface {
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
}

func NewLogger() *Logger {
	if logger_default != nil && takeup == false {
		takeup = true
		return logger_default
	}

	l := new(Logger)
	l.writers = make([]Writer, 0, 2)
	l.tunnel = make(chan *Record, tunnel_size_default)
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

func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

func (l *Logger) Tracef(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(TRACE, fmt, args...)
}

func (l *Logger) Debugf(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

func (l *Logger) Warnf(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

func (l *Logger) Infof(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

func (l *Logger) Errorf(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

func (l *Logger) Fatalf(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}

func (l *Logger) Trace(args ...interface{}) {
	l.deliverRecordToWriter(TRACE, "", args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, "", args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.deliverRecordToWriter(WARNING, "", args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.deliverRecordToWriter(INFO, "", args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.deliverRecordToWriter(ERROR, "", args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.deliverRecordToWriter(FATAL, "", args...)
}

func (l *Logger) Public(args ...interface{}) {
	l.deliverRecordToWriter(PUBLIC, "", args...)
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
	if ok {
		return path.Base(file) + ":" + strconv.Itoa(line)
	}
	return ""
}

func (l *Logger) deliverRecordToWriter(level int, format string, args ...interface{}) {
	var inf string

	if level < l.level {
		return
	}

	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}

	// source code, file and line num
	// format time
	now := time.Now()
	if now.UnixNano() != l.lastTime {
		l.lastTime = now.UnixNano()
		l.lastTimeStr = now.Format(l.layout)
	}
	r := recordPool.Get().(*Record)
	r.info = inf
	r.code = l.CodeLine(3)
	if len(l.skipStr) <= len(r.code) && r.code[:len(l.skipStr)] == l.skipStr {
		r.code = l.CodeLine(4)
	}
	//for i := 0; i < 6; i++ {
	//	fmt.Println(l.CodeLine(i))
	//}
	r.time = l.lastTimeStr
	r.level = level

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
				if r, ok := w.(Rotater); ok {
					if err := r.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
		}
	}
}

// default
var (
	logger_default *Logger
	takeup         = false
)

func GetLogger() *Logger {
	return logger_default
}

func SetLevel(lvl int) {
	logger_default.level = lvl
}

func SetLayout(layout string) {
	logger_default.layout = layout
}

func Tracef(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(TRACE, fmt, args...)
}

func Debugf(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(DEBUG, fmt, args...)
}

func Warnf(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(WARNING, fmt, args...)
}

func Infof(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(INFO, fmt, args...)
}

func Errorf(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(ERROR, fmt, args...)
}

func Fatalf(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(FATAL, fmt, args...)
}

func Publicf(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(PUBLIC, fmt, args...)
}

func Trace(args ...interface{}) {
	logger_default.deliverRecordToWriter(TRACE, "", args...)
}

func Debug(args ...interface{}) {
	logger_default.deliverRecordToWriter(DEBUG, "", args...)
}

func Warn(args ...interface{}) {
	logger_default.deliverRecordToWriter(WARNING, "", args...)
}

func Info(args ...interface{}) {
	logger_default.deliverRecordToWriter(INFO, "", args...)
}

func Error(args ...interface{}) {
	logger_default.deliverRecordToWriter(ERROR, "", args...)
}

func Fatal(args ...interface{}) {
	logger_default.deliverRecordToWriter(FATAL, "", args...)
}

func Public(args ...interface{}) {
	logger_default.deliverRecordToWriter(PUBLIC, "", args...)
}

func Register(w Writer) {
	logger_default.Register(w)
}

func Close() {
	logger_default.Close()
}

func init() {
	logger_default = NewLogger()
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
}
