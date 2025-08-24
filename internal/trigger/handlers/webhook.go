package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/trigger"
)

// WebhookHandler handles HTTP webhook triggers
type WebhookHandler struct {
	manager  *trigger.Manager
	server   *http.Server
	mux      *http.ServeMux
	routes   map[string]string // path+method -> triggerID
	routesMu sync.RWMutex
	serverMu sync.Mutex
}

// WebhookConfig represents webhook configuration
type WebhookConfig struct {
	Path   string `mapstructure:"path"`
	Method string `mapstructure:"method"`
}

// Validate validates the webhook configuration
func (c *WebhookConfig) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("webhook path is required")
	}

	if c.Method == "" {
		c.Method = "POST" // Set default if empty
	}

	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	for _, validMethod := range validMethods {
		if c.Method == validMethod {
			return nil
		}
	}
	
	return fmt.Errorf("invalid HTTP method: %s", c.Method)
}

func (h *WebhookHandler) Initialize(ctx context.Context, manager *trigger.Manager) error {
	h.manager = manager
	h.mux = http.NewServeMux()
	h.routes = make(map[string]string)
	return nil
}

func (h *WebhookHandler) Run(ctx context.Context) error {
	h.serverMu.Lock()
	if h.server != nil {
		h.serverMu.Unlock()
		return nil // Server already running
	}

	h.server = &http.Server{
		Addr:    ":8080",
		Handler: h.mux,
	}
	h.serverMu.Unlock()

	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start webhook server: %v\n", err)
		}
	}()

	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if h.server != nil {
		if err := h.server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("Error shutting down webhook server: %v\n", err)
		}
	}

	return nil
}

func (h *WebhookHandler) Type() trigger.TriggerType {
	return trigger.TypeWebhook
}

func (h *WebhookHandler) Validate(configMap map[string]any) error {
	// Use mapstructure directly for validation
	var config WebhookConfig
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &config,
		TagName: "mapstructure",
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(configMap); err != nil {
		return fmt.Errorf("invalid webhook config: %w", err)
	}

	// Validate using WebhookConfig's Validate method
	if err := config.Validate(); err != nil {
		return err
	}

	// Update configMap with any defaults set by validation
	configMap["method"] = config.Method

	return nil
}

func (h *WebhookHandler) Register(ctx context.Context, instance *trigger.Instance) error {
	if h.mux == nil {
		return fmt.Errorf("webhook handler not initialized")
	}

	var config WebhookConfig
	if err := instance.Bind(&config); err != nil {
		return fmt.Errorf("failed to bind webhook config: %w", err)
	}

	// Create route key
	routeKey := config.Method + ":" + config.Path

	h.routesMu.Lock()
	defer h.routesMu.Unlock()

	if existingTriggerID, exists := h.routes[routeKey]; exists {
		return fmt.Errorf("route '%s %s' is already registered by trigger %s", config.Method, config.Path, existingTriggerID)
	}

	h.routes[routeKey] = instance.ID
	h.mux.HandleFunc(config.Path, h.createHandler(instance.ID, config.Method))

	instance.SetCleanup(func() error {
		h.routesMu.Lock()
		defer h.routesMu.Unlock()

		routeKey := config.Method + ":" + config.Path
		delete(h.routes, routeKey)
		return nil
	})

	return nil
}

func (h *WebhookHandler) Unregister(instance *trigger.Instance) error {
	return nil
}

func (h *WebhookHandler) createHandler(triggerID string, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if method matches
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := map[string]any{
			"method":    r.Method,
			"path":      r.URL.Path,
			"query":     h.parseQuery(r.URL.Query()),
			"headers":   h.parseHeaders(r.Header),
			"timestamp": time.Now(),
		}

		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err == nil && len(body) > 0 {
				data["body"] = string(body)
			}
		}

		if h.manager != nil {
			h.manager.ExecuteTrigger(triggerID, data)
		}

		response := map[string]any{
			"status":     "triggered",
			"trigger_id": triggerID,
			"timestamp":  time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func (h *WebhookHandler) parseQuery(query url.Values) map[string]any {
	result := make(map[string]any)
	for key, values := range query {
		if len(values) == 1 {
			result[key] = values[0]
		} else {
			result[key] = values
		}
	}
	return result
}

func (h *WebhookHandler) parseHeaders(headers http.Header) map[string]any {
	result := make(map[string]any)
	for key, values := range headers {
		if len(values) == 1 {
			result[key] = values[0]
		} else {
			result[key] = values
		}
	}
	return result
}
