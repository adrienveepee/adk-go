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
	"sync"

	"github.com/adrienveepee/adk-go/google/adk/events"
)

// SequentialAgent executes sub-agents sequentially
type SequentialAgent struct {
	*BaseAgent
}

// NewSequentialAgent creates a new sequential agent
func NewSequentialAgent(name string, subAgents []Agent) *SequentialAgent {
	agent := &SequentialAgent{
		BaseAgent: NewBaseAgent(name, "Sequential execution agent"),
	}
	
	for _, subAgent := range subAgents {
		agent.AddSubAgent(subAgent)
	}
	
	return agent
}

// RunAsync executes sub-agents sequentially
func (a *SequentialAgent) RunAsync(ctx context.Context, invocationCtx *InvocationContext) (<-chan *events.Event, error) {
	eventChan := make(chan *events.Event)
	
	go func() {
		defer close(eventChan)
		
		// Execute before agent callback
		if a.BeforeAgentCallback != nil {
			if err := a.BeforeAgentCallback(invocationCtx); err != nil {
				return
			}
		}
		
		// Execute each sub-agent sequentially
		for _, subAgent := range a.SubAgents {
			subEventChan, err := subAgent.RunAsync(ctx, invocationCtx)
			if err != nil {
				// TODO: Better error handling
				continue
			}
			
			// Forward all events from the sub-agent
			for event := range subEventChan {
				eventChan <- event
			}
		}
		
		// Execute after agent callback
		if a.AfterAgentCallback != nil {
			if err := a.AfterAgentCallback(invocationCtx); err != nil {
				return
			}
		}
	}()
	
	return eventChan, nil
}

// ParallelAgent executes sub-agents in parallel
type ParallelAgent struct {
	*BaseAgent
}

// NewParallelAgent creates a new parallel agent
func NewParallelAgent(name string, subAgents []Agent) *ParallelAgent {
	agent := &ParallelAgent{
		BaseAgent: NewBaseAgent(name, "Parallel execution agent"),
	}
	
	for _, subAgent := range subAgents {
		agent.AddSubAgent(subAgent)
	}
	
	return agent
}

// RunAsync executes sub-agents in parallel
func (a *ParallelAgent) RunAsync(ctx context.Context, invocationCtx *InvocationContext) (<-chan *events.Event, error) {
	eventChan := make(chan *events.Event)
	
	go func() {
		defer close(eventChan)
		
		// Execute before agent callback
		if a.BeforeAgentCallback != nil {
			if err := a.BeforeAgentCallback(invocationCtx); err != nil {
				return
			}
		}
		
		// Create a wait group to track sub-agent completion
		var wg sync.WaitGroup
		
		// Execute each sub-agent in parallel
		for _, subAgent := range a.SubAgents {
			wg.Add(1)
			go func(agent Agent) {
				defer wg.Done()
				
				subEventChan, err := agent.RunAsync(ctx, invocationCtx)
				if err != nil {
					return
				}
				
				// Forward all events from the sub-agent
				for event := range subEventChan {
					eventChan <- event
				}
			}(subAgent)
		}
		
		// Wait for all sub-agents to complete
		wg.Wait()
		
		// Execute after agent callback
		if a.AfterAgentCallback != nil {
			if err := a.AfterAgentCallback(invocationCtx); err != nil {
				return
			}
		}
	}()
	
	return eventChan, nil
}

// LoopAgent executes sub-agents in a loop with configurable iterations
type LoopAgent struct {
	*BaseAgent
	MaxIterations int `json:"max_iterations"`
}

// NewLoopAgent creates a new loop agent
func NewLoopAgent(name string, subAgents []Agent, maxIterations int) *LoopAgent {
	agent := &LoopAgent{
		BaseAgent:     NewBaseAgent(name, "Loop execution agent"),
		MaxIterations: maxIterations,
	}
	
	for _, subAgent := range subAgents {
		agent.AddSubAgent(subAgent)
	}
	
	return agent
}

// RunAsync executes sub-agents in a loop
func (a *LoopAgent) RunAsync(ctx context.Context, invocationCtx *InvocationContext) (<-chan *events.Event, error) {
	eventChan := make(chan *events.Event)
	
	go func() {
		defer close(eventChan)
		
		// Execute before agent callback
		if a.BeforeAgentCallback != nil {
			if err := a.BeforeAgentCallback(invocationCtx); err != nil {
				return
			}
		}
		
		// Execute loop iterations
		for iteration := 0; iteration < a.MaxIterations; iteration++ {
			// Check if we should exit the loop early
			if a.shouldExitLoop(invocationCtx) {
				break
			}
			
			// Execute each sub-agent sequentially in this iteration
			for _, subAgent := range a.SubAgents {
				subEventChan, err := subAgent.RunAsync(ctx, invocationCtx)
				if err != nil {
					// TODO: Better error handling
					continue
				}
				
				// Forward all events from the sub-agent
				for event := range subEventChan {
					eventChan <- event
					
					// Check if the event indicates we should exit the loop
					if a.shouldExitLoopFromEvent(event) {
						goto exitLoop
					}
				}
			}
		}
		
	exitLoop:
		// Execute after agent callback
		if a.AfterAgentCallback != nil {
			if err := a.AfterAgentCallback(invocationCtx); err != nil {
				return
			}
		}
	}()
	
	return eventChan, nil
}

// shouldExitLoop checks if the loop should be exited based on session state
func (a *LoopAgent) shouldExitLoop(invocationCtx *InvocationContext) bool {
	// Check for exit loop flag in session state
	if exitFlag, exists := invocationCtx.Session.State.Get("exit_loop"); exists {
		if exit, ok := exitFlag.(bool); ok && exit {
			return true
		}
	}
	return false
}

// shouldExitLoopFromEvent checks if an event indicates the loop should be exited
func (a *LoopAgent) shouldExitLoopFromEvent(event *events.Event) bool {
	// Check if the event has an exit loop action
	return event.Actions.SkipSummarization // Using this as a proxy for exit loop
}