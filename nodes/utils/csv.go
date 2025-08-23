package utils

import (
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

func init() {
	// CSV Parse Node
	node.RegisterNode(&CSVParseUtilNode{
		BaseNode: node.NewBaseNode(
			"teatime.util.csv.parse",
			node.NodeTypeUtil,
			"CSV Parse",
			"CSV 문자열을 파싱하여 배열로 변환합니다.",
			"Table",
			[]node.NodeProperty{
				node.StringProp("csvString", "CSV String",
					node.WithDescription("파싱할 CSV 문자열을 입력하세요"),
					node.WithPlaceholder("name,age,city\nJohn,30,NYC\nJane,25,LA"),
					node.TextArea(8),
					node.Required(),
				),
				node.BoolProp("hasHeader", "Has Header",
					node.WithDescription("첫 번째 행이 헤더인지 여부"),
					node.OptionalWithDefault(true),
				),
				node.StringProp("delimiter", "Delimiter",
					node.WithDescription("구분자를 지정하세요"),
					node.WithPlaceholder(","),
					node.OptionalWithDefault(","),
				),
				node.BoolProp("autoType", "Auto Type Conversion",
					node.WithDescription("숫자와 불린 값을 자동으로 변환할지 여부"),
					node.OptionalWithDefault(true),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.JSON, "rows", "Rows",
					node.WithDescription("파싱된 CSV 행들의 배열입니다."),
				),
				node.OutputProp(node.JSON, "headers", "Headers",
					node.WithDescription("CSV 헤더 배열입니다 (hasHeader가 true인 경우)."),
				),
				node.OutputProp(node.Int64, "rowCount", "Row Count",
					node.WithDescription("총 행 수입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "CSV parsed successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to parse CSV",
				},
			},
		),
	})

	// CSV Generate Node
	node.RegisterNode(&CSVGenerateUtilNode{
		BaseNode: node.NewBaseNode(
			"teatime.util.csv.generate",
			node.NodeTypeUtil,
			"CSV Generate",
			"배열 데이터를 CSV 문자열로 변환합니다.",
			"FileSpreadsheet",
			[]node.NodeProperty{
				node.JSONProp("data", "Data",
					node.WithDescription("CSV로 변환할 배열 데이터를 입력하세요"),
					node.WithPlaceholder(`[{"name":"John","age":30},{"name":"Jane","age":25}]`),
					node.Required(),
				),
				node.JSONProp("headers", "Headers (Optional)",
					node.WithDescription("사용할 헤더 배열 (없으면 첫 번째 객체의 키를 사용)"),
					node.WithPlaceholder(`["name","age","city"]`),
					node.Optional(),
				),
				node.StringProp("delimiter", "Delimiter",
					node.WithDescription("구분자를 지정하세요"),
					node.WithPlaceholder(","),
					node.OptionalWithDefault(","),
				),
				node.BoolProp("includeHeader", "Include Header",
					node.WithDescription("헤더를 포함할지 여부"),
					node.OptionalWithDefault(true),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "csvString", "CSV String",
					node.WithDescription("생성된 CSV 문자열입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "CSV generated successfully",
				},
				{
					ID:          "error",
					Label:       "Error",
					Description: "Failed to generate CSV",
				},
			},
		),
	})
}

// CSVParseUtilNode parses CSV strings into arrays.
type CSVParseUtilNode struct {
	node.BaseNode
}

type csvParseProps struct {
	CSVString string `mapstructure:"csvString"`
	HasHeader bool   `mapstructure:"hasHeader"`
	Delimiter string `mapstructure:"delimiter"`
	AutoType  bool   `mapstructure:"autoType"`
}

func (c *CSVParseUtilNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	var props csvParseProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	if props.CSVString == "" {
		return node.NodeResult{
			Output: map[string]any{
				"rows":     []any{},
				"headers":  []string{},
				"rowCount": int64(0),
			},
			Error:         fmt.Errorf("CSV string is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Set defaults
	if props.Delimiter == "" {
		props.Delimiter = ","
	}

	// Parse CSV
	reader := csv.NewReader(strings.NewReader(props.CSVString))
	reader.Comma = rune(props.Delimiter[0])

	records, err := reader.ReadAll()
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"rows":     []any{},
				"headers":  []string{},
				"rowCount": int64(0),
			},
			Error:         fmt.Errorf("failed to parse CSV: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	if len(records) == 0 {
		return node.NodeResult{
			Output: map[string]any{
				"rows":     []any{},
				"headers":  []string{},
				"rowCount": int64(0),
			},
			Error:         nil,
			Continue:      true,
			OutputHandles: []string{"success"},
		}
	}

	var headers []string
	var dataRows [][]string
	
	if props.HasHeader && len(records) > 0 {
		headers = records[0]
		dataRows = records[1:]
	} else {
		// Generate generic headers
		if len(records) > 0 {
			for i := 0; i < len(records[0]); i++ {
				headers = append(headers, fmt.Sprintf("column_%d", i+1))
			}
		}
		dataRows = records
	}

	// Convert rows to objects if headers are available
	var rows []any
	if len(headers) > 0 {
		for _, row := range dataRows {
			obj := make(map[string]any)
			for i, value := range row {
				if i < len(headers) {
					if props.AutoType {
						obj[headers[i]] = convertValue(value)
					} else {
						obj[headers[i]] = value
					}
				}
			}
			rows = append(rows, obj)
		}
	} else {
		// Return as arrays if no headers
		for _, row := range dataRows {
			if props.AutoType {
				convertedRow := make([]any, len(row))
				for i, value := range row {
					convertedRow[i] = convertValue(value)
				}
				rows = append(rows, convertedRow)
			} else {
				rows = append(rows, row)
			}
		}
	}

	return node.NodeResult{
		Output: map[string]any{
			"rows":     rows,
			"headers":  headers,
			"rowCount": int64(len(dataRows)),
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// CSVGenerateUtilNode converts arrays to CSV strings.
type CSVGenerateUtilNode struct {
	node.BaseNode
}

type csvGenerateProps struct {
	Data          []any    `mapstructure:"data"`
	Headers       []string `mapstructure:"headers"`
	Delimiter     string   `mapstructure:"delimiter"`
	IncludeHeader bool     `mapstructure:"includeHeader"`
}

func (c *CSVGenerateUtilNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	var props csvGenerateProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	if len(props.Data) == 0 {
		return node.NodeResult{
			Output: map[string]any{
				"csvString": "",
			},
			Error:         fmt.Errorf("data array is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Set defaults
	if props.Delimiter == "" {
		props.Delimiter = ","
	}

	var output strings.Builder
	writer := csv.NewWriter(&output)
	writer.Comma = rune(props.Delimiter[0])

	// Determine headers
	var headers []string
	if len(props.Headers) > 0 {
		headers = props.Headers
	} else {
		// Extract headers from first object
		if len(props.Data) > 0 {
			if obj, ok := props.Data[0].(map[string]any); ok {
				for key := range obj {
					headers = append(headers, key)
				}
			}
		}
	}

	// Write header if requested
	if props.IncludeHeader && len(headers) > 0 {
		if err := writer.Write(headers); err != nil {
			return node.NodeResult{
				Output: map[string]any{
					"csvString": "",
				},
				Error:         fmt.Errorf("failed to write CSV header: %w", err),
				Continue:      true,
				OutputHandles: []string{"error"},
			}
		}
	}

	// Write data rows
	for _, item := range props.Data {
		var row []string
		
		switch v := item.(type) {
		case map[string]any:
			// Object - use headers to maintain order
			for _, header := range headers {
				if value, exists := v[header]; exists {
					row = append(row, fmt.Sprintf("%v", value))
				} else {
					row = append(row, "")
				}
			}
		case []any:
			// Array - convert each element
			for _, elem := range v {
				row = append(row, fmt.Sprintf("%v", elem))
			}
		default:
			// Single value
			row = append(row, fmt.Sprintf("%v", v))
		}

		if err := writer.Write(row); err != nil {
			return node.NodeResult{
				Output: map[string]any{
					"csvString": "",
				},
				Error:         fmt.Errorf("failed to write CSV row: %w", err),
				Continue:      true,
				OutputHandles: []string{"error"},
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"csvString": "",
			},
			Error:         fmt.Errorf("failed to finalize CSV: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	return node.NodeResult{
		Output: map[string]any{
			"csvString": output.String(),
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// Helper function to convert string values to appropriate types
func convertValue(value string) any {
	// Try to convert to number
	if intVal, err := strconv.Atoi(value); err == nil {
		return int64(intVal)
	}
	
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}
	
	// Try to convert to boolean
	if boolVal, err := strconv.ParseBool(value); err == nil {
		return boolVal
	}
	
	// Return as string if no conversion possible
	return value
}