/**
 * @Author: wenliangzhang
 * @Description:
 * @File: ip
 * @Version: 1.0.0
 * @Date: 2022/4/19 11:23 AM
 */
package common

import (
	"bytes"
	"fmt"
	"github.com/golib/v3/env"
	"net"
	"runtime"
	"strconv"
	"strings"
)

func InternalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if inter.Flags&net.FlagUp != 0 && !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

var httpAddr string

func GetHttpAddr() string {
	if httpAddr == "" {
		httpAddr = fmt.Sprintf("http://%s:%s", InternalIP(), env.HttpPort())
	}
	return httpAddr
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
