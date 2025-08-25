package trigger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
	"github.com/ironpark/teatime/internal/trigger/handlers"
)

func init() {
	node.RegisterNode(&CommandTriggerNode{
		BaseNode: node.NewBaseNode(
			"teatime.trigger.command",
			node.NodeTypeTrigger,
			"Command",
			"명령어를 통해 워크플로우를 실행합니다.",
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
				node.JSONProp("args", "Arguments",
					node.WithDescription("명령어 인자 정의 [{이름, 필수여부, 리스트여부, 설명}, ...]"),
					node.WithDefault([]string{}),
					node.ArgList(),
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

// CommandArg represents a command line argument definition
type CommandArg struct {
	Name        string `json:"name" mapstructure:"name"`               // 옵션이름
	Required    bool   `json:"required" mapstructure:"required"`       // 필수여부
	List        bool   `json:"list" mapstructure:"list"`               // 리스트여부
	Description string `json:"description" mapstructure:"description"` // 설명
}

type commandTriggerProps struct {
	Description string       `mapstructure:"description"`
	Args        []CommandArg `mapstructure:"args"`
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
	if argsValue, ok := ctx["args"]; ok {
		var commandArgs []CommandArg
		
		// Handle different types of args values
		switch v := argsValue.(type) {
		case []CommandArg:
			commandArgs = v
		case []any:
			// Convert from any slice (from JSON unmarshaling)
			for _, item := range v {
				if argMap, ok := item.(map[string]any); ok {
					arg := CommandArg{
						Name:        getStringFromMap(argMap, "name"),
						Required:    getBoolFromMap(argMap, "required"),
						List:        getBoolFromMap(argMap, "list"),
						Description: getStringFromMap(argMap, "description"),
					}
					commandArgs = append(commandArgs, arg)
				}
			}
		}

		// Create output properties for each configured argument
		for _, arg := range commandArgs {
			if arg.Name == "" {
				continue // Skip empty argument definitions
			}

			// Determine property type based on argument configuration
			var propType node.PropertyType
			if arg.List {
				propType = node.StringArray // List arguments are string arrays
			} else {
				propType = node.String // Single arguments are strings
			}

			// Create output property for this argument
			argOutput := node.OutputProp(propType, arg.Name, arg.Name,
				node.WithDescription(fmt.Sprintf("CLI argument: %s - %s", arg.Name, arg.Description)),
			)
			output = append(output, argOutput)
		}
	}

	return output
}

// Helper functions for type conversion from map[string]any
func getStringFromMap(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getBoolFromMap(m map[string]any, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// Run executes the command trigger logic.
// This is called when: teatime run recipe-file --arg1 value1 --arg2 value2
func (c *CommandTriggerNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states *node.WorkflowState) node.NodeResult {
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

	// Extract CLI context from execution context
	var cmdContext handlers.CommandContext
	if err := states.BindExecContext(&cmdContext); err != nil {
		return node.NodeResult{
			Error:    fmt.Errorf("failed to bind execution context: %w", err),
			Continue: false,
		}
	}

	recipeFile := cmdContext.RecipeFile
	cliArgs := cmdContext.CliArgs

	// Merge CLI arguments with defined arguments from props
	finalArgs := make(map[string]any)

	// Start with defined arguments (set defaults based on configuration)
	if props.Args != nil {
		for _, arg := range props.Args {
			if arg.Name != "" {
				// Set default value based on argument type
				if arg.List {
					finalArgs[arg.Name] = []string{} // Default empty array for list arguments
				} else {
					finalArgs[arg.Name] = "" // Default empty string for single arguments
				}
			}
		}
	}

	// Override with CLI arguments from execution context
	for k, v := range cliArgs {
		finalArgs[k] = v
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
func (c *CommandTriggerNode) GetExpectedArgs() []CommandArg {
	props := c.GetProperties(node.PropertyContext{})
	for _, prop := range props {
		if prop.Key == "args" {
			// Handle different types of args values
			switch v := prop.Value.(type) {
			case []CommandArg:
				return v
			case []any:
				// Convert from any slice (from JSON unmarshaling)
				var commandArgs []CommandArg
				for _, item := range v {
					if argMap, ok := item.(map[string]any); ok {
						arg := CommandArg{
							Name:        getStringFromMap(argMap, "name"),
							Required:    getBoolFromMap(argMap, "required"),
							List:        getBoolFromMap(argMap, "list"),
							Description: getStringFromMap(argMap, "description"),
						}
						commandArgs = append(commandArgs, arg)
					}
				}
				return commandArgs
			}
		}
	}
	return []CommandArg{}
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
