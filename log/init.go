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
	"os"

	"github.com/BurntSushi/toml"

	"github.com/golib/v2/env"
)

var defaultLogConfigPath = "conf/log.toml"

type LogConfig struct {
	AppName     string `toml:"appname"`
	RotateUnit  int    `toml:"rotateunit"`
	RotateCount int    `toml:"rotatecount"`
	Level       string `toml:"level"`
	Stdout      bool   `toml:"stdout"`
}

var globalConfig LogConfig

var defaultLoggerConfig = LogConfig{
	"singularity",
	1,
	48,
	"warn",
	false,
}

func Init(filePath string) {
	if filePath == "" {
		filePath = defaultLogConfigPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("log.toml not exist")
		globalConfig = defaultLoggerConfig
	} else {
		if _, err = toml.DecodeFile(filePath, &globalConfig); err != nil {
			fmt.Sprintf("Can't load config file, %s", err.Error())
			globalConfig = defaultLoggerConfig
		}
	}

	// 初始化日志目录
	initLogDir(env.LogRootPath())
	if env.AppName() == "unknown" {
		env.SetAppName(globalConfig.AppName)
	}
	initGlobalLogger(globalConfig)
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
