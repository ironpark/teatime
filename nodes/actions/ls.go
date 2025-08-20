package actions

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&LsActionNode{
		BaseNode: *node.NewBaseNode(
			"teatime.action.ls",
			node.NodeTypeAction,
			"ls",
			"디렉토리의 파일 및 하위 디렉토리 목록을 조회하는 액션 노드입니다.",
			"FolderOpen",
		),
	})
}

type LsActionNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *LsActionNode) Output() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Files",
			Description: "파일 목록입니다.",
			Key:         "files",
			Value:       "",
			Type:        node.StringArray,
		},
		{
			Name:        "Directories",
			Description: "디렉토리 목록입니다.",
			Key:         "directories",
			Value:       "",
			Type:        node.StringArray,
		},
		{
			Name:        "Total Count",
			Description: "전체 항목 수입니다.",
			Key:         "totalCount",
			Value:       "",
			Type:        node.Int64,
		},
		{
			Name:        "Error",
			Description: "오류 메시지입니다.",
			Key:         "error",
			Value:       "",
			Type:        node.String,
			Optional:    true,
		},
	}
}

func (r *LsActionNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "Directory Path",
			Description: "조회할 디렉토리 경로를 입력하세요",
			Optional:    false,
			Key:         "path",
			Value:       ".",
			Type:        node.String,
		},
		{
			Name:        "Recursive",
			Description: "하위 디렉토리를 포함하여 파일 목록을 출력할지 여부를 선택하세요",
			Optional:    true,
			Key:         "recursive",
			Value:       "false",
			Type:        node.Bool,
		},
		{
			Name:        "Filter Type",
			Description: "필터링할 타입을 선택하세요",
			Optional:    true,
			Key:         "filterType",
			Value:       "all",
			Type:        node.String,
			Options:     []string{"all", "files", "directories"},
		},
		{
			Name:        "Pattern",
			Description: "파일 이름 패턴 (glob 패턴 지원)",
			Optional:    true,
			Key:         "pattern",
			Value:       "*",
			Type:        node.String,
		},
		{
			Name:        "Show Hidden",
			Description: "숨김 파일 표시 여부",
			Optional:    true,
			Key:         "showHidden",
			Value:       "false",
			Type:        node.Bool,
		},
		{
			Name:        "Sort By",
			Description: "정렬 기준",
			Optional:    true,
			Key:         "sortBy",
			Value:       "name",
			Type:        node.String,
			Options:     []string{"name", "size", "modified", "created"},
		},
		{
			Name:        "Sort Order",
			Description: "정렬 순서",
			Optional:    true,
			Key:         "sortOrder",
			Value:       "asc",
			Type:        node.String,
			Options:     []string{"asc", "desc"},
		},
		{
			Name:        "Max Depth",
			Description: "최대 탐색 깊이 (재귀 탐색 시)",
			Optional:    true,
			Key:         "maxDepth",
			Value:       "1",
			Type:        node.Int64,
		},
	}
}

func (r *LsActionNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *LsActionNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
