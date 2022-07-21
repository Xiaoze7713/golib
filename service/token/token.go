/**
 * @Author: wenliangzhang
 * @Description:
 * @File: token
 * @Version: 1.0.0
 * @Date: 2022/4/20 2:05 PM
 */
package token

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"

	"git.singularity-ai.com/backend/library/arch/web"
	"git.singularity-ai.com/backend/library/log"
)

type TokenConfig struct {
	Name string `toml:"name"`
	Host string `toml:"host"`
}

type TokenData struct {
	UserId  string `json:"userId"`
	RobotId string `json:"robotId"`
}

type TokenResp struct {
	Code     int       `json:"code"`
	CodeMsg  string    `json:"code_msg"`
	RespData TokenData `json:"resp_data"`
}

var defaultTokenConfigPath = "conf/service/token.toml"

var tokenConfig TokenConfig

func InitToken(filePath string) {
	if filePath == "" {
		filePath = defaultTokenConfigPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Infoln("token.toml not exist")
		return
	}
	if _, err := toml.DecodeFile(filePath, &tokenConfig); err != nil {
		panic(fmt.Sprintf("Can't load config file, %s", err.Error()))
	}
}

func VerifyToken(ctx *web.WebContext, token string) (TokenData, error) {

	u := "https://" + tokenConfig.Host + "/token/verify"

	body := struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}{
		Data: struct {
			Token string `json:"token"`
		}{
			token,
		},
	}
	bodyByte, _ := json.Marshal(body)

	request, _ := http.NewRequest("POST", u, bytes.NewReader(bodyByte))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("trace_id", ctx.SerializeSpanContext())

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		ctx.Errorf("Verify Token failed, err:%v", err)
		return TokenData{}, err
	}
	if resp.StatusCode != 200 {
		ctx.Errorf("Verify Token failed, errno:%v,errmsg:%v", resp.StatusCode, resp.Status)
		return TokenData{}, err
	}

	b, _ := io.ReadAll(resp.Body)
	var respBody TokenResp
	json.Unmarshal(b, &respBody)
	if respBody.Code != 200 {
		ctx.Warnf("Verify Token failed, errno:%v,errmsg:%v", respBody.Code, respBody.CodeMsg)
		return TokenData{}, errors.New(respBody.CodeMsg)
	}
	return respBody.RespData, nil
}
