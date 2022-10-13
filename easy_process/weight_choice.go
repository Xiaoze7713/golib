package easy_process

import (
	"errors"
	"math/rand"
	"sort"
	"time"
)

type DataAny interface {
	any
}

type ProbEvent[T DataAny] struct {
	Prob float64     `json:"prob"`
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

func (m *ProbEvent[T]) Hit() (prob float64, hit bool) {
	rand.Seed(time.Now().UnixMilli())
	p := rand.Float64()
	return p, p <= m.Prob
}

func (m *ProbEvent[T]) GetData() (data T) {
	return m.Data.(T)
}

func NewProbEvent[T DataAny](prob float64, name string, data interface{}) (pe *ProbEvent[T]) {
	return &ProbEvent[T]{
		Prob: prob,
		Name: name,
		Data: data,
	}
}

type Reply interface {
	Weight() int64
	String() string
}

type SimpleReply struct {
	weight  int64
	content string
}

func NewSimpleReply(weight int64, content string) *SimpleReply {
	return &SimpleReply{
		weight:  weight,
		content: content,
	}
}

func (m *SimpleReply) String() string {
	return m.content
}

func (m *SimpleReply) Weight() int64 {
	return m.weight
}

type AddedWeightReply struct {
	Left  int64
	Right int64
	Ele   Reply
}

func (m *AddedWeightReply) Contain(value int64) bool {
	return value >= m.Left && value < m.Right
}

func (m *AddedWeightReply) Less(value int64) bool {
	return m.Right <= value
}

func (m *AddedWeightReply) More(value int64) bool {
	return m.Left > value
}

type WeightChoice struct {
	choices     []*AddedWeightReply
	totalWeight int64
}

func (m *WeightChoice) Load(replyList []Reply) {
	var tmpChoice []*AddedWeightReply
	for _, reply := range replyList {
		tmpChoice = append(tmpChoice, &AddedWeightReply{
			Left:  m.totalWeight,
			Right: m.totalWeight + reply.Weight(),
			Ele:   reply,
		})
		m.totalWeight += reply.Weight()
	}
	sort.Slice(tmpChoice, func(i, j int) bool {
		return tmpChoice[i].Left < tmpChoice[j].Left
	})
	m.choices = tmpChoice
}

func (m *WeightChoice) find(roll int64) Reply {
	for _, choice := range m.choices {
		if choice.Contain(roll) {
			//fmt.Printf("%v %v %v", roll, choice.Left, choice.Right)
			return choice.Ele
		}
	}
	return nil
}

func (m *WeightChoice) Choice() (res Reply, err error) {
	roll := rand.Int63n(m.totalWeight)
	reply := m.find(roll)
	if reply == nil {
		return nil, errors.New("value not found")
	}
	return reply, nil
}
