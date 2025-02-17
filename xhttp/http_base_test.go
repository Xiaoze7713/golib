package xhttp

import (
	"encoding/json"
	"fmt"
	"github.com/golib/v3/utils"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	var err error
	a := struct {
		Data map[string]interface{} `json:"data"`
	}{map[string]interface{}{"float64": float64(1.0), "int64": int64(1 << 62)}}
	s := utils.MustJson(a)
	println(utils.MustJson(a))
	b := map[string]interface{}{}
	err = jsoniter.UnmarshalFromString(s, &b)
	if err != nil {
		println(utils.MustJson(a), err)
	}
	println(fmt.Sprintf("%+v", b))
	c := map[string]interface{}{}
	err = json.Unmarshal([]byte(s), &c)
	if err != nil {
		println(utils.MustJson(a), err)
	}
	println(fmt.Sprintf("%+v", c))
	d := map[string]interface{}{}
	decoder := jsoniter.NewDecoder(strings.NewReader(s))
	decoder.UseNumber()
	err = decoder.Decode(&d)
	if err != nil {
		println(utils.MustJson(a), err)
	}
	println(fmt.Sprintf("%+v", d))
	dc := json.NewDecoder(strings.NewReader(s))
	dc.UseNumber()
	dc.More()
}
