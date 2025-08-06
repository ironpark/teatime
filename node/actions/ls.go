package trigger

import (
	"github.com/ironpark/teatime/node"
	"github.com/ironpark/teatime/node/types"
)

func init() {
	node.RegisterNode(&LsActionNode{})
}

type LsActionNode struct {
	customParams []types.NodeProperty
}

func (r *LsActionNode) Name() string {
	return "ls"
}

func (r *LsActionNode) Description() string {
	return "디렉토리의 파일 및 하위 디렉토리 목록을 조회하는 액션 노드입니다."
}

func (r *LsActionNode) Type() types.NodeType {
	return types.NodeTypeAction
}

func (r *LsActionNode) ID() string {
	return "teatime.action.ls"
}

func (r *LsActionNode) Output() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Files",
			Description: "파일 목록입니다.",
			Key:         "files",
			Value:       "",
			Type:        types.JSONArray,
		},
		{
			Name:        "Directories",
			Description: "디렉토리 목록입니다.",
			Key:         "directories",
			Value:       "",
			Type:        types.JSONArray,
		},
		{
			Name:        "Total Count",
			Description: "전체 항목 수입니다.",
			Key:         "totalCount",
			Value:       "",
			Type:        types.Int64,
		},
		{
			Name:        "Error",
			Description: "오류 메시지입니다.",
			Key:         "error",
			Value:       "",
			Type:        types.Text,
			Optional:    true,
		},
	}
}

func (r *LsActionNode) Properties() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "Directory Path",
			Description: "조회할 디렉토리 경로를 입력하세요",
			Optional:    false,
			Key:         "path",
			Value:       ".",
			Type:        types.String,
		},
		{
			Name:        "Recursive",
			Description: "하위 디렉토리를 포함하여 파일 목록을 출력할지 여부를 선택하세요",
			Optional:    true,
			Key:         "recursive",
			Value:       "false",
			Type:        types.Bool,
		},
		{
			Name:        "Filter Type",
			Description: "필터링할 타입을 선택하세요",
			Optional:    true,
			Key:         "filterType",
			Value:       "all",
			Type:        types.String,
			Options:     []string{"all", "files", "directories"},
		},
		{
			Name:        "Pattern",
			Description: "파일 이름 패턴 (glob 패턴 지원)",
			Optional:    true,
			Key:         "pattern",
			Value:       "*",
			Type:        types.String,
		},
		{
			Name:        "Show Hidden",
			Description: "숨김 파일 표시 여부",
			Optional:    true,
			Key:         "showHidden",
			Value:       "false",
			Type:        types.Bool,
		},
		{
			Name:        "Sort By",
			Description: "정렬 기준",
			Optional:    true,
			Key:         "sortBy",
			Value:       "name",
			Type:        types.String,
			Options:     []string{"name", "size", "modified", "created"},
		},
		{
			Name:        "Sort Order",
			Description: "정렬 순서",
			Optional:    true,
			Key:         "sortOrder",
			Value:       "asc",
			Type:        types.String,
			Options:     []string{"asc", "desc"},
		},
		{
			Name:        "Max Depth",
			Description: "최대 탐색 깊이 (재귀 탐색 시)",
			Optional:    true,
			Key:         "maxDepth",
			Value:       "1",
			Type:        types.Int64,
		},
	}
}

func (r *LsActionNode) CustomParams() []types.NodeProperty {
	return r.customParams
}

func (r *LsActionNode) AddCustomParam(param types.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
