package trie

import (
	"git.singularity-ai.com/backend/library/utils"
	"sort"
	"unicode"
)

type Node struct {
	isLeaf   bool
	leafData interface{}
	tail     string
	nextMap  map[string]*Node
}

func (m *Node) IsLeaf() bool {
	return m.isLeaf
}

func (m *Node) SetLeaf() {
	m.isLeaf = true
}

func (m *Node) HasNext(w string) bool {
	if _, ok := m.nextMap[w]; ok {
		return true
	}
	return false
}

func (m *Node) AddNode(w string, node *Node) {
	m.nextMap[w] = node
}

func (m *Node) GetNode(w string) (n *Node) {
	n, ok := m.nextMap[w]
	if ok {
		return n
	}
	return nil
}

func (m *Node) GetTail() string {
	return m.tail
}

func (m *Node) GetData() interface{} {
	return m.leafData
}

func (m *Node) SetData(data interface{}) {
	m.leafData = data
}

type DATrie struct {
	root *Node
}

func NewTrie() *DATrie {
	return &DATrie{root: &Node{nextMap: map[string]*Node{}}}
}

type Result struct {
	Str  string
	Data interface{}
}

func (m *DATrie) Insert(word string, data interface{}) {
	idx := 0
	node := m.root
	s := []rune(word)
	for idx < len(s) {
		c := s[idx]
		isLeaf := idx == len(s)-1
		if !node.HasNext(string(c)) {
			nodeData := data
			if !isLeaf {
				nodeData = nil
			}
			node.AddNode(string(c), &Node{isLeaf: isLeaf, leafData: nodeData, nextMap: map[string]*Node{}})
		}
		node = node.GetNode(string(c))
		idx += 1
	}
	if !node.IsLeaf() {
		node.SetLeaf()
		node.SetData(data)
	}
}

func (m *DATrie) FindStrict(word string) interface{} {
	idx := 0
	node := m.root
	s := []rune(word)
	for node != nil && idx < len(s) {
		c := s[idx]
		if !node.HasNext(string(c)) {
			return nil
		}
		node = node.GetNode(string(c))
		idx += 1
	}
	if node.IsLeaf() {
		return node.GetData()
	}
	return nil
}

func (m *DATrie) prefix(word string) (results []*Result) {
	idx := 0
	node := m.root
	results = []*Result{}
	s := []rune(word)
	for node != nil && idx < len(s) {
		c := s[idx]
		if !node.HasNext(string(c)) {
			return results
		}
		node = node.GetNode(string(c))
		if node.IsLeaf() {
			results = append(results, &Result{
				Str:  string(s[:idx+1]),
				Data: node.GetData(),
			})
		}
		idx += 1
	}
	return results
}

func (m *DATrie) deepSearch(node *Node, sentence string, startIdx int, idx int, ignoreCase bool) (resultList []*Result) {
	if node == nil {
		return nil
	}
	runeSentence := []rune(sentence)
	//println("sentence", string(runeSentence), startIdx, idx, string(runeSentence[startIdx:idx+1]), utils.MustJson(node.nextMap), "leaf", node.isLeaf, "data", utils.MustJson(node.leafData))
	if node.IsLeaf() {
		resultList = append(resultList, &Result{
			Str:  string(runeSentence[startIdx:idx]),
			Data: node.GetData(),
		})
	}
	if idx >= len([]rune(sentence)) {
		return resultList
	}
	c := runeSentence[idx]
	nodeKey := string(c)
	//println("node [", nodeKey, "]")
	//println(nodeKey)
	// 优先保持大小写
	if !ignoreCase {
		if !node.HasNext(nodeKey) {
			return resultList
		}
		searchList := m.deepSearch(node.GetNode(nodeKey), sentence, startIdx, idx+1, ignoreCase)
		if len(searchList) > 0 {
			resultList = append(resultList, searchList...)
		}
	} else if ignoreCase {
		upKey := string(unicode.ToUpper(c))
		lowerKey := string(unicode.ToLower(c))
		//println(upKey, lowerKey)
		if node.HasNext(upKey) {
			//println("upper", upKey)
			searchList := m.deepSearch(node.GetNode(upKey), sentence, startIdx, idx+1, ignoreCase)
			if len(searchList) > 0 {
				resultList = append(resultList, searchList...)
			}
		}
		if node.HasNext(lowerKey) {
			//println("lower", lowerKey)
			searchList := m.deepSearch(node.GetNode(lowerKey), sentence, startIdx, idx+1, ignoreCase)
			if len(searchList) > 0 {
				resultList = append(resultList, searchList...)
			}
		}
	}
	return resultList
}

func (m *DATrie) maxPrefix(word string, startIdx int, ignoreCase bool) (result *Result) {
	idx := startIdx
	node := m.root
	results := []*Result{nil}
	s := []rune(word)
	println(string(s[idx:]))
	for node != nil && idx < len(s) {
		// c sentence word
		c := s[idx]
		nodeKey := string(c)
		//println("node [", nodeKey, "]")
		//println(nodeKey)
		println(utils.MustJson(node.nextMap))
		if ignoreCase {
			// 优先保持大小写
			if node.HasNext(nodeKey) {

			} else {
				upKey := string(unicode.ToUpper(c))
				lowerKey := string(unicode.ToLower(c))
				if node.HasNext(upKey) {
					nodeKey = upKey
				} else if node.HasNext(lowerKey) {
					nodeKey = lowerKey
				}
			}
		}
		println(nodeKey)
		if !node.HasNext(nodeKey) {
			return results[len(results)-1]
		}
		//println(nodeKey)
		node = node.GetNode(nodeKey)
		if node.IsLeaf() {
			results = append(results, &Result{
				Str:  string(s[startIdx : idx+1]),
				Data: node.GetData(),
			})
		}
		idx += 1
	}
	return results[len(results)-1]
}

type ScoreItem interface {
	Score() float64
}

func (m *DATrie) maxScore(word string, startIdx int) (result *Result) {
	idx := startIdx
	node := m.root
	results := []*Result{nil}
	s := []rune(word)
	for node != nil && idx < len(s) {
		c := s[idx]
		if !node.HasNext(string(c)) {
			return results[len(results)-1]
		}
		node = node.GetNode(string(c))
		if node.IsLeaf() {
			results = append(results, &Result{
				Str:  string(s[:idx+1]),
				Data: node.GetData(),
			})
		}
		idx += 1
	}
	sort.Slice(results, func(i, j int) bool {
		item1, ok1 := results[i].Data.(ScoreItem)
		s1 := 0.0
		if ok1 {
			s1 = item1.Score()
		}
		s2 := 0.0
		item2, ok2 := results[j].Data.(ScoreItem)
		if ok2 {
			s2 = item2.Score()
		}
		return s1 < s2
	})
	return results[len(results)-1]
}

func (m *DATrie) Match3(content string) (result []*Result) {
	runeContent := []rune(content)
	l := len(runeContent)
	i := 0
	resultList := []*Result{}
	for i < l {
		matchResults := m.deepSearch(m.root, content, i, i, false)
		if len(matchResults) > 0 {
			sort.Slice(matchResults, func(i, j int) bool {
				return len([]rune(matchResults[i].Str)) < len([]rune(matchResults[j].Str))
			})
			res := matchResults[len(matchResults)-1]
			//println("match", res.Str, utils.MustJson(res.RespData))
			resultList = append(resultList, res)
			i += len([]rune(res.Str))
		} else {
			resultList = append(resultList, &Result{
				Str:  string(runeContent[i]),
				Data: nil,
			})
			i += 1
		}
	}
	return resultList
}
func (m *DATrie) MatchIgc(content string) (result []*Result) {
	runeContent := []rune(content)
	l := len(runeContent)
	i := 0
	resultList := []*Result{}
	for i < l {
		//println("igc start ", content, i)
		matchResults := m.deepSearch(m.root, content, i, i, true)
		if len(matchResults) > 0 {
			sort.Slice(matchResults, func(i, j int) bool {
				return len([]rune(matchResults[i].Str)) < len([]rune(matchResults[j].Str))
			})
			res := matchResults[len(matchResults)-1]
			//println("match", res.Str, utils.MustJson(res.RespData))
			resultList = append(resultList, res)
			i += len([]rune(res.Str))
		} else {
			resultList = append(resultList, &Result{
				Str:  string(runeContent[i]),
				Data: nil,
			})
			i += 1
		}
	}
	return resultList
}
