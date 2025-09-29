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

package tools

import (
	"context"
	"fmt"
	"reflect"

	"github.com/adrienveepee/adk-go/google/adk/events"
	"github.com/adrienveepee/adk-go/google/adk/models"
)

// ToolContext provides the context for tool execution
type ToolContext struct {
	InvocationContext interface{} // Will be *agents.InvocationContext, using interface{} to avoid import cycle
	FunctionCallID    string
	EventActions      *events.EventActions
}

// NewToolContext creates a new tool context
func NewToolContext(invocationCtx interface{}, functionCallID string) *ToolContext {
	return &ToolContext{
		InvocationContext: invocationCtx,
		FunctionCallID:    functionCallID,
		EventActions:      &events.EventActions{},
	}
}

// Tool is the interface that all tools must implement
type Tool interface {
	// GetName returns the tool's name
	GetName() string
	
	// GetDescription returns the tool's description
	GetDescription() string
	
	// IsLongRunning returns whether the tool is long running
	IsLongRunning() bool
	
	// RunAsync executes the tool asynchronously
	RunAsync(ctx context.Context, args map[string]interface{}, toolCtx *ToolContext) (interface{}, error)
	
	// ProcessLLMRequest processes the outgoing LLM request for this tool
	ProcessLLMRequest(toolCtx *ToolContext, llmRequest *models.LLMRequest) error
}

// BaseTool provides the base implementation for all tools
type BaseTool struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	IsLongRunning_  bool   `json:"is_long_running"`
}

// NewBaseTool creates a new base tool
func NewBaseTool(name, description string, isLongRunning bool) *BaseTool {
	return &BaseTool{
		Name:            name,
		Description:     description,
		IsLongRunning_:  isLongRunning,
	}
}

// GetName returns the tool's name
func (t *BaseTool) GetName() string {
	return t.Name
}

// GetDescription returns the tool's description
func (t *BaseTool) GetDescription() string {
	return t.Description
}

// IsLongRunning returns whether the tool is long running
func (t *BaseTool) IsLongRunning() bool {
	return t.IsLongRunning_
}

// RunAsync is the base implementation - to be overridden by concrete tools
func (t *BaseTool) RunAsync(ctx context.Context, args map[string]interface{}, toolCtx *ToolContext) (interface{}, error) {
	return nil, fmt.Errorf("RunAsync not implemented for tool: %s", t.Name)
}

// ProcessLLMRequest is the base implementation for processing LLM requests
func (t *BaseTool) ProcessLLMRequest(toolCtx *ToolContext, llmRequest *models.LLMRequest) error {
	// Base implementation adds the tool to the LLM request
	// TODO: Convert tool to appropriate format for LLM
	return nil
}

// FunctionTool wraps a Go function as a tool
type FunctionTool struct {
	*BaseTool
	Function interface{} `json:"-"`
}

// NewFunctionTool creates a new function tool from a Go function
func NewFunctionTool(fn interface{}) (*FunctionTool, error) {
	fnType := reflect.TypeOf(fn)
	
	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("expected function, got %s", fnType.Kind())
	}
	
	// Extract name from function (simplified - could be improved)
	name := "function_tool"
	description := "A function tool"
	
	return &FunctionTool{
		BaseTool: NewBaseTool(name, description, false),
		Function: fn,
	}, nil
}

// RunAsync executes the wrapped function
func (ft *FunctionTool) RunAsync(ctx context.Context, args map[string]interface{}, toolCtx *ToolContext) (interface{}, error) {
	fnValue := reflect.ValueOf(ft.Function)
	fnType := reflect.TypeOf(ft.Function)
	
	// Create input values for function call
	inputs := make([]reflect.Value, fnType.NumIn())
	
	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)
		paramName := fmt.Sprintf("param%d", i) // Simplified parameter naming
		
		if argValue, exists := args[paramName]; exists {
			inputs[i] = reflect.ValueOf(argValue)
		} else {
			// Use zero value for missing parameters
			inputs[i] = reflect.Zero(paramType)
		}
	}
	
	// Call the function
	results := fnValue.Call(inputs)
	
	// Return the first result (simplified)
	if len(results) > 0 {
		return results[0].Interface(), nil
	}
	
	return nil, nil
}

// ExampleTool adds examples to the LLM request
type ExampleTool struct {
	*BaseTool
	Examples []interface{} `json:"examples"`
}

// NewExampleTool creates a new example tool
func NewExampleTool(examples []interface{}) *ExampleTool {
	return &ExampleTool{
		BaseTool: NewBaseTool("example_tool", "A tool that adds examples to the LLM request", false),
		Examples: examples,
	}
}

// ProcessLLMRequest adds examples to the LLM request
func (et *ExampleTool) ProcessLLMRequest(toolCtx *ToolContext, llmRequest *models.LLMRequest) error {
	// TODO: Add examples to the LLM request in appropriate format
	return nil
}

// AgentTool wraps an agent as a tool for delegation
type AgentTool struct {
	*BaseTool
	Agent interface{} `json:"-"` // Will be agents.Agent, using interface{} to avoid import cycle
}

// NewAgentTool creates a new agent tool
func NewAgentTool(agent interface{}) *AgentTool {
	// Using reflection to get name and description to avoid import cycle
	agentValue := reflect.ValueOf(agent)
	nameMethod := agentValue.MethodByName("GetName")
	descMethod := agentValue.MethodByName("GetDescription")
	
	var name, description string
	if nameMethod.IsValid() {
		nameResults := nameMethod.Call(nil)
		if len(nameResults) > 0 {
			name = "transfer_to_" + nameResults[0].String()
		}
	}
	if descMethod.IsValid() {
		descResults := descMethod.Call(nil)
		if len(descResults) > 0 {
			description = "Transfer to " + descResults[0].String()
		}
	}
	
	return &AgentTool{
		BaseTool: NewBaseTool(name, description, false),
		Agent:    agent,
	}
}

// RunAsync delegates to the wrapped agent
func (at *AgentTool) RunAsync(ctx context.Context, args map[string]interface{}, toolCtx *ToolContext) (interface{}, error) {
	// Using reflection to call RunAsync on the agent to avoid import cycle
	agentValue := reflect.ValueOf(at.Agent)
	runAsyncMethod := agentValue.MethodByName("RunAsync")
	
	if !runAsyncMethod.IsValid() {
		return nil, fmt.Errorf("agent does not have RunAsync method")
	}
	
	// Call the method with context and invocation context
	ctxValue := reflect.ValueOf(ctx)
	invCtxValue := reflect.ValueOf(toolCtx.InvocationContext)
	
	results := runAsyncMethod.Call([]reflect.Value{ctxValue, invCtxValue})
	
	if len(results) < 2 {
		return nil, fmt.Errorf("unexpected return values from RunAsync")
	}
	
	// Check for error
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}
	
	// For simplicity, return a placeholder result
	// In a real implementation, we would properly handle the event channel
	return "Agent execution completed", nil
}

// Built-in tools

// ExitLoop is a built-in tool for exiting loops
func ExitLoop(toolCtx *ToolContext) error {
	// Set the exit loop flag in event actions
	toolCtx.EventActions.SkipSummarization = true
	return nil
}

// TransferToAgent is a built-in tool for transferring to another agent
func TransferToAgent(agentName string, toolCtx *ToolContext) error {
	toolCtx.EventActions.TransferToAgent = agentName
	return nil
}