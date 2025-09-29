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

package examples

import (
	"context"
	"sync"

	"github.com/adrienveepee/adk-go/google/adk/events"
)

// Example represents a single example with input and output
type Example struct {
	Input  *events.Content `json:"input"`
	Output *events.Content `json:"output"`
}

// NewExample creates a new example
func NewExample(input, output *events.Content) *Example {
	return &Example{
		Input:  input,
		Output: output,
	}
}

// ExampleProvider is the interface that all example providers must implement
type ExampleProvider interface {
	// GetExamples retrieves examples for the given context
	GetExamples(ctx context.Context, query string) ([]*Example, error)
}

// BaseExampleProvider provides the base implementation for all example providers
type BaseExampleProvider struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewBaseExampleProvider creates a new base example provider
func NewBaseExampleProvider(name, description string) *BaseExampleProvider {
	return &BaseExampleProvider{
		Name:        name,
		Description: description,
	}
}

// GetExamples is the base implementation - to be overridden by concrete providers
func (p *BaseExampleProvider) GetExamples(ctx context.Context, query string) ([]*Example, error) {
	return []*Example{}, nil
}

// InMemoryExampleProvider provides an in-memory example store
type InMemoryExampleProvider struct {
	*BaseExampleProvider
	mu       sync.RWMutex
	examples map[string][]*Example // Key: category/query, Value: examples
}

// NewInMemoryExampleProvider creates a new in-memory example provider
func NewInMemoryExampleProvider() *InMemoryExampleProvider {
	return &InMemoryExampleProvider{
		BaseExampleProvider: NewBaseExampleProvider("in_memory_examples", "In-memory example storage"),
		examples:           make(map[string][]*Example),
	}
}

// AddExample adds an example to the provider
func (p *InMemoryExampleProvider) AddExample(category string, example *Example) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.examples[category] == nil {
		p.examples[category] = make([]*Example, 0)
	}
	
	p.examples[category] = append(p.examples[category], example)
}

// AddExamples adds multiple examples to the provider
func (p *InMemoryExampleProvider) AddExamples(category string, examples []*Example) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.examples[category] == nil {
		p.examples[category] = make([]*Example, 0)
	}
	
	p.examples[category] = append(p.examples[category], examples...)
}

// GetExamples retrieves examples for the given query/category
func (p *InMemoryExampleProvider) GetExamples(ctx context.Context, query string) ([]*Example, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	if examples, exists := p.examples[query]; exists {
		// Return a copy of the examples
		result := make([]*Example, len(examples))
		copy(result, examples)
		return result, nil
	}
	
	return []*Example{}, nil
}

// GetAllCategories returns all available categories
func (p *InMemoryExampleProvider) GetAllCategories() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	categories := make([]string, 0, len(p.examples))
	for category := range p.examples {
		categories = append(categories, category)
	}
	
	return categories
}

// VertexAiExampleStore provides a Vertex AI example store
type VertexAiExampleStore struct {
	*BaseExampleProvider
	// TODO: Add Vertex AI specific fields
}

// NewVertexAiExampleStore creates a new Vertex AI example store
func NewVertexAiExampleStore() *VertexAiExampleStore {
	return &VertexAiExampleStore{
		BaseExampleProvider: NewBaseExampleProvider("vertex_ai_examples", "Vertex AI example storage"),
	}
}

// GetExamples retrieves examples from Vertex AI
func (v *VertexAiExampleStore) GetExamples(ctx context.Context, query string) ([]*Example, error) {
	// TODO: Implement Vertex AI example retrieval
	return []*Example{}, nil
}

// ExampleManager manages multiple example providers
type ExampleManager struct {
	mu        sync.RWMutex
	providers map[string]ExampleProvider
}

// NewExampleManager creates a new example manager
func NewExampleManager() *ExampleManager {
	return &ExampleManager{
		providers: make(map[string]ExampleProvider),
	}
}

// RegisterProvider registers an example provider
func (m *ExampleManager) RegisterProvider(name string, provider ExampleProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.providers[name] = provider
}

// GetProvider gets an example provider by name
func (m *ExampleManager) GetProvider(name string) (ExampleProvider, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	provider, exists := m.providers[name]
	return provider, exists
}

// GetExamples retrieves examples from all providers
func (m *ExampleManager) GetExamples(ctx context.Context, query string) ([]*Example, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var allExamples []*Example
	
	for _, provider := range m.providers {
		examples, err := provider.GetExamples(ctx, query)
		if err != nil {
			continue // Skip providers that error
		}
		allExamples = append(allExamples, examples...)
	}
	
	return allExamples, nil
}