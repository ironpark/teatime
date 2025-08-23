package actions

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&CommandActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.command",
			node.NodeTypeAction,
			"Command",
			"시스템 명령어를 실행합니다.",
			"Terminal",
			[]node.NodeProperty{
				node.StringProp("command", "Command",
					node.WithDescription("실행할 명령어를 입력하세요"),
					node.Required(),
				),
				node.StringProp("workdir", "Working Directory",
					node.WithDescription("작업 디렉토리를 입력하세요"),
					node.Optional(),
				),
				node.IntProp("timeout", "Timeout",
					node.WithDescription("타임아웃 시간(초)"),
					node.OptionalWithDefault(int64(30)),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "output", "Output",
					node.WithDescription("명령어 실행 결과입니다."),
				),
				node.OutputProp(node.Int64, "exitCode", "Exit Code",
					node.WithDescription("명령어 종료 코드입니다."),
				),
				node.OutputProp(node.String, "error", "Error",
					node.WithDescription("오류 메시지입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "Success output handle",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Error output handle",
				},
			},
		),
	})
}

// CommandActionNode executes system commands and captures their output.
type CommandActionNode struct {
	node.BaseNode
}

// Run executes the configured system command and returns the result.
func (c *CommandActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	// Extract command parameters
	command, ok := resolvedProps["command"].(string)
	if !ok || command == "" {
		return node.NodeResult{
			Error:         fmt.Errorf("command is required"),
			Continue:      false,
			OutputHandles: []string{"default"},
		}
	}

	workdir, _ := resolvedProps["workdir"].(string)
	timeoutSecs, _ := resolvedProps["timeout"].(int64)
	if timeoutSecs <= 0 {
		timeoutSecs = 30
	}

	// Create command context with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	// Execute command
	cmd := exec.CommandContext(cmdCtx, "sh", "-c", command)
	if workdir != "" {
		cmd.Dir = workdir
	}

	output, err := cmd.CombinedOutput()
	exitCode := 0
	errorMsg := ""
	var handles []string
	if err != nil {
		errorMsg = err.Error()
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
		handles = []string{"error"}
	} else {
		handles = []string{"success"}
	}
	return node.NodeResult{
		Output: map[string]any{
			"output":   string(output),
			"exitCode": int64(exitCode),
			"error":    errorMsg,
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: handles,
	}
}
