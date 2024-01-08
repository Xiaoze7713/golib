package es

import (
	"github.com/Xiaoze7713/golib/v3/utils"
	"github.com/Xiaoze7713/golib/v3/xContext"
	"github.com/Xiaoze7713/golib/v3/xContext/base_if/xtrace_base"
	"github.com/Xiaoze7713/golib/v3/xContext/loggers/xlog"
	"github.com/Xiaoze7713/golib/v3/xContext/metrics/null_metric"
	"io/ioutil"
	"path"
	"src/entity"
	"src/logic/dict"
	"strings"
	"sync"
	"testing"
	"time"
)

func InitAll() {
	xlog.SetupLogDefault()
	xContext.Init(xlog.GetLogger(), xtrace_base.XTraceNoop{}, null_metric.MetricsNull{}, nil, nil, nil, nil, nil)
	err := Init([]string{"http://47.92.76.225:9200"})
	if err != nil {
		xlog.Error(err)
	}
}

func TestName(t *testing.T) {
	ctx := xContext.NewXContext("es_start")
	//err = Handler.WordIndex(ctx)
	//if err != nil {
	//	xlog.Error(err)
	//}
	InitAll()
	err := dict.Handler.Insert(ctx, &entity.WordType{
		Tags: []entity.Tag{"食物", "水果", "植物"},
		Word: "苹果",
	})
	if err != nil {
		xlog.Error(err)
	}
	err = dict.Handler.UpdateWord(ctx, &entity.WordType{
		Tags: []entity.Tag{"食物", "水果", "植物"},
		Word: "苹果",
	})
	if err != nil {
		xlog.Error(err)
	}
	time.Sleep(time.Second * 3)
}

func TestLoadData(t *testing.T) {
	InitAll()
	dataPath := "/Users/cangxiaoze/user/git/thesaurus/data/"
	idx, err := ioutil.ReadFile(path.Join(dataPath, "index.idx"))
	if err != nil {
		xlog.Error(err)
		return
	}
	w := sync.WaitGroup{}
	for _, line := range strings.Split(string(idx), "\n") {
		nsLine := strings.TrimSpace(line)
		if nsLine == "" {
			continue
		}
		items := strings.Split(nsLine, " ")
		xlog.Debug(utils.MustJson(items))
		fileName := items[0]
		tag := items[1]
		wordDict, err := ioutil.ReadFile(path.Join(dataPath, fileName))
		if err != nil {
			xlog.Debug(err)
			continue
		}
		w.Add(1)
		go func(data []byte) {
			a := 0
			for _, line2 := range strings.Split(string(data), "\n") {
				wordLine := strings.TrimSpace(line2)
				word := strings.Split(wordLine, "\t")[0]
				err = func() (err error) {
					ctx := xContext.NewXContext(word + tag)
					newWord := &entity.WordType{
						Tags: []entity.Tag{entity.Tag(tag)},
						Word: strings.TrimSpace(word),
					}
					err = dict.Handler.Insert(ctx, newWord)
					defer ctx.Fin()
					return nil
				}()
				a += 1
				if a%100 == 0 {
					xlog.Debug(a, fileName, word)
				}
				if err != nil {
					xlog.Debug(err)
					continue
				}
			}
		}(wordDict)
	}
	w.Wait()
}
