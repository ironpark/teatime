package actions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	// File Read Node
	node.RegisterNode(&FileReadActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.file.read",
			node.NodeTypeAction,
			"File Read",
			"파일을 읽어서 내용을 반환하는 액션 노드입니다.",
			"FileText",
			[]node.NodeProperty{
				node.StringProp("path", "File Path",
					node.WithDescription("읽을 파일의 경로를 입력하세요"),
					node.WithPlaceholder("/path/to/file.txt"),
					node.Required(),
				),
				node.SelectProp("encoding", "Encoding", []string{"utf-8", "ascii", "latin1"},
					node.WithDescription("파일 인코딩을 선택하세요"),
					node.OptionalWithDefault("utf-8"),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "content", "Content",
					node.WithDescription("파일의 내용입니다."),
				),
				node.OutputProp(node.String, "filename", "Filename",
					node.WithDescription("파일명입니다."),
				),
				node.OutputProp(node.Int64, "size", "File Size",
					node.WithDescription("파일 크기(바이트)입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "File read successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to read file",
				},
			},
		),
	})

	// File Write Node
	node.RegisterNode(&FileWriteActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.file.write",
			node.NodeTypeAction,
			"File Write",
			"내용을 파일에 저장하는 액션 노드입니다.",
			"Save",
			[]node.NodeProperty{
				node.StringProp("path", "File Path",
					node.WithDescription("저장할 파일의 경로를 입력하세요"),
					node.WithPlaceholder("/path/to/output.txt"),
					node.Required(),
				),
				node.StringProp("content", "Content",
					node.WithDescription("저장할 내용을 입력하세요"),
					node.TextArea(10),
					node.Required(),
				),
				node.SelectProp("mode", "Write Mode", []string{"overwrite", "append"},
					node.WithDescription("파일 쓰기 모드를 선택하세요"),
					node.OptionalWithDefault("overwrite"),
				),
				node.BoolProp("createDirs", "Create Directories",
					node.WithDescription("필요한 디렉토리를 자동으로 생성할지 여부"),
					node.OptionalWithDefault(false),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "path", "File Path",
					node.WithDescription("저장된 파일의 경로입니다."),
				),
				node.OutputProp(node.Int64, "bytesWritten", "Bytes Written",
					node.WithDescription("저장된 바이트 수입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "File written successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to write file",
				},
			},
		),
	})
}

// FileReadActionNode reads files and returns their content.
type FileReadActionNode struct {
	node.BaseNode
}

type fileReadProps struct {
	Path     string `mapstructure:"path"`
	Encoding string `mapstructure:"encoding"`
}

func (f *FileReadActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	var props fileReadProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	if props.Path == "" {
		return node.NodeResult{
			Output: map[string]any{
				"content":  "",
				"filename": "",
				"size":     int64(0),
			},
			Error:         fmt.Errorf("file path is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Set defaults
	if props.Encoding == "" {
		props.Encoding = "utf-8"
	}

	// Read file
	content, err := os.ReadFile(props.Path)
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"content":  "",
				"filename": filepath.Base(props.Path),
				"size":     int64(0),
			},
			Error:         fmt.Errorf("failed to read file: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Get file info
	fileInfo, err := os.Stat(props.Path)
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"content":  string(content),
				"filename": filepath.Base(props.Path),
				"size":     int64(len(content)),
			},
			Error:         fmt.Errorf("failed to get file info: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	return node.NodeResult{
		Output: map[string]any{
			"content":  string(content),
			"filename": fileInfo.Name(),
			"size":     fileInfo.Size(),
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// FileWriteActionNode writes content to files.
type FileWriteActionNode struct {
	node.BaseNode
}

type fileWriteProps struct {
	Path       string `mapstructure:"path"`
	Content    string `mapstructure:"content"`
	Mode       string `mapstructure:"mode"`
	CreateDirs bool   `mapstructure:"createDirs"`
}

func (f *FileWriteActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	var props fileWriteProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	if props.Path == "" {
		return node.NodeResult{
			Output: map[string]any{
				"path":         "",
				"bytesWritten": int64(0),
			},
			Error:         fmt.Errorf("file path is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Set defaults
	if props.Mode == "" {
		props.Mode = "overwrite"
	}

	// Create directories if requested
	if props.CreateDirs {
		dir := filepath.Dir(props.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return node.NodeResult{
				Output: map[string]any{
					"path":         props.Path,
					"bytesWritten": int64(0),
				},
				Error:         fmt.Errorf("failed to create directories: %w", err),
				Continue:      true,
				OutputHandles: []string{"error"},
			}
		}
	}

	var err error
	var bytesWritten int

	// Write file based on mode
	if props.Mode == "append" {
		file, err := os.OpenFile(props.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return node.NodeResult{
				Output: map[string]any{
					"path":         props.Path,
					"bytesWritten": int64(0),
				},
				Error:         fmt.Errorf("failed to open file for append: %w", err),
				Continue:      true,
				OutputHandles: []string{"error"},
			}
		}
		defer file.Close()

		bytesWritten, err = file.WriteString(props.Content)
		if err != nil {
			return node.NodeResult{
				Output: map[string]any{
					"path":         props.Path,
					"bytesWritten": int64(0),
				},
				Error:         fmt.Errorf("failed to write to file: %w", err),
				Continue:      true,
				OutputHandles: []string{"error"},
			}
		}
	} else {
		// Overwrite mode
		err = os.WriteFile(props.Path, []byte(props.Content), 0644)
		if err != nil {
			return node.NodeResult{
				Output: map[string]any{
					"path":         props.Path,
					"bytesWritten": int64(0),
				},
				Error:         fmt.Errorf("failed to write file: %w", err),
				Continue:      true,
				OutputHandles: []string{"error"},
			}
		}
		bytesWritten = len(props.Content)
	}

	return node.NodeResult{
		Output: map[string]any{
			"path":         props.Path,
			"bytesWritten": int64(bytesWritten),
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}