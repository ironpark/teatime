package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/mapstructure/v2"
	"github.com/ironpark/teatime/internal/trigger"
)

// FileWatchHandler handles file system events
type FileWatchHandler struct {
	manager *trigger.Manager
	watcher *fsnotify.Watcher
	watches map[string]string
	mu      sync.RWMutex
}

// FileWatchConfig represents file watch configuration
type FileWatchConfig struct {
	Path string `mapstructure:"path"`
}

// Validate validates the file watch configuration
func (c *FileWatchConfig) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("file path is required")
	}

	absPath, err := filepath.Abs(c.Path)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}
	
	c.Path = absPath
	return nil
}

func (h *FileWatchHandler) Initialize(ctx context.Context, manager *trigger.Manager) error {
	h.manager = manager
	h.watches = make(map[string]string)
	
	var err error
	h.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}
	
	return nil
}

func (h *FileWatchHandler) Run(ctx context.Context) error {
	if h.watcher == nil {
		return fmt.Errorf("file watcher not initialized")
	}
	
	go h.watchEvents()
	<-ctx.Done()
	
	if h.watcher != nil {
		h.watcher.Close()
	}
	
	return nil
}

func (h *FileWatchHandler) Type() trigger.TriggerType {
	return trigger.TypeFileWatch
}

func (h *FileWatchHandler) Validate(configMap map[string]any) error {
	// Use mapstructure directly for validation
	var config FileWatchConfig
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &config,
		TagName: "mapstructure",
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(configMap); err != nil {
		return fmt.Errorf("invalid file watch config: %w", err)
	}

	// Validate using FileWatchConfig's Validate method
	if err := config.Validate(); err != nil {
		return err
	}

	// Update configMap with any defaults set by validation
	configMap["path"] = config.Path

	return nil
}

func (h *FileWatchHandler) Register(ctx context.Context, instance *trigger.Instance) error {
	if h.watcher == nil {
		return fmt.Errorf("file watcher not initialized")
	}

	var config FileWatchConfig
	if err := instance.Bind(&config); err != nil {
		return fmt.Errorf("failed to bind file watch config: %w", err)
	}
	
	h.mu.Lock()
	defer h.mu.Unlock()

	if existingTriggerID, exists := h.watches[config.Path]; exists {
		return fmt.Errorf("path '%s' is already being watched by trigger %s", config.Path, existingTriggerID)
	}

	err := h.watcher.Add(config.Path)
	if err != nil {
		return fmt.Errorf("failed to watch path '%s': %w", config.Path, err)
	}

	h.watches[config.Path] = instance.ID

	instance.SetCleanup(func() error {
		h.mu.Lock()
		defer h.mu.Unlock()
		
		if h.watcher != nil {
			h.watcher.Remove(config.Path)
		}
		delete(h.watches, config.Path)
		return nil
	})

	fmt.Printf("Watching file path: %s\n", config.Path)
	return nil
}

func (h *FileWatchHandler) Unregister(instance *trigger.Instance) error {
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

	data := map[string]any{
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

