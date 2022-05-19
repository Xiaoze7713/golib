/**
 * @Author: wenliangzhang
 * @Description:
 * @File: invoke
 * @Version: 1.0.0
 * @Date: 2022/5/19 7:34 PM
 */
package common

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"strings"
)

// BuildCmd scheme协议头拼接
func BuildCmd(strPath string, arrParams map[string]interface{}, arrPd map[string]string) string {

	// 统计字段 tag tab source
	if len(arrPd) == 0 {
		arrPd = make(map[string]string)
		arrPd["tab"] = ""
		arrPd["tag"] = ""
		arrPd["source"] = ""
	}

	var strParams string
	if len(arrParams) != 0 {
		jsonParams, _ := jsoniter.Marshal(arrParams)
		strParams = strings.Replace(string(jsonParams), "\\u0026", "&", -1)
		strParams = strings.Replace(url.QueryEscape(strParams), "+", "%20", -1)
	} else {
		strParams = ""
	}

	var strCmd string
	strCmdTpl := "singularity://singularity-ai.com%s?params=%s&tab=%s&tag=%s&source=%s"
	strCmd = fmt.Sprintf(strCmdTpl, strPath, strParams, arrPd["tab"], arrPd["tag"], arrPd["source"])
	return strCmd
}
