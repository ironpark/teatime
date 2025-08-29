package agent

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// InMemoryStore is a simple in-memory implementation of Memory
type InMemoryStore struct {
	conversations map[string]*Conversation
	messages      map[string][]Message  // conversationID -> messages
	messageIndex  map[string]Message    // messageID -> message
	cache         map[string]cacheEntry // simple cache
	mu            sync.RWMutex
	cacheMu       sync.RWMutex
}

type cacheEntry struct {
	value     any
	expiresAt time.Time
}

// NewInMemoryStore creates a new in-memory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		conversations: make(map[string]*Conversation),
		messages:      make(map[string][]Message),
		messageIndex:  make(map[string]Message),
		cache:         make(map[string]cacheEntry),
	}
}

// CreateConversation creates a new conversation
func (s *InMemoryStore) CreateConversation(ctx context.Context, title string) (*Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conversationID := uuid.New().String()
	now := time.Now()
	
	conv := &Conversation{
		ID:           conversationID,
		Title:        title,
		CreatedAt:    now,
		UpdatedAt:    now,
		MessageCount: 0,
		Metadata:     make(map[string]any),
		Tags:         make([]string, 0),
	}

	s.conversations[conversationID] = conv
	s.messages[conversationID] = make([]Message, 0)

	return conv, nil
}

// GetConversation retrieves a conversation by ID
func (s *InMemoryStore) GetConversation(ctx context.Context, conversationID string) (*Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv, exists := s.conversations[conversationID]
	if !exists {
		return nil, fmt.Errorf("conversation %s not found", conversationID)
	}

	// Return a copy to avoid external modifications
	result := *conv
	return &result, nil
}

// UpdateConversation updates conversation metadata
func (s *InMemoryStore) UpdateConversation(ctx context.Context, conversationID string, updates map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, exists := s.conversations[conversationID]
	if !exists {
		return fmt.Errorf("conversation %s not found", conversationID)
	}

	// Update fields
	if title, ok := updates["title"].(string); ok {
		conv.Title = title
	}
	
	if tags, ok := updates["tags"].([]string); ok {
		conv.Tags = tags
	}
	
	if metadata, ok := updates["metadata"].(map[string]any); ok {
		for k, v := range metadata {
			conv.Metadata[k] = v
		}
	}

	conv.UpdatedAt = time.Now()
	return nil
}

// DeleteConversation removes a conversation and all its messages
func (s *InMemoryStore) DeleteConversation(ctx context.Context, conversationID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.conversations[conversationID]; !exists {
		return fmt.Errorf("conversation %s not found", conversationID)
	}

	// Remove messages from index
	if messages, exists := s.messages[conversationID]; exists {
		for _, msg := range messages {
			delete(s.messageIndex, msg.ID)
		}
	}

	delete(s.conversations, conversationID)
	delete(s.messages, conversationID)

	return nil
}

// ListConversations returns a list of conversations with pagination
func (s *InMemoryStore) ListConversations(ctx context.Context, limit, offset int) ([]Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var conversations []Conversation
	for _, conv := range s.conversations {
		conversations = append(conversations, *conv)
	}

	// Sort by updated time (most recent first)
	sort.Slice(conversations, func(i, j int) bool {
		return conversations[i].UpdatedAt.After(conversations[j].UpdatedAt)
	})

	// Apply pagination
	if offset >= len(conversations) {
		return []Conversation{}, nil
	}

	end := offset + limit
	if end > len(conversations) {
		end = len(conversations)
	}

	return conversations[offset:end], nil
}

// AddMessage adds a message to a conversation
func (s *InMemoryStore) AddMessage(ctx context.Context, conversationID string, message Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, exists := s.conversations[conversationID]
	if !exists {
		return fmt.Errorf("conversation %s not found", conversationID)
	}

	// Ensure message has an ID and timestamp
	if message.ID == "" {
		message.ID = uuid.New().String()
	}
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}
	if message.Metadata == nil {
		message.Metadata = make(map[string]any)
	}

	// Add to messages and index
	s.messages[conversationID] = append(s.messages[conversationID], message)
	s.messageIndex[message.ID] = message

	// Update conversation
	conv.MessageCount++
	conv.UpdatedAt = time.Now()

	return nil
}

// GetMessages retrieves messages from a conversation with limit
func (s *InMemoryStore) GetMessages(ctx context.Context, conversationID string, limit int) ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messages, exists := s.messages[conversationID]
	if !exists {
		return []Message{}, nil
	}

	// Sort by timestamp
	sortedMessages := make([]Message, len(messages))
	copy(sortedMessages, messages)
	sort.Slice(sortedMessages, func(i, j int) bool {
		return sortedMessages[i].Timestamp.Before(sortedMessages[j].Timestamp)
	})

	if limit <= 0 || limit > len(sortedMessages) {
		return sortedMessages, nil
	}

	// Return the most recent messages
	start := len(sortedMessages) - limit
	return sortedMessages[start:], nil
}

// GetMessagesByRange retrieves messages within a time range
func (s *InMemoryStore) GetMessagesByRange(ctx context.Context, conversationID string, startTime, endTime time.Time) ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messages, exists := s.messages[conversationID]
	if !exists {
		return []Message{}, nil
	}

	var filtered []Message
	for _, msg := range messages {
		if msg.Timestamp.After(startTime) && msg.Timestamp.Before(endTime) {
			filtered = append(filtered, msg)
		}
	}

	// Sort by timestamp
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	return filtered, nil
}

// UpdateMessage updates a message
func (s *InMemoryStore) UpdateMessage(ctx context.Context, messageID string, updates map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, exists := s.messageIndex[messageID]
	if !exists {
		return fmt.Errorf("message %s not found", messageID)
	}

	// Update message fields
	if content, ok := updates["content"].(string); ok {
		msg.Content = content
	}
	
	if metadata, ok := updates["metadata"].(map[string]any); ok {
		if msg.Metadata == nil {
			msg.Metadata = make(map[string]any)
		}
		for k, v := range metadata {
			msg.Metadata[k] = v
		}
	}

	// Update in index and messages slice
	s.messageIndex[messageID] = msg
	
	// Find and update in conversation messages
	for convID, messages := range s.messages {
		for i, m := range messages {
			if m.ID == messageID {
				s.messages[convID][i] = msg
				break
			}
		}
	}

	return nil
}

// DeleteMessage removes a message
func (s *InMemoryStore) DeleteMessage(ctx context.Context, messageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.messageIndex[messageID]
	if !exists {
		return fmt.Errorf("message %s not found", messageID)
	}

	// Remove from index
	delete(s.messageIndex, messageID)

	// Remove from conversation messages
	for convID, messages := range s.messages {
		for i, m := range messages {
			if m.ID == messageID {
				s.messages[convID] = append(messages[:i], messages[i+1:]...)
				
				// Update conversation count
				if conv, exists := s.conversations[convID]; exists {
					conv.MessageCount--
					conv.UpdatedAt = time.Now()
				}
				break
			}
		}
	}

	return nil
}

// AddMessages adds multiple messages in batch
func (s *InMemoryStore) AddMessages(ctx context.Context, conversationID string, messages []Message) error {
	for _, msg := range messages {
		if err := s.AddMessage(ctx, conversationID, msg); err != nil {
			return err
		}
	}
	return nil
}

// ClearConversation removes all messages from a conversation
func (s *InMemoryStore) ClearConversation(ctx context.Context, conversationID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, exists := s.conversations[conversationID]
	if !exists {
		return fmt.Errorf("conversation %s not found", conversationID)
	}

	// Remove messages from index
	if messages, exists := s.messages[conversationID]; exists {
		for _, msg := range messages {
			delete(s.messageIndex, msg.ID)
		}
	}

	s.messages[conversationID] = make([]Message, 0)
	conv.MessageCount = 0
	conv.UpdatedAt = time.Now()

	return nil
}

// SearchMessages searches for messages based on query
func (s *InMemoryStore) SearchMessages(ctx context.Context, query SearchQuery) (*SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []Message
	
	for _, messages := range s.messages {
		for _, msg := range messages {
			if s.matchesMessageQuery(msg, query) {
				results = append(results, msg)
			}
		}
	}

	// Sort by timestamp (most recent first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.After(results[j].Timestamp)
	})

	// Apply pagination
	total := len(results)
	offset := query.Offset
	limit := query.Limit
	
	if limit == 0 {
		limit = 50 // default limit
	}

	if offset >= len(results) {
		results = []Message{}
	} else {
		end := offset + limit
		if end > len(results) {
			end = len(results)
		}
		results = results[offset:end]
	}

	return &SearchResult{
		Messages: results,
		Total:    total,
		HasMore:  offset+len(results) < total,
	}, nil
}

// SearchConversations searches for conversations based on query
func (s *InMemoryStore) SearchConversations(ctx context.Context, query SearchQuery) (*SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []Conversation
	
	for _, conv := range s.conversations {
		if s.matchesConversationQuery(*conv, query) {
			results = append(results, *conv)
		}
	}

	// Sort by updated time (most recent first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdatedAt.After(results[j].UpdatedAt)
	})

	// Apply pagination
	total := len(results)
	offset := query.Offset
	limit := query.Limit
	
	if limit == 0 {
		limit = 20 // default limit
	}

	if offset >= len(results) {
		results = []Conversation{}
	} else {
		end := offset + limit
		if end > len(results) {
			end = len(results)
		}
		results = results[offset:end]
	}

	return &SearchResult{
		Conversations: results,
		Total:         total,
		HasMore:       offset+len(results) < total,
	}, nil
}

// matchesMessageQuery checks if a message matches the search query
func (s *InMemoryStore) matchesMessageQuery(msg Message, query SearchQuery) bool {
	// Text search
	if query.Query != "" {
		queryLower := strings.ToLower(query.Query)
		if !strings.Contains(strings.ToLower(msg.Content), queryLower) {
			return false
		}
	}

	// Role filter
	if len(query.Roles) > 0 {
		found := false
		for _, role := range query.Roles {
			if msg.Role == role {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Time range filter
	if query.StartTime != nil && msg.Timestamp.Before(*query.StartTime) {
		return false
	}
	if query.EndTime != nil && msg.Timestamp.After(*query.EndTime) {
		return false
	}

	return true
}

// matchesConversationQuery checks if a conversation matches the search query
func (s *InMemoryStore) matchesConversationQuery(conv Conversation, query SearchQuery) bool {
	// Text search in title
	if query.Query != "" {
		queryLower := strings.ToLower(query.Query)
		if !strings.Contains(strings.ToLower(conv.Title), queryLower) {
			return false
		}
	}

	// Tags filter
	if len(query.Tags) > 0 {
		for _, queryTag := range query.Tags {
			found := false
			for _, convTag := range conv.Tags {
				if convTag == queryTag {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	// Time range filter
	if query.StartTime != nil && conv.CreatedAt.Before(*query.StartTime) {
		return false
	}
	if query.EndTime != nil && conv.CreatedAt.After(*query.EndTime) {
		return false
	}

	return true
}