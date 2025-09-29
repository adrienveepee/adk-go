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

// Package main demonstrates the ADK Go SDK functionality
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/adrienveepee/adk-go/google/adk/agents"
	"github.com/adrienveepee/adk-go/google/adk/events"
	"github.com/adrienveepee/adk-go/google/adk/runners"
	"github.com/adrienveepee/adk-go/google/adk/sessions"
	"github.com/adrienveepee/adk-go/google/adk/tools"
)

// Example function that can be used as a tool
func getCapitalCity(country string) string {
	capitals := map[string]string{
		"france":        "Paris",
		"japan":         "Tokyo",
		"canada":        "Ottawa",
		"united states": "Washington, D.C.",
		"germany":       "Berlin",
	}
	
	if capital, exists := capitals[country]; exists {
		return capital
	}
	return fmt.Sprintf("Sorry, I don't know the capital of %s", country)
}

func main() {
	ctx := context.Background()
	
	// Create a function tool from the getCapitalCity function
	capitalTool, err := tools.NewFunctionTool(getCapitalCity)
	if err != nil {
		log.Fatalf("Failed to create capital tool: %v", err)
	}
	capitalTool.Name = "get_capital_city"
	capitalTool.Description = "Get the capital city of a country"
	
	// Create an LLM agent
	agent := agents.NewAgent(
		"capital_assistant",
		"gemini-2.0-flash",
		"You are a helpful assistant that can provide capital cities of countries. Use the get_capital_city tool when users ask about capitals.",
	).SetDescription("An assistant that provides capital city information").
		SetTools([]tools.Tool{capitalTool}).
		SetOutputKey("capital_response")
	
	// Create a session service
	sessionService := sessions.NewInMemorySessionService()
	
	// Create a runner
	runner := runners.NewRunner(agent, "capital_app", sessionService)
	
	// Create a user message
	userMessage := &events.Content{
		Role: "user",
		Parts: []events.Part{
			{Text: "What is the capital of France?"},
		},
	}
	
	// Run the agent
	fmt.Println("Running ADK Go SDK example...")
	fmt.Printf("User: %s\n", userMessage.Parts[0].Text)
	
	eventChan, err := runner.RunAsync(ctx, "user123", "session456", userMessage)
	if err != nil {
		log.Fatalf("Failed to run agent: %v", err)
	}
	
	// Process events
	for event := range eventChan {
		if event.Content != nil && event.Content.Role == "model" {
			fmt.Printf("Assistant: %s\n", event.Content.Parts[0].Text)
		}
	}
	
	// Demonstrate workflow agents
	fmt.Println("\nDemonstrating workflow agents...")
	
	// Create sub-agents for workflow
	greeter := agents.NewAgent(
		"greeter",
		"gemini-2.0-flash",
		"You are a friendly greeter. Say hello to the user.",
	).SetOutputKey("greeting")
	
	taskExecutor := agents.NewAgent(
		"task_executor", 
		"gemini-2.0-flash",
		"You are a task executor. Help the user with their request.",
	).SetOutputKey("task_result")
	
	// Create a sequential agent
	sequentialAgent := agents.NewSequentialAgent(
		"sequential_workflow",
		[]agents.Agent{greeter, taskExecutor},
	)
	
	// Create a parallel agent
	parallelAgent := agents.NewParallelAgent(
		"parallel_workflow",
		[]agents.Agent{greeter, taskExecutor},
	)
	
	// Create a loop agent
	loopAgent := agents.NewLoopAgent(
		"loop_workflow",
		[]agents.Agent{greeter},
		2, // Max 2 iterations
	)
	
	// Test each workflow agent
	workflows := []struct {
		name  string
		agent agents.Agent
	}{
		{"Sequential", sequentialAgent},
		{"Parallel", parallelAgent}, 
		{"Loop", loopAgent},
	}
	
	for _, workflow := range workflows {
		fmt.Printf("\nTesting %s Workflow:\n", workflow.name)
		
		workflowRunner := runners.NewRunner(workflow.agent, "workflow_app", sessionService)
		
		workflowMessage := &events.Content{
			Role: "user",
			Parts: []events.Part{
				{Text: "Hello, please help me with a task."},
			},
		}
		
		eventChan, err := workflowRunner.RunAsync(ctx, "user123", "workflow_session", workflowMessage)
		if err != nil {
			log.Printf("Failed to run %s workflow: %v", workflow.name, err)
			continue
		}
		
		for event := range eventChan {
			if event.Content != nil && event.Content.Role == "model" {
				fmt.Printf("  %s: %s\n", event.Author, event.Content.Parts[0].Text)
			}
		}
	}
	
	fmt.Println("\nADK Go SDK demonstration complete!")
}