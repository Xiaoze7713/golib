package xlog

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type ConfFileWriter struct {
	On                  bool    `json:"On"`
	Fmt                 FmtType `json:"fmt"`
	LogPath             string  `json:"LogPath"`
	RotateLogPath       string  `json:"RotateLogPath"`
	WfLogPath           string  `json:"WfLogPath"`
	RotateWfLogPath     string  `json:"RotateWfLogPath"`
	PublicLogPath       string  `json:"PublicLogPath"`
	RotatePublicLogPath string  `json:"RotatePublicLogPath"`
}

type ConfConsoleWriter struct {
	On  bool    `json:"On"`
	Fmt FmtType `json:"fmt"`
}

type LogConfig struct {
	Level string            `json:"LogLevel"`
	FW    ConfFileWriter    `json:"FileWriter"`
	CW    ConfConsoleWriter `json:"ConsoleWriter"`
}

func SetupLogDefault() {
	var lc LogConfig
	lc.Level = "DEBUG"
	lc.CW = ConfConsoleWriter{
		On:  true,
		Fmt: FmtTypeRaw,
	}
	SetupLogWithConf(lc)
}

func SetupLogWithConfFile(file string) (err error) {
	var lc LogConfig
	cnt, err := ioutil.ReadFile(file)

	if err = json.Unmarshal(cnt, &lc); err != nil {
		return
	}
	return SetupLogWithConf(lc)

}
func SetupLogWithConf(lc LogConfig) (err error) {
	if lc.FW.On {
		if len(lc.FW.LogPath) > 0 {
			w := NewFileWriter()
			w.SetFileName(lc.FW.LogPath)
			w.SetPathPattern(lc.FW.RotateLogPath)
			w.SetLogLevelFloor(TRACE)
			if len(lc.FW.WfLogPath) > 0 {
				w.SetLogLevelCeil(PUBLIC)
			} else {
				w.SetLogLevelCeil(FATAL)
			}
			w.SetFmt(lc.FW.Fmt)
			Register(w)
		}

		if len(lc.FW.WfLogPath) > 0 {
			wfw := NewFileWriter()
			wfw.SetFileName(lc.FW.WfLogPath)
			wfw.SetPathPattern(lc.FW.RotateWfLogPath)
			wfw.SetLogLevelFloor(WARNING)
			wfw.SetLogLevelCeil(FATAL)
			wfw.SetFmt(lc.FW.Fmt)
			Register(wfw)
		}

		if len(lc.FW.PublicLogPath) > 0 {
			wfp := NewFileWriter()
			wfp.SetFileName(lc.FW.PublicLogPath)
			wfp.SetPathPattern(lc.FW.RotatePublicLogPath)
			wfp.SetLogLevelFloor(PUBLIC)
			wfp.SetLogLevelCeil(PUBLIC)
			wfp.SetFmt(lc.FW.Fmt)
			Register(wfp)
		}
	}

	if lc.CW.On {
		w := NewConsoleWriter()
		w.SetFmt(lc.CW.Fmt)
		Register(w)
	}

	switch lc.Level {
	case "trace":
		SetLevel(TRACE)

	case "debug":
		SetLevel(DEBUG)

	case "info":
		SetLevel(INFO)

	case "warning":
		SetLevel(WARNING)

	case "error":
		SetLevel(ERROR)

	case "fatal":
		SetLevel(FATAL)

	default:
		err = errors.New("invalid log level")
	}
	return
}
