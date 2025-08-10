package actions

import (
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&SaveActionNode{
		BaseNode: *node.NewBaseNode("teatime.action.save", node.NodeTypeAction, "Save", "파일을 저장하는 액션 노드입니다.", "Save"),
	})
}

// 파일을 저장하는 액션 노드
type SaveActionNode struct {
	node.BaseNode
	customParams []node.NodeProperty
}

func (r *SaveActionNode) Output() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "File Path",
			Description: "저장된 파일 경로입니다.",
			Key:         "savedPath",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "File Size",
			Description: "파일 크기(바이트)입니다.",
			Key:         "fileSize",
			Value:       "",
			Type:        node.Int64,
		},
		{
			Name:        "Success",
			Description: "저장 성공 여부입니다.",
			Key:         "success",
			Value:       "",
			Type:        node.Bool,
		},
	}
}

func (r *SaveActionNode) Properties() []node.NodeProperty {
	return []node.NodeProperty{
		{
			Name:        "File Path",
			Description: "저장할 파일 경로를 입력하세요",
			Optional:    false,
			Key:         "filePath",
			Value:       "",
			Type:        node.String,
		},
		{
			Name:        "Content",
			Description: "저장할 내용을 입력하세요",
			Optional:    false,
			Key:         "content",
			Value:       "",
			Type:        node.Text,
		},
		{
			Name:        "Encoding",
			Description: "파일 인코딩",
			Optional:    true,
			Key:         "encoding",
			Value:       "utf-8",
			Type:        node.String,
			Options:     []string{"utf-8", "utf-16", "ascii", "base64"},
		},
		{
			Name:        "Append",
			Description: "파일에 추가할지 여부",
			Optional:    true,
			Key:         "append",
			Value:       "false",
			Type:        node.Bool,
		},
		{
			Name:        "Create Directory",
			Description: "디렉토리가 없으면 생성할지 여부",
			Optional:    true,
			Key:         "createDir",
			Value:       "true",
			Type:        node.Bool,
		},
	}
}

func (r *SaveActionNode) CustomParams() []node.NodeProperty {
	return r.customParams
}

func (r *SaveActionNode) AddCustomParam(param node.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
