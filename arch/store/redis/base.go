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

	"git.singularity-ai.com/backend/library/arch/web"
)

func sampleDo(ctx *web.WebContext, method string, args ...interface{}) (reply interface{}, err error) {
	conn := client.C.Get()
	defer func(conn rds.Conn) {
		err := conn.Close()
		if err != nil {
			ctx.Fatalf("close conn err=%v", err)
		}
	}(conn)
	return conn.Do(method, args...)
}

func sampleDoInt(ctx *web.WebContext, method string, args ...interface{}) (reply int, err error) {
	reply, err = rds.Int(sampleDo(ctx, method, args...))
	if err == rds.ErrNil {
		return 0, nil
	}
	if err != nil {
		ctx.Fatalf("redis err=%v", err)
	}
	return
}

func sampleDoInt64(ctx *web.WebContext, method string, args ...interface{}) (reply int64, err error) {
	reply, err = rds.Int64(sampleDo(ctx, method, args...))
	if err == rds.ErrNil {
		return 0, nil
	}
	if err != nil {
		ctx.Fatalf("redis err=%v", err)
	}
	return
}

func sampleDoString(ctx *web.WebContext, method string, args ...interface{}) (reply string, err error) {
	reply, err = rds.String(sampleDo(ctx, method, args...))
	if err == rds.ErrNil {
		return "", nil
	}
	if err != nil {
		ctx.Fatalf("redis err=%v", err)
	}
	return
}
