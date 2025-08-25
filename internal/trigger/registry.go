package trigger

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/ironpark/teatime/internal/node"
)

// Registry manages trigger handler plugins and provides dynamic registration.
// It stores singleton handler instances and handles their initialization.
type Registry struct {
	handlers map[string]Handler
	mu       sync.RWMutex
}

// NewRegistry creates a new empty handler registry.
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]Handler),
	}
}

// Register registers a handler instance with the registry and initializes it.
// Returns an error if the handler type is already registered or if handler is not a pointer.
func (r *Registry) Register(ctx context.Context, handler Handler, eventCh chan<- Event) error {
	// Validate that handler is a pointer
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}
	
	handlerValue := reflect.ValueOf(handler)
	if handlerValue.Kind() != reflect.Pointer {
		return fmt.Errorf("handler must be a pointer type, got %T", handler)
	}
	
	if handlerValue.IsNil() {
		return fmt.Errorf("handler pointer cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	nodeRef := handler.NodeRef()
	
	// Validate that the node reference exists in the node registry
	if !node.TriggerExists(nodeRef) {
		return fmt.Errorf("trigger node '%s' is not registered in the node registry", nodeRef)
	}
	
	if _, exists := r.handlers[nodeRef]; exists {
		return fmt.Errorf("handler for node '%s' already registered", nodeRef)
	}

	// Initialize the handler with event channel
	if ctx != nil {
		if err := handler.Initialize(ctx, eventCh); err != nil {
			return fmt.Errorf("failed to initialize handler %s: %w", nodeRef, err)
		}
	}

	r.handlers[nodeRef] = handler
	return nil
}

// UnregisterHandler removes a handler from the registry.
func (r *Registry) UnregisterHandler(nodeRef string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.handlers, nodeRef)
}

// GetHandler returns the registered handler instance.
func (r *Registry) GetHandler(nodeRef string) (Handler, error) {
	r.mu.RLock()
	handler, exists := r.handlers[nodeRef]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown handler for node: %s", nodeRef)
	}

	return handler, nil
}

// List returns metadata for all registered handler types.
func (r *Registry) List() []HandlerInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handlers := make([]HandlerInfo, 0, len(r.handlers))
	for nodeRef, handler := range r.handlers {
		handlers = append(handlers, HandlerInfo{
			NodeRef:     nodeRef,
			Name:        handler.Name(),
			Description: handler.Description(),
		})
	}
	return handlers
}

// HandlerInfo contains metadata for a trigger handler.
type HandlerInfo struct {
	NodeRef     string // The node reference this handler manages
	Name        string // Human-readable name
	Description string // Description of the handler's functionality
}

// SupportsNode checks if a node reference is supported
func (r *Registry) SupportsNode(nodeRef string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.handlers[nodeRef]
	return exists
}

// DefaultRegistry is the global registry instance.
// Use this for simple cases where a single registry is sufficient.
var DefaultRegistry = NewRegistry()