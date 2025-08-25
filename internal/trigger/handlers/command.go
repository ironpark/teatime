package handlers

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/trigger"
)

// CommandHandler handles command-based triggers
type CommandHandler struct {
	manager            *trigger.Manager
	registeredCommands map[string]string
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

func (h *CommandHandler) Initialize(manager *trigger.Manager) error {
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

func (h *CommandHandler) Validate(configMap map[string]any) error {
	// Use mapstructure directly for validation
	var config CommandConfig
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &config,
		TagName: "mapstructure",
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(configMap); err != nil {
		return fmt.Errorf("invalid command config: %w", err)
	}

	// Validate using CommandConfig's Validate method
	return config.Validate()
}

func (h *CommandHandler) Register(instance *trigger.Instance) error {
	var config CommandConfig
	if err := instance.Bind(&config); err != nil {
		return fmt.Errorf("failed to bind command config: %w", err)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if existingTriggerID, exists := h.registeredCommands[config.Command]; exists {
		return fmt.Errorf("command '%s' is already registered by trigger %s", config.Command, existingTriggerID)
	}

	h.registeredCommands[config.Command] = instance.ID

	instance.SetCleanup(func() error {
		h.mu.Lock()
		defer h.mu.Unlock()
		delete(h.registeredCommands, config.Command)
		return nil
	})

	fmt.Printf("Registered command: %s (global: %v)\n", config.Command, config.Global)

	return nil
}

func (h *CommandHandler) Unregister(instance *trigger.Instance) error {
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

	if h.manager != nil {
		h.manager.ExecuteTrigger(ctx, triggerID, data)
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
