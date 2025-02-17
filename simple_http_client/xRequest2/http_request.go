package xRequest2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golib/v3/utils"
	"github.com/golib/v3/xContext"
	"github.com/go-resty/resty/v2"
	"net/http"
	url2 "net/url"
	"time"
)

type HttpBase struct {
	Cli         *resty.Client
	Url         string
	skipReqLog  map[string]bool
	skipRespLog map[string]bool
}

func (m *HttpBase) Post(ctx *xContext.XContext, url string, req interface{}, resp interface{}, headersMap ...map[string]string) (err error) {
	return m.Request(ctx, url, req, resp, http.MethodPost, headersMap...)
}

func (m *HttpBase) Get(ctx *xContext.XContext, url string, req interface{}, resp interface{}, headersMap ...map[string]string) (err error) {
	return m.Request(ctx, url, req, resp, http.MethodGet, headersMap...)
}

func (m *HttpBase) Request(ctx *xContext.XContext, url string, req interface{}, resp interface{}, httpMethod string, headersMap ...map[string]string) (err error) {
	path := ""
	host := ""
	u, err := url2.Parse(url)
	if err == nil {
		path = u.Path
		host = u.Host
	}
	ctx = xContext.NewChildXContext(ctx, path+"[REQ]")
	defer ctx.Fin()
	ctx.LogFields("url", url)
	if ok, exist := m.skipReqLog[path]; !exist || !ok {
		ctx.Infof("[%v][%v][REQ]%v", host, path, utils.MustJson(req))
		ctx.LogFields(string(xContext.ReqBody), utils.MustJson(req))
	} else {
		ctx.Infof("[%v][%v][REQ]%v", host, path, "SkipDumpBody")
		ctx.LogFields(string(xContext.RespBody), "skip body field")
	}
	ctx.SetTag("trace", ctx.SerializeSpanContext())
	//ctx.SetTag("url", m.Url)
	httpRequest := m.Cli.R().
		SetBody(req).
		SetHeader("Content-Type", "application/json").
		SetHeader("trace_id", ctx.SerializeSpanContext()).
		SetHeader("time", time.Now().String())
	if len(headersMap) > 0 {
		for _, headers := range headersMap {
			httpRequest.SetHeaders(headers)
		}
	}
	ctx.LogFields("headers", utils.MustJson(httpRequest.Header))
	var httpResp *resty.Response
	switch httpMethod {
	case http.MethodPost:
		httpResp, err = httpRequest.
			Post(url)
	case http.MethodGet:
		httpResp, err = httpRequest.
			Get(url)
	default:
		httpResp, err = httpRequest.
			Get(url)
	}
	if err != nil {
		ctx.SetTag("err", err.Error())
		return err
	}
	defer ctx.SetTag("http_status", httpResp.StatusCode())
	if httpResp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("exception http code %v", httpResp.StatusCode()))
	}
	respBody := httpResp.Body()
	//ctx.SetTag(string(xContext.RespBody), string(respBody))
	//ctx.Debugf("resp %v", string(respBody))
	if ok, exist := m.skipRespLog[path]; !exist || !ok {
		ctx.Infof("[%v][%v][RESP] %v", host, path, string(respBody))
		ctx.LogFields(string(xContext.RespBody), string(respBody))
	} else {
		ctx.Infof("[%v][%v][RESP] %v", host, path, "SkipDumpBody")
		ctx.LogFields(string(xContext.RespBody), "skip body field")
	}
	ctx.LogFields(string(xContext.HttpRespHeader), utils.MustJson(httpResp.Header()))
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
