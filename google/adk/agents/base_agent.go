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

	"github.com/adrienveepee/adk-go/google/adk/events"
	"github.com/adrienveepee/adk-go/google/adk/sessions"
)

// InvocationContext provides the context for agent invocation
type InvocationContext struct {
	Session sessions.Session
	// Add other context fields as needed
}

// Agent is the interface that all agents must implement
type Agent interface {
	// GetName returns the agent's name
	GetName() string
	
	// GetDescription returns the agent's description
	GetDescription() string
	
	// RunAsync executes the agent asynchronously and returns events
	RunAsync(ctx context.Context, invocationCtx *InvocationContext) (<-chan *events.Event, error)
	
	// RunLive executes the agent in live mode (bidi-streaming)
	RunLive(ctx context.Context, invocationCtx *InvocationContext) (<-chan *events.Event, error)
	
	// GetSubAgents returns the list of sub-agents
	GetSubAgents() []Agent
	
	// FindAgent finds an agent by name in the hierarchy
	FindAgent(name string) Agent
	
	// FindSubAgent finds a direct sub-agent by name
	FindSubAgent(name string) Agent
	
	// GetParentAgent returns the parent agent
	GetParentAgent() Agent
	
	// SetParentAgent sets the parent agent
	SetParentAgent(parent Agent)
	
	// GetRootAgent returns the root agent in the hierarchy
	GetRootAgent() Agent
}

// BaseAgent provides the base implementation for all agents
type BaseAgent struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	SubAgents    []Agent `json:"sub_agents,omitempty"`
	ParentAgent  Agent   `json:"-"`
	
	// Callbacks
	BeforeAgentCallback func(ctx *InvocationContext) error `json:"-"`
	AfterAgentCallback  func(ctx *InvocationContext) error `json:"-"`
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(name, description string) *BaseAgent {
	return &BaseAgent{
		Name:        name,
		Description: description,
		SubAgents:   make([]Agent, 0),
	}
}

// GetName returns the agent's name
func (a *BaseAgent) GetName() string {
	return a.Name
}

// GetDescription returns the agent's description
func (a *BaseAgent) GetDescription() string {
	return a.Description
}

// GetSubAgents returns the list of sub-agents
func (a *BaseAgent) GetSubAgents() []Agent {
	return a.SubAgents
}

// AddSubAgent adds a sub-agent
func (a *BaseAgent) AddSubAgent(subAgent Agent) {
	a.SubAgents = append(a.SubAgents, subAgent)
	subAgent.SetParentAgent(a)
}

// FindAgent finds an agent by name in the hierarchy
func (a *BaseAgent) FindAgent(name string) Agent {
	if a.Name == name {
		return a
	}
	
	for _, subAgent := range a.SubAgents {
		if found := subAgent.FindAgent(name); found != nil {
			return found
		}
	}
	
	return nil
}

// FindSubAgent finds a direct sub-agent by name
func (a *BaseAgent) FindSubAgent(name string) Agent {
	for _, subAgent := range a.SubAgents {
		if subAgent.GetName() == name {
			return subAgent
		}
	}
	return nil
}

// GetParentAgent returns the parent agent
func (a *BaseAgent) GetParentAgent() Agent {
	return a.ParentAgent
}

// SetParentAgent sets the parent agent
func (a *BaseAgent) SetParentAgent(parent Agent) {
	a.ParentAgent = parent
}

// GetRootAgent returns the root agent in the hierarchy
func (a *BaseAgent) GetRootAgent() Agent {
	if a.ParentAgent == nil {
		return a
	}
	return a.ParentAgent.GetRootAgent()
}

// RunAsync is the base implementation - to be overridden by concrete agents
func (a *BaseAgent) RunAsync(ctx context.Context, invocationCtx *InvocationContext) (<-chan *events.Event, error) {
	// Base implementation returns empty channel
	eventChan := make(chan *events.Event)
	close(eventChan)
	return eventChan, nil
}

// RunLive is the base implementation for live mode
func (a *BaseAgent) RunLive(ctx context.Context, invocationCtx *InvocationContext) (<-chan *events.Event, error) {
	// Default implementation delegates to RunAsync
	return a.RunAsync(ctx, invocationCtx)
}