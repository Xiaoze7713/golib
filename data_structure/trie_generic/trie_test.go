package trie_generic_test

import (
	"git.singularity-ai.com/backend/library/v2/data_structure/trie_generic"
	"git.singularity-ai.com/backend/library/v2/easy_process"
	"git.singularity-ai.com/backend/library/v2/utils"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestTrie(t *testing.T) {
	r := trie_generic.NewTrie[string](nil)
	r.Insert("河北", &trie_generic.DataType[string]{Data: "河北"})
	r.Insert("湖南", &trie_generic.DataType[string]{Data: "湖南"})
	r.Insert("湖北", &trie_generic.DataType[string]{Data: "湖北省"})
	results := r.Match3("我从湖北来")
	println(utils.MustJson(results))
	sb := strings.Builder{}
	for _, r := range results {
		if r.Data == nil {
			sb.WriteString(r.Str)
			continue
		} else {
			sb.WriteString(r.Data.Data)
		}
	}
	println(sb.String())
}

type StrList []string

func TestTrieIgc(t *testing.T) {
	r := trie_generic.NewTrie[StrList](nil)
	rand.Seed(time.Now().UnixMilli())
	r.Insert("Hi", &trie_generic.DataType[StrList]{Data: StrList{"哈哈哈哈", "你好啊"}})
	r.Insert("heLLo", &trie_generic.DataType[StrList]{Data: StrList{"hello", "hi"}})
	r.Insert("hehe", &trie_generic.DataType[StrList]{Data: StrList{"呵呵呵", "嘿嘿"}})
	//r.Insert("湖南", "湖南")
	//r.Insert("湖北", "湖北省")
	results := r.MatchIgc("hi hello hehello")
	println(utils.MustJson(results))
	sb := strings.Builder{}
	for _, r := range results {
		if r.Data == nil {
			sb.WriteString(r.Str)
			continue
		}
		ds := r.Data.Data
		sb.WriteString(easy_process.ChoiceOne(ds))
	}
	println(sb.String())
}

func TestCut(t *testing.T) {
	sentence := "你好么？我叫干不倒"
	r, err := regexp.Compile("[^.。?？!！]+[.。?？!！]?")
	if err != nil {
		return
	}
	println(r.FindString(sentence))
}
