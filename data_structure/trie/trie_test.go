package trie

import (
	"git.singularity-ai.com/backend/library/utils"
	"regexp"
	"src/common/easy_process"
	"strings"
	"testing"
)

func TestTrie(t *testing.T) {
	r := NewTrie()
	r.Insert("河北", "河北")
	r.Insert("湖南", "湖南")
	r.Insert("湖北", "湖北省")
	results := r.Match3("我从湖北来")
	println(utils.MustJson(results))
	sb := strings.Builder{}
	for _, r := range results {
		if r.Data == nil {
			sb.WriteString(r.Str)
			continue
		}
		ds, ok := r.Data.(string)
		if !ok {
			sb.WriteString(r.Str)
		}
		sb.WriteString(ds)
	}
	println(sb.String())
}

func TestTrieIgc(t *testing.T) {
	r := NewTrie()
	r.Insert("Hi", []string{"哈哈哈"})
	r.Insert("heLLo", []string{"你好"})
	r.Insert("hehe", []string{"呵呵"})
	//r.Insert("湖南", "湖南")
	//r.Insert("湖北", "湖北省")
	results := r.MatchIgc("hi hello")
	println(utils.MustJson(results))
	sb := strings.Builder{}
	for _, r := range results {
		if r.Data == nil {
			sb.WriteString(r.Str)
			continue
		}
		ds, ok := r.Data.([]string)
		if !ok {
			sb.WriteString(r.Str)
		}
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
