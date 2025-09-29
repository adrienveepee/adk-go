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

package evaluation

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrienveepee/adk-go/google/adk/agents"
	"github.com/adrienveepee/adk-go/google/adk/events"
	"github.com/adrienveepee/adk-go/google/adk/runners"
	"github.com/adrienveepee/adk-go/google/adk/sessions"
)

// EvaluationResult represents the result of evaluating an agent
type EvaluationResult struct {
	TestCase    string                 `json:"test_case"`
	Input       *events.Content        `json:"input"`
	Expected    *events.Content        `json:"expected,omitempty"`
	Actual      *events.Content        `json:"actual"`
	Score       float64                `json:"score"`
	Passed      bool                   `json:"passed"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	ExecutionTime time.Duration        `json:"execution_time"`
	Error       string                 `json:"error,omitempty"`
}

// EvaluationSet represents a set of test cases for evaluation
type EvaluationSet struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	TestCases   []*EvaluationTestCase  `json:"test_cases"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// EvaluationTestCase represents a single test case
type EvaluationTestCase struct {
	Name        string                 `json:"name"`
	Input       *events.Content        `json:"input"`
	Expected    *events.Content        `json:"expected,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// EvaluationConfig represents configuration for evaluation
type EvaluationConfig struct {
	MaxConcurrency int                    `json:"max_concurrency,omitempty"`
	Timeout        time.Duration          `json:"timeout,omitempty"`
	Metrics        []string               `json:"metrics,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// EvaluationReport represents the overall evaluation report
type EvaluationReport struct {
	AgentName       string                 `json:"agent_name"`
	EvaluationSet   string                 `json:"evaluation_set"`
	StartTime       time.Time              `json:"start_time"`
	EndTime         time.Time              `json:"end_time"`
	TotalTestCases  int                    `json:"total_test_cases"`
	PassedTestCases int                    `json:"passed_test_cases"`
	FailedTestCases int                    `json:"failed_test_cases"`
	AverageScore    float64                `json:"average_score"`
	Results         []*EvaluationResult    `json:"results"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// AgentEvaluator provides evaluation functionality for agents
type AgentEvaluator struct {
	Config *EvaluationConfig
}

// NewAgentEvaluator creates a new agent evaluator
func NewAgentEvaluator(config *EvaluationConfig) *AgentEvaluator {
	if config == nil {
		config = &EvaluationConfig{
			MaxConcurrency: 5,
			Timeout:        30 * time.Second,
			Metrics:        []string{"accuracy"},
		}
	}
	return &AgentEvaluator{
		Config: config,
	}
}

// Evaluate evaluates an agent against an evaluation set
func (e *AgentEvaluator) Evaluate(ctx context.Context, agent agents.Agent, evalSet *EvaluationSet) (*EvaluationReport, error) {
	report := &EvaluationReport{
		AgentName:     agent.GetName(),
		EvaluationSet: evalSet.Name,
		StartTime:     time.Now(),
		Results:       make([]*EvaluationResult, 0),
		Metadata:      make(map[string]interface{}),
	}
	
	// Create session service and runner
	sessionService := sessions.NewInMemorySessionService()
	runner := runners.NewRunner(agent, "evaluation_app", sessionService)
	
	// Evaluate each test case
	for i, testCase := range evalSet.TestCases {
		result := e.evaluateTestCase(ctx, runner, testCase, i)
		report.Results = append(report.Results, result)
		
		if result.Passed {
			report.PassedTestCases++
		} else {
			report.FailedTestCases++
		}
	}
	
	report.EndTime = time.Now()
	report.TotalTestCases = len(evalSet.TestCases)
	
	// Calculate average score
	totalScore := 0.0
	for _, result := range report.Results {
		totalScore += result.Score
	}
	if len(report.Results) > 0 {
		report.AverageScore = totalScore / float64(len(report.Results))
	}
	
	return report, nil
}

// evaluateTestCase evaluates a single test case
func (e *AgentEvaluator) evaluateTestCase(ctx context.Context, runner *runners.Runner, testCase *EvaluationTestCase, index int) *EvaluationResult {
	result := &EvaluationResult{
		TestCase: testCase.Name,
		Input:    testCase.Input,
		Expected: testCase.Expected,
		Metadata: make(map[string]interface{}),
	}
	
	startTime := time.Now()
	
	// Create a timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, e.Config.Timeout)
	defer cancel()
	
	// Generate unique session ID for this test case
	sessionID := fmt.Sprintf("eval_session_%d_%d", time.Now().Unix(), index)
	
	// Run the agent
	eventChan, err := runner.RunAsync(timeoutCtx, "eval_user", sessionID, testCase.Input)
	if err != nil {
		result.Error = err.Error()
		result.Score = 0.0
		result.Passed = false
		result.ExecutionTime = time.Since(startTime)
		return result
	}
	
	// Collect the agent's response
	var finalEvent *events.Event
	for event := range eventChan {
		if event.IsFinalResponse {
			finalEvent = event
		}
	}
	
	result.ExecutionTime = time.Since(startTime)
	
	if finalEvent == nil || finalEvent.Content == nil {
		result.Error = "No final response from agent"
		result.Score = 0.0
		result.Passed = false
		return result
	}
	
	result.Actual = finalEvent.Content
	
	// Evaluate the response
	score, passed := e.evaluateResponse(testCase, finalEvent.Content)
	result.Score = score
	result.Passed = passed
	
	return result
}

// evaluateResponse evaluates the agent's response against the expected response
func (e *AgentEvaluator) evaluateResponse(testCase *EvaluationTestCase, actual *events.Content) (float64, bool) {
	// If no expected response, consider it a pass with full score
	if testCase.Expected == nil {
		return 1.0, true
	}
	
	// Simple text comparison (could be enhanced with more sophisticated metrics)
	if len(testCase.Expected.Parts) > 0 && len(actual.Parts) > 0 {
		expectedText := testCase.Expected.Parts[0].Text
		actualText := actual.Parts[0].Text
		
		// Simple substring matching
		if strings.Contains(strings.ToLower(actualText), strings.ToLower(expectedText)) {
			return 1.0, true
		}
		
		// Partial credit based on similarity (simplified)
		similarity := e.calculateSimilarity(expectedText, actualText)
		return similarity, similarity >= 0.7
	}
	
	return 0.0, false
}

// calculateSimilarity calculates a simple similarity score between two strings
func (e *AgentEvaluator) calculateSimilarity(text1, text2 string) float64 {
	// Simple Jaccard similarity based on words
	words1 := strings.Fields(strings.ToLower(text1))
	words2 := strings.Fields(strings.ToLower(text2))
	
	if len(words1) == 0 && len(words2) == 0 {
		return 1.0
	}
	
	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}
	
	// Create sets
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)
	
	for _, word := range words1 {
		set1[word] = true
	}
	
	for _, word := range words2 {
		set2[word] = true
	}
	
	// Calculate intersection and union
	intersection := 0
	union := len(set1)
	
	for word := range set2 {
		if set1[word] {
			intersection++
		} else {
			union++
		}
	}
	
	if union == 0 {
		return 1.0
	}
	
	return float64(intersection) / float64(union)
}

// FindConfigForTestFile finds evaluation configuration for a test file
func (e *AgentEvaluator) FindConfigForTestFile(testFilePath string) (*EvaluationConfig, error) {
	// Look for config file in the same directory
	dir := filepath.Dir(testFilePath)
	_ = dir // TODO: Use this to load actual config file
	
	// TODO: Implement config file loading
	// For now, return default config
	return &EvaluationConfig{
		MaxConcurrency: 5,
		Timeout:        30 * time.Second,
		Metrics:        []string{"accuracy"},
	}, nil
}

// LoadEvaluationSet loads an evaluation set from a file
func LoadEvaluationSet(filePath string) (*EvaluationSet, error) {
	// TODO: Implement evaluation set loading from file
	// This would typically load from JSON, YAML, or other formats
	
	// For now, return a mock evaluation set
	return &EvaluationSet{
		Name:        "Mock Evaluation Set",
		Description: "A mock evaluation set for testing",
		TestCases: []*EvaluationTestCase{
			{
				Name: "Basic Test",
				Input: &events.Content{
					Role: "user",
					Parts: []events.Part{
						{Text: "Hello, what is 2+2?"},
					},
				},
				Expected: &events.Content{
					Role: "assistant",
					Parts: []events.Part{
						{Text: "4"},
					},
				},
			},
		},
	}, nil
}