package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/ironpark/teatime/internal/trigger"
)

// FileWatchHandler handles file system events
type FileWatchHandler struct {
	manager *trigger.Manager
	watcher *fsnotify.Watcher
	watches map[string]string
	mu      sync.RWMutex
	eventCh chan<- trigger.Event
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

func (h *FileWatchHandler) Initialize(ctx context.Context, eventCh chan<- trigger.Event) error {
	h.watches = make(map[string]string)
	h.eventCh = eventCh

	var err error
	h.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	return nil
}

func (h *FileWatchHandler) Start(ctx context.Context) error {
	if h.watcher == nil {
		return fmt.Errorf("file watcher not initialized")
	}

	h.watchEvents(ctx)
	if h.watcher != nil {
		h.watcher.Close()
	}

	return nil
}

func (h *FileWatchHandler) NodeRef() string {
	return "teatime.trigger.filewatch"
}

func (h *FileWatchHandler) Name() string {
	return "File Watcher"
}

func (h *FileWatchHandler) Description() string {
	return "Triggers workflows on file system events"
}


func (h *FileWatchHandler) Register(ctx context.Context, id string, configMap map[string]any) error {
	if h.watcher == nil {
		return fmt.Errorf("file watcher not initialized")
	}

	var config FileWatchConfig
	if err := trigger.BindAndValidate(&config, configMap); err != nil {
		return fmt.Errorf("failed to validate file watch config: %w", err)
	}

	// Update configMap with any defaults set by validation
	configMap["path"] = config.Path

	h.mu.Lock()
	defer h.mu.Unlock()

	if existingTriggerID, exists := h.watches[config.Path]; exists {
		return fmt.Errorf("path '%s' is already being watched by trigger %s", config.Path, existingTriggerID)
	}

	err := h.watcher.Add(config.Path)
	if err != nil {
		return fmt.Errorf("failed to watch path '%s': %w", config.Path, err)
	}

	h.watches[config.Path] = id

	fmt.Printf("Watching file path: %s\n", config.Path)
	return nil
}

func (h *FileWatchHandler) Unregister(ctx context.Context, id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Find and remove the path associated with this trigger ID
	var pathToRemove string
	for path, triggerID := range h.watches {
		if triggerID == id {
			pathToRemove = path
			break
		}
	}

	if pathToRemove != "" {
		if h.watcher != nil {
			h.watcher.Remove(pathToRemove)
		}
		delete(h.watches, pathToRemove)
		fmt.Printf("Stopped watching file path: %s\n", pathToRemove)
	}

	return nil
}

func (h *FileWatchHandler) watchEvents(ctx context.Context) {
	if h.watcher == nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-h.watcher.Events:
			if !ok {
				return
			}
			h.handleFileEvent(ctx, event)

		case err, ok := <-h.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("File watcher error: %v\n", err)
		}
	}
}

func (h *FileWatchHandler) handleFileEvent(_ context.Context, event fsnotify.Event) {
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

	if h.eventCh != nil {
		event := trigger.Event{
			TriggerID:   triggerID,
			Data:        data,
			TriggeredAt: time.Now(),
		}
		select {
		case h.eventCh <- event:
		default:
			fmt.Printf("Warning: event channel full for trigger %s\n", triggerID)
		}
	}
}
