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

func (h *FileWatchHandler) Register(ctx context.Context, instance *trigger.Instance) error {
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

