package actions

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&LsActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.ls",
			node.NodeTypeAction,
			"ls",
			"디렉토리의 파일 및 하위 디렉토리 목록을 조회하는 액션 노드입니다.",
			"FolderOpen",
			[]node.NodeProperty{
				node.StringProp("path", "Directory Path",
					node.WithDescription("조회할 디렉토리 경로를 입력하세요"),
					node.RequiredWithDefault("."),
				),
				node.BoolProp("recursive", "Recursive",
					node.WithDescription("하위 디렉토리를 포함하여 파일 목록을 출력할지 여부를 선택하세요"),
					node.OptionalWithDefault(false),
				),
				node.SelectProp("filterType", "Filter Type", []string{"all", "files", "directories"},
					node.WithDescription("필터링할 타입을 선택하세요"),
					node.OptionalWithDefault("all"),
				),
				node.StringProp("pattern", "Pattern",
					node.WithDescription("파일 이름 패턴 (glob 패턴 지원)"),
					node.OptionalWithDefault("*"),
				),
				node.BoolProp("showHidden", "Show Hidden",
					node.WithDescription("숨김 파일 표시 여부"),
					node.OptionalWithDefault(false),
				),
				node.SelectProp("sortBy", "Sort By", []string{"name", "size", "modified", "created"},
					node.WithDescription("정렬 기준"),
					node.OptionalWithDefault("name"),
				),
				node.SelectProp("sortOrder", "Sort Order", []string{"asc", "desc"},
					node.WithDescription("정렬 순서"),
					node.OptionalWithDefault("asc"),
				),
				node.IntProp("maxDepth", "Max Depth",
					node.WithDescription("최대 탐색 깊이 (재귀 탐색 시)"),
					node.OptionalWithDefault(int64(1)),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.StringArray, "files", "Files",
					node.WithDescription("파일 목록입니다."),
				),
				node.OutputProp(node.StringArray, "directories", "Directories",
					node.WithDescription("디렉토리 목록입니다."),
				),
				node.OutputProp(node.Int64, "totalCount", "Total Count",
					node.WithDescription("전체 항목 수입니다."),
				),
				node.OutputProp(node.String, "error", "Error",
					node.WithDescription("오류 메시지입니다."),
				),
			},
			nil, // Use default output handle
		),
	})
}

type LsActionNode struct {
	node.BaseNode
}