package actions

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

var (
	httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	contentTypes = []string{"application/json", "application/x-www-form-urlencoded", "text/plain", "application/xml", "multipart/form-data"}
)

func init() {
	node.RegisterNode(&HTTPRequestActionNode{
		BaseNode: node.NewBaseNode(
			"teatime.action.http",
			node.NodeTypeAction,
			"HTTP Request",
			"HTTP 요청을 보내고 응답을 받는 액션 노드입니다.",
			"Globe",
			[]node.NodeProperty{
				node.StringProp("url", "URL",
					node.WithDescription("요청할 URL을 입력하세요"),
					node.WithPlaceholder("https://api.example.com/data"),
					node.Required(),
				),
				node.SelectProp("method", "Method", httpMethods,
					node.WithDescription("HTTP 메서드를 선택하세요"),
					node.RequiredWithDefault("GET"),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.String, "response", "Response",
					node.WithDescription("HTTP 응답 본문입니다."),
				),
				node.OutputProp(node.Int64, "statusCode", "Status Code",
					node.WithDescription("HTTP 상태 코드입니다."),
				),
				node.OutputProp(node.JSON, "headers", "Headers",
					node.WithDescription("응답 헤더들입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Success",
					Description: "HTTP request completed successfully",
				},
				{
					ID:          "error",
					Label:       "Error", 
					Description: "HTTP request failed",
				},
			},
		),
	})
}

// HTTPRequestActionNode sends HTTP requests and returns responses.
type HTTPRequestActionNode struct {
	node.BaseNode
}

// GetProperties returns dynamic properties based on the selected method.
func (h *HTTPRequestActionNode) GetProperties(ctx node.PropertyContext) []node.NodeProperty {
	baseProps := h.BaseNode.GetProperties(ctx)
	props := make([]node.NodeProperty, 0, len(baseProps)+10)
	props = append(props, baseProps...)

	// Get current method selection
	method, _ := ctx["method"].(string)
	if method == "" {
		method = "GET"
	}

	// Add headers property
	props = append(props,
		node.JSONProp("headers", "Headers",
			node.WithDescription("요청 헤더를 JSON 형식으로 입력하세요"),
			node.WithPlaceholder(`{"Authorization": "Bearer token", "Content-Type": "application/json"}`),
			node.Optional(),
		),
	)

	// Add body property for methods that support it
	if method != "GET" && method != "HEAD" && method != "OPTIONS" {
		props = append(props,
			node.SelectProp("contentType", "Content Type", contentTypes,
				node.WithDescription("요청 본문의 Content-Type을 선택하세요"),
				node.OptionalWithDefault("application/json"),
			),
			node.StringProp("body", "Body",
				node.WithDescription("요청 본문을 입력하세요"),
				node.WithPlaceholder(`{"key": "value"}`),
				node.TextArea(5),
				node.Optional(),
			),
		)
	}

	// Add additional options
	props = append(props,
		node.IntProp("timeout", "Timeout (seconds)",
			node.WithDescription("요청 타임아웃 시간(초)"),
			node.WithRange(1, 300, 1),
			node.OptionalWithDefault(int64(30)),
		),
		node.BoolProp("followRedirects", "Follow Redirects",
			node.WithDescription("리디렉션을 자동으로 따를지 여부"),
			node.OptionalWithDefault(true),
		),
	)

	return props
}

type httpActionProps struct {
	URL             string            `mapstructure:"url"`
	Method          string            `mapstructure:"method"`
	Headers         map[string]string `mapstructure:"headers"`
	ContentType     string            `mapstructure:"contentType"`
	Body            string            `mapstructure:"body"`
	Timeout         int64             `mapstructure:"timeout"`
	FollowRedirects bool              `mapstructure:"followRedirects"`
}

// Run executes the HTTP request and returns the response.
func (h *HTTPRequestActionNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	var props httpActionProps
	if err := mapstructure.Decode(resolvedProps, &props); err != nil {
		return node.NodeResult{
			Error:         err,
			Continue:      false,
			OutputHandles: []string{"error"},
		}
	}

	if props.URL == "" {
		return node.NodeResult{
			Output: map[string]any{
				"response":   "",
				"statusCode": 0,
				"headers":    map[string]string{},
			},
			Error:         fmt.Errorf("URL is required"),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Set defaults
	if props.Method == "" {
		props.Method = "GET"
	}
	if props.Timeout == 0 {
		props.Timeout = 30
	}
	if props.ContentType == "" {
		props.ContentType = "application/json"
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(props.Timeout) * time.Second,
	}

	// Don't follow redirects if specified
	if !props.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Prepare request body
	var bodyReader io.Reader
	if props.Body != "" && props.Method != "GET" && props.Method != "HEAD" && props.Method != "OPTIONS" {
		bodyReader = strings.NewReader(props.Body)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, props.Method, props.URL, bodyReader)
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"response":   "",
				"statusCode": 0,
				"headers":    map[string]string{},
			},
			Error:         fmt.Errorf("failed to create request: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Set content type if body is present
	if bodyReader != nil && props.ContentType != "" {
		req.Header.Set("Content-Type", props.ContentType)
	}

	// Set custom headers
	for key, value := range props.Headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"response":   "",
				"statusCode": 0,
				"headers":    map[string]string{},
			},
			Error:         fmt.Errorf("request failed: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return node.NodeResult{
			Output: map[string]any{
				"response":   "",
				"statusCode": int64(resp.StatusCode),
				"headers":    convertHeaders(resp.Header),
			},
			Error:         fmt.Errorf("failed to read response: %w", err),
			Continue:      true,
			OutputHandles: []string{"error"},
		}
	}

	// Determine output handle based on status code
	outputHandle := "success"
	if resp.StatusCode >= 400 {
		outputHandle = "error"
	}

	return node.NodeResult{
		Output: map[string]any{
			"response":   string(responseBody),
			"statusCode": int64(resp.StatusCode),
			"headers":    convertHeaders(resp.Header),
		},
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{outputHandle},
	}
}

// convertHeaders converts http.Header to map[string]string
func convertHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = values[0] // Take first value if multiple exist
		}
	}
	return result
}