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
	"context"
	"testing"

	"github.com/adrienveepee/adk-go/google/adk/events"
)

func TestNewBaseLLM(t *testing.T) {
	modelName := "test-model"
	llm := NewBaseLLM(modelName)
	
	if llm.GetModelName() != modelName {
		t.Errorf("Expected model name to be %s, got %s", modelName, llm.GetModelName())
	}
}

func TestNewGeminiLLM(t *testing.T) {
	modelName := "gemini-2.0-flash"
	llm := NewGeminiLLM(modelName)
	
	if llm.GetModelName() != modelName {
		t.Errorf("Expected model name to be %s, got %s", modelName, llm.GetModelName())
	}
	
	if llm.BaseLLM == nil {
		t.Error("BaseLLM should be initialized")
	}
}

func TestGeminiLLMConnect(t *testing.T) {
	llm := NewGeminiLLM("gemini-2.0-flash")
	
	err := llm.Connect(context.Background())
	if err != nil {
		t.Errorf("Connect should not return error: %v", err)
	}
}

func TestGeminiLLMGenerateContentAsync(t *testing.T) {
	llm := NewGeminiLLM("gemini-2.0-flash")
	
	request := &LLMRequest{
		Contents: []*events.Content{
			{
				Role: "user",
				Parts: []events.Part{
					{Text: "Hello"},
				},
			},
		},
	}
	
	eventChan, err := llm.GenerateContentAsync(context.Background(), request)
	if err != nil {
		t.Errorf("GenerateContentAsync should not return error: %v", err)
	}
	
	// Should receive at least one event
	eventCount := 0
	for event := range eventChan {
		eventCount++
		
		if event.Author != llm.GetModelName() {
			t.Errorf("Expected event author to be %s, got %s", llm.GetModelName(), event.Author)
		}
		
		if event.Content == nil {
			t.Error("Event should have content")
		}
		
		if !event.IsFinalResponse {
			t.Error("Event should be marked as final response")
		}
	}
	
	if eventCount == 0 {
		t.Error("Should receive at least one event")
	}
}

func TestGeminiLLMSupportedModels(t *testing.T) {
	llm := NewGeminiLLM("gemini-2.0-flash")
	
	models := llm.SupportedModels()
	
	if len(models) == 0 {
		t.Error("Should have supported models")
	}
	
	// Check for expected models
	expectedModels := []string{
		"gemini-2.0-flash",
		"gemini-1.5-pro",
		"gemini-1.5-flash",
		"gemini-1.0-pro",
	}
	
	for _, expected := range expectedModels {
		found := false
		for _, model := range models {
			if model == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected model %s not found in supported models", expected)
		}
	}
}

func TestLLMRegistry(t *testing.T) {
	// Test NewLLM with Gemini model
	llm, err := NewLLM("gemini-2.0-flash")
	if err != nil {
		t.Errorf("NewLLM should not return error for gemini model: %v", err)
	}
	
	if llm == nil {
		t.Error("NewLLM should return an LLM instance")
	}
	
	if llm.GetModelName() != "gemini-2.0-flash" {
		t.Errorf("Expected model name to be 'gemini-2.0-flash', got %s", llm.GetModelName())
	}
	
	// Test NewLLM with unsupported model
	_, err = NewLLM("unsupported-model")
	if err == nil {
		t.Error("NewLLM should return error for unsupported model")
	}
}

func TestResolve(t *testing.T) {
	// Test Resolve (alias for NewLLM)
	llm, err := Resolve("gemini-1.5-pro")
	if err != nil {
		t.Errorf("Resolve should not return error for gemini model: %v", err)
	}
	
	if llm == nil {
		t.Error("Resolve should return an LLM instance")
	}
	
	if llm.GetModelName() != "gemini-1.5-pro" {
		t.Errorf("Expected model name to be 'gemini-1.5-pro', got %s", llm.GetModelName())
	}
}

func TestGenerateContentConfig(t *testing.T) {
	temperature := float32(0.7)
	maxTokens := 1000
	topP := float32(0.9)
	topK := 40
	
	config := &GenerateContentConfig{
		Temperature:     &temperature,
		MaxOutputTokens: &maxTokens,
		TopP:            &topP,
		TopK:            &topK,
	}
	
	if *config.Temperature != temperature {
		t.Errorf("Expected temperature to be %f, got %f", temperature, *config.Temperature)
	}
	
	if *config.MaxOutputTokens != maxTokens {
		t.Errorf("Expected max output tokens to be %d, got %d", maxTokens, *config.MaxOutputTokens)
	}
	
	if *config.TopP != topP {
		t.Errorf("Expected topP to be %f, got %f", topP, *config.TopP)
	}
	
	if *config.TopK != topK {
		t.Errorf("Expected topK to be %d, got %d", topK, *config.TopK)
	}
}

func TestLLMRequest(t *testing.T) {
	config := &GenerateContentConfig{
		Temperature: func() *float32 { v := float32(0.8); return &v }(),
	}
	
	contents := []*events.Content{
		{
			Role: "user",
			Parts: []events.Part{
				{Text: "Hello"},
			},
		},
	}
	
	tools := []interface{}{
		"test_tool",
	}
	
	request := &LLMRequest{
		Contents: contents,
		Config:   config,
		Tools:    tools,
	}
	
	if len(request.Contents) != 1 {
		t.Errorf("Expected 1 content, got %d", len(request.Contents))
	}
	
	if request.Contents[0].Role != "user" {
		t.Errorf("Expected content role to be 'user', got %s", request.Contents[0].Role)
	}
	
	if request.Config.Temperature == nil || *request.Config.Temperature != 0.8 {
		t.Error("Config temperature should be set correctly")
	}
	
	if len(request.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(request.Tools))
	}
}

func TestLLMResponse(t *testing.T) {
	content := &events.Content{
		Role: "model",
		Parts: []events.Part{
			{Text: "Hello, how can I help you?"},
		},
	}
	
	response := &LLMResponse{
		Content: content,
	}
	
	if response.Content.Role != "model" {
		t.Errorf("Expected content role to be 'model', got %s", response.Content.Role)
	}
	
	if len(response.Content.Parts) != 1 {
		t.Errorf("Expected 1 content part, got %d", len(response.Content.Parts))
	}
	
	if response.Content.Parts[0].Text != "Hello, how can I help you?" {
		t.Errorf("Expected content text to match")
	}
}

func TestRegistryInstanceCaching(t *testing.T) {
	// Create the same model twice
	llm1, err1 := NewLLM("gemini-2.0-flash")
	llm2, err2 := NewLLM("gemini-2.0-flash")
	
	if err1 != nil || err2 != nil {
		t.Errorf("NewLLM should not return errors: %v, %v", err1, err2)
	}
	
	// Both should be the same instance due to caching
	if llm1 != llm2 {
		t.Error("Registry should return the same instance for the same model")
	}
}