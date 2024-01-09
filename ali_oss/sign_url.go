/**
 * @Author: wenliangzhang
 * @Description:
 * @File: sign_url
 * @Version: 1.0.0
 * @Date: 2022/4/26 11:46 AM
 */
package ali_oss

import (
	"github.com/Xiaoze7713/golib/v3/xContext/loggers/xlog"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func (b *Bucket) SignObjectURL(object, method string, expire int64, options ...oss.Option) (string, error) {
	url, err := b.SignURL(object, oss.HTTPMethod(method), expire, options...)
	if err != nil {
		xlog.Errorf("SignURL failed,err=%v", err.Error())
		return "", err
	}
	return url, nil
}
