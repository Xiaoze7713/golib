package trie

import (
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
			node.AddNode(string(c), &Node{isLeaf: isLeaf, leafData: data, nextMap: map[string]*Node{}})
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

func (m *DATrie) maxPrefix(word string, startIdx int, ignoreCase bool) (result *Result) {
	idx := startIdx
	node := m.root
	results := []*Result{nil}
	s := []rune(word)
	for node != nil && idx < len(s) {
		c := s[idx]
		nodeKey := string(c)
		//println(nodeKey)
		//println(utils.MustJson(node.nextMap))
		if ignoreCase {
			if node.HasNext(string(unicode.ToUpper(c))) {
				nodeKey = string(unicode.ToUpper(c))
			} else if node.HasNext(string(unicode.ToLower(c))) {
				nodeKey = string(unicode.ToLower(c))
			} else {
				return results[len(results)-1]
			}
		} else {
			if !node.HasNext(nodeKey) {
				return results[len(results)-1]
			}
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
		matchResult := m.maxPrefix(content, i, false)
		if matchResult != nil {
			resultList = append(resultList, matchResult)
			i += len([]rune(matchResult.Str))
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
		matchResult := m.maxPrefix(content, i, true)
		if matchResult != nil {
			resultList = append(resultList, matchResult)
			i += len([]rune(matchResult.Str))
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
