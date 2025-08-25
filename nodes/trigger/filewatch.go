package trigger

import (
	"context"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
	"github.com/ironpark/teatime/internal/trigger/handlers"
)

func init() {
	node.RegisterNode(&FileWatchTriggerNode{
		BaseNode: node.NewBaseNode(
			"teatime.trigger.filewatch",
			node.NodeTypeTrigger,
			"File Watcher",
			"파일 변경사항을 감지하여 워크플로우를 실행합니다.",
			"FolderOpen",
			[]node.NodeProperty{
				node.StringProp("path", "File Path",
					node.WithDescription("감시할 파일 또는 디렉토리 경로"),
					node.Required(),
				),
				node.BoolProp("watchCreate", "Watch Create",
					node.WithDescription("파일 생성 이벤트를 감지합니다."),
					node.OptionalWithDefault(true),
				),
				node.BoolProp("watchModify", "Watch Modify",
					node.WithDescription("파일 수정 이벤트를 감지합니다."),
					node.OptionalWithDefault(true),
				),
				node.BoolProp("watchRemove", "Watch Remove",
					node.WithDescription("파일 삭제 이벤트를 감지합니다."),
					node.OptionalWithDefault(true),
				),
				node.BoolProp("watchRename", "Watch Rename",
					node.WithDescription("파일 이름 변경 이벤트를 감지합니다."),
					node.OptionalWithDefault(true),
				),
				node.BoolProp("watchChmod", "Watch Permission Change",
					node.WithDescription("파일 권한 변경 이벤트를 감지합니다."),
					node.OptionalWithDefault(false),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.Date, "timestamp", "Timestamp",
					node.WithDescription("파일 변경 시점의 날짜와 시간입니다."),
				),
				node.OutputProp(node.String, "path", "File Path",
					node.WithDescription("변경된 파일의 경로입니다."),
				),
				node.OutputProp(node.String, "operation", "Operation",
					node.WithDescription("파일 작업 유형 (CREATE|WRITE|REMOVE|RENAME|CHMOD)입니다."),
				),
				node.OutputProp(node.Bool, "created", "Created",
					node.WithDescription("파일이 생성되었는지 여부입니다."),
				),
				node.OutputProp(node.Bool, "modified", "Modified",
					node.WithDescription("파일이 수정되었는지 여부입니다."),
				),
				node.OutputProp(node.Bool, "removed", "Removed",
					node.WithDescription("파일이 삭제되었는지 여부입니다."),
				),
				node.OutputProp(node.Bool, "renamed", "Renamed",
					node.WithDescription("파일이 이름이 변경되었는지 여부입니다."),
				),
				node.OutputProp(node.Bool, "chmod", "Permission Changed",
					node.WithDescription("파일 권한이 변경되었는지 여부입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "File Changed",
					Description: "File system event detected",
				},
			},
		),
	})
}


// FileWatchTriggerNode triggers workflow execution on file system events.
type FileWatchTriggerNode struct {
	node.BaseNode
}

// Run executes the file watch trigger logic.
// This is called when a file system event is detected.
func (f *FileWatchTriggerNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	// Extract parameters using mapstructure
	var config handlers.FileWatchConfig
	if err := mapstructure.Decode(resolvedProps, &config); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to decode properties: %w", err),
			Continue:      false,
			OutputHandles: []string{"success"},
		}
	}

	// Extract file event information from execution context
	var filewatchContext handlers.FilewatchContext
	if err := states.BindExecContext(&filewatchContext); err != nil {
		return node.NodeResult{
			Error:    fmt.Errorf("failed to bind execution context: %w", err),
			Continue: false,
		}
	}

	// Extract individual fields for convenience
	filePath := filewatchContext.Path
	operation := filewatchContext.Operation
	timestamp := filewatchContext.Timestamp
	created := filewatchContext.Created
	modified := filewatchContext.Modified
	removed := filewatchContext.Removed
	renamed := filewatchContext.Renamed
	chmod := filewatchContext.Chmod

	// Event filtering is already done at the handler level based on configuration
	// so we can proceed with processing

	// Build output data
	output := map[string]any{
		"timestamp": timestamp,
		"path":      filePath,
		"operation": operation,
		"created":   created,
		"modified":  modified,
		"removed":   removed,
		"renamed":   renamed,
		"chmod":     chmod,
	}

	return node.NodeResult{
		Output:        output,
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// GetWatchPath returns the configured watch path.
func (f *FileWatchTriggerNode) GetWatchPath() string {
	props := f.GetProperties(node.PropertyContext{})
	for _, prop := range props {
		if prop.Key == "path" {
			if path, ok := prop.Value.(string); ok {
				return path
			}
		}
	}
	return ""
}