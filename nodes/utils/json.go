package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	// JSON Parse Node
	node.RegisterNode(&JSONParseUtilNode{
		BaseNode: node.NewBaseNode(
			"teatime.util.json.parse",
			node.NodeTypeUtil,
			"JSON Parse",
			"JSON 문자열을 파싱하여 객체로 변환합니다.",
			"Braces",
			[]node.NodeProperty{
				node.StringProp("jsonString", "JSON String",
					node.WithDescription("파싱할 JSON 문자열을 입력하세요"),
					node.WithPlaceholder(`{"key": "value", "number": 123}`),
					node.TextArea(5),
					node.Required(),
				),
				node.StringProp("path", "JSON Path (Optional)",
					node.WithDescription("특정 경로의 값을 추출하려면 JSON 경로를 입력하세요 (예: data.items[0].name)"),
					node.WithPlaceholder("data.items[0].name"),
					node.Optional(),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.JSON, "parsed", "Parsed Object",
					node.WithDescription("파싱된 JSON 객체입니다."),
				),
				node.OutputProp(node.String, "type", "Type",
					node.WithDescription("파싱된 값의 타입입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "JSON parsed successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to parse JSON",
				},
			},
		),
	})

	// JSON Stringify Node
	node.RegisterNode(&JSONStringifyUtilNode{
		BaseNode: node.NewBaseNode(
			"teatime.util.json.stringify",
			node.NodeTypeUtil,
			"JSON Stringify",
			"객체를 JSON 문자열로 변환합니다.",
			"Code",
			[]node.NodeProperty{
				node.JSONProp("object", "Object",
					node.WithDescription("JSON으로 변환할 객체를 입력하세요"),
					node.Required(),
				),
				node.BoolProp("pretty", "Pretty Print",
					node.WithDescription("가독성을 위해 들여쓰기를 추가할지 여부"),
					node.OptionalWithDefault(false),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "jsonString", "JSON String",
					node.WithDescription("변환된 JSON 문자열입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "Object converted to JSON successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to convert object to JSON",
				},
			},
		),
	})
}

// JSONParseUtilNode parses JSON strings into objects.
type JSONParseUtilNode struct {
	node.BaseNode
}

type jsonParseProps struct {
	JSONString string `mapstructure:"jsonString"`
	Path       string `mapstructure:"path"`
}

func (j *JSONParseUtilNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states *node.WorkflowState) node.NodeResult {
	var props jsonParseProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	if props.JSONString == "" {
		return node.NodeResult{
			Output: map[string]any{
				"parsed": nil,
				"type":   "null",
			},
			Error:         fmt.Errorf("JSON string is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Parse JSON
	var parsed any
	if err := json.Unmarshal([]byte(props.JSONString), &parsed); err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"parsed": nil,
				"type":   "null",
			},
			Error:         fmt.Errorf("failed to parse JSON: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	result := parsed
	resultType := getValueType(parsed)

	// Extract specific path if provided
	if props.Path != "" {
		extracted, err := extractJSONPath(parsed, props.Path)
		if err != nil {
			return node.NodeResult{
				Output: map[string]any{
					"parsed": nil,
					"type":   "null",
				},
				Error:         fmt.Errorf("failed to extract path %s: %w", props.Path, err),
				Continue:      true,
				OutputHandles: []string{"error"},
			}
		}
		result = extracted
		resultType = getValueType(extracted)
	}

	return node.NodeResult{
		Output: map[string]any{
			"parsed": result,
			"type":   resultType,
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// JSONStringifyUtilNode converts objects to JSON strings.
type JSONStringifyUtilNode struct {
	node.BaseNode
}

type jsonStringifyProps struct {
	Object any  `mapstructure:"object"`
	Pretty bool `mapstructure:"pretty"`
}

func (j *JSONStringifyUtilNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states *node.WorkflowState) node.NodeResult {
	var props jsonStringifyProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	var jsonBytes []byte
	var err error

	if props.Pretty {
		jsonBytes, err = json.MarshalIndent(props.Object, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(props.Object)
	}

	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"jsonString": "",
			},
			Error:         fmt.Errorf("failed to convert to JSON: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	return node.NodeResult{
		Output: map[string]any{
			"jsonString": string(jsonBytes),
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// Helper functions

func getValueType(value any) string {
	if value == nil {
		return "null"
	}

	switch value.(type) {
	case bool:
		return "boolean"
	case float64, int, int64:
		return "number"
	case string:
		return "string"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return reflect.TypeOf(value).String()
	}
}

func extractJSONPath(data any, path string) (any, error) {
	// Simple JSON path extraction (supports dot notation and array indices)
	// Example: "data.items[0].name"
	
	parts := parseJSONPath(path)
	current := data

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]any:
			if part.isArray {
				return nil, fmt.Errorf("cannot use array index on object at path part: %s", part.key)
			}
			var ok bool
			current, ok = v[part.key]
			if !ok {
				return nil, fmt.Errorf("key not found: %s", part.key)
			}
		case []any:
			if !part.isArray {
				return nil, fmt.Errorf("cannot use object key on array at path part: %s", part.key)
			}
			if part.index < 0 || part.index >= len(v) {
				return nil, fmt.Errorf("array index out of bounds: %d", part.index)
			}
			current = v[part.index]
		default:
			return nil, fmt.Errorf("cannot navigate further into non-object/non-array value")
		}
	}

	return current, nil
}

type pathPart struct {
	key     string
	isArray bool
	index   int
}

func parseJSONPath(path string) []pathPart {
	var parts []pathPart
	segments := strings.Split(path, ".")

	for _, segment := range segments {
		if strings.Contains(segment, "[") {
			// Handle array notation: "items[0]"
			openBracket := strings.Index(segment, "[")
			closeBracket := strings.Index(segment, "]")
			
			if openBracket > 0 {
				// Object key first, then array index
				key := segment[:openBracket]
				parts = append(parts, pathPart{key: key, isArray: false})
			}
			
			if closeBracket > openBracket {
				indexStr := segment[openBracket+1 : closeBracket]
				if index, err := strconv.Atoi(indexStr); err == nil {
					parts = append(parts, pathPart{isArray: true, index: index})
				}
			}
		} else {
			// Simple object key
			if segment != "" {
				parts = append(parts, pathPart{key: segment, isArray: false})
			}
		}
	}

	return parts
}