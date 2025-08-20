// Package runner executes recipe workflows by managing node execution order,
// dependencies, and state propagation between nodes.
package runner

import (
	"context"
	"fmt"
	"sync"

	n "github.com/ironpark/teatime/internal/node"
	"github.com/ironpark/teatime/internal/recipe"
)

type NodeExecutionStatus string

// Node execution states.
const (
	// StateWait indicates a node is waiting for dependencies.
	StateWait NodeExecutionStatus = "wait"
	// StateRun indicates a node is currently executing.
	StateRun NodeExecutionStatus = "run"
	// StateDone indicates a node has completed successfully.
	StateDone NodeExecutionStatus = "done"
	// StateError indicates a node encountered an error.
	StateError NodeExecutionStatus = "error"
)

// Callback is called during node execution to report state changes.
// It receives the recipe, current state, node being executed, output data,
// and any error that occurred.
type Callback func(recipe *recipe.Recipe, state NodeExecutionStatus, node recipe.Node, output map[string]any, err error)

// Run executes a recipe workflow starting from the specified node.
// It manages concurrent node execution while respecting dependencies.
//
// The execution flow:
//  1. Validates and starts the trigger node
//  2. Executes nodes concurrently as dependencies are met
//  3. Propagates outputs between connected nodes
//  4. Reports state changes via callback
//
// Execution stops if any node encounters an error or the context is cancelled.
func Run(ctx context.Context, target *recipe.Recipe, startNodeId string, properties map[string]any, callback Callback) error {
	triggerNode, err := target.GetNodeById(startNodeId)
	if err != nil {
		return err
	}
	if callback == nil {
		return fmt.Errorf("callback is nil")
	}
	runState := &runState{
		recipe:       target,
		states:       map[string]any{},
		nodeExecuted: map[string]bool{},
		nodeResults:  map[string]chan error{},
		lock:         sync.RWMutex{},
		waitGroup:    sync.WaitGroup{},
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = startNode(ctx, runState, triggerNode, properties, func(recipe *recipe.Recipe, state NodeExecutionStatus, node recipe.Node, output map[string]any, err error) {
		callback(recipe, state, node, output, err)
		if err != nil {
			cancel()
		}
	})
	if err != nil {
		return err
	}
	runState.waitGroup.Wait()
	return nil
}

// startNode initiates execution of a single node if not already executed.
// It validates inputs, marks the node as executed, and launches executeNode
// in a goroutine for concurrent execution.
func startNode(ctx context.Context, state *runState, node recipe.Node, properties map[string]any, callback Callback) error {
	if !state.setNodeExecuted(node.Id) {
		// node already executed
		return nil
	}
	rawNode := node.GetRawNode()
	// validate node input
	propCtx := n.PropertyContext(properties)
	err := n.ValidateProperties(rawNode.GetProperties(propCtx), properties)
	if err != nil {
		callback(state.recipe, StateError, node, nil, err)
		return err
	}
	err = rawNode.ValidateProperties(properties)
	if err != nil {
		callback(state.recipe, StateError, node, nil, err)
		return err
	}

	state.setNodeInput(node.Id, properties)
	state.waitGroup.Add(1)
	// run node in goroutine
	go executeNode(ctx, state, node, callback)
	return nil
}

// executeNode runs a node's logic after dependencies are satisfied.
// It waits for dependencies, executes the node, stores outputs,
// and triggers execution of connected downstream nodes.
func executeNode(ctx context.Context, state *runState, node recipe.Node, callback Callback) {
	defer state.waitGroup.Done()
	callback(state.recipe, StateWait, node, nil, nil)
	// wait dependencies
	err := state.waitDependencies(node.Id)
	if err != nil {
		callback(state.recipe, StateError, node, nil, err)
		return
	}
	rawNode := node.GetRawNode()
	if rawNode == nil {
		callback(state.recipe, StateError, node, nil, fmt.Errorf("node %s has no raw node", node.Id))
		return
	}

	resultChannel := state.getNodeResultChannel(node.Id)
	callback(state.recipe, StateRun, node, nil, nil)
	// run node
	result := rawNode.Run(ctx, state.states)
	if result.Error != nil {
		callback(state.recipe, StateError, node, nil, result.Error)
		resultChannel <- result.Error
		return
	}
	// if node is not allowed to continue flow, stop execution
	if !result.Continue {
		callback(state.recipe, StateDone, node, result.Output, nil)
		return
	}
	// set node output to states
	state.setNodeOutput(node.Id, result.Output)
	callback(state.recipe, StateDone, node, result.Output, nil)
	// run next nodes
	nextNodes, err := state.recipe.GetConnectedNodes(node.Id)
	if err != nil {
		callback(state.recipe, StateError, node, nil, err)
		resultChannel <- err
		return
	}
	resultChannel <- nil
	for _, nextNode := range nextNodes {
		rawNextNode := nextNode.GetRawNode()
		nextResultChannel := state.getNodeResultChannel(nextNode.Id)
		if rawNextNode == nil {
			callback(state.recipe, StateError, nextNode, nil, fmt.Errorf("next node %s has no raw node", nextNode.Id))
			nextResultChannel <- fmt.Errorf("next node %s has no raw node", nextNode.Id)
			return
		}
		nextInput, err := n.ResolveInput(rawNextNode, state.states)
		if err != nil {
			callback(state.recipe, StateError, nextNode, nil, err)
			nextResultChannel <- err
			return
		}
		startNode(ctx, state, nextNode, nextInput, callback)
	}
}
