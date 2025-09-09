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

package models

import (
	"time"
)

// Agent represents an AI agent in the ADK system
type Agent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Model       string    `json:"model"`
	Instruction string    `json:"instruction"`
	Description string    `json:"description"`
	Tools       []Tool    `json:"tools"`
	SubAgents   []Agent   `json:"sub_agents,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tool represents a tool that agents can use
type Tool struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Session represents a conversation session
type Session struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Message represents a message in a session
type Message struct {
	ID        string                 `json:"id"`
	Role      string                 `json:"role"` // "user", "assistant", "system"
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// AgentConfig represents configuration for creating an agent
type AgentConfig struct {
	Name        string   `json:"name" binding:"required"`
	Model       string   `json:"model" binding:"required"`
	Instruction string   `json:"instruction"`
	Description string   `json:"description"`
	Tools       []string `json:"tools,omitempty"`
}

// SessionConfig represents configuration for creating a session
type SessionConfig struct {
	AgentID string `json:"agent_id" binding:"required"`
}

// MessageRequest represents a request to send a message
type MessageRequest struct {
	Content  string                 `json:"content" binding:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}