package stores

import (
	"github.com/google/uuid"
	node2 "github.com/ironpark/teatime/internal/node"
)

type NodeInfo struct {
	Id          string
	Name        string
	Type        string
	Description string
}

type Node struct {
	Id          string
	LocalId     string
	Name        string
	Description string
	Type        string
	Properties  []node2.NodeProperty
	Output      []node2.NodeProperty
}

type nodeStore struct {
}

func NewNodeStore() *nodeStore {
	return &nodeStore{}
}

func (t *nodeStore) GetNodeInfos() []NodeInfo {
	nodes := node2.GetNodes()
	nodeInfos := make([]NodeInfo, len(nodes))
	for i, node := range nodes {
		nodeInfos[i] = NodeInfo{
			Id:          node.Ref(),
			Name:        node.Name(),
			Type:        string(node.Type()),
			Description: node.Description(),
		}
	}
	return nodeInfos
}

func (t *nodeStore) GetNodeInfosByType(nodeType string) []NodeInfo {
	nodes := node2.GetNodesByType(node2.NodeType(nodeType))
	nodeInfos := make([]NodeInfo, len(nodes))
	for i, node := range nodes {
		nodeInfos[i] = NodeInfo{
			Id:          node.Ref(),
			Name:        node.Name(),
			Type:        string(node.Type()),
			Description: node.Description(),
		}
	}
	return nodeInfos
}

func (t *nodeStore) GetNodeInfo(id string) NodeInfo {
	node, err := node2.GetNodeByRef(id)
	if err != nil {
		return NodeInfo{}
	}
	return NodeInfo{
		Id:          node.Ref(),
		Name:        node.Name(),
		Type:        string(node.Type()),
		Description: node.Description(),
	}
}

func (t *nodeStore) CreateNode(nodeId string) Node {
	node, err := node2.GetNodeByRef(nodeId)
	if err != nil {
		return Node{}
	}
	return Node{
		Id:          node.Ref(),
		LocalId:     uuid.New().String(),
		Name:        node.Name(),
		Type:        string(node.Type()),
		Description: node.Description(),
		Properties:  node.Properties(),
		Output:      node.Output(),
	}
}
