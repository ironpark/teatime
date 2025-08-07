package node

type NodeType string

const (
	NodeTypeTrigger NodeType = "trigger"
	NodeTypeBranch  NodeType = "branch"
	NodeTypeAction  NodeType = "action"
	NodeTypeUtil    NodeType = "util"
)

type BaseNode struct {
	ref         string
	nodeType    NodeType
	name        string
	description string
}

func NewBaseNode(ref string, nodeType NodeType, name string, description string) *BaseNode {
	return &BaseNode{
		ref:         ref,
		nodeType:    nodeType,
		name:        name,
		description: description,
	}
}

type NodeInfo struct {
	Ref         string   `json:"ref"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        NodeType `json:"type"`
}

func (r *BaseNode) Ref() string {
	return r.ref
}

func (r *BaseNode) Name() string {
	return r.name
}

func (r *BaseNode) Type() NodeType {
	return r.nodeType
}

func (r *BaseNode) Description() string {
	return r.description
}

func (r *BaseNode) Info() NodeInfo {
	return NodeInfo{
		Ref:         r.ref,
		Name:        r.name,
		Description: r.description,
		Type:        r.nodeType,
	}
}

// Node 모든 노드가 구현해야 하는 기본 인터페이스
type Node interface {
	// 노드의 고유 식별자
	Ref() string
	// 노드 이름
	Name() string
	// 노드 타입 반환
	Type() NodeType
	// 노드 설명
	Description() string
	// 노드의 입력 속성 정의
	Properties() []NodeProperty
	// 노드의 출력 속성 정의
	Output() []NodeProperty
	// 노드 정보 반환
	Info() NodeInfo
}
