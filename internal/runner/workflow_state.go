// Package runner provides workflow execution state management.
package runner

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/node"
)

// workflowState implements the node.WorkflowState interface.
// It stores all node inputs and outputs with structured keys for easy access.
type WorkflowState struct {
	data    map[string]any
	execCtx map[string]any
	history []node.StateChange
	mu      sync.RWMutex
}

// NewWorkflowState creates a new WorkflowState instance.
func NewWorkflowState() *WorkflowState {
	return &WorkflowState{
		data:    make(map[string]any),
		execCtx: make(map[string]any),
		history: make([]node.StateChange, 0),
	}
}

func (ws *WorkflowState) BindExecContext(value any) error {
	if value == nil {
		return nil
	}
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	// decode execCtx to value
	err := mapstructure.Decode(ws.execCtx, value)
	if err != nil {
		return err
	}
	return nil
}

// SetExecContext sets execution context data from a map.
func (ws *WorkflowState) SetExecContext(contextMap map[string]any) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for key, value := range contextMap {
		ws.execCtx[key] = value
	}
}

func (ws *WorkflowState) ExecContext() map[string]any {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.execCtx
}

// GetOutput retrieves an output value from a specific node.
func (ws *WorkflowState) GetOutput(nodeId string, key string) any {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.data[fmt.Sprintf("%s.output.%s", nodeId, key)]
}

// GetInput retrieves an input value from a specific node.
func (ws *WorkflowState) GetInput(nodeId string, key string) any {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.data[fmt.Sprintf("%s.input.%s", nodeId, key)]
}

// SetOutput sets an output value for a specific node.
func (ws *WorkflowState) SetOutput(nodeId string, key string, value any) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	stateKey := fmt.Sprintf("%s.output.%s", nodeId, key)
	oldValue := ws.data[stateKey]
	ws.data[stateKey] = value

	ws.history = append(ws.history, node.StateChange{
		Key:       stateKey,
		OldValue:  oldValue,
		NewValue:  value,
		Timestamp: time.Now(),
	})
}

// SetOutputs sets multiple output values for a specific node.
func (ws *WorkflowState) SetOutputs(nodeId string, outputs map[string]any) {
	for key, value := range outputs {
		ws.SetOutput(nodeId, key, value)
	}
}

// SetInput sets an input value for a specific node.
func (ws *WorkflowState) SetInput(nodeId string, key string, value any) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	stateKey := fmt.Sprintf("%s.input.%s", nodeId, key)
	oldValue := ws.data[stateKey]
	ws.data[stateKey] = value

	ws.history = append(ws.history, node.StateChange{
		Key:       stateKey,
		OldValue:  oldValue,
		NewValue:  value,
		Timestamp: time.Now(),
	})
}

// SetInputs sets multiple input values for a specific node.
func (ws *WorkflowState) SetInputs(nodeId string, inputs map[string]any) {
	for key, value := range inputs {
		ws.SetInput(nodeId, key, value)
	}
}

// Get retrieves a value by key (for backward compatibility).
func (ws *WorkflowState) Get(key string) any {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.data[key]
}

// Set sets a value by key (for backward compatibility).
func (ws *WorkflowState) Set(key string, value any) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	oldValue := ws.data[key]
	ws.data[key] = value

	ws.history = append(ws.history, node.StateChange{
		Key:       key,
		OldValue:  oldValue,
		NewValue:  value,
		Timestamp: time.Now(),
	})
}

// GetHistory returns a copy of the state change history.
func (ws *WorkflowState) GetHistory() []node.StateChange {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	history := make([]node.StateChange, len(ws.history))
	copy(history, ws.history)
	return history
}

// ToMap returns a copy of the internal data map (for backward compatibility).
func (ws *WorkflowState) ToMap() map[string]any {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	result := make(map[string]any)
	for k, v := range ws.data {
		result[k] = v
	}
	return result
}
