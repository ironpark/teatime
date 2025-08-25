package handlers

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ironpark/teatime/internal/trigger"
)

// CommandHandler handles command-based triggers
type CommandHandler struct {
	manager            *trigger.Manager
	registeredCommands map[string]string
	eventCh            chan<- trigger.Event
	mu                 sync.RWMutex
}

// CommandConfig represents command configuration
type CommandConfig struct {
	Command string `mapstructure:"command"`
	Global  bool   `mapstructure:"global"`
}

// Validate validates the command configuration
func (c *CommandConfig) Validate() error {
	if c.Command == "" {
		return fmt.Errorf("command name is required")
	}

	if len(c.Command) < 1 {
		return fmt.Errorf("command name must be at least 1 character")
	}

	return nil
}

func (h *CommandHandler) Initialize(ctx context.Context, eventCh chan<- trigger.Event) error {
	h.registeredCommands = make(map[string]string)
	h.eventCh = eventCh
	return nil
}

func (h *CommandHandler) Start(ctx context.Context) error {
	// Command handler doesn't need a background process
	return nil
}

func (h *CommandHandler) NodeRef() string {
	return "teatime.trigger.command"
}

func (h *CommandHandler) Name() string {
	return "Manual Command"
}

func (h *CommandHandler) Description() string {
	return "Triggers workflows via manual command execution"
}


func (h *CommandHandler) Register(ctx context.Context, id string, configMap map[string]any) error {
	var config CommandConfig
	if err := trigger.BindAndValidate(&config, configMap); err != nil {
		return fmt.Errorf("failed to validate command config: %w", err)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if existingTriggerID, exists := h.registeredCommands[config.Command]; exists {
		return fmt.Errorf("command '%s' is already registered by trigger %s", config.Command, existingTriggerID)
	}

	h.registeredCommands[config.Command] = id

	fmt.Printf("Registered command: %s (global: %v)\n", config.Command, config.Global)

	return nil
}

func (h *CommandHandler) Unregister(ctx context.Context, id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Find and remove the command associated with this trigger ID
	for command, triggerID := range h.registeredCommands {
		if triggerID == id {
			delete(h.registeredCommands, command)
			fmt.Printf("Unregistered command: %s\n", command)
			break
		}
	}
	
	return nil
}

func (h *CommandHandler) ExecuteCommand(ctx context.Context, command string, args map[string]any) error {
	h.mu.RLock()
	triggerID, exists := h.registeredCommands[command]
	h.mu.RUnlock()

	if !exists {
		return fmt.Errorf("command not found: %s", command)
	}

	data := map[string]any{
		"command":   command,
		"args":      args,
		"timestamp": time.Now(),
		"workdir":   getCurrentWorkingDirectory(),
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
			// Channel is full, log warning
			fmt.Printf("Warning: event channel full for trigger %s\n", triggerID)
		}
	}

	return nil
}

func (h *CommandHandler) GetRegisteredCommands() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	commands := make([]string, 0, len(h.registeredCommands))
	for command := range h.registeredCommands {
		commands = append(commands, command)
	}

	return commands
}

func getCurrentWorkingDirectory() string {
	dir, _ := os.Getwd()
	return dir
}
