package stores

import (
	"github.com/google/uuid"
	node2 "github.com/ironpark/teatime/internal/node"
	"github.com/samber/lo"
)

// Node represents a workflow node instance with its configuration and runtime state.
// It embeds NodeInfo for type information and adds instance-specific data like
// unique identifiers and property values.
type Node struct {
	node2.NodeInfo
	Id            string               `json:"id"`            // Reference ID of the node type
	LocalId       string               `json:"localId"`       // Unique instance ID for this node
	Properties    []node2.NodeProperty `json:"properties"`    // Input properties with current values
	Output        []node2.NodeProperty `json:"output"`        // Output properties with current values
	OutputHandles []node2.OutputHandle `json:"outputHandles"` // Output handles with current values
}

// nodeStore manages workflow node definitions and instances.
// It provides access to registered node types and creates node instances
// with unique identifiers and default property values.
//
// The store is safe for concurrent use as it only reads from the global
// node registry maintained by the node package.
type nodeStore struct {
}

// NewNodeStore creates a new node store for managing workflow nodes.
// The returned store provides access to all registered node types and
// can create node instances with unique identifiers.
func NewNodeStore() *nodeStore {
	return &nodeStore{}
}

// GetNodeInfos returns information about all registered node types.
// This includes node names, descriptions, categories, input/output definitions,
// and other metadata needed for the workflow editor UI.
//
// The returned slice contains read-only information and is safe to modify.
func (t *nodeStore) GetNodeInfos() []node2.NodeInfo {
	nodes := node2.GetNodes()
	nodeInfos := make([]node2.NodeInfo, len(nodes))
	for i, node := range nodes {
		nodeInfos[i] = node.Info()
	}
	return nodeInfos
}

// GetNodeInfosByType returns information about nodes of the specified type.
// The nodeType parameter should be one of the valid node types (e.g., "action", "trigger", "branch").
//
// Returns an empty slice if the nodeType is not recognized or if no nodes
// of that type are registered.
func (t *nodeStore) GetNodeInfosByType(nodeType string) []node2.NodeInfo {
	nodes := node2.GetNodesByType(node2.NodeType(nodeType))
	return lo.Map(nodes, func(node node2.Node, _ int) node2.NodeInfo {
		return node.Info()
	})
}

// GetNodeInfo retrieves information for a specific node type by its reference ID.
// The id parameter should be the unique reference identifier for the node type.
//
// Returns an empty NodeInfo struct if the node ID is not found or invalid.
// This method does not return an error; callers should check if the returned
// NodeInfo has valid data.
func (t *nodeStore) GetNodeInfo(id string) node2.NodeInfo {
	node, err := node2.GetNodeByRef(id)
	if err != nil {
		return node2.NodeInfo{}
	}
	return node.Info()
}

// CreateNode creates a new node instance from the specified node type.
// It generates a unique local ID for the instance and initializes the node
// with default property values from the node type definition.
//
// The nodeId parameter should be a valid node type reference ID.
// Returns an empty Node struct if the nodeId is not found or invalid.
// Each call generates a new unique LocalId for the instance.
func (t *nodeStore) CreateNode(nodeId string) Node {
	node, err := node2.GetNodeByRef(nodeId)
	if err != nil {
		return Node{}
	}
	return Node{
		Id:            node.Ref(),
		LocalId:       uuid.New().String(),
		NodeInfo:      node.Info(),
		Properties:    node.GetProperties(node2.PropertyContext{}),
		Output:        node.GetOutput(node2.PropertyContext{}),
		OutputHandles: node.GetOutputHandles(node2.PropertyContext{}),
	}
}
