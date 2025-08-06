package trigger

import (
	"github.com/ironpark/teatime/node"
	"github.com/ironpark/teatime/node/types"
)

func init() {
	node.RegisterNode(&SaveActionNode{})
}

// 파일을 저장하는 액션 노드
type SaveActionNode struct {
	customParams []types.NodeProperty
}

func (r *SaveActionNode) Name() string {
	return "Save"
}

func (r *SaveActionNode) Description() string {
	return "파일을 저장하는 액션 노드입니다."
}

func (r *SaveActionNode) Type() types.NodeType {
	return types.NodeTypeAction
}

func (r *SaveActionNode) ID() string {
	return "teatime.action.save"
}

func (r *SaveActionNode) Output() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "File Path",
			Description: "저장된 파일 경로입니다.",
			Key:         "savedPath",
			Value:       "",
			Type:        types.String,
		},
		{
			Name:        "File Size",
			Description: "파일 크기(바이트)입니다.",
			Key:         "fileSize",
			Value:       "",
			Type:        types.Int64,
		},
		{
			Name:        "Success",
			Description: "저장 성공 여부입니다.",
			Key:         "success",
			Value:       "",
			Type:        types.Bool,
		},
	}
}

func (r *SaveActionNode) Properties() []types.NodeProperty {
	return []types.NodeProperty{
		{
			Name:        "File Path",
			Description: "저장할 파일 경로를 입력하세요",
			Optional:    false,
			Key:         "filePath",
			Value:       "",
			Type:        types.String,
		},
		{
			Name:        "Content",
			Description: "저장할 내용을 입력하세요",
			Optional:    false,
			Key:         "content",
			Value:       "",
			Type:        types.Text,
		},
		{
			Name:        "Encoding",
			Description: "파일 인코딩",
			Optional:    true,
			Key:         "encoding",
			Value:       "utf-8",
			Type:        types.String,
			Options:     []string{"utf-8", "utf-16", "ascii", "base64"},
		},
		{
			Name:        "Append",
			Description: "파일에 추가할지 여부",
			Optional:    true,
			Key:         "append",
			Value:       "false",
			Type:        types.Bool,
		},
		{
			Name:        "Create Directory",
			Description: "디렉토리가 없으면 생성할지 여부",
			Optional:    true,
			Key:         "createDir",
			Value:       "true",
			Type:        types.Bool,
		},
	}
}

func (r *SaveActionNode) CustomParams() []types.NodeProperty {
	return r.customParams
}

func (r *SaveActionNode) AddCustomParam(param types.NodeProperty) {
	r.customParams = append(r.customParams, param)
}
