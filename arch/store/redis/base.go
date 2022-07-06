/**
 * @Author: wenliangzhang
 * @Description:
 * @File: base
 * @Version: 1.0.0
 * @Date: 2022/7/5 7:59 PM
 */
package redis

import (
	rds "github.com/gomodule/redigo/redis"

	"git.singularity-ai.com/backend/library/log"
)

func sampleDo(method string, args ...interface{}) (reply interface{}, err error) {
	conn := client.C.Get()
	defer func(conn rds.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("close conn err=%v", err)
		}
	}(conn)
	return conn.Do(method, args...)
}

func sampleDoInt(method string, args ...interface{}) (reply int, err error) {
	reply, err = rds.Int(sampleDo(method, args...))
	if err == rds.ErrNil {
		return 0, nil
	}
	if err != nil {
		log.Fatalf("redis err=%v", err)
	}
	return
}

func sampleDoInt64(method string, args ...interface{}) (reply int64, err error) {
	reply, err = rds.Int64(sampleDo(method, args...))
	if err == rds.ErrNil {
		return 0, nil
	}
	if err != nil {
		log.Fatalf("redis err=%v", err)
	}
	return
}

func sampleDoString(method string, args ...interface{}) (reply string, err error) {
	reply, err = rds.String(sampleDo(method, args...))
	if err == rds.ErrNil {
		return "", nil
	}
	if err != nil {
		log.Fatalf("redis err=%v", err)
	}
	return
}
