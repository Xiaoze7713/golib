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

func NewBucketClient(client *oss.Client, bucket string) (*oss.Bucket, error) {
	b, _ := client.Bucket(bucket)
	return b, nil
}
