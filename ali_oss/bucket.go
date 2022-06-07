/**
 * @Author: wenliangzhang
 * @Description:
 * @File: bucket
 * @Version: 1.0.0
 * @Date: 2022/4/26 11:07 AM
 */
package ali_oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Bucket struct {
	Config *Config
	*oss.Bucket
}

func BucketClient(client *oss.Client, c *Config) (*Bucket, error) {
	if c == nil {
		c = config
	}
	b, _ := client.Bucket(c.BucketName)
	return &Bucket{config, b}, nil
}
