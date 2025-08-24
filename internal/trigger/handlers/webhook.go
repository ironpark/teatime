package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ironpark/teatime/internal/trigger"
)

// WebhookHandler handles HTTP webhook triggers
type WebhookHandler struct {
	manager *trigger.Manager
	router  *gin.Engine
}

func (h *WebhookHandler) Initialize(ctx context.Context, manager *trigger.Manager) error {
	h.manager = manager
	return nil
}

func (h *WebhookHandler) Run(ctx context.Context) error {
	// Webhook handler doesn't need a background process
	return nil
}

func (h *WebhookHandler) Type() trigger.TriggerType {
	return trigger.TypeWebhook
}

func (h *WebhookHandler) Validate(config map[string]interface{}) error {
	path, ok := config["path"].(string)
	if !ok || path == "" {
		return fmt.Errorf("webhook path is required")
	}

	method, ok := config["method"].(string)
	if !ok {
		config["method"] = "POST" // default
	} else {
		validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
		valid := false
		for _, validMethod := range validMethods {
			if method == validMethod {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid HTTP method: %s", method)
		}
	}

	return nil
}

func (h *WebhookHandler) Register(ctx context.Context, instance *trigger.Instance) error {
	if h.router == nil {
		h.router = gin.New()
		h.router.Use(gin.Recovery())
		
		go func() {
			if err := h.router.Run(":8080"); err != nil {
				fmt.Printf("Failed to start webhook server: %v\n", err)
			}
		}()
	}

	path := instance.Config["path"].(string)
	method := instance.Config["method"].(string)

	h.router.Handle(method, path, h.createHandler(instance.ID))

	instance.SetCleanup(func() error {
		return nil
	})

	return nil
}

func (h *WebhookHandler) Unregister(instance *trigger.Instance) error {
	return nil
}

func (h *WebhookHandler) createHandler(triggerID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]interface{}{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"query":     c.Request.URL.Query(),
			"headers":   c.Request.Header,
			"timestamp": time.Now(),
		}

		if c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil && len(body) > 0 {
				data["body"] = string(body)
			}
		}

		if len(c.Params) > 0 {
			params := make(map[string]string)
			for _, param := range c.Params {
				params[param.Key] = param.Value
			}
			data["params"] = params
		}

		if h.manager != nil {
			h.manager.ExecuteTrigger(triggerID, data)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "triggered",
			"trigger_id":  triggerID,
			"timestamp":   time.Now(),
		})
	}
}