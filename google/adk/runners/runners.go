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

package runners

import (
	"context"
	"fmt"

	"github.com/adrienveepee/adk-go/google/adk/agents"
	"github.com/adrienveepee/adk-go/google/adk/artifacts"
	"github.com/adrienveepee/adk-go/google/adk/events"
	"github.com/adrienveepee/adk-go/google/adk/memory"
	"github.com/adrienveepee/adk-go/google/adk/sessions"
)

// Runner orchestrates agent execution with sessions and services
type Runner struct {
	Agent           agents.Agent
	AppName         string
	SessionService  sessions.SessionService
	MemoryService   memory.MemoryService
	ArtifactService artifacts.ArtifactService
}

// NewRunner creates a new runner instance
func NewRunner(agent agents.Agent, appName string, sessionService sessions.SessionService) *Runner {
	return &Runner{
		Agent:          agent,
		AppName:        appName,
		SessionService: sessionService,
		// Use default in-memory services if not provided
		MemoryService:   memory.NewInMemoryMemoryService(),
		ArtifactService: artifacts.NewInMemoryArtifactService(),
	}
}

// Run executes an agent synchronously and returns the final response
func (r *Runner) Run(ctx context.Context, userID, sessionID string, newMessage *events.Content) (*events.Event, error) {
	eventChan, err := r.RunAsync(ctx, userID, sessionID, newMessage)
	if err != nil {
		return nil, err
	}
	
	var finalEvent *events.Event
	for event := range eventChan {
		finalEvent = event
	}
	
	return finalEvent, nil
}

// RunAsync executes an agent asynchronously and returns a channel of events
func (r *Runner) RunAsync(ctx context.Context, userID, sessionID string, newMessage *events.Content) (<-chan *events.Event, error) {
	// Get or create session
	session, err := r.getOrCreateSession(userID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create session: %w", err)
	}
	
	// Add user message to session if provided
	if newMessage != nil {
		userEvent := events.NewEvent()
		userEvent.Author = "user"
		userEvent.Content = newMessage
		session.AddEvent(userEvent)
		
		// Persist the event
		r.SessionService.AppendEvent(r.AppName, userID, sessionID, userEvent)
	}
	
	// Create invocation context
	invocationCtx := &agents.InvocationContext{
		Session: *session,
	}
	
	// Execute agent
	eventChan, err := r.Agent.RunAsync(ctx, invocationCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to run agent: %w", err)
	}
	
	// Create output channel that persists events
	outputChan := make(chan *events.Event)
	
	go func() {
		defer close(outputChan)
		
		for event := range eventChan {
			// Persist event to session
			r.SessionService.AppendEvent(r.AppName, userID, sessionID, event)
			
			// Forward event to output channel
			outputChan <- event
		}
	}()
	
	return outputChan, nil
}

// RunLive executes an agent in live mode (bidi-streaming)
func (r *Runner) RunLive(ctx context.Context, userID, sessionID string) (<-chan *events.Event, error) {
	// Get or create session
	session, err := r.getOrCreateSession(userID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create session: %w", err)
	}
	
	// Create invocation context
	invocationCtx := &agents.InvocationContext{
		Session: *session,
	}
	
	// Execute agent in live mode
	eventChan, err := r.Agent.RunLive(ctx, invocationCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to run agent in live mode: %w", err)
	}
	
	// Create output channel that persists events
	outputChan := make(chan *events.Event)
	
	go func() {
		defer close(outputChan)
		
		for event := range eventChan {
			// Persist event to session
			r.SessionService.AppendEvent(r.AppName, userID, sessionID, event)
			
			// Forward event to output channel
			outputChan <- event
		}
	}()
	
	return outputChan, nil
}

// CloseSession closes a session
func (r *Runner) CloseSession(userID, sessionID string) error {
	return r.SessionService.CloseSession(r.AppName, userID, sessionID)
}

// getOrCreateSession retrieves an existing session or creates a new one
func (r *Runner) getOrCreateSession(userID, sessionID string) (*sessions.Session, error) {
	// Try to get existing session
	session, err := r.SessionService.GetSession(r.AppName, userID, sessionID)
	if err != nil {
		return nil, err
	}
	
	// If session doesn't exist, create a new one
	if session == nil {
		session, err = r.SessionService.CreateSession(r.AppName, userID, sessionID, nil)
		if err != nil {
			return nil, err
		}
	}
	
	return session, nil
}

// InMemoryRunner provides an in-memory implementation of Runner
type InMemoryRunner struct {
	*Runner
}

// NewInMemoryRunner creates a new in-memory runner
func NewInMemoryRunner(agent agents.Agent, appName string) *InMemoryRunner {
	sessionService := sessions.NewInMemorySessionService()
	runner := NewRunner(agent, appName, sessionService)
	
	return &InMemoryRunner{
		Runner: runner,
	}
}