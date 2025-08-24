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
	manager           *trigger.Manager
	registeredCommands map[string]string
	mu                sync.RWMutex
}

func (h *CommandHandler) Initialize(ctx context.Context, manager *trigger.Manager) error {
	h.manager = manager
	h.registeredCommands = make(map[string]string)
	return nil
}

func (h *CommandHandler) Run(ctx context.Context) error {
	// Command handler doesn't need a background process
	return nil
}

func (h *CommandHandler) Type() trigger.TriggerType {
	return trigger.TypeCommand
}

func (h *CommandHandler) Validate(config map[string]interface{}) error {
	command, ok := config["command"].(string)
	if !ok || command == "" {
		return fmt.Errorf("command name is required")
	}

	if len(command) < 1 {
		return fmt.Errorf("command name must be at least 1 character")
	}

	return nil
}

func (h *CommandHandler) Register(ctx context.Context, instance *trigger.Instance) error {
	command := instance.Config["command"].(string)
	global, _ := instance.Config["global"].(bool)

	h.mu.Lock()
	defer h.mu.Unlock()

	if existingTriggerID, exists := h.registeredCommands[command]; exists {
		return fmt.Errorf("command '%s' is already registered by trigger %s", command, existingTriggerID)
	}

	h.registeredCommands[command] = instance.ID

	instance.SetCleanup(func() error {
		h.mu.Lock()
		defer h.mu.Unlock()
		delete(h.registeredCommands, command)
		return nil
	})

	fmt.Printf("Registered command: %s (global: %v)\n", command, global)

	return nil
}

func (h *CommandHandler) Unregister(instance *trigger.Instance) error {
	return nil
}

func (h *CommandHandler) ExecuteCommand(command string, args map[string]interface{}) error {
	h.mu.RLock()
	triggerID, exists := h.registeredCommands[command]
	h.mu.RUnlock()

	if !exists {
		return fmt.Errorf("command not found: %s", command)
	}

	data := map[string]interface{}{
		"command":   command,
		"args":      args,
		"timestamp": time.Now(),
		"workdir":   getCurrentWorkingDirectory(),
	}

	if h.manager != nil {
		h.manager.ExecuteTrigger(triggerID, data)
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