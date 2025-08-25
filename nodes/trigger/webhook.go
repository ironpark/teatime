package trigger

import (
	"context"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
	"github.com/ironpark/teatime/internal/trigger/handlers"
)

func init() {
	node.RegisterNode(&WebhookTriggerNode{
		BaseNode: node.NewBaseNode(
			"teatime.trigger.webhook",
			node.NodeTypeTrigger,
			"Webhook",
			"Webhook을 통해 워크플로우를 실행합니다.",
			"Webhook",
			[]node.NodeProperty{
				node.StringProp("path", "Path",
					node.WithDescription("Webhook 경로 (예: /webhook/my-trigger)"),
					node.Required(),
				),
				node.SelectProp("method", "HTTP Method", []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
					node.WithDescription("허용할 HTTP 메서드"),
					node.OptionalWithDefault("POST"),
				),
				node.BoolProp("requireAuth", "Require Authentication",
					node.WithDescription("인증 필요 여부"),
					node.OptionalWithDefault(false),
				),
				node.StringProp("secret", "Secret",
					node.WithDescription("Webhook 검증을 위한 시크릿 키"),
					node.Optional(),
				),
			},
			[]node.NodeProperty{
				node.OutputProp(node.Date, "timestamp", "Timestamp",
					node.WithDescription("호출시점의 날짜와 시간입니다."),
				),
				node.OutputProp(node.String, "method", "HTTP Method",
					node.WithDescription("요청 HTTP 메서드입니다."),
				),
				node.OutputProp(node.String, "path", "Request Path",
					node.WithDescription("요청 경로입니다."),
				),
				node.OutputProp(node.JSON, "headers", "Headers",
					node.WithDescription("요청 헤더들입니다."),
				),
				node.OutputProp(node.JSON, "query", "Query Parameters",
					node.WithDescription("쿼리 파라미터들입니다."),
				),
				node.OutputProp(node.JSON, "body", "Request Body",
					node.WithDescription("요청 바디입니다."),
				),
				node.OutputProp(node.String, "remoteAddr", "Remote Address",
					node.WithDescription("요청자 IP 주소입니다."),
				),
			},
			[]node.OutputHandle{
				{
					ID:          "success",
					Label:       "Triggered",
					Description: "Webhook successfully triggered",
				},
				{
					ID:          "unauthorized",
					Label:       "Unauthorized",
					Description: "Authentication failed",
				},
			},
		),
	})
}


// WebhookTriggerNode triggers workflow execution via HTTP webhooks.
type WebhookTriggerNode struct {
	node.BaseNode
}

// Run executes the webhook trigger logic.
// This is called when an HTTP request is received on the configured path.
func (w *WebhookTriggerNode) Run(ctx context.Context, resolvedProps node.PropertyContext, states node.WorkflowState) node.NodeResult {
	// Extract parameters using mapstructure
	var config handlers.WebhookConfig
	if err := mapstructure.Decode(resolvedProps, &config); err != nil {
		return node.NodeResult{
			Error:         fmt.Errorf("failed to decode properties: %w", err),
			Continue:      false,
			OutputHandles: []string{"success"},
		}
	}

	// Extract request information from execution context
	var requestContext handlers.WebhookContext

	if err := states.BindExecContext(&requestContext); err != nil {
		return node.NodeResult{
			Error:    fmt.Errorf("failed to bind execution context: %w", err),
			Continue: false,
		}
	}

	// Check authentication if required
	if config.RequireAuth {
		authenticated := false

		// Check for secret in different places
		if config.Secret != "" {
			// Check Authorization header
			if auth, ok := requestContext.Headers["authorization"].(string); ok && auth == "Bearer "+config.Secret {
				authenticated = true
			}
			// Check X-Secret header
			if secret, ok := requestContext.Headers["x-secret"].(string); ok && secret == config.Secret {
				authenticated = true
			}
			// Check secret query parameter
			if secret, ok := requestContext.Query["secret"].(string); ok && secret == config.Secret {
				authenticated = true
			}
		}

		if !authenticated {
			return node.NodeResult{
				Output:        states.ExecContext(),
				Error:         fmt.Errorf("authentication failed"),
				Continue:      true,
				OutputHandles: []string{"unauthorized"},
			}
		}
	}

	// Check if method is allowed
	if config.Method != "" && requestContext.Method != "" && requestContext.Method != config.Method {
		return node.NodeResult{
			Output:        states.ExecContext(),
			Error:         fmt.Errorf("method %s not allowed, expected %s", requestContext.Method, config.Method),
			Continue:      true,
			OutputHandles: []string{"unauthorized"},
		}
	}

	return node.NodeResult{
		Output:        states.ExecContext(),
		Error:         nil,
		Continue:      true,
		OutputHandles: []string{"success"},
	}
}

// GetWebhookPath returns the webhook path for HTTP server registration.
func (w *WebhookTriggerNode) GetWebhookPath() string {
	props := w.GetProperties(node.PropertyContext{})
	for _, prop := range props {
		if prop.Key == "path" {
			if path, ok := prop.Value.(string); ok {
				return path
			}
		}
	}
	return ""
}

// GetAllowedMethod returns the allowed HTTP method.
func (w *WebhookTriggerNode) GetAllowedMethod() string {
	props := w.GetProperties(node.PropertyContext{})
	for _, prop := range props {
		if prop.Key == "method" {
			if method, ok := prop.Value.(string); ok {
				return method
			}
		}
	}
	return "POST"
}
