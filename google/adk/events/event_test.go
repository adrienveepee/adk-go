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

package events

import (
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	event := NewEvent()
	
	if event.ID == "" {
		t.Error("Event ID should not be empty")
	}
	
	if event.Timestamp.IsZero() {
		t.Error("Event timestamp should not be zero")
	}
	
	if event.Actions.StateDelta == nil {
		event.Actions.StateDelta = make(map[string]interface{})
	}
}

func TestNewEventWithID(t *testing.T) {
	testID := "test-event-id"
	event := NewEventWithID(testID)
	
	if event.ID != testID {
		t.Errorf("Expected event ID to be %s, got %s", testID, event.ID)
	}
	
	if event.Timestamp.IsZero() {
		t.Error("Event timestamp should not be zero")
	}
}

func TestEventIsFinalized(t *testing.T) {
	event := NewEvent()
	
	// Test default state
	if event.IsFinalized() {
		t.Error("Event should not be finalized by default")
	}
	
	// Test finalized state
	event.IsFinalResponse = true
	if !event.IsFinalized() {
		t.Error("Event should be finalized when IsFinalResponse is true")
	}
}

func TestEventContent(t *testing.T) {
	event := NewEvent()
	
	content := &Content{
		Role: "user",
		Parts: []Part{
			{Text: "Hello, world!"},
		},
	}
	
	event.Content = content
	
	if event.Content.Role != "user" {
		t.Errorf("Expected content role to be 'user', got %s", event.Content.Role)
	}
	
	if len(event.Content.Parts) != 1 {
		t.Errorf("Expected 1 content part, got %d", len(event.Content.Parts))
	}
	
	if event.Content.Parts[0].Text != "Hello, world!" {
		t.Errorf("Expected content text to be 'Hello, world!', got %s", event.Content.Parts[0].Text)
	}
}

func TestEventActions(t *testing.T) {
	event := NewEvent()
	
	// Test transfer to agent
	event.Actions.TransferToAgent = "test_agent"
	if event.Actions.TransferToAgent != "test_agent" {
		t.Errorf("Expected transfer to agent to be 'test_agent', got %s", event.Actions.TransferToAgent)
	}
	
	// Test escalate
	event.Actions.Escalate = true
	if !event.Actions.Escalate {
		t.Error("Expected escalate to be true")
	}
	
	// Test state delta
	event.Actions.StateDelta = map[string]interface{}{
		"key": "value",
	}
	if event.Actions.StateDelta["key"] != "value" {
		t.Errorf("Expected state delta key to be 'value', got %v", event.Actions.StateDelta["key"])
	}
}

func TestEventTimestamp(t *testing.T) {
	before := time.Now()
	event := NewEvent()
	after := time.Now()
	
	if event.Timestamp.Before(before) || event.Timestamp.After(after) {
		t.Error("Event timestamp should be between before and after times")
	}
}

func TestPartText(t *testing.T) {
	part := Part{Text: "Test text"}
	
	if part.Text != "Test text" {
		t.Errorf("Expected part text to be 'Test text', got %s", part.Text)
	}
}

func TestContentMultipleParts(t *testing.T) {
	content := &Content{
		Role: "model",
		Parts: []Part{
			{Text: "Part 1"},
			{Text: "Part 2"},
		},
	}
	
	if content.Role != "model" {
		t.Errorf("Expected content role to be 'model', got %s", content.Role)
	}
	
	if len(content.Parts) != 2 {
		t.Errorf("Expected 2 content parts, got %d", len(content.Parts))
	}
	
	if content.Parts[0].Text != "Part 1" {
		t.Errorf("Expected first part text to be 'Part 1', got %s", content.Parts[0].Text)
	}
	
	if content.Parts[1].Text != "Part 2" {
		t.Errorf("Expected second part text to be 'Part 2', got %s", content.Parts[1].Text)
	}
}