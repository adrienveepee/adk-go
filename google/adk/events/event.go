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
	"time"

	"github.com/google/uuid"
)

// Content represents message content with role and parts
type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

// Part represents a part of content (text, image, etc.)
type Part struct {
	Text string `json:"text,omitempty"`
	// Could be extended for images, files, etc.
}

// EventActions represents actions that can be taken with an event
type EventActions struct {
	TransferToAgent         string                 `json:"transfer_to_agent,omitempty"`
	Escalate                bool                   `json:"escalate,omitempty"`
	SkipSummarization       bool                   `json:"skip_summarization,omitempty"`
	StateDelta              map[string]interface{} `json:"state_delta,omitempty"`
	ArtifactDelta           map[string]interface{} `json:"artifact_delta,omitempty"`
	RequestedAuthConfigs    []interface{}          `json:"requested_auth_configs,omitempty"`
}

// Event represents a single event in the ADK system
type Event struct {
	ID                    string       `json:"id"`
	InvocationID          string       `json:"invocation_id"`
	Timestamp             time.Time    `json:"timestamp"`
	Author                string       `json:"author"`
	Content               *Content     `json:"content,omitempty"`
	Branch                string       `json:"branch,omitempty"`
	IsFinalResponse       bool         `json:"is_final_response"`
	Actions               EventActions `json:"actions,omitempty"`
	LongRunningToolIDs    []string     `json:"long_running_tool_ids,omitempty"`
}

// NewEvent creates a new event with a unique ID and current timestamp
func NewEvent() *Event {
	return &Event{
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Actions:   EventActions{},
	}
}

// NewEventWithID creates a new event with the given ID
func NewEventWithID(id string) *Event {
	return &Event{
		ID:        id,
		Timestamp: time.Now(),
		Actions:   EventActions{},
	}
}

// IsFinalized returns whether this event represents a final response
func (e *Event) IsFinalized() bool {
	return e.IsFinalResponse
}

// GetFunctionCalls extracts function calls from the event content
func (e *Event) GetFunctionCalls() []interface{} {
	// TODO: Implement function call extraction
	return nil
}

// GetFunctionResponses extracts function responses from the event content
func (e *Event) GetFunctionResponses() []interface{} {
	// TODO: Implement function response extraction
	return nil
}

// HasTrailingCodeExecutionResult checks if the event has code execution results
func (e *Event) HasTrailingCodeExecutionResult() bool {
	// TODO: Implement code execution result detection
	return false
}