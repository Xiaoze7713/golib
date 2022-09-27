package xRequest2

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.singularity-ai.com/backend/library/utils"
	"git.singularity-ai.com/backend/library/xContext"
	"github.com/go-resty/resty/v2"
	"net/url"
	"time"
)

type HttpBase struct {
	Url string
	Cli *resty.Client
}

func (m *HttpBase) Post(ctx *xContext.XContext, req interface{}, resp interface{}) (err error) {
	method := ""
	u, err := url.Parse(m.Url)
	if err == nil {
		method = u.Path
	}

	ctx = xContext.NewChildXContext(ctx, method+"[REQ]")
	defer ctx.Fin()
	ctx.LogFields("request", m.Url)
	ctx.Infof("[%v] req %v", method, utils.MustJson(req))
	//ctx.SetTag(string(xContext.ReqBody), utils.MustJson(req))
	ctx.SetTag("trace", ctx.SerializeSpanContext())
	ctx.SetTag("url", m.Url)
	httpResp, err := m.Cli.R().
		SetBody(req).
		SetHeader("Content-Type", "application/json").
		SetHeader("trace_id", ctx.SerializeSpanContext()).
		SetHeader("time", time.Now().String()).
		Post(m.Url)
	if err != nil {
		ctx.SetTag("err", err.Error())
		return err
	}
	defer ctx.SetTag("http_status", httpResp.StatusCode)
	if httpResp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("exception http code %v", httpResp.StatusCode()))
	}
	respBody := httpResp.Body()
	//ctx.SetTag(string(xContext.RespBody), string(respBody))
	//ctx.Debugf("resp %v", string(respBody))
	ctx.Infof("[%v] resp %v", method, string(respBody))
	ctx.LogFields("resp", string(respBody))
	err = json.Unmarshal(respBody, resp)
	if err != nil {
		return err
	}
	return nil
}

//func (m *HttpBase) Get(ctx *xContext.XContext, req interface{}, resp interface{}) (err error) {
//	bodyBytes, err := json.Marshal(req)
//	if err != nil {
//		return err
//	}
//	body := bytes.NewReader(bodyBytes)
//	request, err := http.NewRequest(http.MethodGet, m.Url, body)
//	if err != nil {
//		return err
//	}
//	ctx.SetTag(string(xContext.ReqBody), string(bodyBytes))
//	request.Header.Add("contentType", "application/json")
//	ctx.SetTag("do_request", ctx.SerializeSpanContext())
//	request.Header.Add("trace_id", ctx.SerializeSpanContext())
//	httpResp, err := m.Cli.Do(request)
//	if err != nil {
//		return err
//	}
//	ctx.LogFields("request", m.Url)
//	ctx.SetTag("url", m.Url)
//	defer httpResp.Body.Close()
//	ctx.SetTag("http_status", httpResp.StatusCode)
//	if httpResp.StatusCode != 200 {
//		return errors.New("except err != 200")
//	}
//	respBody, err := ioutil.ReadAll(httpResp.Body)
//	if err != nil {
//		ctx.SetTag(string(xContext.Error), "read resp body err")
//		return errors.New("read resp body err")
//	}
//	ctx.SetTag(string(xContext.RespBody), string(respBody))
//	ctx.Debugf("resp=%v", string(respBody))
//	err = json.Unmarshal(respBody, resp)
//	if err != nil {
//		return err
//	}
//	return nil
//}
