package stores

type Store struct {
	NodeStore
}

func NewStore() *Store {
	return &Store{
		NodeStore: NewNodeStore(),
	}
}

type NodeStore interface {
	GetNodeInfos() []NodeInfo
	GetNodeInfosByType(nodeType string) []NodeInfo
	GetNodeInfo(id string) NodeInfo
	CreateNode(nodeId string) Node
}
