package node

import (
	"sync"

	"github.com/ironpark/teatime/node/types"
)

var (
	TriggerNodes = []types.Node{}
	ActionNodes  = []types.Node{}
	BranchNodes  = []types.Node{}
	UtilNodes    = []types.Node{}

	lock sync.RWMutex
)

func RegisterNode(node types.Node) {
	lock.Lock()
	defer lock.Unlock()

	switch node.Type() {
	case types.NodeTypeTrigger:
		TriggerNodes = append(TriggerNodes, node)
	case types.NodeTypeAction:
		ActionNodes = append(ActionNodes, node)
	case types.NodeTypeBranch:
		BranchNodes = append(BranchNodes, node)
	case types.NodeTypeUtil:
		UtilNodes = append(UtilNodes, node)
	}
}

func GetNodes() []types.Node {
	lock.RLock()
	defer lock.RUnlock()

	return append(TriggerNodes, append(ActionNodes, append(BranchNodes, UtilNodes...)...)...)
}

func GetNodesByType(nodeType types.NodeType) []types.Node {
	lock.RLock()
	defer lock.RUnlock()
	switch nodeType {
	case types.NodeTypeTrigger:
		return TriggerNodes
	case types.NodeTypeAction:
		return ActionNodes
	case types.NodeTypeBranch:
		return BranchNodes
	case types.NodeTypeUtil:
		return UtilNodes
	}
	return []types.Node{}
}

func GetNode(id string) types.Node {
	lock.RLock()
	defer lock.RUnlock()

	for _, node := range GetNodes() {
		if node.ID() == id {
			return node
		}
	}
	return nil
}
