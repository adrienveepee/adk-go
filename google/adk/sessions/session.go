// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sessions

import (
	"sync"
	"time"

	"github.com/adrienveepee/adk-go/google/adk/events"
	"github.com/google/uuid"
)

// State represents session state with key-value storage
type State struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// NewState creates a new state instance
func NewState() *State {
	return &State{
		data: make(map[string]interface{}),
	}
}

// NewStateWithData creates a new state instance with initial data
func NewStateWithData(data map[string]interface{}) *State {
	stateCopy := make(map[string]interface{})
	for k, v := range data {
		stateCopy[k] = v
	}
	return &State{
		data: stateCopy,
	}
}

// Get retrieves a value from the state
func (s *State) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.data[key]
	return value, exists
}

// Set stores a value in the state
func (s *State) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Update updates the state with a map of key-value pairs
func (s *State) Update(updates map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, value := range updates {
		s.data[key] = value
	}
}

// ToDict returns a copy of the state as a map
func (s *State) ToDict() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]interface{})
	for k, v := range s.data {
		result[k] = v
	}
	return result
}

// HasDelta checks if the state has changes
func (s *State) HasDelta() bool {
	// TODO: Implement delta tracking
	return false
}

// Session represents a user session
type Session struct {
	ID             string         `json:"id"`
	AppName        string         `json:"app_name"`
	UserID         string         `json:"user_id"`
	State          *State         `json:"state"`
	Events         []*events.Event `json:"events"`
	LastUpdateTime time.Time      `json:"last_update_time"`
}

// NewSession creates a new session
func NewSession(appName, userID, sessionID string, initialState map[string]interface{}) *Session {
	if sessionID == "" {
		sessionID = uuid.New().String()
	}
	
	var state *State
	if initialState != nil {
		state = NewStateWithData(initialState)
	} else {
		state = NewState()
	}
	
	return &Session{
		ID:             sessionID,
		AppName:        appName,
		UserID:         userID,
		State:          state,
		Events:         make([]*events.Event, 0),
		LastUpdateTime: time.Now(),
	}
}

// AddEvent adds an event to the session
func (s *Session) AddEvent(event *events.Event) {
	s.Events = append(s.Events, event)
	s.LastUpdateTime = time.Now()
}

// GetEvents returns all events in the session
func (s *Session) GetEvents() []*events.Event {
	return s.Events
}

// SessionService interface for managing sessions
type SessionService interface {
	// CreateSession creates a new session
	CreateSession(appName, userID, sessionID string, initialState map[string]interface{}) (*Session, error)
	
	// GetSession retrieves a session by ID
	GetSession(appName, userID, sessionID string) (*Session, error)
	
	// DeleteSession deletes a session
	DeleteSession(appName, userID, sessionID string) error
	
	// ListSessions lists all sessions for a user
	ListSessions(appName, userID string) ([]*Session, error)
	
	// AppendEvent adds an event to a session
	AppendEvent(appName, userID, sessionID string, event *events.Event) error
	
	// ListEvents lists events for a session
	ListEvents(appName, userID, sessionID string) ([]*events.Event, error)
	
	// CloseSession closes a session
	CloseSession(appName, userID, sessionID string) error
}

// InMemorySessionService provides an in-memory implementation of SessionService
type InMemorySessionService struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// NewInMemorySessionService creates a new in-memory session service
func NewInMemorySessionService() *InMemorySessionService {
	return &InMemorySessionService{
		sessions: make(map[string]*Session),
	}
}

// sessionKey creates a unique key for a session
func (s *InMemorySessionService) sessionKey(appName, userID, sessionID string) string {
	return appName + ":" + userID + ":" + sessionID
}

// CreateSession creates a new session
func (s *InMemorySessionService) CreateSession(appName, userID, sessionID string, initialState map[string]interface{}) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	session := NewSession(appName, userID, sessionID, initialState)
	key := s.sessionKey(appName, userID, sessionID)
	s.sessions[key] = session
	
	return session, nil
}

// GetSession retrieves a session by ID
func (s *InMemorySessionService) GetSession(appName, userID, sessionID string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	key := s.sessionKey(appName, userID, sessionID)
	session, exists := s.sessions[key]
	if !exists {
		return nil, nil // Session not found
	}
	
	return session, nil
}

// DeleteSession deletes a session
func (s *InMemorySessionService) DeleteSession(appName, userID, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	key := s.sessionKey(appName, userID, sessionID)
	delete(s.sessions, key)
	
	return nil
}

// ListSessions lists all sessions for a user
func (s *InMemorySessionService) ListSessions(appName, userID string) ([]*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	prefix := appName + ":" + userID + ":"
	var sessions []*Session
	
	for key, session := range s.sessions {
		if len(key) > len(prefix) && key[:len(prefix)] == prefix {
			sessions = append(sessions, session)
		}
	}
	
	return sessions, nil
}

// AppendEvent adds an event to a session
func (s *InMemorySessionService) AppendEvent(appName, userID, sessionID string, event *events.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	key := s.sessionKey(appName, userID, sessionID)
	session, exists := s.sessions[key]
	if !exists {
		return nil // Session not found
	}
	
	session.AddEvent(event)
	return nil
}

// ListEvents lists events for a session
func (s *InMemorySessionService) ListEvents(appName, userID, sessionID string) ([]*events.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	key := s.sessionKey(appName, userID, sessionID)
	session, exists := s.sessions[key]
	if !exists {
		return nil, nil // Session not found
	}
	
	return session.GetEvents(), nil
}

// CloseSession closes a session
func (s *InMemorySessionService) CloseSession(appName, userID, sessionID string) error {
	// For in-memory implementation, we don't need to do anything special
	return nil
}