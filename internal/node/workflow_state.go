// Package node provides the core node system for workflow automation.
package node

import "time"

// WorkflowState represents the global state of a workflow execution.
// It provides thread-safe access to node inputs/outputs and execution context.
type WorkflowState interface {
	// Execution context methods
	BindExecContext(value any) error
	ExecContext() map[string]any

	// Node I/O methods
	GetOutput(nodeId string, key string) any
	GetInput(nodeId string, key string) any
	SetOutput(nodeId string, key string, value any)
	SetOutputs(nodeId string, outputs map[string]any)
	SetInput(nodeId string, key string, value any)
	SetInputs(nodeId string, inputs map[string]any)

	// General state methods
	Get(key string) any
	Set(key string, value any)
	ToMap() map[string]any

	// History tracking
	GetHistory() []StateChange
}

// StateChange represents a change in the workflow state for history tracking.
type StateChange struct {
	Key       string    `json:"key"`
	OldValue  any       `json:"oldValue"`
	NewValue  any       `json:"newValue"`
	Timestamp time.Time `json:"timestamp"`
}
