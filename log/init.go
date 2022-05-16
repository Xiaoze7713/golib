/**
 * @Author: wenliangzhang
 * @Description:
 * @File: init
 * @Version: 1.0.0
 * @Date: 2022/4/21 6:33 PM
 */
package log

import (
	"fmt"
	"git.singularity-ai.com/backend/library/env"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var defaultLogConfigPath = "conf/log.toml"

type LogConfig struct {
	AppName     string `toml:"appname"`
	RotateUnit  int    `toml:"rotateunit"`
	RotateCount int    `toml:"rotatecount"`
	Level       string `toml:"level"`
	Stdout      bool   `toml:"stdout"`
}

var loggerDef *Logger
var loggerWf *Logger

func Init(filePath string) {
	var config LogConfig
	if filePath == "" {
		filePath = defaultLogConfigPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Println("log.toml not exist")
		return
	}
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		panic(fmt.Sprintf("Can't load config file, %s", err.Error()))
	}
	// 初始化日志目录
	initLogDir(env.LogRootPath())
	if env.AppName() == "" {
		env.SetAppName(config.AppName)
	}

	loggerDef = &Logger{
		NewLogger(config, ""),
	}

	loggerWf = &Logger{
		NewLogger(config, ".wf"),
	}

}

func initLogDir(path string) error {
	if path == "" {
		path = "./"
	}
	_, staterr := os.Stat(path)
	if os.IsNotExist(staterr) {
		// 创建目录
		cterr := os.MkdirAll(path, 0777)
		if cterr != nil {
			return fmt.Errorf("log conf err: create log dir '%s' error: %s", path, cterr)
		}
	}
	return nil
}
