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

package agents

import (
	"context"
	"testing"

	"github.com/adrienveepee/adk-go/internal/models"
)

func TestCreateAgent(t *testing.T) {
	service := NewService()
	config := models.AgentConfig{
		Name:        "Test Agent",
		Model:       "gemini-2.0-flash",
		Instruction: "You are a helpful assistant",
		Description: "A test agent",
	}

	agent, err := service.CreateAgent(context.Background(), config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if agent.Name != config.Name {
		t.Errorf("Expected agent name %s, got %s", config.Name, agent.Name)
	}

	if agent.Model != config.Model {
		t.Errorf("Expected agent model %s, got %s", config.Model, agent.Model)
	}

	if agent.ID == "" {
		t.Error("Expected agent ID to be set")
	}
}

func TestGetAgent(t *testing.T) {
	service := NewService()
	config := models.AgentConfig{
		Name:  "Test Agent",
		Model: "gemini-2.0-flash",
	}

	// Create agent first
	createdAgent, err := service.CreateAgent(context.Background(), config)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	// Retrieve agent
	retrievedAgent, err := service.GetAgent(context.Background(), createdAgent.ID)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if retrievedAgent.ID != createdAgent.ID {
		t.Errorf("Expected agent ID %s, got %s", createdAgent.ID, retrievedAgent.ID)
	}
}

func TestGetAgentNotFound(t *testing.T) {
	service := NewService()

	_, err := service.GetAgent(context.Background(), "non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent agent, got nil")
	}
}

func TestRunAgent(t *testing.T) {
	service := NewService()
	config := models.AgentConfig{
		Name:  "Test Agent",
		Model: "gemini-2.0-flash",
	}

	// Create agent first
	agent, err := service.CreateAgent(context.Background(), config)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	// Run agent
	input := "Hello, world!"
	result, err := service.RunAgent(context.Background(), agent.ID, input)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}