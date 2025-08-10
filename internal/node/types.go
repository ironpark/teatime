package node

type NodeType string

const (
	NodeTypeTrigger NodeType = "trigger"
	NodeTypeBranch  NodeType = "branch"
	NodeTypeAction  NodeType = "action"
	NodeTypeUtil    NodeType = "util"
)

type NodeInfo struct {
	Ref         string   `json:"ref"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        NodeType `json:"type"`
	Icon        string   `json:"icon"`
}

type BaseNode struct {
	nodeInfo NodeInfo
}

func getDefaultIcon(nodeType NodeType) string {
	switch nodeType {
	case NodeTypeTrigger:
		return "Zap"
	case NodeTypeBranch:
		return "GitBranch"
	case NodeTypeAction:
		return "Play"
	case NodeTypeUtil:
		return "Settings"
	default:
		return "Activity"
	}
}
func NewBaseNode(ref string, nodeType NodeType, name string, description string, icon string) *BaseNode {
	nodeInfo := NodeInfo{
		Ref:         ref,
		Type:        nodeType,
		Name:        name,
		Description: description,
		Icon:        getDefaultIcon(nodeType),
	}
	if icon != "" {
		nodeInfo.Icon = icon
	}
	return &BaseNode{
		nodeInfo: nodeInfo,
	}
}

func (r *BaseNode) Ref() string {
	return r.nodeInfo.Ref
}

func (r *BaseNode) Name() string {
	return r.nodeInfo.Name
}

func (r *BaseNode) Icon() string {
	if r.nodeInfo.Icon == "" {
		switch r.nodeInfo.Type {
		case NodeTypeTrigger:
			return "Zap"
		case NodeTypeBranch:
			return "GitBranch"
		case NodeTypeAction:
			return "Play"
		case NodeTypeUtil:
			return "Settings"
		default:
			return "Activity"
		}
	}
	return r.nodeInfo.Icon
}

func (r *BaseNode) Type() NodeType {
	return r.nodeInfo.Type
}

func (r *BaseNode) Description() string {
	return r.nodeInfo.Description
}

func (r *BaseNode) Info() NodeInfo {
	return r.nodeInfo
}

// Node 모든 노드가 구현해야 하는 기본 인터페이스
type Node interface {
	// 노드의 고유 식별자
	Ref() string
	// 노드 이름
	Name() string
	// 노드 타입 반환
	Type() NodeType
	// 노드 아이콘 반환
	Icon() string
	// 노드 설명
	Description() string
	// 노드의 입력 속성 정의
	Properties() []NodeProperty
	// 노드의 출력 속성 정의
	Output() []NodeProperty
	// 노드 정보 반환
	Info() NodeInfo
}
