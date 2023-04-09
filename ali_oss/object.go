/**
 * @Author: wenliangzhang
 * @Description:
 * @File: object
 * @Version: 1.0.0
 * @Date: 2022/4/26 11:18 AM
 */
package ali_oss

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"git.singularity-ai.com/backend/library/v2/log"
)

func (b *Bucket) DownloadObject(object string) ([]byte, error) {
	body, err := b.GetObject(object)
	if err != nil {
		log.Errorf("GetObject failed, err=%v", err.Error())
		return nil, err
	}

	data, err := ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		log.Errorf("GetObject failed, err=%v", err.Error())
		return nil, err
	}
	return data, nil
}

func (b *Bucket) DownloadObjectToFile(object, file string) error {
	// check 文件是否已存在
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// 文件不存在则新建
		fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0660)
		if err != nil {
			log.Errorln("open file failed,err=%v", err.Error())
			return err
		}
		defer fd.Close()

		body, err := b.GetObject(object)
		if err != nil {
			log.Errorln("GetObject failed,err=%v", err.Error())
			return err
		}
		io.Copy(fd, body)
		body.Close()

	} else {
		err := b.GetObjectToFile(object, file)
		if err != nil {
			log.Errorln("GetObjectToFile failed,err=%v", err.Error())
			return err
		}
	}

	return nil
}

func (b *Bucket) UploadObject(object string, data []byte, options ...oss.Option) error {
	if b == nil {
		log.Errorln("bucket client is nil")
		return errors.New("bucket client is nil")
	}
	err := b.PutObject(object, bytes.NewReader(data), options...)
	if err != nil {
		log.Errorf("PutObject failed, err=%v", err.Error())
		return err
	}
	return nil
}

func (b *Bucket) UploadObjectFromFile(object string, filePath string) error {
	if b == nil {
		log.Errorln("bucket client is nil")
		return errors.New("bucket client is nil")
	}
	err := b.PutObjectFromFile(object, filePath)
	if err != nil {
		log.Errorf("PutObject failed, err=%v", err.Error())
		return err
	}
	return nil
}

func (b *Bucket) GetObjectUrl(object string) string {
	return strings.Join([]string{b.Config.DomainName, object}, "/")
}
