package actions

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&SaveActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.save",
			node.NodeTypeAction,
			"Save",
			"파일을 저장하는 액션 노드입니다.",
			"Save",
			[]node.NodeProperty{
				node.StringProp("filePath", "File Path",
					node.WithDescription("저장할 파일 경로를 입력하세요"),
					node.Required(),
				),
				node.StringProp("content", "Content",
					node.WithDescription("저장할 내용을 입력하세요"),
					node.TextArea(10),
					node.Required(),
				),
				node.SelectProp("encoding", "Encoding", []string{"utf-8", "utf-16", "ascii", "base64"},
					node.WithDescription("파일 인코딩"),
					node.OptionalWithDefault("utf-8"),
				),
				node.BoolProp("append", "Append",
					node.WithDescription("파일에 추가할지 여부"),
					node.OptionalWithDefault(false),
				),
				node.BoolProp("createDir", "Create Directory",
					node.WithDescription("디렉토리가 없으면 생성할지 여부"),
					node.OptionalWithDefault(true),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "savedPath", "File Path",
					node.WithDescription("저장된 파일 경로입니다."),
				),
				node.OutputProp(node.Int64, "fileSize", "File Size",
					node.WithDescription("파일 크기(바이트)입니다."),
				),
				node.OutputProp(node.Bool, "success", "Success",
					node.WithDescription("저장 성공 여부입니다."),
				),
			},
			nil, // Use default output handle
		),
	})
}

// 파일을 저장하는 액션 노드
type SaveActionNode struct {
	node.BaseNode
}
