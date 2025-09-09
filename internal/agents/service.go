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
	"fmt"
	"time"

	"github.com/adrienveepee/adk-go/internal/models"
	"github.com/google/uuid"
)

// Service represents the agent service
type Service struct {
	agents map[string]*models.Agent
}

// NewService creates a new agent service
func NewService() *Service {
	return &Service{
		agents: make(map[string]*models.Agent),
	}
}

// CreateAgent creates a new agent
func (s *Service) CreateAgent(ctx context.Context, config models.AgentConfig) (*models.Agent, error) {
	agent := &models.Agent{
		ID:          uuid.New().String(),
		Name:        config.Name,
		Model:       config.Model,
		Instruction: config.Instruction,
		Description: config.Description,
		Tools:       []models.Tool{}, // TODO: Convert tool names to tool objects
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.agents[agent.ID] = agent
	return agent, nil
}

// GetAgent retrieves an agent by ID
func (s *Service) GetAgent(ctx context.Context, id string) (*models.Agent, error) {
	agent, exists := s.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent with ID %s not found", id)
	}
	return agent, nil
}

// RunAgent executes an agent with the given input
func (s *Service) RunAgent(ctx context.Context, id string, input string) (string, error) {
	agent, err := s.GetAgent(ctx, id)
	if err != nil {
		return "", err
	}

	// TODO: Implement actual agent execution logic
	// This is a placeholder implementation
	response := fmt.Sprintf("Agent '%s' with model '%s' processed input: %s", agent.Name, agent.Model, input)
	return response, nil
}

// ListAgents returns all agents
func (s *Service) ListAgents(ctx context.Context) ([]*models.Agent, error) {
	agents := make([]*models.Agent, 0, len(s.agents))
	for _, agent := range s.agents {
		agents = append(agents, agent)
	}
	return agents, nil
}