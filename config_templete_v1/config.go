package config_templete_v1

import (
	"fmt"
)

type SimpleHttpServerConfig struct {
	ServerURL string   `toml:"server_url"`
	TimeOut   int64    `toml:"time_out"`
	Hosts     []string `toml:"hosts"`
	PathList  []string `toml:"path_list"`
}

type SimpleRPCServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func (m *SimpleRPCServerConfig) ToAddress() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

type ServiceInfo struct {
	NetName     string `toml:"net_name"`
	EnvName     string `toml:"env_name"`
	ServiceName string `toml:"service_name"`
	IDC         string `toml:"idc"`
}

type MysqlCommon struct {
	User         string `toml:"user"`
	Password     string `toml:"password"`
	Host         string `toml:"host"`
	CharSet      string `toml:"charset"`
	ParseTime    bool   `toml:"parserTime"`
	Loc          string `toml:"loc"`
	DatabaseName string `toml:"database"`
}

func (m *MysqlCommon) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local", m.User, m.Password, m.Host, m.DatabaseName, m.CharSet, m.ParseTime)
}

type EsConfig struct {
}

type Neo4jConfig struct {
	Url      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type RedisConfig struct {
	Host       string   `toml:"host"`
	Password   string   `toml:"password"`
	DBNumber   int64    `toml:"db_no"`
	Mode       string   `toml:"mode"`
	MultiHosts []string `toml:"multi_hosts"`
	UserName   string   `toml:"user_name"`
}

func (m *RedisConfig) IsCluster() bool {
	return m.Mode == "cluster"
}
