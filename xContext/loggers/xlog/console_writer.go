package xlog

import (
	"fmt"
	"os"
)

type colorRecord Record

type ConsoleWriter struct {
	fmt FmtType
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

func (w *ConsoleWriter) Write(r *Record) error {
	data := ""
	switch w.Fmt() {
	case FmtTypeRaw:
		data = r.RawString()
	case FmtTypeJson:
		data = r.JsonString()
	case FmtTypeColorRaw:
		data = r.ColorString()
	default:
		data = r.JsonString()
	}
	_, err := fmt.Fprint(os.Stdout, data)
	return err
}

func (w *ConsoleWriter) Init() error {
	return nil
}

func (w *ConsoleWriter) SetJson() {
	w.fmt = FmtTypeJson
}

func (w *ConsoleWriter) SetRaw() {
	w.fmt = FmtTypeRaw
}

func (w *ConsoleWriter) SetColorRaw() {
	w.fmt = FmtTypeColorRaw
}

func (w *ConsoleWriter) SetFmt(fmt FmtType) {
	w.fmt = fmt
}

func (w *ConsoleWriter) Fmt() FmtType {
	return w.fmt
}
