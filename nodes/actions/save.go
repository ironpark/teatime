package actions

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-viper/mapstructure/v2"
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
				node.SelectProp("encoding", "Encoding", []string{"utf-8", "utf-16"},
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
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "File saved successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to save file",
				},
			},
		),
	})
}

type saveActionProps struct {
	FilePath  string `mapstructure:"filePath"`
	Content   string `mapstructure:"content"`
	Encoding  string `mapstructure:"encoding"`
	Append    bool   `mapstructure:"append"`
	CreateDir bool   `mapstructure:"createDir"`
}

// SaveActionNode saves content to a file with various options.
type SaveActionNode struct {
	node.BaseNode
}

// Run executes the file save operation.
func (s *SaveActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	// Extract parameters using mapstructure
	var props saveActionProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to decode properties: %w", err),
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	// Validate required parameters
	if props.FilePath == "" {
		return node.NodeResult{
			Output: map[string]any{
				"savedPath": "",
				"fileSize":  int64(0),
				"success":   false,
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

	// Create directory if needed
	if props.CreateDir {
		dir := filepath.Dir(props.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return node.NodeResult{
				Output: map[string]any{
					"savedPath": "",
					"fileSize":  int64(0),
					"success":   false,
				},
				Error:         fmt.Errorf("failed to create directory %s: %w", dir, err),
				Continue:      true,
				OutputHandles: []string{"error"},
			}
		}
	}

	// Determine file open flags
	flags := os.O_CREATE | os.O_WRONLY
	if props.Append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	// Open file
	file, err := os.OpenFile(props.FilePath, flags, 0644)
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"savedPath": "",
				"fileSize":  int64(0),
				"success":   false,
			},
			Error:         fmt.Errorf("failed to open file %s: %w", props.FilePath, err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}
	defer file.Close()

	// Write content
	var bytesWritten int64
	var writeErr error

	switch props.Encoding {
	case "utf-8":
		n, err := io.WriteString(file, props.Content)
		bytesWritten = int64(n)
		writeErr = err
	case "utf-16":
		// Convert to UTF-16 (simplified - using UTF-16LE)
		content := []byte(props.Content)
		utf16Content := make([]byte, 0, len(content)*2+2)
		// Add BOM for UTF-16LE
		utf16Content = append(utf16Content, 0xFF, 0xFE)
		for _, b := range content {
			utf16Content = append(utf16Content, b, 0x00)
		}
		n, err := file.Write(utf16Content)
		bytesWritten = int64(n)
		writeErr = err
	default:
		return node.NodeResult{
			Output: map[string]any{
				"savedPath": "",
				"fileSize":  int64(0),
				"success":   false,
			},
			Error:         fmt.Errorf("unsupported encoding: %s", props.Encoding),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	if writeErr != nil {
		return node.NodeResult{
			Output: map[string]any{
				"savedPath": "",
				"fileSize":  int64(0),
				"success":   false,
			},
			Error:         fmt.Errorf("failed to write to file %s: %w", props.FilePath, writeErr),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Get absolute path for output
	absPath, err := filepath.Abs(props.FilePath)
	if err != nil {
		absPath = props.FilePath // fallback to original path
	}

	return node.NodeResult{
		Output: map[string]any{
			"savedPath": absPath,
			"fileSize":  bytesWritten,
			"success":   true,
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}
