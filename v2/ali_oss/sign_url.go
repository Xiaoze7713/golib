/**
 * @Author: wenliangzhang
 * @Description:
 * @File: sign_url
 * @Version: 1.0.0
 * @Date: 2022/4/26 11:46 AM
 */
package ali_oss

import (
	"github.com/golib/v2/log"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func (b *Bucket) SignObjectURL(object, method string, expire int64, options ...oss.Option) (string, error) {
	url, err := b.SignURL(object, oss.HTTPMethod(method), expire, options...)
	if err != nil {
		log.Errorf("SignURL failed,err=%v", err.Error())
		return "", err
	}
	return url, nil
}
