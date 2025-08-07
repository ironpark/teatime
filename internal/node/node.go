package node

import (
	"errors"
	"sync"
)

var (
	triggerNodes = []Node{}
	actionNodes  = []Node{}
	branchNodes  = []Node{}
	utilNodes    = []Node{}
	refToNode    = make(map[string]Node)
	lock         sync.RWMutex
)

func RegisterNode(node Node) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := refToNode[node.Ref()]; ok {
		return
	}
	switch node.Type() {
	case NodeTypeTrigger:
		triggerNodes = append(triggerNodes, node)
	case NodeTypeAction:
		actionNodes = append(actionNodes, node)
	case NodeTypeBranch:
		branchNodes = append(branchNodes, node)
	case NodeTypeUtil:
		utilNodes = append(utilNodes, node)
	}
	refToNode[node.Ref()] = node
}

func GetNodes() []Node {
	lock.RLock()
	defer lock.RUnlock()
	return append(triggerNodes, append(actionNodes, append(branchNodes, utilNodes...)...)...)
}

func GetNodesByType(nodeType NodeType) []Node {
	lock.RLock()
	defer lock.RUnlock()
	switch nodeType {
	case NodeTypeTrigger:
		return triggerNodes
	case NodeTypeAction:
		return actionNodes
	case NodeTypeBranch:
		return branchNodes
	case NodeTypeUtil:
		return utilNodes
	}
	return []Node{}
}

func GetNodeByRef(ref string) (Node, error) {
	lock.RLock()
	defer lock.RUnlock()

	if node, ok := refToNode[ref]; ok {
		return node, nil
	}
	return nil, errors.New("node not found")
}
