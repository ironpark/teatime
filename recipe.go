package main

type NodeType string
type PropertyType int

const (
	// Basic Types
	Invalid PropertyType = iota
	Bool
	Int64
	Uint64
	Float64
	String
	// Special Types
	JSON
	XML
	Date
	// Arrays
	TextArray
	NumberArray
	BooleanArray
	JSONArray
	XMLArray
)

const (
	NodeTypeTrigger NodeType = "trigger"
	NodeTypeBranch  NodeType = "branch"
	NodeTypeAction  NodeType = "action"
	NodeTypeUtil    NodeType = "util"
)

// 레시피는 로컬 자동화 워크플로우 프로젝트를 의미한다.
type Recipe struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// 레시피에 포함된 노드들
	Nodes []Node `json:"nodes"`
	// 노드 간 연결 관계
	Edges []Edge `json:"edges"`
}

type Node struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data struct {
		Label      string         `json:"label"`
		NodeType   string         `json:"nodeType"`
		Properties []NodeProperty `json:"properties"`
	} `json:"data"`
}

type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type NodeProperty struct {
	// Type of the property (string, number, boolean, array, json, xml)
	Type PropertyType `json:"type"`
	// Name of the property
	Name string `json:"name"`
	// Description of the property
	Description string `json:"description"`
	// Whether the property is required
	Required bool `json:"required"`
	// Key of the property
	Key string `json:"key"`
	// Value of the property
	Value any `json:"value"`
	// not editable
	ReadOnly bool `json:"readOnly"`
}

func (r *Recipe) GetTriggerNodes() (triggerNodes []Node) {
	triggerNodes = []Node{}
	for _, node := range r.Nodes {
		if node.Type == string(NodeTypeTrigger) {
			triggerNodes = append(triggerNodes, node)
		}
	}
	return triggerNodes
}
