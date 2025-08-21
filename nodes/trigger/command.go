package trigger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	node.RegisterNode(&CommandTriggerNode{
		BaseNode: node.NewBaseNode(
			"teatime.trigger.command",
			node.NodeTypeTrigger,
			"Command",
			"명령어를 통해 워크플로우를 실행하는 트리거 노드입니다.",
			"Zap",
			[]node.NodeProperty{
				node.StringProp("command", "Command",
					node.WithDescription("명령어의 이름을 입력하세요"),
					node.Required(),
				),
				node.BoolProp("global", "Global",
					node.WithDescription("명령어를 전역으로 등록할지 여부"),
					node.Optional(),
					node.WithDefault(false),
				),
				node.StringArrayProp("args", "Arguments",
					node.WithDescription("명령어 인자 정의"),
					node.WithDefault([]string{}),
					node.DynamicList(true),
					node.Optional(),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "recipeFile", "Recipe File",
					node.WithDescription("실행된 레시피 파일 경로입니다."),
				),
				node.OutputProp(node.JSON, "args", "Arguments",
					node.WithDescription("명령어 인자들입니다."),
				),
				node.OutputProp(node.String, "workdir", "Working Directory",
					node.WithDescription("명령어가 호출된 디렉토리입니다."),
				),
				node.OutputProp(node.Date, "timestamp", "Timestamp",
					node.WithDescription("호출시점의 날짜와 시간입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Start",
					Description: "Command trigger activated",
				},
			},
		),
	})
}

type commandTriggerProps struct {
	Description string         `mapstructure:"description"`
	Args        map[string]any `mapstructure:"args"`
}

// CommandTriggerNode triggers workflow execution via CLI commands.
type CommandTriggerNode struct {
	node.BaseNode
}

// GetOutput returns dynamic output properties based on configured arguments.
func (c *CommandTriggerNode) GetOutput(ctx node.PropertyContext) []node.NodeProperty {
	// Get base output properties
	baseOutput := c.BaseNode.GetOutput(ctx)
	output := make([]node.NodeProperty, 0, len(baseOutput)+10)
	output = append(output, baseOutput...)

	// Get args from context to generate dynamic outputs
	if args, ok := ctx["args"].(map[string]any); ok {
		for argKey, argValue := range args {
			// Determine the property type based on the value type
			var propType node.PropertyType
			switch argValue.(type) {
			case bool:
				propType = node.Bool
			case int, int32, int64, float32, float64:
				propType = node.Float64
			case string:
				propType = node.String
			case []interface{}, []string:
				propType = node.StringArray
			default:
				propType = node.JSON
			}

			// Create output property for this argument
			argOutput := node.OutputProp(propType, argKey, argKey,
				node.WithDescription(fmt.Sprintf("CLI argument: %s", argKey)),
			)
			output = append(output, argOutput)
		}
	}

	return output
}

// Run executes the command trigger logic.
// This is called when: teatime run recipe-file --arg1 value1 --arg2 value2
func (c *CommandTriggerNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	// Extract parameters using mapstructure
	var props commandTriggerProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"success"},
		}
	}

	// Get current working directory
	workdir, err := os.Getwd()
	if err != nil {
		workdir = "."
	}

	// Get current timestamp
	timestamp := time.Now()

	// Extract recipe file and CLI arguments from states
	// These would be populated by the CLI system when running: teatime run recipe-file --arg1 value1
	recipeFile := ""
	if rf, ok := states["recipeFile"].(string); ok {
		recipeFile = rf
	}

	// Merge CLI arguments with defined arguments from props
	finalArgs := make(map[string]any)

	// Start with defined arguments (defaults from node configuration)
	if props.Args != nil {
		for k, v := range props.Args {
			finalArgs[k] = v
		}
	}

	// Override with CLI arguments from states
	if cliArgs, ok := states["cliArgs"].(map[string]any); ok {
		for k, v := range cliArgs {
			finalArgs[k] = v
		}
	}

	// Build output including individual argument outputs
	output := map[string]any{
		"recipeFile": recipeFile,
		"args":       finalArgs,
		"workdir":    workdir,
		"timestamp":  timestamp,
	}

	// Add individual argument values as separate outputs
	for argKey, argValue := range finalArgs {
		output[argKey] = argValue
	}

	return node.NodeResult{
		Output:        output,
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// GetExpectedArgs returns the argument definitions for CLI help and validation.
func (c *CommandTriggerNode) GetExpectedArgs() map[string]any {
	props := c.GetProperties(node.PropertyContext{})
	for _, prop := range props {
		if prop.Key == "args" {
			if args, ok := prop.Value.(map[string]any); ok {
				return args
			}
		}
	}
	return make(map[string]any)
}

// GetDescription returns the command description.
func (c *CommandTriggerNode) GetDescription() string {
	props := c.GetProperties(node.PropertyContext{})
	for _, prop := range props {
		if prop.Key == "description" {
			if desc, ok := prop.Value.(string); ok {
				return desc
			}
		}
	}
	return "Command trigger for workflow execution"
}
