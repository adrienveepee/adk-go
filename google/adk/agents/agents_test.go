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

	"github.com/adrienveepee/adk-go/google/adk/events"
	"github.com/adrienveepee/adk-go/google/adk/sessions"
)

func TestNewBaseAgent(t *testing.T) {
	name := "test_agent"
	description := "A test agent"
	
	agent := NewBaseAgent(name, description)
	
	if agent.GetName() != name {
		t.Errorf("Expected agent name to be %s, got %s", name, agent.GetName())
	}
	
	if agent.GetDescription() != description {
		t.Errorf("Expected agent description to be %s, got %s", description, agent.GetDescription())
	}
	
	if agent.GetSubAgents() == nil {
		t.Error("SubAgents should be initialized")
	}
	
	if len(agent.GetSubAgents()) != 0 {
		t.Errorf("Expected 0 sub-agents, got %d", len(agent.GetSubAgents()))
	}
}

func TestBaseAgentAddSubAgent(t *testing.T) {
	parentAgent := NewBaseAgent("parent", "Parent agent")
	childAgent := NewBaseAgent("child", "Child agent")
	
	parentAgent.AddSubAgent(childAgent)
	
	if len(parentAgent.GetSubAgents()) != 1 {
		t.Errorf("Expected 1 sub-agent, got %d", len(parentAgent.GetSubAgents()))
	}
	
	if parentAgent.GetSubAgents()[0].GetName() != "child" {
		t.Errorf("Expected sub-agent name to be 'child', got %s", parentAgent.GetSubAgents()[0].GetName())
	}
	
	if childAgent.GetParentAgent() != parentAgent {
		t.Error("Child agent parent should be set to parent agent")
	}
}

func TestBaseAgentFindAgent(t *testing.T) {
	rootAgent := NewBaseAgent("root", "Root agent")
	childAgent := NewBaseAgent("child", "Child agent")
	grandchildAgent := NewBaseAgent("grandchild", "Grandchild agent")
	
	rootAgent.AddSubAgent(childAgent)
	childAgent.AddSubAgent(grandchildAgent)
	
	// Test finding self
	found := rootAgent.FindAgent("root")
	if found != rootAgent {
		t.Error("Should find self")
	}
	
	// Test finding direct child
	found = rootAgent.FindAgent("child")
	if found != childAgent {
		t.Error("Should find direct child")
	}
	
	// Test finding grandchild
	found = rootAgent.FindAgent("grandchild")
	if found != grandchildAgent {
		t.Error("Should find grandchild")
	}
	
	// Test finding non-existent agent
	found = rootAgent.FindAgent("non_existent")
	if found != nil {
		t.Error("Should not find non-existent agent")
	}
}

func TestBaseAgentFindSubAgent(t *testing.T) {
	parentAgent := NewBaseAgent("parent", "Parent agent")
	childAgent := NewBaseAgent("child", "Child agent")
	grandchildAgent := NewBaseAgent("grandchild", "Grandchild agent")
	
	parentAgent.AddSubAgent(childAgent)
	childAgent.AddSubAgent(grandchildAgent)
	
	// Test finding direct sub-agent
	found := parentAgent.FindSubAgent("child")
	if found != childAgent {
		t.Error("Should find direct sub-agent")
	}
	
	// Test not finding grandchild (not direct sub-agent)
	found = parentAgent.FindSubAgent("grandchild")
	if found != nil {
		t.Error("Should not find grandchild as direct sub-agent")
	}
	
	// Test finding non-existent agent
	found = parentAgent.FindSubAgent("non_existent")
	if found != nil {
		t.Error("Should not find non-existent agent")
	}
}

func TestBaseAgentGetRootAgent(t *testing.T) {
	rootAgent := NewBaseAgent("root", "Root agent")
	childAgent := NewBaseAgent("child", "Child agent")
	grandchildAgent := NewBaseAgent("grandchild", "Grandchild agent")
	
	rootAgent.AddSubAgent(childAgent)
	childAgent.AddSubAgent(grandchildAgent)
	
	// Test root agent returns self
	if rootAgent.GetRootAgent() != rootAgent {
		t.Error("Root agent should return self as root")
	}
	
	// Test child agent returns root
	if childAgent.GetRootAgent() != rootAgent {
		t.Error("Child agent should return root agent")
	}
	
	// Test grandchild agent returns root
	if grandchildAgent.GetRootAgent() != rootAgent {
		t.Error("Grandchild agent should return root agent")
	}
}

func TestBaseAgentRunAsync(t *testing.T) {
	agent := NewBaseAgent("test", "Test agent")
	
	session := sessions.NewSession("app", "user", "session", nil)
	invocationCtx := &InvocationContext{Session: *session}
	
	eventChan, err := agent.RunAsync(context.Background(), invocationCtx)
	if err != nil {
		t.Errorf("RunAsync should not return error: %v", err)
	}
	
	// Base implementation should return empty channel
	eventCount := 0
	for range eventChan {
		eventCount++
	}
	
	if eventCount != 0 {
		t.Errorf("Base agent should return 0 events, got %d", eventCount)
	}
}

func TestNewLlmAgent(t *testing.T) {
	name := "test_llm_agent"
	model := "gemini-2.0-flash"
	instruction := "You are a test assistant"
	
	agent := NewLlmAgent(name, model, instruction)
	
	if agent.GetName() != name {
		t.Errorf("Expected agent name to be %s, got %s", name, agent.GetName())
	}
	
	if agent.Model != model {
		t.Errorf("Expected model to be %s, got %s", model, agent.Model)
	}
	
	if agent.Instruction != instruction {
		t.Errorf("Expected instruction to be %s, got %s", instruction, agent.Instruction)
	}
	
	if agent.IncludeContents != IncludeContentsDefault {
		t.Errorf("Expected default include contents, got %s", agent.IncludeContents)
	}
}

func TestNewAgent(t *testing.T) {
	// Test that NewAgent is an alias for NewLlmAgent
	name := "test_agent"
	model := "gemini-2.0-flash"
	instruction := "You are a test assistant"
	
	agent := NewAgent(name, model, instruction)
	
	if agent.GetName() != name {
		t.Errorf("Expected agent name to be %s, got %s", name, agent.GetName())
	}
	
	if agent.Model != model {
		t.Errorf("Expected model to be %s, got %s", model, agent.Model)
	}
}

func TestLlmAgentSetters(t *testing.T) {
	agent := NewLlmAgent("test", "gemini-2.0-flash", "test instruction")
	
	// Test SetDescription
	description := "Test description"
	agent.SetDescription(description)
	if agent.GetDescription() != description {
		t.Errorf("Expected description to be %s, got %s", description, agent.GetDescription())
	}
	
	// Test SetOutputKey
	outputKey := "test_output"
	agent.SetOutputKey(outputKey)
	if agent.OutputKey != outputKey {
		t.Errorf("Expected output key to be %s, got %s", outputKey, agent.OutputKey)
	}
	
	// Test SetIncludeContents
	agent.SetIncludeContents(IncludeContentsNone)
	if agent.IncludeContents != IncludeContentsNone {
		t.Errorf("Expected include contents to be %s, got %s", IncludeContentsNone, agent.IncludeContents)
	}
}

func TestLlmAgentGetCanonicalInstruction(t *testing.T) {
	agent := NewLlmAgent("test", "gemini-2.0-flash", "Main instruction")
	
	// Test without global instruction
	if agent.GetCanonicalInstruction() != "Main instruction" {
		t.Errorf("Expected canonical instruction to be 'Main instruction', got %s", agent.GetCanonicalInstruction())
	}
	
	// Test with global instruction
	agent.GlobalInstruction = "Global instruction"
	expected := "Global instruction\n\nMain instruction"
	if agent.GetCanonicalInstruction() != expected {
		t.Errorf("Expected canonical instruction to be %s, got %s", expected, agent.GetCanonicalInstruction())
	}
}

func TestSequentialAgent(t *testing.T) {
	agent1 := NewBaseAgent("agent1", "First agent")
	agent2 := NewBaseAgent("agent2", "Second agent")
	
	sequentialAgent := NewSequentialAgent("sequential", []Agent{agent1, agent2})
	
	if sequentialAgent.GetName() != "sequential" {
		t.Errorf("Expected sequential agent name to be 'sequential', got %s", sequentialAgent.GetName())
	}
	
	if len(sequentialAgent.GetSubAgents()) != 2 {
		t.Errorf("Expected 2 sub-agents, got %d", len(sequentialAgent.GetSubAgents()))
	}
}

func TestParallelAgent(t *testing.T) {
	agent1 := NewBaseAgent("agent1", "First agent")
	agent2 := NewBaseAgent("agent2", "Second agent")
	
	parallelAgent := NewParallelAgent("parallel", []Agent{agent1, agent2})
	
	if parallelAgent.GetName() != "parallel" {
		t.Errorf("Expected parallel agent name to be 'parallel', got %s", parallelAgent.GetName())
	}
	
	if len(parallelAgent.GetSubAgents()) != 2 {
		t.Errorf("Expected 2 sub-agents, got %d", len(parallelAgent.GetSubAgents()))
	}
}

func TestLoopAgent(t *testing.T) {
	agent1 := NewBaseAgent("agent1", "First agent")
	maxIterations := 3
	
	loopAgent := NewLoopAgent("loop", []Agent{agent1}, maxIterations)
	
	if loopAgent.GetName() != "loop" {
		t.Errorf("Expected loop agent name to be 'loop', got %s", loopAgent.GetName())
	}
	
	if len(loopAgent.GetSubAgents()) != 1 {
		t.Errorf("Expected 1 sub-agent, got %d", len(loopAgent.GetSubAgents()))
	}
	
	if loopAgent.MaxIterations != maxIterations {
		t.Errorf("Expected max iterations to be %d, got %d", maxIterations, loopAgent.MaxIterations)
	}
}

func TestLoopAgentShouldExitLoop(t *testing.T) {
	agent1 := NewBaseAgent("agent1", "First agent")
	loopAgent := NewLoopAgent("loop", []Agent{agent1}, 3)
	
	session := sessions.NewSession("app", "user", "session", nil)
	invocationCtx := &InvocationContext{Session: *session}
	
	// Test default state (should not exit)
	if loopAgent.shouldExitLoop(invocationCtx) {
		t.Error("Should not exit loop by default")
	}
	
	// Test with exit flag set
	session.State.Set("exit_loop", true)
	if !loopAgent.shouldExitLoop(invocationCtx) {
		t.Error("Should exit loop when exit_loop flag is true")
	}
	
	// Test with false exit flag
	session.State.Set("exit_loop", false)
	if loopAgent.shouldExitLoop(invocationCtx) {
		t.Error("Should not exit loop when exit_loop flag is false")
	}
}

func TestLoopAgentShouldExitLoopFromEvent(t *testing.T) {
	agent1 := NewBaseAgent("agent1", "First agent")
	loopAgent := NewLoopAgent("loop", []Agent{agent1}, 3)
	
	event := events.NewEvent()
	
	// Test default event (should not exit)
	if loopAgent.shouldExitLoopFromEvent(event) {
		t.Error("Should not exit loop from default event")
	}
	
	// Test with skip summarization set (proxy for exit loop)
	event.Actions.SkipSummarization = true
	if !loopAgent.shouldExitLoopFromEvent(event) {
		t.Error("Should exit loop when SkipSummarization is true")
	}
}