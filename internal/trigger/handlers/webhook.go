package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"

	"github.com/ironpark/teatime/internal/trigger"
)

// WebhookContext represents the HTTP request context for webhook triggers.
type WebhookContext struct {
	Method     string         `mapstructure:"method"`
	Path       string         `mapstructure:"path"`
	Headers    map[string]any `mapstructure:"headers"`
	Query      map[string]any `mapstructure:"query"`
	Body       map[string]any `mapstructure:"body"`
	RemoteAddr string         `mapstructure:"remoteAddr"`
}

// WebhookConfig represents webhook configuration
type WebhookConfig struct {
	Path        string `mapstructure:"path"`
	Method      string `mapstructure:"method"`
	RequireAuth bool   `mapstructure:"requireAuth"`
	Secret      string `mapstructure:"secret"`
}

// route represents a registered webhook route
type route struct {
	pattern   string
	method    string
	handler   http.HandlerFunc
	triggerID string
}

// WebhookHandler handles HTTP webhook triggers with proper route cleanup
type WebhookHandler struct {
	server   *http.Server
	routes   []route
	routesMu sync.RWMutex
	serverMu sync.Mutex
	eventCh  chan<- trigger.Event
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
	if !slices.Contains(validMethods, c.Method) {
		return fmt.Errorf("invalid HTTP method: %s", c.Method)
	}

	return nil
}

func (h *WebhookHandler) Initialize(ctx context.Context, eventCh chan<- trigger.Event) error {
	h.routes = make([]route, 0)
	h.eventCh = eventCh
	return nil
}

func (h *WebhookHandler) Start(ctx context.Context) error {
	h.serverMu.Lock()
	if h.server != nil {
		h.serverMu.Unlock()
		return nil // Server already running
	}

	h.server = &http.Server{
		Addr:    ":8080",
		Handler: h,
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

// ServeHTTP implements http.Handler interface for custom routing
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.routesMu.RLock()
	defer h.routesMu.RUnlock()

	for _, route := range h.routes {
		if route.method == r.Method && route.pattern == r.URL.Path {
			route.handler(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func (h *WebhookHandler) NodeRef() string {
	return "teatime.trigger.webhook"
}

func (h *WebhookHandler) Name() string {
	return "HTTP Webhook"
}

func (h *WebhookHandler) Description() string {
	return "Triggers workflows via HTTP webhook endpoints"
}

func (h *WebhookHandler) Register(ctx context.Context, id string, configMap map[string]any) error {
	var config WebhookConfig
	if err := trigger.BindAndValidate(&config, configMap); err != nil {
		return fmt.Errorf("failed to validate webhook config: %w", err)
	}

	// Update configMap with any defaults set by validation
	configMap["method"] = config.Method

	h.routesMu.Lock()
	defer h.routesMu.Unlock()

	// Check for existing route
	for _, existingRoute := range h.routes {
		if existingRoute.method == config.Method && existingRoute.pattern == config.Path {
			return fmt.Errorf("route '%s %s' is already registered by trigger %s",
				config.Method, config.Path, existingRoute.triggerID)
		}
	}

	// Add new route
	newRoute := route{
		pattern:   config.Path,
		method:    config.Method,
		handler:   h.createHandler(id, config.Method),
		triggerID: id,
	}
	h.routes = append(h.routes, newRoute)

	return nil
}

func (h *WebhookHandler) Unregister(ctx context.Context, id string) error {
	h.routesMu.Lock()
	defer h.routesMu.Unlock()

	// Remove route by filtering
	filteredRoutes := make([]route, 0, len(h.routes))
	for _, route := range h.routes {
		if route.triggerID != id {
			filteredRoutes = append(filteredRoutes, route)
		}
	}
	h.routes = filteredRoutes

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

		if h.eventCh != nil {
			event := trigger.Event{
				TriggerID:   triggerID,
				Data:        data,
				TriggeredAt: time.Now(),
			}
			select {
			case h.eventCh <- event:
			default:
				fmt.Printf("Warning: event channel full for trigger %s\n", triggerID)
			}
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
