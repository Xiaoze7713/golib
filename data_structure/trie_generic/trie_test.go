package trie_generic_test

import (
	"github.com/golib/v3/data_structure/trie_generic"
	"github.com/golib/v3/utils"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
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

type DictItem struct {
	Str string
	ID  int
}

func TestTrieIgc(t *testing.T) {
	// sampleS := "日本人在二战期间曾经对中国和其他国家犯下了许多残忍的罪行，包括对平民的大规模屠杀、掠夺资源、进行细菌战等。这些罪行导致了大量的死亡和伤害，对受害者和他们的家庭造成了极大的痛苦和伤害。\n因此，一些人可能会对日本人感到恐惧或厌恶，尤其是在听到他们曾经的暴行时。这种情感是很正常的，但也应该注意到历史是过去的事实，我们不能让历史的伤痛和痛苦影响我们今天的生活和关系。我们应该以平等和包容的态度对待不同的文化和国家，尊重和平等对待每个人。"
	targetCodeList := []int{16286, 639, 1620, 9064, 3171, 7180, 16532, 2049, 8699, 1327, 10235, 12291, 4169, 535, 606, 21619, 17268, 4502, 20374, 220, 1117, 2748, 2049, 14135, 4502, 12849, 13844, 226, 15587, 21966, 226, 22432, 20146, 2634, 4887, 221, 22386, 20374, 13625, 606, 12864, 4502, 17257, 1327, 9518, 220, 2049, 11560, 5307, 1327, 658, 670, 4502, 13520, 22635, 606, 16840, 4502, 18983, 1327, 9518, 221, 252, 12243, 220, 7932, 639, 11682, 691, 2049, 16286, 639, 14929, 14801, 2633, 11360, 220, 13762, 3083, 1620, 11920, 658, 670, 16532, 4502, 3136, 5986, 3060, 221, 6688, 4769, 14857, 3083, 2383, 17203, 4502, 220, 710, 595, 2277, 6287, 17687, 1026, 11344, 3083, 22312, 4502, 9031, 220, 15008, 8489, 6241, 11344, 4502, 697, 4442, 1327, 18983, 14512, 15008, 9318, 4502, 18775, 1327, 10211, 221, 15008, 2277, 6287, 667, 14140, 1327, 10994, 4502, 14718, 2049, 2381, 8385, 4502, 16012, 1327, 12291, 220, 13656, 12008, 4887, 2049, 2381, 3542, 8670, 221, 1, 6, 7, 730, 16411, 6237, 568, 730, 4502, 1257, 1957, 3083, 667, 8699, 17857, 12784, 8796, 4502, 1308, 225, 2}
	sampleS := "[SEP][USER]给我10个威猛狗狗的名字？[SEP][BOT]当然可以！以下是一些威猛狗狗的名字：\n\n1. 大狗（Big Dog）\n2. 狂野（Wild）\n3. 力量（Powerful）\n4. 勇气（Brave）\n5. 勇猛（Fearless）\n6. 无敌（Invincible）\n7. 无畏（Unafraid）\n8. 忠诚（Loyal）\n9. 勇敢（Courageous）\n10. 坚韧不拔"
	targetCodeList = []int{2, 7, 20213, 49970, 559, 1872, 4201, 18547, 4502, 1257, 1957, 225, 2, 8, 2355, 4060, 11657, 222, 9377, 3083, 7932, 1872, 4201, 18547, 4502, 1257, 1957, 224, 252, 108, 63041, 1778, 4181, 432, 52870, 51605, 433, 248, 109, 63041, 4175, 6896, 432, 57712, 433, 248, 110, 63041, 10866, 432, 53098, 25102, 433, 248, 111, 63041, 10974, 432, 61374, 24518, 433, 248, 112, 63041, 10975, 432, 53432, 51243, 25042, 433, 248, 113, 63041, 16215, 432, 55635, 25308, 27312, 433, 248, 114, 63041, 16229, 432, 50368, 24086, 29238, 433, 248, 115, 63041, 14695, 432, 58207, 24560, 154, 433, 248, 116, 63041, 10973, 432, 51273, 26210, 25122, 163, 161, 433, 248, 49970, 63041, 12448, 536, 2731}
	content, _ := ioutil.ReadFile("/Users/cangxiaoze/user/git/token_dict/vocab.json")
	all := map[string]int{}
	id2Str := map[int]string{}
	err := jsoniter.Unmarshal(content, &all)
	if err != nil {
		t.Error(err)
	}
	r := trie_generic.NewTrie[int](nil)
	rand.Seed(time.Now().UnixMilli())
	for k, v := range all {
		r.Insert(k, &trie_generic.DataType[int]{Data: v})
		id2Str[v] = k
	}
	sp := trie_generic.NewTrie[int](nil)
	spToken := []string{"[SEP]"}
	for _, token := range spToken {
		val, ok := all[token]
		if ok {
			sp.Insert(token, &trie_generic.DataType[int]{Data: val})
		}
	}
	results := sp.Match3(sampleS)
	tokenStr := ""
	ids := []int{}
	unk, ok := all["[UNK]"]
	if !ok {
		return
	}
	for _, result := range results {
		if result.Data != nil {
			subResultList := r.Match3(tokenStr)
			for _, subResult := range subResultList {
				if subResult.Data != nil {
					ids = append(ids, subResult.Data.Data)
				} else {
					ids = append(ids, unk)
					//fmt.Println("fatal con`t match")
				}
			}
			ids = append(ids, result.Data.Data)
			tokenStr = ""
		} else {
			tokenStr += result.Str
		}
	}
	if tokenStr != "" {
		subResultList := r.Match3(tokenStr)
		for _, subResult := range subResultList {
			if subResult.Data != nil {
				ids = append(ids, subResult.Data.Data)
			} else {
				ids = append(ids, unk)
				//fmt.Println("fatal con`t match")
			}
		}
	}
	println(utils.MustJson(ids))
	decode := ""
	decode2 := ""
	for i, id := range ids {
		decode += id2Str[id]
		id2 := targetCodeList[i]
		decode2 += id2Str[id2]
		println("xz", id, id2Str[id], "zb", id2, id2Str[id2])
	}
	println(decode)
	println(decode2)
	//decode = ""
	//for _, id := range targetCodeList {
	//	decode += id2Str[id]
	//	println("zb", id, id2Str[id])
	//}
	//println(decode)
	time.Sleep(time.Second * 4)
}

func TestCut(t *testing.T) {
	sentence := "你好么？我叫干不倒"
	r, err := regexp.Compile("[^.。?？!！]+[.。?？!！]?")
	if err != nil {
		return
	}
	println(r.FindString(sentence))
}
