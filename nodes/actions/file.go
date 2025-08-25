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
	// File Delete Node
	node.RegisterNode(&FileDeleteActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.file.delete",
			node.NodeTypeAction,
			"File Delete",
			"파일 또는 디렉토리를 삭제합니다.",
			"Trash2",
			[]node.NodeProperty{
				node.StringProp("path", "Path",
					node.WithDescription("삭제할 파일 또는 디렉토리의 경로를 입력하세요"),
					node.WithPlaceholder("/path/to/file.txt or /path/to/directory"),
					node.Required(),
				),
				node.BoolProp("recursive", "Recursive",
					node.WithDescription("디렉토리를 재귀적으로 삭제할지 여부"),
					node.OptionalWithDefault(false),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "deletedPath", "Deleted Path",
					node.WithDescription("삭제된 파일 또는 디렉토리의 경로입니다."),
				),
				node.OutputProp(node.Bool, "wasDirectory", "Was Directory",
					node.WithDescription("삭제된 경로가 디렉토리였는지 여부입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "File/directory deleted successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to delete file/directory",
				},
			},
		),
	})

	// File Move Node
	node.RegisterNode(&FileMoveActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.file.move",
			node.NodeTypeAction,
			"File Move",
			"파일 또는 디렉토리를 이동하거나 이름을 변경합니다.",
			"Move",
			[]node.NodeProperty{
				node.StringProp("sourcePath", "Source Path",
					node.WithDescription("이동할 파일 또는 디렉토리의 경로를 입력하세요"),
					node.WithPlaceholder("/path/to/source"),
					node.Required(),
				),
				node.StringProp("destinationPath", "Destination Path",
					node.WithDescription("목적지 경로를 입력하세요"),
					node.WithPlaceholder("/path/to/destination"),
					node.Required(),
				),
				node.BoolProp("createDirs", "Create Directories",
					node.WithDescription("목적지 디렉토리가 없을 경우 자동으로 생성할지 여부"),
					node.OptionalWithDefault(false),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "sourcePath", "Source Path",
					node.WithDescription("원본 경로입니다."),
				),
				node.OutputProp(node.String, "destinationPath", "Destination Path",
					node.WithDescription("목적지 경로입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "File/directory moved successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to move file/directory",
				},
			},
		),
	})

	// File Read Node
	node.RegisterNode(&FileReadActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.file.read",
			node.NodeTypeAction,
			"File Read",
			"파일을 읽어서 내용을 반환합니다.",
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
			"내용을 파일에 저장합니다.",
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

func (f *FileReadActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states *node.WorkflowState) node.NodeResult {
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

func (f *FileWriteActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states *node.WorkflowState) node.NodeResult {
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

// FileDeleteActionNode deletes files and directories.
type FileDeleteActionNode struct {
	node.BaseNode
}

type fileDeleteProps struct {
	Path      string `mapstructure:"path"`
	Recursive bool   `mapstructure:"recursive"`
}

func (f *FileDeleteActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states *node.WorkflowState) node.NodeResult {
	var props fileDeleteProps
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
				"deletedPath":  "",
				"wasDirectory": false,
			},
			Error:         fmt.Errorf("path is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Check if path exists and if it's a directory
	fileInfo, err := os.Stat(props.Path)
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"deletedPath":  props.Path,
				"wasDirectory": false,
			},
			Error:         fmt.Errorf("failed to stat path: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	isDirectory := fileInfo.IsDir()

	// Delete the file or directory
	if isDirectory && props.Recursive {
		err = os.RemoveAll(props.Path)
	} else if isDirectory && !props.Recursive {
		err = os.Remove(props.Path) // Will fail if directory is not empty
	} else {
		err = os.Remove(props.Path)
	}

	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"deletedPath":  props.Path,
				"wasDirectory": isDirectory,
			},
			Error:         fmt.Errorf("failed to delete: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	return node.NodeResult{
		Output: map[string]any{
			"deletedPath":  props.Path,
			"wasDirectory": isDirectory,
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// FileMoveActionNode moves/renames files and directories.
type FileMoveActionNode struct {
	node.BaseNode
}

type fileMoveProps struct {
	SourcePath      string `mapstructure:"sourcePath"`
	DestinationPath string `mapstructure:"destinationPath"`
	CreateDirs      bool   `mapstructure:"createDirs"`
}

func (f *FileMoveActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states *node.WorkflowState) node.NodeResult {
	var props fileMoveProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	if props.SourcePath == "" {
		return node.NodeResult{
			Output: map[string]any{
				"sourcePath":      "",
				"destinationPath": props.DestinationPath,
			},
			Error:         fmt.Errorf("source path is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	if props.DestinationPath == "" {
		return node.NodeResult{
			Output: map[string]any{
				"sourcePath":      props.SourcePath,
				"destinationPath": "",
			},
			Error:         fmt.Errorf("destination path is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Create destination directories if requested
	if props.CreateDirs {
		destDir := filepath.Dir(props.DestinationPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return node.NodeResult{
				Output: map[string]any{
					"sourcePath":      props.SourcePath,
					"destinationPath": props.DestinationPath,
				},
				Error:         fmt.Errorf("failed to create destination directories: %w", err),
				Continue:      true,
				OutputHandles: []string{"error"},
			}
		}
	}

	// Check if source exists
	if _, err := os.Stat(props.SourcePath); err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"sourcePath":      props.SourcePath,
				"destinationPath": props.DestinationPath,
			},
			Error:         fmt.Errorf("source path does not exist: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Move/rename the file or directory
	err := os.Rename(props.SourcePath, props.DestinationPath)
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"sourcePath":      props.SourcePath,
				"destinationPath": props.DestinationPath,
			},
			Error:         fmt.Errorf("failed to move: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	return node.NodeResult{
		Output: map[string]any{
			"sourcePath":      props.SourcePath,
			"destinationPath": props.DestinationPath,
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}