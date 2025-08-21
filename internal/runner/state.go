package runner

import (
	"fmt"
	"sync"

	"github.com/ironpark/teatime/internal/recipe"
)

// runState maintains the execution state of a recipe run.
// It tracks node inputs/outputs, execution status, and synchronization.
type runState struct {
	recipe *recipe.Recipe
	// states stores all node I/O values with keys:
	// "{nodeId}.input.{propertyKey}" for node inputs
	// "{nodeId}.output.{outputKey}" for node outputs
	states       map[string]any
	nodeExecuted map[string]bool
	// nodeResults contains channels for signaling node completion
	nodeResults map[string]chan error
	lock        sync.RWMutex
	waitGroup   sync.WaitGroup
}

// setNodeOutput stores a node's output values in the state map.
// Keys are formatted as "{nodeId}.output.{key}".
func (state *runState) setNodeOutput(nodeId string, output map[string]any) {
	state.lock.Lock()
	defer state.lock.Unlock()
	for key, value := range output {
		state.states[fmt.Sprintf("%s.output.%s", nodeId, key)] = value
	}
}

// setNodeInput stores a node's input values in the state map.
// Keys are formatted as "{nodeId}.input.{key}".
func (state *runState) setNodeInput(nodeId string, input map[string]any) {
	state.lock.Lock()
	defer state.lock.Unlock()
	for key, value := range input {
		state.states[fmt.Sprintf("%s.input.%s", nodeId, key)] = value
	}
}

// getNodeResultChannel returns the result channel for a node,
// creating it if it doesn't exist. Used for dependency synchronization.
func (state *runState) getNodeResultChannel(nodeId string) chan error {
	state.lock.Lock()
	defer state.lock.Unlock()

	if state.nodeResults[nodeId] == nil {
		state.nodeResults[nodeId] = make(chan error, 1)
	}
	return state.nodeResults[nodeId]
}

// setNodeExecuted marks a node as executed.
// Returns false if the node was already executed, true otherwise.
func (state *runState) setNodeExecuted(nodeId string) bool {
	state.lock.Lock()
	defer state.lock.Unlock()
	if state.nodeExecuted[nodeId] {
		return false
	}
	state.nodeExecuted[nodeId] = true
	return true
}

// waitDependencies blocks until all dependency nodes complete.
// Returns immediately if no dependencies exist.
// Returns an error if any dependency fails.
func (state *runState) waitDependencies(nodeId string) error {
	// check dependencies
	dependencies, err := state.recipe.GetNodeDependencies(nodeId)
	if err != nil {
		return err
	}
	if len(dependencies) > 0 {
		return nil
	}
	for _, dependency := range dependencies {
		err := <-state.getNodeResultChannel(dependency)
		if err != nil {
			return err
		}
	}
	return nil
}
