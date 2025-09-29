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

package code_executors

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// CodeExecutorContext provides context for code execution
type CodeExecutorContext struct {
	mu                   sync.RWMutex
	executionID          string
	inputFiles           []string
	processedFileNames   []string
	errorCount           int
	stateDelta           map[string]interface{}
	codeExecutionResult  string
}

// NewCodeExecutorContext creates a new code executor context
func NewCodeExecutorContext() *CodeExecutorContext {
	return &CodeExecutorContext{
		inputFiles:         make([]string, 0),
		processedFileNames: make([]string, 0),
		stateDelta:         make(map[string]interface{}),
	}
}

// SetExecutionID sets the execution ID
func (c *CodeExecutorContext) SetExecutionID(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.executionID = id
}

// GetExecutionID gets the execution ID
func (c *CodeExecutorContext) GetExecutionID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.executionID
}

// AddInputFiles adds input files
func (c *CodeExecutorContext) AddInputFiles(files []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.inputFiles = append(c.inputFiles, files...)
}

// GetInputFiles gets input files
func (c *CodeExecutorContext) GetInputFiles() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]string, len(c.inputFiles))
	copy(result, c.inputFiles)
	return result
}

// ClearInputFiles clears input files
func (c *CodeExecutorContext) ClearInputFiles() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.inputFiles = c.inputFiles[:0]
}

// AddProcessedFileNames adds processed file names
func (c *CodeExecutorContext) AddProcessedFileNames(names []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.processedFileNames = append(c.processedFileNames, names...)
}

// GetProcessedFileNames gets processed file names
func (c *CodeExecutorContext) GetProcessedFileNames() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]string, len(c.processedFileNames))
	copy(result, c.processedFileNames)
	return result
}

// IncrementErrorCount increments the error count
func (c *CodeExecutorContext) IncrementErrorCount() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errorCount++
}

// GetErrorCount gets the error count
func (c *CodeExecutorContext) GetErrorCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.errorCount
}

// ResetErrorCount resets the error count
func (c *CodeExecutorContext) ResetErrorCount() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errorCount = 0
}

// GetStateDelta gets the state delta
func (c *CodeExecutorContext) GetStateDelta() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[string]interface{})
	for k, v := range c.stateDelta {
		result[k] = v
	}
	return result
}

// UpdateCodeExecutionResult updates the code execution result
func (c *CodeExecutorContext) UpdateCodeExecutionResult(result string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.codeExecutionResult = result
}

// CodeExecutor is the interface that all code executors must implement
type CodeExecutor interface {
	// ExecuteCode executes the given code and returns the result
	ExecuteCode(ctx context.Context, code, language string, execCtx *CodeExecutorContext) (string, error)
}

// BaseCodeExecutor provides the base implementation for all code executors
type BaseCodeExecutor struct {
	CodeBlockDelimiters       []string `json:"code_block_delimiters"`
	ExecutionResultDelimiters []string `json:"execution_result_delimiters"`
	ErrorRetryAttempts        int      `json:"error_retry_attempts"`
	OptimizeDataFile          bool     `json:"optimize_data_file"`
	Stateful                  bool     `json:"stateful"`
}

// NewBaseCodeExecutor creates a new base code executor
func NewBaseCodeExecutor() *BaseCodeExecutor {
	return &BaseCodeExecutor{
		CodeBlockDelimiters:       []string{"```", "```"},
		ExecutionResultDelimiters: []string{"<execution_result>", "</execution_result>"},
		ErrorRetryAttempts:        3,
		OptimizeDataFile:          true,
		Stateful:                  true,
	}
}

// ExecuteCode is the base implementation - to be overridden by concrete executors
func (e *BaseCodeExecutor) ExecuteCode(ctx context.Context, code, language string, execCtx *CodeExecutorContext) (string, error) {
	return "", fmt.Errorf("ExecuteCode not implemented for base executor")
}

// UnsafeLocalCodeExecutor executes code locally (unsafe for production)
type UnsafeLocalCodeExecutor struct {
	*BaseCodeExecutor
}

// NewUnsafeLocalCodeExecutor creates a new unsafe local code executor
func NewUnsafeLocalCodeExecutor() *UnsafeLocalCodeExecutor {
	return &UnsafeLocalCodeExecutor{
		BaseCodeExecutor: NewBaseCodeExecutor(),
	}
}

// ExecuteCode executes code locally
func (e *UnsafeLocalCodeExecutor) ExecuteCode(ctx context.Context, code, language string, execCtx *CodeExecutorContext) (string, error) {
	switch strings.ToLower(language) {
	case "python", "python3":
		return e.executePython(ctx, code, execCtx)
	case "bash", "shell", "sh":
		return e.executeBash(ctx, code, execCtx)
	case "javascript", "js", "node":
		return e.executeJavaScript(ctx, code, execCtx)
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
}

// executePython executes Python code
func (e *UnsafeLocalCodeExecutor) executePython(ctx context.Context, code string, execCtx *CodeExecutorContext) (string, error) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "adk_python_*.py")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Write code to file
	if _, err := tempFile.WriteString(code); err != nil {
		return "", fmt.Errorf("failed to write code to temp file: %w", err)
	}
	tempFile.Close()
	
	// Execute Python
	cmd := exec.CommandContext(ctx, "python3", tempFile.Name())
	output, err := cmd.CombinedOutput()
	
	result := string(output)
	if err != nil {
		execCtx.IncrementErrorCount()
		return result, fmt.Errorf("python execution failed: %w", err)
	}
	
	execCtx.UpdateCodeExecutionResult(result)
	return result, nil
}

// executeBash executes Bash code
func (e *UnsafeLocalCodeExecutor) executeBash(ctx context.Context, code string, execCtx *CodeExecutorContext) (string, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", code)
	output, err := cmd.CombinedOutput()
	
	result := string(output)
	if err != nil {
		execCtx.IncrementErrorCount()
		return result, fmt.Errorf("bash execution failed: %w", err)
	}
	
	execCtx.UpdateCodeExecutionResult(result)
	return result, nil
}

// executeJavaScript executes JavaScript code using Node.js
func (e *UnsafeLocalCodeExecutor) executeJavaScript(ctx context.Context, code string, execCtx *CodeExecutorContext) (string, error) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "adk_js_*.js")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	
	// Write code to file
	if _, err := tempFile.WriteString(code); err != nil {
		return "", fmt.Errorf("failed to write code to temp file: %w", err)
	}
	tempFile.Close()
	
	// Execute Node.js
	cmd := exec.CommandContext(ctx, "node", tempFile.Name())
	output, err := cmd.CombinedOutput()
	
	result := string(output)
	if err != nil {
		execCtx.IncrementErrorCount()
		return result, fmt.Errorf("javascript execution failed: %w", err)
	}
	
	execCtx.UpdateCodeExecutionResult(result)
	return result, nil
}

// ContainerCodeExecutor executes code in a container
type ContainerCodeExecutor struct {
	*BaseCodeExecutor
	Image      string `json:"image"`
	DockerPath string `json:"docker_path"`
	BaseURL    string `json:"base_url"`
}

// NewContainerCodeExecutor creates a new container code executor
func NewContainerCodeExecutor(image string) *ContainerCodeExecutor {
	executor := &ContainerCodeExecutor{
		BaseCodeExecutor: NewBaseCodeExecutor(),
		Image:            image,
		DockerPath:       "docker",
	}
	return executor
}

// ExecuteCode executes code in a container
func (e *ContainerCodeExecutor) ExecuteCode(ctx context.Context, code, language string, execCtx *CodeExecutorContext) (string, error) {
	// Create temporary directory for code execution
	tempDir, err := os.MkdirTemp("", "adk_container_exec_")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create code file
	var fileName string
	switch strings.ToLower(language) {
	case "python", "python3":
		fileName = "code.py"
	case "javascript", "js", "node":
		fileName = "code.js"
	case "bash", "shell", "sh":
		fileName = "code.sh"
	default:
		fileName = "code.txt"
	}
	
	codeFile := filepath.Join(tempDir, fileName)
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		return "", fmt.Errorf("failed to write code file: %w", err)
	}
	
	// Execute in container
	var cmd *exec.Cmd
	switch strings.ToLower(language) {
	case "python", "python3":
		cmd = exec.CommandContext(ctx, e.DockerPath, "run", "--rm", "-v", tempDir+":/workspace", e.Image, "python3", "/workspace/"+fileName)
	case "javascript", "js", "node":
		cmd = exec.CommandContext(ctx, e.DockerPath, "run", "--rm", "-v", tempDir+":/workspace", e.Image, "node", "/workspace/"+fileName)
	case "bash", "shell", "sh":
		cmd = exec.CommandContext(ctx, e.DockerPath, "run", "--rm", "-v", tempDir+":/workspace", e.Image, "bash", "/workspace/"+fileName)
	default:
		return "", fmt.Errorf("unsupported language for container execution: %s", language)
	}
	
	output, err := cmd.CombinedOutput()
	result := string(output)
	
	if err != nil {
		execCtx.IncrementErrorCount()
		return result, fmt.Errorf("container execution failed: %w", err)
	}
	
	execCtx.UpdateCodeExecutionResult(result)
	return result, nil
}

// VertexAiCodeExecutor executes code using Vertex AI
type VertexAiCodeExecutor struct {
	*BaseCodeExecutor
	ResourceName string `json:"resource_name"`
}

// NewVertexAiCodeExecutor creates a new Vertex AI code executor
func NewVertexAiCodeExecutor(resourceName string) *VertexAiCodeExecutor {
	return &VertexAiCodeExecutor{
		BaseCodeExecutor: NewBaseCodeExecutor(),
		ResourceName:     resourceName,
	}
}

// ExecuteCode executes code using Vertex AI
func (e *VertexAiCodeExecutor) ExecuteCode(ctx context.Context, code, language string, execCtx *CodeExecutorContext) (string, error) {
	// TODO: Implement Vertex AI code execution
	// This would involve calling the Vertex AI API to execute code in a managed environment
	
	// For now, return a mock result
	result := fmt.Sprintf("Mock Vertex AI execution result for %s code:\n%s", language, code)
	execCtx.UpdateCodeExecutionResult(result)
	
	return result, nil
}