package node

import (
	"errors"
	"sync"
)

var (
	// triggerNodes stores all registered trigger nodes.
	triggerNodes = []Node{}
	// actionNodes stores all registered action nodes.
	actionNodes  = []Node{}
	// branchNodes stores all registered branch nodes.
	branchNodes  = []Node{}
	// utilNodes stores all registered utility nodes.
	utilNodes    = []Node{}
	// refToNode maps node reference IDs to their implementations.
	refToNode    = make(map[string]Node)
	// lock protects concurrent access to the node registry.
	lock         sync.RWMutex
)

// RegisterNode registers a new node implementation in the global registry.
// Nodes are categorized by type and stored in type-specific slices for efficient lookup.
// Node implementations should call this during package initialization.
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

// GetNodes returns all registered nodes regardless of type.
// The returned slice combines all node types in order: trigger, action, branch, util.
func GetNodes() []Node {
	lock.RLock()
	defer lock.RUnlock()
	return append(triggerNodes, append(actionNodes, append(branchNodes, utilNodes...)...)...)
}

// GetNodesByType returns all registered nodes of the specified type.
// Returns an empty slice if the node type is not recognized.
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

// GetNodeByRef retrieves a registered node by its reference ID.
// Returns an error if no node with the specified reference exists.
func GetNodeByRef(ref string) (Node, error) {
	lock.RLock()
	defer lock.RUnlock()

	if node, ok := refToNode[ref]; ok {
		return node, nil
	}
	return nil, errors.New("node not found")
}
