package trigger

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// WebhookHandler handles HTTP webhook triggers
type WebhookHandler struct {
	manager *Manager
	router  *gin.Engine
}

func (h *WebhookHandler) SetManager(manager interface{}) {
	if m, ok := manager.(*Manager); ok {
		h.manager = m
	}
}

func (h *WebhookHandler) Type() TriggerType {
	return TypeWebhook
}

func (h *WebhookHandler) Validate(config map[string]interface{}) error {
	path, ok := config["path"].(string)
	if !ok || path == "" {
		return fmt.Errorf("webhook path is required")
	}

	method, ok := config["method"].(string)
	if !ok {
		config["method"] = "POST" // default
	} else {
		validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
		valid := false
		for _, validMethod := range validMethods {
			if method == validMethod {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid HTTP method: %s", method)
		}
	}

	return nil
}

func (h *WebhookHandler) Register(ctx context.Context, instance *Instance) error {
	if h.router == nil {
		h.router = gin.New()
		h.router.Use(gin.Recovery())
		
		go func() {
			if err := h.router.Run(":8080"); err != nil {
				fmt.Printf("Failed to start webhook server: %v\n", err)
			}
		}()
	}

	path := instance.Config["path"].(string)
	method := instance.Config["method"].(string)

	h.router.Handle(method, path, h.createHandler(instance.ID))

	instance.SetCleanup(func() error {
		return nil
	})

	return nil
}

func (h *WebhookHandler) Unregister(instance *Instance) error {
	return nil
}

func (h *WebhookHandler) createHandler(triggerID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]interface{}{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"query":     c.Request.URL.Query(),
			"headers":   c.Request.Header,
			"timestamp": time.Now(),
		}

		if c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil && len(body) > 0 {
				data["body"] = string(body)
			}
		}

		if len(c.Params) > 0 {
			params := make(map[string]string)
			for _, param := range c.Params {
				params[param.Key] = param.Value
			}
			data["params"] = params
		}

		if h.manager != nil {
			h.manager.ExecuteTrigger(triggerID, data)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "triggered",
			"trigger_id":  triggerID,
			"timestamp":   time.Now(),
		})
	}
}

// ScheduleHandler handles cron-based schedule triggers
type ScheduleHandler struct {
	manager   *Manager
	scheduler *cron.Cron
}

func (h *ScheduleHandler) SetManager(manager interface{}) {
	if m, ok := manager.(*Manager); ok {
		h.manager = m
		h.scheduler = cron.New(cron.WithSeconds())
		h.scheduler.Start()
	}
}

func (h *ScheduleHandler) Type() TriggerType {
	return TypeSchedule
}

func (h *ScheduleHandler) Validate(config map[string]interface{}) error {
	cronExpr, ok := config["cron"].(string)
	if !ok || cronExpr == "" {
		return fmt.Errorf("cron expression is required")
	}

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := parser.Parse(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression '%s': %w", cronExpr, err)
	}

	return nil
}

func (h *ScheduleHandler) Register(ctx context.Context, instance *Instance) error {
	if h.scheduler == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	cronExpr := instance.Config["cron"].(string)
	
	entryID, err := h.scheduler.AddFunc(cronExpr, func() {
		if h.manager != nil {
			data := map[string]interface{}{
				"timestamp": time.Now(),
				"cron":      cronExpr,
				"scheduled": true,
			}
			h.manager.ExecuteTrigger(instance.ID, data)
		}
	})
	
	if err != nil {
		return fmt.Errorf("failed to schedule cron job: %w", err)
	}

	instance.SetCleanup(func() error {
		if h.scheduler != nil {
			h.scheduler.Remove(entryID)
		}
		return nil
	})

	instance.Config["entryID"] = int(entryID)

	fmt.Printf("Scheduled cron job: %s (ID: %d)\n", cronExpr, entryID)
	return nil
}

func (h *ScheduleHandler) Unregister(instance *Instance) error {
	return nil
}

func (h *ScheduleHandler) Shutdown() {
	if h.scheduler != nil {
		ctx := h.scheduler.Stop()
		<-ctx.Done()
	}
}

// CommandHandler handles command-based triggers
type CommandHandler struct {
	manager           *Manager
	registeredCommands map[string]string
	mu                sync.RWMutex
}

func (h *CommandHandler) SetManager(manager interface{}) {
	if m, ok := manager.(*Manager); ok {
		h.manager = m
		h.registeredCommands = make(map[string]string)
	}
}

func (h *CommandHandler) Type() TriggerType {
	return TypeCommand
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

func (h *CommandHandler) Register(ctx context.Context, instance *Instance) error {
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

func (h *CommandHandler) Unregister(instance *Instance) error {
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

// FileWatchHandler handles file system events
type FileWatchHandler struct {
	manager *Manager
	watcher *fsnotify.Watcher
	watches map[string]string
	mu      sync.RWMutex
}

func (h *FileWatchHandler) SetManager(manager interface{}) {
	if m, ok := manager.(*Manager); ok {
		h.manager = m
		h.watches = make(map[string]string)
		
		var err error
		h.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			fmt.Printf("Failed to create file watcher: %v\n", err)
			return
		}
		
		go h.watchEvents()
	}
}

func (h *FileWatchHandler) Type() TriggerType {
	return TypeFileWatch
}

func (h *FileWatchHandler) Validate(config map[string]interface{}) error {
	path, ok := config["path"].(string)
	if !ok || path == "" {
		return fmt.Errorf("file path is required")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}
	
	config["path"] = absPath
	return nil
}

func (h *FileWatchHandler) Register(ctx context.Context, instance *Instance) error {
	if h.watcher == nil {
		return fmt.Errorf("file watcher not initialized")
	}

	path := instance.Config["path"].(string)
	
	h.mu.Lock()
	defer h.mu.Unlock()

	if existingTriggerID, exists := h.watches[path]; exists {
		return fmt.Errorf("path '%s' is already being watched by trigger %s", path, existingTriggerID)
	}

	err := h.watcher.Add(path)
	if err != nil {
		return fmt.Errorf("failed to watch path '%s': %w", path, err)
	}

	h.watches[path] = instance.ID

	instance.SetCleanup(func() error {
		h.mu.Lock()
		defer h.mu.Unlock()
		
		if h.watcher != nil {
			h.watcher.Remove(path)
		}
		delete(h.watches, path)
		return nil
	})

	fmt.Printf("Watching file path: %s\n", path)
	return nil
}

func (h *FileWatchHandler) Unregister(instance *Instance) error {
	return nil
}

func (h *FileWatchHandler) watchEvents() {
	if h.watcher == nil {
		return
	}

	for {
		select {
		case event, ok := <-h.watcher.Events:
			if !ok {
				return
			}
			h.handleFileEvent(event)

		case err, ok := <-h.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("File watcher error: %v\n", err)
		}
	}
}

func (h *FileWatchHandler) handleFileEvent(event fsnotify.Event) {
	h.mu.RLock()
	triggerID, exists := h.watches[event.Name]
	h.mu.RUnlock()

	if !exists {
		return
	}

	data := map[string]interface{}{
		"path":      event.Name,
		"operation": event.Op.String(),
		"timestamp": time.Now(),
	}

	data["created"] = event.Op&fsnotify.Create != 0
	data["modified"] = event.Op&fsnotify.Write != 0
	data["removed"] = event.Op&fsnotify.Remove != 0
	data["renamed"] = event.Op&fsnotify.Rename != 0
	data["chmod"] = event.Op&fsnotify.Chmod != 0

	if h.manager != nil {
		h.manager.ExecuteTrigger(triggerID, data)
	}
}

func (h *FileWatchHandler) Shutdown() {
	if h.watcher != nil {
		h.watcher.Close()
	}
}
