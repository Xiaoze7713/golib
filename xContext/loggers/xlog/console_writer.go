package xlog

import (
	"fmt"
	"os"
)

type colorRecord Record

func (r *colorRecord) String() string {
	content := r.content
	if r.spanID != "" {
		content = fmt.Sprintf("[span_id:%v]", r.spanID) + content
	}
	if r.traceID != "" {
		content = fmt.Sprintf("[trace_id:%v]", r.traceID) + content
	}
	//return fmt.Sprintf("[%s][%s][%s] %s\n", LevelFlags[r.level], r.time, r.code, r.content)
	switch r.level {
	case TRACE:
		return fmt.Sprintf("\033[96m%s\033[0m [\033[96m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LevelFlags[r.level], r.code, content)
	case DEBUG:
		return fmt.Sprintf("\033[96m%s\033[0m [\033[94m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LevelFlags[r.level], r.code, content)

	case INFO:
		return fmt.Sprintf("\033[96m%s\033[0m [\033[92m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LevelFlags[r.level], r.code, content)

	case WARNING:
		return fmt.Sprintf("\033[96m%s\033[0m [\033[103m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LevelFlags[r.level], r.code, content)

	case ERROR:
		return fmt.Sprintf("\033[96m%s\033[0m [\033[101m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LevelFlags[r.level], r.code, content)

	case FATAL:
		return fmt.Sprintf("\033[96m%s\033[0m [\033[105m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LevelFlags[r.level], r.code, content)
	case PUBLIC:
		return fmt.Sprintf("\033[96m%s\033[0m [\033[96m%s\033[0m] \033[47;30m%s\033[0m \033[36m%s\033[0m\n",
			r.time, LevelFlags[r.level], r.code, content)
	}
	return ""
}

type ConsoleWriter struct {
	color bool
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

func (w *ConsoleWriter) Write(r *Record) error {
	if w.color {
		fmt.Fprint(os.Stdout, ((*colorRecord)(r)).String())
	} else {
		fmt.Fprint(os.Stdout, r.String())
	}
	return nil
}

func (w *ConsoleWriter) Init() error {
	return nil
}

func (w *ConsoleWriter) SetColor(c bool) {
	w.color = c
}
