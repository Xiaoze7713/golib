package xRequest

import (
	"bytes"
	"encoding/json"
	"errors"
	"git.singularity-ai.com/backend/library/xContext"
	"io/ioutil"
	"net/http"
)

type HttpBase struct {
	Url string
	Cli *http.Client
}

func (m *HttpBase) Post(ctx *xContext.XContext, req interface{}, resp interface{}) (err error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	body := bytes.NewReader(bodyBytes)
	request, err := http.NewRequest(http.MethodPost, m.Url, body)
	if err != nil {
		return err
	}
	ctx.SetTag(string(xContext.ReqBody), string(bodyBytes))
	request.Header.Add("contentType", "application/json")
	ctx.SetTag("do_request", ctx.SerializeSpanContext())
	request.Header.Add("trace_id", ctx.SerializeSpanContext())
	httpResp, err := m.Cli.Do(request)
	if err != nil {
		return err
	}
	ctx.LogFields("request", m.Url)
	ctx.SetTag("url", m.Url)
	defer httpResp.Body.Close()
	ctx.SetTag("http_status", httpResp.StatusCode)
	if httpResp.StatusCode != 200 {
		return errors.New("except err != 200")
	}
	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		ctx.SetTag(string(xContext.Error), "read resp body err")
		return errors.New("read resp body err")
	}
	ctx.SetTag(string(xContext.RespBody), string(respBody))
	ctx.Debugf("resp=%v", string(respBody))
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
