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
	"testing"

	"github.com/adrienveepee/adk-go/google/adk/events"
)

func TestNewState(t *testing.T) {
	state := NewState()
	
	if state == nil {
		t.Error("NewState should not return nil")
	}
	
	if state.data == nil {
		t.Error("State data should be initialized")
	}
}

func TestNewStateWithData(t *testing.T) {
	initialData := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	
	state := NewStateWithData(initialData)
	
	if state == nil {
		t.Error("NewStateWithData should not return nil")
	}
	
	value1, exists1 := state.Get("key1")
	if !exists1 || value1 != "value1" {
		t.Errorf("Expected key1 to be 'value1', got %v", value1)
	}
	
	value2, exists2 := state.Get("key2")
	if !exists2 || value2 != 42 {
		t.Errorf("Expected key2 to be 42, got %v", value2)
	}
}

func TestStateGetSet(t *testing.T) {
	state := NewState()
	
	// Test setting and getting a value
	state.Set("test_key", "test_value")
	
	value, exists := state.Get("test_key")
	if !exists {
		t.Error("Key should exist after setting")
	}
	
	if value != "test_value" {
		t.Errorf("Expected value to be 'test_value', got %v", value)
	}
	
	// Test getting a non-existent key
	_, exists = state.Get("non_existent_key")
	if exists {
		t.Error("Non-existent key should not exist")
	}
}

func TestStateUpdate(t *testing.T) {
	state := NewState()
	
	// Set initial value
	state.Set("key1", "initial")
	
	// Update with multiple values
	updates := map[string]interface{}{
		"key1": "updated",
		"key2": 123,
		"key3": true,
	}
	
	state.Update(updates)
	
	// Verify updates
	value1, _ := state.Get("key1")
	if value1 != "updated" {
		t.Errorf("Expected key1 to be 'updated', got %v", value1)
	}
	
	value2, _ := state.Get("key2")
	if value2 != 123 {
		t.Errorf("Expected key2 to be 123, got %v", value2)
	}
	
	value3, _ := state.Get("key3")
	if value3 != true {
		t.Errorf("Expected key3 to be true, got %v", value3)
	}
}

func TestStateToDict(t *testing.T) {
	state := NewState()
	
	state.Set("key1", "value1")
	state.Set("key2", 42)
	
	dict := state.ToDict()
	
	if len(dict) != 2 {
		t.Errorf("Expected dict to have 2 entries, got %d", len(dict))
	}
	
	if dict["key1"] != "value1" {
		t.Errorf("Expected dict[key1] to be 'value1', got %v", dict["key1"])
	}
	
	if dict["key2"] != 42 {
		t.Errorf("Expected dict[key2] to be 42, got %v", dict["key2"])
	}
}

func TestNewSession(t *testing.T) {
	appName := "test_app"
	userID := "test_user"
	sessionID := "test_session"
	initialState := map[string]interface{}{
		"initial_key": "initial_value",
	}
	
	session := NewSession(appName, userID, sessionID, initialState)
	
	if session.ID != sessionID {
		t.Errorf("Expected session ID to be %s, got %s", sessionID, session.ID)
	}
	
	if session.AppName != appName {
		t.Errorf("Expected app name to be %s, got %s", appName, session.AppName)
	}
	
	if session.UserID != userID {
		t.Errorf("Expected user ID to be %s, got %s", userID, session.UserID)
	}
	
	value, exists := session.State.Get("initial_key")
	if !exists || value != "initial_value" {
		t.Errorf("Expected initial state to be preserved")
	}
	
	if len(session.Events) != 0 {
		t.Errorf("Expected empty events list, got %d events", len(session.Events))
	}
}

func TestSessionAddEvent(t *testing.T) {
	session := NewSession("app", "user", "session", nil)
	
	event := events.NewEvent()
	event.Author = "test_author"
	
	session.AddEvent(event)
	
	if len(session.Events) != 1 {
		t.Errorf("Expected 1 event, got %d events", len(session.Events))
	}
	
	if session.Events[0].Author != "test_author" {
		t.Errorf("Expected event author to be 'test_author', got %s", session.Events[0].Author)
	}
}

func TestInMemorySessionService(t *testing.T) {
	service := NewInMemorySessionService()
	
	if service == nil {
		t.Error("NewInMemorySessionService should not return nil")
	}
	
	if service.sessions == nil {
		t.Error("Session service sessions map should be initialized")
	}
}

func TestSessionServiceCreateSession(t *testing.T) {
	service := NewInMemorySessionService()
	
	appName := "test_app"
	userID := "test_user"
	sessionID := "test_session"
	initialState := map[string]interface{}{
		"key": "value",
	}
	
	session, err := service.CreateSession(appName, userID, sessionID, initialState)
	if err != nil {
		t.Errorf("CreateSession should not return error: %v", err)
	}
	
	if session == nil {
		t.Error("CreateSession should return a session")
	}
	
	if session.ID != sessionID {
		t.Errorf("Expected session ID to be %s, got %s", sessionID, session.ID)
	}
}

func TestSessionServiceGetSession(t *testing.T) {
	service := NewInMemorySessionService()
	
	appName := "test_app"
	userID := "test_user"
	sessionID := "test_session"
	
	// Create a session
	createdSession, _ := service.CreateSession(appName, userID, sessionID, nil)
	
	// Get the session
	retrievedSession, err := service.GetSession(appName, userID, sessionID)
	if err != nil {
		t.Errorf("GetSession should not return error: %v", err)
	}
	
	if retrievedSession == nil {
		t.Error("GetSession should return the session")
	}
	
	if retrievedSession.ID != createdSession.ID {
		t.Errorf("Retrieved session ID should match created session ID")
	}
	
	// Test getting non-existent session
	nonExistentSession, err := service.GetSession(appName, userID, "non_existent")
	if err != nil {
		t.Errorf("GetSession should not return error for non-existent session: %v", err)
	}
	
	if nonExistentSession != nil {
		t.Error("GetSession should return nil for non-existent session")
	}
}

func TestSessionServiceDeleteSession(t *testing.T) {
	service := NewInMemorySessionService()
	
	appName := "test_app"
	userID := "test_user"
	sessionID := "test_session"
	
	// Create a session
	service.CreateSession(appName, userID, sessionID, nil)
	
	// Delete the session
	err := service.DeleteSession(appName, userID, sessionID)
	if err != nil {
		t.Errorf("DeleteSession should not return error: %v", err)
	}
	
	// Verify session is deleted
	session, _ := service.GetSession(appName, userID, sessionID)
	if session != nil {
		t.Error("Session should be deleted")
	}
}

func TestSessionServiceAppendEvent(t *testing.T) {
	service := NewInMemorySessionService()
	
	appName := "test_app"
	userID := "test_user"
	sessionID := "test_session"
	
	// Create a session
	service.CreateSession(appName, userID, sessionID, nil)
	
	// Create an event
	event := events.NewEvent()
	event.Author = "test_author"
	
	// Append the event
	err := service.AppendEvent(appName, userID, sessionID, event)
	if err != nil {
		t.Errorf("AppendEvent should not return error: %v", err)
	}
	
	// Verify event was appended
	session, _ := service.GetSession(appName, userID, sessionID)
	if len(session.Events) != 1 {
		t.Errorf("Expected 1 event, got %d events", len(session.Events))
	}
	
	if session.Events[0].Author != "test_author" {
		t.Errorf("Expected event author to be 'test_author', got %s", session.Events[0].Author)
	}
}

func TestSessionServiceListSessions(t *testing.T) {
	service := NewInMemorySessionService()
	
	appName := "test_app"
	userID := "test_user"
	
	// Create multiple sessions
	service.CreateSession(appName, userID, "session1", nil)
	service.CreateSession(appName, userID, "session2", nil)
	service.CreateSession(appName, "other_user", "session3", nil) // Different user
	
	// List sessions for the user
	sessions, err := service.ListSessions(appName, userID)
	if err != nil {
		t.Errorf("ListSessions should not return error: %v", err)
	}
	
	if len(sessions) != 2 {
		t.Errorf("Expected 2 sessions for user, got %d sessions", len(sessions))
	}
}