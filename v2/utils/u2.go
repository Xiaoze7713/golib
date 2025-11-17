package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	tokenKey = "token"
)

func UTimeMs() int64 {
	return time.Now().UnixMilli()
}

func SecondToMs(v int64) int64 {
	return v * 1000
}

func MultiDur(dur time.Duration, times int64) time.Duration {
	return time.Duration(times) * dur
}

func FileMd5(filePath string) (md5String string, err error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	m := md5.New()
	m.Write(file)
	md5String = hex.EncodeToString(m.Sum(nil))
	return
}

func MaxInt32(ns ...int32) (m int32) {
	m = ns[0]

	for i := 1; i < len(ns); i++ {
		if ns[i] > m {
			m = ns[i]
		}
	}
	return m
}

func BytesMd5(bytes []byte) (md5String string, err error) {
	m := md5.New()
	m.Write(bytes)
	md5String = hex.EncodeToString(m.Sum(nil))
	return
}

func StringsMd5(strList []string) (md5String string, err error) {
	m := md5.New()
	for _, s := range strList {
		m.Write([]byte(s))
	}
	md5String = hex.EncodeToString(m.Sum(nil))
	return
}

func MinInt32(ns ...int32) (m int32) {
	m = ns[0]

	for i := 1; i < len(ns); i++ {
		if ns[i] < m {
			m = ns[i]
		}
	}
	return m
}

func APPVersionCode() (code int32, err error) {
	t := time.Now()
	ts := fmt.Sprintf("%s", t.Format("20060102"))
	c, err := strconv.Atoi(ts)
	if err != nil {
		return
	}
	code = int32(c)
	return
}

func UTimeMSDur(start int64) int64 {
	return time.Now().UnixMilli() - start
}

func UTimeMSDurF(start int64) float64 {
	return float64((time.Now().UnixNano() - start) / 1000000)
}

func Find(patten, s string) (idxList [][]int, err error) {
	reCmp, err := regexp.Compile(patten)
	if err != nil {
		return nil, err
	}
	idxList = reCmp.FindAllStringIndex(s, -1)
	return idxList, nil
}

func UnMarshalToString(v interface{}) (s string) {
	if reflect.TypeOf(v).Kind() == reflect.String {
		s = v.(string)
		return s
	}
	return string(v.([]byte))
}

func UnMarshalToInterface(v interface{}, out interface{}) (err error) {
	//	s := string(v.([]byte))
	var sb []byte
	if reflect.TypeOf(v).Kind() != reflect.String {
		sb = v.([]byte)
	} else {
		sb = []byte(v.(string))
	}
	if out == nil {
		return errors.New("nil out")
	}
	if reflect.TypeOf(out).Kind() == reflect.Ptr {
		err = json.Unmarshal(sb, out)
	} else {
		err = json.Unmarshal(sb, &out)
	}
	return err

}

func MarshalToString(v interface{}) (s string, err error) {
	sb, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	s = string(sb)
	return s, err
}

func MarshalToStringNoErr(v interface{}) (s string) {
	sb, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	s = string(sb)
	return s
}
func FuncName(v interface{}) (s string) {
	if reflect.TypeOf(v).Kind() != reflect.Func {
		s = "NOT_FUNC"
	}
	funcName := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
	l1 := strings.Split(funcName, ".")
	s1 := l1[len(l1)-1]
	l2 := strings.Split(s1, "-fm")
	s = l2[0]
	return s
}

func MustJson(v interface{}) (s string) {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		return reflect.TypeOf(v).Kind().String() + err.Error()
	}
	return s
}

func JsonUnMarshalString(s string, objPtr interface{}) (err error) {
	decoder := jsoniter.NewDecoder(strings.NewReader(s))
	decoder.UseNumber()
	err = decoder.Decode(objPtr)
	return err
}

func JsonMarshalString(v interface{}) (s string, err error) {
	s, err = jsoniter.MarshalToString(v)
	if err != nil {
		return "", errors.New(reflect.TypeOf(v).Kind().String() + err.Error())
	}
	return s, err
}
