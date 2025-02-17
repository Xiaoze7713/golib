/**
 * @Author: wenliangzhang
 * @Description:
 * @File: init
 * @Version: 1.0.0
 * @Date: 2022/5/9 7:54 PM
 */
package mysql

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/golib/v2/log"
)

type ConfigItem struct {
	User         string `toml:"user"`
	Password     string `toml:"password"`
	Host         string `toml:"host"`
	CharSet      string `toml:"charset"`
	ParseTime    bool   `toml:"parserTime"`
	Loc          string `toml:"loc"`
	DatabaseName string `toml:"database"`
}

type GormConnector struct {
	DB *gorm.DB
}

var defaultMysqlConfigPath = "conf/service/mysql.toml"

var dbConnectors map[string]GormConnector

func Init(filePath string) error {

	if filePath == "" {
		filePath = defaultMysqlConfigPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Println("mysql.toml not exist")
		return nil
	}

	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("read mysql.toml failed,err=%s", err.Error())
		return err
	}
	confStrRaw := string(bs)
	var mysqlConf map[string]ConfigItem
	if _, err := toml.Decode(confStrRaw, &mysqlConf); err != nil {
		panic(fmt.Sprintf("Can't load config file, %s", err.Error()))
	}

	dbConnectors = make(map[string]GormConnector)

	for dbName, configItem := range mysqlConf {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
			configItem.User, configItem.Password, configItem.Host, configItem.DatabaseName, configItem.CharSet, configItem.ParseTime)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		db.Logger.LogMode(0)
		if err != nil {
			log.Errorln("mysql open failed,db=%s,err=%s", dbName, err.Error())
			continue
		}
		dbConnectors[dbName] = GormConnector{db}
	}
	return nil
}
