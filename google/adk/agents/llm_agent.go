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
	"github.com/adrienveepee/adk-go/google/adk/models"
	"github.com/adrienveepee/adk-go/google/adk/tools"
)

// IncludeContents determines how conversation history is included
type IncludeContents string

const (
	IncludeContentsDefault IncludeContents = "default"
	IncludeContentsNone    IncludeContents = "none"
)

// LlmAgent represents an agent powered by a Large Language Model
type LlmAgent struct {
	*BaseAgent
	
	// Core LLM configuration
	Model                 string                        `json:"model"`
	Instruction           string                        `json:"instruction"`
	GlobalInstruction     string                        `json:"global_instruction,omitempty"`
	GenerateContentConfig *models.GenerateContentConfig `json:"generate_content_config,omitempty"`
	
	// Input/Output configuration
	InputSchema    interface{}     `json:"input_schema,omitempty"`
	OutputSchema   interface{}     `json:"output_schema,omitempty"`
	OutputKey      string          `json:"output_key,omitempty"`
	IncludeContents IncludeContents `json:"include_contents,omitempty"`
	
	// Tools and capabilities
	Tools        []tools.Tool `json:"tools,omitempty"`
	Examples     []interface{} `json:"examples,omitempty"`
	CodeExecutor interface{}   `json:"code_executor,omitempty"`
	Planner      interface{}   `json:"planner,omitempty"`
	
	// Transfer control settings
	DisallowTransferToParent bool `json:"disallow_transfer_to_parent,omitempty"`
	DisallowTransferToPeers  bool `json:"disallow_transfer_to_peers,omitempty"`
	
	// Callbacks
	BeforeModelCallback func(*InvocationContext) error `json:"-"`
	AfterModelCallback  func(*InvocationContext) error `json:"-"`
	BeforeToolCallback  func(*InvocationContext) error `json:"-"`
	AfterToolCallback   func(*InvocationContext) error `json:"-"`
	
	// Internal
	llm models.LLM `json:"-"`
}

// NewLlmAgent creates a new LLM agent
func NewLlmAgent(name, model, instruction string) *LlmAgent {
	return &LlmAgent{
		BaseAgent:       NewBaseAgent(name, ""),
		Model:           model,
		Instruction:     instruction,
		IncludeContents: IncludeContentsDefault,
		Tools:           make([]tools.Tool, 0),
		Examples:        make([]interface{}, 0),
	}
}

// NewAgent is an alias for NewLlmAgent for convenience
func NewAgent(name, model, instruction string) *LlmAgent {
	return NewLlmAgent(name, model, instruction)
}

// SetDescription sets the agent description
func (a *LlmAgent) SetDescription(description string) *LlmAgent {
	a.Description = description
	return a
}

// SetTools sets the agent tools
func (a *LlmAgent) SetTools(tools []tools.Tool) *LlmAgent {
	a.Tools = tools
	return a
}

// AddTool adds a tool to the agent
func (a *LlmAgent) AddTool(tool tools.Tool) *LlmAgent {
	a.Tools = append(a.Tools, tool)
	return a
}

// SetOutputKey sets the output key for storing results in session state
func (a *LlmAgent) SetOutputKey(outputKey string) *LlmAgent {
	a.OutputKey = outputKey
	return a
}

// SetOutputSchema sets the output schema for structured responses
func (a *LlmAgent) SetOutputSchema(schema interface{}) *LlmAgent {
	a.OutputSchema = schema
	return a
}

// SetInputSchema sets the input schema for structured inputs
func (a *LlmAgent) SetInputSchema(schema interface{}) *LlmAgent {
	a.InputSchema = schema
	return a
}

// SetGenerateContentConfig sets the content generation configuration
func (a *LlmAgent) SetGenerateContentConfig(config *models.GenerateContentConfig) *LlmAgent {
	a.GenerateContentConfig = config
	return a
}

// SetIncludeContents sets how conversation history is included
func (a *LlmAgent) SetIncludeContents(includeContents IncludeContents) *LlmAgent {
	a.IncludeContents = includeContents
	return a
}

// GetCanonicalModel returns the resolved LLM model
func (a *LlmAgent) GetCanonicalModel() models.LLM {
	if a.llm == nil {
		llm, err := models.Resolve(a.Model)
		if err != nil {
			// TODO: Better error handling
			return nil
		}
		a.llm = llm
	}
	return a.llm
}

// GetCanonicalInstruction returns the complete instruction including global instruction
func (a *LlmAgent) GetCanonicalInstruction() string {
	instruction := a.Instruction
	if a.GlobalInstruction != "" {
		instruction = a.GlobalInstruction + "\n\n" + instruction
	}
	return instruction
}

// GetCanonicalTools returns all tools including sub-agent tools
func (a *LlmAgent) GetCanonicalTools() []tools.Tool {
	allTools := make([]tools.Tool, len(a.Tools))
	copy(allTools, a.Tools)
	
	// Add sub-agent tools for delegation
	for _, subAgent := range a.SubAgents {
		agentTool := tools.NewAgentTool(subAgent)
		allTools = append(allTools, agentTool)
	}
	
	return allTools
}

// RunAsync executes the LLM agent asynchronously
func (a *LlmAgent) RunAsync(ctx context.Context, invocationCtx *InvocationContext) (<-chan *events.Event, error) {
	eventChan := make(chan *events.Event)
	
	go func() {
		defer close(eventChan)
		
		// Execute before agent callback
		if a.BeforeAgentCallback != nil {
			if err := a.BeforeAgentCallback(invocationCtx); err != nil {
				// TODO: Better error handling
				return
			}
		}
		
		// Execute before model callback
		if a.BeforeModelCallback != nil {
			if err := a.BeforeModelCallback(invocationCtx); err != nil {
				// TODO: Better error handling
				return
			}
		}
		
		// Get LLM model
		llm := a.GetCanonicalModel()
		if llm == nil {
			// TODO: Better error handling
			return
		}
		
		// Build LLM request
		request := a.buildLLMRequest(invocationCtx)
		
		// Generate content
		responseEventChan, err := llm.GenerateContentAsync(ctx, request)
		if err != nil {
			// TODO: Better error handling
			return
		}
		
		// Process events and handle tool calls
		for event := range responseEventChan {
			// Process tool calls if any
			if a.hasToolCalls(event) {
				toolEvents, err := a.processToolCalls(ctx, event, invocationCtx)
				if err != nil {
					// TODO: Better error handling
					continue
				}
				
				// Forward tool events
				for _, toolEvent := range toolEvents {
					eventChan <- toolEvent
				}
			} else {
				// Handle output key storage
				if a.OutputKey != "" && event.IsFinalResponse && event.Content != nil {
					a.storeOutputInSession(event, invocationCtx)
				}
				
				// Forward the event
				eventChan <- event
			}
		}
		
		// Execute after model callback
		if a.AfterModelCallback != nil {
			if err := a.AfterModelCallback(invocationCtx); err != nil {
				// TODO: Better error handling
				return
			}
		}
		
		// Execute after agent callback
		if a.AfterAgentCallback != nil {
			if err := a.AfterAgentCallback(invocationCtx); err != nil {
				// TODO: Better error handling
				return
			}
		}
	}()
	
	return eventChan, nil
}

// buildLLMRequest builds the LLM request from the agent configuration
func (a *LlmAgent) buildLLMRequest(invocationCtx *InvocationContext) *models.LLMRequest {
	request := &models.LLMRequest{
		Config: a.GenerateContentConfig,
	}
	
	// Add conversation history if requested
	if a.IncludeContents == IncludeContentsDefault {
		request.Contents = a.buildContents(invocationCtx)
	} else {
		// Include only the instruction as system message
		request.Contents = []*events.Content{
			{
				Role: "system",
				Parts: []events.Part{
					{Text: a.GetCanonicalInstruction()},
				},
			},
		}
	}
	
	// Add tools
	canonicalTools := a.GetCanonicalTools()
	request.Tools = make([]interface{}, len(canonicalTools))
	for i, tool := range canonicalTools {
		request.Tools[i] = tool
	}
	
	return request
}

// buildContents builds the conversation contents from session history
func (a *LlmAgent) buildContents(invocationCtx *InvocationContext) []*events.Content {
	contents := make([]*events.Content, 0)
	
	// Add system instruction
	if instruction := a.GetCanonicalInstruction(); instruction != "" {
		contents = append(contents, &events.Content{
			Role: "system",
			Parts: []events.Part{
				{Text: instruction},
			},
		})
	}
	
	// Add session events as conversation history
	for _, event := range invocationCtx.Session.Events {
		if event.Content != nil {
			contents = append(contents, event.Content)
		}
	}
	
	return contents
}

// hasToolCalls checks if an event contains tool calls
func (a *LlmAgent) hasToolCalls(event *events.Event) bool {
	// TODO: Implement tool call detection
	return false
}

// processToolCalls processes tool calls in an event
func (a *LlmAgent) processToolCalls(ctx context.Context, event *events.Event, invocationCtx *InvocationContext) ([]*events.Event, error) {
	// TODO: Implement tool call processing
	// For now, create a mock tool context
	toolCtx := tools.NewToolContext(invocationCtx, "mock-function-call-id")
	_ = toolCtx // Use the variable to avoid unused variable error
	
	return []*events.Event{}, nil
}

// storeOutputInSession stores the agent output in session state
func (a *LlmAgent) storeOutputInSession(event *events.Event, invocationCtx *InvocationContext) {
	if event.Content != nil && len(event.Content.Parts) > 0 {
		text := event.Content.Parts[0].Text
		invocationCtx.Session.State.Set(a.OutputKey, text)
	}
}