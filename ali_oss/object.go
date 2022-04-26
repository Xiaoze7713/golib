/**
 * @Author: wenliangzhang
 * @Description:
 * @File: object
 * @Version: 1.0.0
 * @Date: 2022/4/26 11:18 AM
 */
package ali_oss

import (
	"errors"
	"git.singularity-ai.com/backend/ws_service/library/log"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"io/ioutil"
	"os"
)

func DownloadObject(bClient *oss.Bucket, object string) ([]byte, error) {
	if bClient == nil {
		log.Errorln("bucket client is nil")
		return nil, errors.New("bucket client is nil")
	}
	body, err := bClient.GetObject(object)
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

func DownloadObjectToFile(bClient *oss.Bucket, object, file string) error {
	if bClient == nil {
		log.Errorln("bucket client is nil")
		return errors.New("bucket client is nil")
	}

	// check 文件是否已存在
	if _, err := os.Stat(file); os.IsNotExist(err) {
		//文件不存在则新建
		fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0660)
		if err != nil {
			log.Errorln("open file failed,err=%v", err.Error())
			return err
		}
		defer fd.Close()

		body, err := bClient.GetObject(object)
		if err != nil {
			log.Errorln("GetObject failed,err=%v", err.Error())
			return err
		}
		io.Copy(fd, body)
		body.Close()

	} else {
		err := bClient.GetObjectToFile(object, file)
		if err != nil {
			log.Errorln("GetObjectToFile failed,err=%v", err.Error())
			return err
		}
	}

	return nil
}
