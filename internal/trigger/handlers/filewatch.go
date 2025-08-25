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

// FilewatchContext represents the file system event context for filewatch triggers.
type FilewatchContext struct {
	Path      string    `mapstructure:"path"`
	Operation string    `mapstructure:"operation"`
	Timestamp time.Time `mapstructure:"timestamp"`
	Created   bool      `mapstructure:"created"`
	Modified  bool      `mapstructure:"modified"`
	Removed   bool      `mapstructure:"removed"`
	Renamed   bool      `mapstructure:"renamed"`
	Chmod     bool      `mapstructure:"chmod"`
}

type FileWatchConfig struct {
	Path        string `mapstructure:"path"`
	WatchCreate bool   `mapstructure:"watchCreate"`
	WatchModify bool   `mapstructure:"watchModify"`
	WatchRemove bool   `mapstructure:"watchRemove"`
	WatchRename bool   `mapstructure:"watchRename"`
	WatchChmod  bool   `mapstructure:"watchChmod"`
}

// FileWatchHandler handles file system events
type FileWatchHandler struct {
	watcher *fsnotify.Watcher
	watches map[string]string          // path -> triggerID
	configs map[string]FileWatchConfig // triggerID -> config
	mu      sync.RWMutex
	eventCh chan<- trigger.Event
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
	h.configs = make(map[string]FileWatchConfig)
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
	h.configs[id] = config

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

	// Remove config
	delete(h.configs, id)

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
	config, configExists := h.configs[triggerID]
	h.mu.RUnlock()

	if !exists || !configExists {
		return
	}

	// Check event type flags
	created := event.Op&fsnotify.Create != 0
	modified := event.Op&fsnotify.Write != 0
	removed := event.Op&fsnotify.Remove != 0
	renamed := event.Op&fsnotify.Rename != 0
	chmod := event.Op&fsnotify.Chmod != 0

	// Filter events based on configuration
	shouldProcess := (created && config.WatchCreate) ||
		(modified && config.WatchModify) ||
		(removed && config.WatchRemove) ||
		(renamed && config.WatchRename) ||
		(chmod && config.WatchChmod)

	// Skip if this event type is not configured to be watched
	if !shouldProcess {
		return
	}

	data := map[string]any{
		"path":      event.Name,
		"operation": event.Op.String(),
		"timestamp": time.Now(),
		"created":   created,
		"modified":  modified,
		"removed":   removed,
		"renamed":   renamed,
		"chmod":     chmod,
	}

	if h.eventCh != nil {
		triggerEvent := trigger.Event{
			TriggerID:   triggerID,
			Data:        data,
			TriggeredAt: time.Now(),
		}
		select {
		case h.eventCh <- triggerEvent:
		default:
			fmt.Printf("Warning: event channel full for trigger %s\n", triggerID)
		}
	}
}
