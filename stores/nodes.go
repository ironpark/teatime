package stores

import (
	"github.com/google/uuid"
	"github.com/ironpark/teatime/node"
	"github.com/ironpark/teatime/node/types"
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
	Properties  []types.NodeProperty
	Output      []types.NodeProperty
}

type nodeStore struct {
}

func NewNodeStore() *nodeStore {
	return &nodeStore{}
}

func (t *nodeStore) GetNodeInfos() []NodeInfo {
	nodes := node.GetNodes()
	nodeInfos := make([]NodeInfo, len(nodes))
	for i, node := range nodes {
		nodeInfos[i] = NodeInfo{
			Id:          node.ID(),
			Name:        node.Name(),
			Type:        string(node.Type()),
			Description: node.Description(),
		}
	}
	return nodeInfos
}

func (t *nodeStore) GetNodeInfosByType(nodeType string) []NodeInfo {
	nodes := node.GetNodesByType(types.NodeType(nodeType))
	nodeInfos := make([]NodeInfo, len(nodes))
	for i, node := range nodes {
		nodeInfos[i] = NodeInfo{
			Id:          node.ID(),
			Name:        node.Name(),
			Type:        string(node.Type()),
			Description: node.Description(),
		}
	}
	return nodeInfos
}

func (t *nodeStore) GetNodeInfo(id string) NodeInfo {
	node := node.GetNode(id)
	return NodeInfo{
		Id:          node.ID(),
		Name:        node.Name(),
		Type:        string(node.Type()),
		Description: node.Description(),
	}
}

func (t *nodeStore) CreateNode(nodeId string) Node {
	node := node.GetNode(nodeId)
	return Node{
		Id:          node.ID(),
		LocalId:     uuid.New().String(),
		Name:        node.Name(),
		Type:        string(node.Type()),
		Description: node.Description(),
		Properties:  node.Properties(),
		Output:      node.Output(),
	}
}
