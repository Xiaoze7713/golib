package trie

import (
	"git.singularity-ai.com/backend/library/utils"
	"regexp"
	"testing"
)

func TestTrie(t *testing.T) {
	r := NewTrie()
	r.Insert("河北", "河北")
	r.Insert("湖南", "湖南")
	r.Insert("湖北", "湖北省")
	results := r.Match3("我从湖北省来")
	println(utils.MustJson(results))
}

func TestCut(t *testing.T) {
	sentence := "你好么？我叫干不倒"
	r, err := regexp.Compile("[^.。?？!！]+[.。?？!！]?")
	if err != nil {
		return
	}
	println(r.FindString(sentence))
}
