package stores

import (
	"github.com/google/uuid"
	node2 "github.com/ironpark/teatime/internal/node"
	"github.com/samber/lo"
)

type Node struct {
	node2.NodeInfo
	Id         string               `json:"id"`
	LocalId    string               `json:"localId"`
	Properties []node2.NodeProperty `json:"properties"`
	Output     []node2.NodeProperty `json:"output"`
}

type nodeStore struct {
}

func NewNodeStore() *nodeStore {
	return &nodeStore{}
}

func (t *nodeStore) GetNodeInfos() []node2.NodeInfo {
	nodes := node2.GetNodes()
	nodeInfos := make([]node2.NodeInfo, len(nodes))
	for i, node := range nodes {
		nodeInfos[i] = node.Info()
	}
	return nodeInfos
}

func (t *nodeStore) GetNodeInfosByType(nodeType string) []node2.NodeInfo {
	nodes := node2.GetNodesByType(node2.NodeType(nodeType))
	nodeInfos := make([]node2.NodeInfo, len(nodes))
	for i, node := range nodes {
		nodeInfos[i] = node.Info()
	}
	return lo.Map(nodes, func(node node2.Node, _ int) node2.NodeInfo {
		return node.Info()
	})
}

func (t *nodeStore) GetNodeInfo(id string) node2.NodeInfo {
	node, err := node2.GetNodeByRef(id)
	if err != nil {
		return node2.NodeInfo{}
	}
	return node.Info()
}

func (t *nodeStore) CreateNode(nodeId string) Node {
	node, err := node2.GetNodeByRef(nodeId)
	if err != nil {
		return Node{}
	}
	return Node{
		Id:         node.Ref(),
		LocalId:    uuid.New().String(),
		NodeInfo:   node.Info(),
		Properties: node.Properties(),
		Output:     node.Output(),
	}
}
