package graph

import (
	"errors"
	"fmt"
)

type Node struct {
	ID    int64       `json:"id"`
	Data  interface{} `json:"data"`
	Value float64     `json:"value"`
}

type Edge struct {
	ID    int64       `json:"id"`
	S     *Node       `json:"s"`
	E     *Node       `json:"e"`
	Value float64     `json:"value"`
	Data  interface{} `json:"data"`
}

type Graph struct {
	NodeList  []*Node           `json:"nodeList"`
	NodeMap   map[int64]*Node   `json:"nodeMap"`
	EdgeMap   map[int64][]*Edge `json:"edgeMap"`
	Node2Node map[string]bool   `json:"node2Node"`
}

func (m *Graph) nodePairStr(nodeID1, nodeID2 int64) string {
	return fmt.Sprintf("%d-%d", nodeID1, nodeID2)
}

func (m *Graph) Marked(nodeID1, nodeID2 int64) bool {
	_, ok := m.Node2Node[m.nodePairStr(nodeID1, nodeID2)]
	return ok
}

func (m *Graph) Mark(nodeID1, nodeID2 int64) {
	m.Node2Node[m.nodePairStr(nodeID1, nodeID2)] = true
}

func (m *Graph) AddNode(id int64, data interface{}, value float64) (err error) {
	if len(m.NodeList) == 0 {
		m.NodeList = make([]*Node, 0, 7)
		m.NodeMap = make(map[int64]*Node, 7)
	}
	m.NodeList = append(m.NodeList, &Node{
		ID:    id,
		Data:  data,
		Value: value,
	})
	m.NodeMap[id] = m.NodeList[len(m.NodeList)-1]
	return nil
}

func (m *Graph) AddEdge(nodeID1, nodeID2 int64, data interface{}, value float64, undirected bool) (err error) {
	node1, ok1 := m.NodeMap[nodeID1]
	node2, ok2 := m.NodeMap[nodeID2]
	if !ok1 || !ok2 {
		return errors.New("error node")
	}
	edge := &Edge{
		S:     node1,
		E:     node2,
		Value: value,
		Data:  data,
	}
	edgeList, ok := m.EdgeMap[nodeID1]
	if !ok {
		edgeList = []*Edge{
			edge,
		}
		m.EdgeMap[nodeID1] = edgeList
	} else {
		m.EdgeMap[nodeID1] = append(m.EdgeMap[nodeID1], edge)
	}
	if undirected {
		edge = &Edge{
			S:     node2,
			E:     node1,
			Value: value,
			Data:  data,
		}
		edgeList, ok = m.EdgeMap[nodeID1]
		if !ok {
			edgeList = []*Edge{
				edge,
			}
			m.EdgeMap[nodeID1] = edgeList
		} else {
			m.EdgeMap[nodeID1] = append(m.EdgeMap[nodeID1], edge)
		}
	}
	return nil
}

func (m *Graph) NextEdges(nodeID int64) (edgeList []*Edge) {
	curNode, ok := m.NodeMap[nodeID]
	if !ok {
		return nil
	}
	nextEdges, ok := m.EdgeMap[curNode.ID]
	if !ok || len(nextEdges) == 0 {
		return nil
	}
	return nextEdges
}
