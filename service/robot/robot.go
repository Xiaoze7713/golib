/**
 * @Author: wenliangzhang
 * @Description:
 * @File: robot
 * @Version: 1.0.0
 * @Date: 2022/7/15 11:35 AM
 */
package robot

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Xiaoze7713/golib/v3/arch/web"
)

type ZhCn struct {
	Title   string      `json:"title"`
	Content interface{} `json:"content"`
}

// 飞书自定义机器人开发文档
// https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN

type Post struct {
	ZhCn ZhCn `json:"zh_cn"`
}

type Content struct {
	Post Post `json:"post,omitempty"`
}

type Msg struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type Response struct {
	Extra         interface{} `json:"Extra"`
	StatusCode    int         `json:"StatusCode"`
	StatusMessage string      `json:"StatusMessage"`
}

func SendMsg(ctx *web.WebContext, webhook string, title string, content interface{}) error {

	msg := Msg{
		"post",
		Content{
			Post{
				ZhCn{
					title,
					content,
				},
			},
		},
	}
	bodyByte, _ := json.Marshal(msg)
	resp, err := http.Post(webhook, "application/json", bytes.NewReader(bodyByte))
	if err != nil {
		ctx.Errorf("Send Robot Msg failed, err:%v", err)
		return err
	}
	if resp.StatusCode != 200 {
		ctx.Errorf("Send Robot Msg failed, errno:%v,errmsg:%v", resp.StatusCode, resp.Status)
		return err
	}
	b, _ := io.ReadAll(resp.Body)
	var respBody Response
	json.Unmarshal(b, &respBody)
	if respBody.StatusCode != 0 {
		ctx.Warnf("Send Robot Msg failed, errno:%v,errmsg:%v", respBody.StatusCode, respBody.StatusMessage)
		return errors.New(respBody.StatusMessage)
	}
	return nil
}
