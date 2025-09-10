# ADK Go SDK

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Go Tests](https://github.com/adrienveepee/adk-go/actions/workflows/go-tests.yml/badge.svg)](https://github.com/adrienveepee/adk-go/actions/workflows/go-tests.yml)

## Overview

The Agent Development Kit (ADK) Go SDK is a comprehensive, code-first Go toolkit for building, evaluating, and deploying sophisticated AI agents. This Go implementation provides full compatibility with the Python ADK architecture while leveraging Go's performance, concurrency, and type safety advantages.

The ADK Go SDK provides:
- **Rich Agent Types**: LLM agents, workflow agents (sequential, parallel, loop), and custom agents
- **Comprehensive Tool System**: Function tools, agent tools, example tools, and extensible tool framework
- **Advanced Features**: Planning, code execution, memory management, and evaluation
- **Production Ready**: Session management, artifact storage, and robust error handling
- **High Performance**: Built with Go's concurrency and performance characteristics

## Key Features

### Core Agent Types

#### LLM Agents
Powered by Large Language Models with advanced capabilities:

```go
agent := agents.NewAgent(
    "assistant",
    "gemini-2.0-flash",
    "You are a helpful assistant that can search and analyze information.",
).SetDescription("An intelligent assistant").
  SetTools([]tools.Tool{searchTool}).
  SetOutputKey("response")
```

#### Workflow Agents
Orchestrate complex processes with deterministic patterns:

```go
// Sequential execution
sequential := agents.NewSequentialAgent("workflow", []agents.Agent{
    dataAgent, analysisAgent, reportAgent,
})

// Parallel execution
parallel := agents.NewParallelAgent("parallel_tasks", []agents.Agent{
    taskA, taskB, taskC,
})

// Loop execution
loop := agents.NewLoopAgent("iterative_process", []agents.Agent{
    processAgent,
}, 5) // Max 5 iterations
```

#### Custom Agents
Implement unique logic by extending `BaseAgent`:

```go
type CustomAgent struct {
    *agents.BaseAgent
    // Custom fields
}

func (a *CustomAgent) RunAsync(ctx context.Context, invocationCtx *agents.InvocationContext) (<-chan *events.Event, error) {
    // Custom implementation
}
```

### Tool System

#### Function Tools
Wrap Go functions as agent tools:

```go
func getWeather(location string) string {
    // Implementation
    return "Sunny, 25°C"
}

tool, _ := tools.NewFunctionTool(getWeather)
tool.Name = "get_weather"
tool.Description = "Get current weather for a location"
```

#### Agent Tools
Enable agent-to-agent delegation:

```go
expertAgent := agents.NewAgent("expert", "gemini-2.0-flash", "I am a domain expert")
agentTool := tools.NewAgentTool(expertAgent)
```

### Advanced Features

#### Session Management
Persistent conversation and state management:

```go
sessionService := sessions.NewInMemorySessionService()
session, _ := sessionService.CreateSession("app", "user123", "session456", nil)
```

#### Memory Services
Long-term memory and retrieval:

```go
memoryService := memory.NewInMemoryMemoryService()
// Or use Vertex AI RAG
// memoryService := memory.NewVertexAiRagMemoryService()
```

#### Code Execution
Safe code execution in multiple languages:

```go
executor := code_executors.NewUnsafeLocalCodeExecutor()
result, _ := executor.ExecuteCode(ctx, "print('Hello, World!')", "python", execCtx)

// Or use containerized execution
containerExec := code_executors.NewContainerCodeExecutor("python:3.9")
```

#### Planning System
Multi-step reasoning and planning:

```go
planner := planners.NewBuiltInPlanner()
// Or use ReAct planner
// planner := planners.NewPlanReActPlanner(5)
```

#### Evaluation Framework
Comprehensive agent testing and evaluation:

```go
evaluator := evaluation.NewAgentEvaluator(nil)
evalSet := &evaluation.EvaluationSet{
    Name: "Test Suite",
    TestCases: []*evaluation.EvaluationTestCase{
        {
            Name: "Basic Test",
            Input: &events.Content{
                Role: "user",
                Parts: []events.Part{{Text: "What is 2+2?"}},
            },
            Expected: &events.Content{
                Role: "assistant",
                Parts: []events.Part{{Text: "4"}},
            },
        },
    },
}

report, _ := evaluator.Evaluate(ctx, agent, evalSet)
```

## Getting Started

### Installation

```bash
go mod init your-project
go get github.com/adrienveepee/adk-go
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/adrienveepee/adk-go/google/adk/agents"
    "github.com/adrienveepee/adk-go/google/adk/events"
    "github.com/adrienveepee/adk-go/google/adk/runners"
    "github.com/adrienveepee/adk-go/google/adk/sessions"
)

func main() {
    ctx := context.Background()
    
    // Create an agent
    agent := agents.NewAgent(
        "assistant",
        "gemini-2.0-flash",
        "You are a helpful assistant.",
    )
    
    // Create session service and runner
    sessionService := sessions.NewInMemorySessionService()
    runner := runners.NewRunner(agent, "my_app", sessionService)
    
    // Create user message
    message := &events.Content{
        Role: "user",
        Parts: []events.Part{{Text: "Hello, how are you?"}},
    }
    
    // Run the agent
    eventChan, err := runner.RunAsync(ctx, "user123", "session456", message)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process events
    for event := range eventChan {
        if event.Content != nil && event.Content.Role == "model" {
            fmt.Printf("Assistant: %s\n", event.Content.Parts[0].Text)
        }
    }
}
```

### Advanced Example

```go
// Create a complex multi-agent system
func createAdvancedSystem() agents.Agent {
    // Create specialized agents
    researcher := agents.NewAgent(
        "researcher",
        "gemini-2.0-flash",
        "You research topics thoroughly and provide detailed information.",
    ).SetTools([]tools.Tool{searchTool, webScrapeTool})
    
    analyzer := agents.NewAgent(
        "analyzer",
        "gemini-2.0-flash", 
        "You analyze data and provide insights.",
    ).SetTools([]tools.Tool{analysisTool})
    
    writer := agents.NewAgent(
        "writer",
        "gemini-2.0-flash",
        "You write clear, comprehensive reports.",
    ).SetOutputKey("final_report")
    
    // Create workflow
    workflow := agents.NewSequentialAgent(
        "research_workflow",
        []agents.Agent{researcher, analyzer, writer},
    )
    
    // Create coordinator agent
    coordinator := agents.NewAgent(
        "coordinator",
        "gemini-2.0-flash",
        "You coordinate complex research tasks.",
    ).SetDescription("Research coordination system")
    
    coordinator.AddSubAgent(workflow)
    
    return coordinator
}
```

## Architecture

The ADK Go SDK follows a modular architecture with clear separation of concerns:

```
google/adk/
├── agents/           # Agent types and interfaces
├── tools/            # Tool system and implementations
├── models/           # LLM model interfaces and implementations
├── sessions/         # Session and state management
├── runners/          # Agent execution orchestration
├── events/           # Event system for agent communication
├── memory/           # Memory services and storage
├── artifacts/        # Artifact management
├── evaluation/       # Agent evaluation framework
├── examples/         # Example providers
├── planners/         # Planning and reasoning systems
└── code_executors/   # Code execution environments
```

## Model Support

The SDK currently supports:
- **Gemini 2.0 Flash** - Latest high-performance model
- **Gemini 1.5 Pro** - Advanced reasoning and long context
- **Gemini 1.5 Flash** - Fast and efficient
- **Gemini 1.0 Pro** - Stable baseline model

Additional model providers can be easily added through the LLM interface.

## Development

### Running Tests

```bash
go test ./google/adk/...
```

### Building Examples

```bash
go build ./examples/main.go
./main
```

### Contributing

This Go SDK aims for complete feature parity with the Python ADK implementation. Contributions are welcome to:

- Add missing features from the Python version
- Implement additional model providers
- Enhance tool integrations
- Improve performance and reliability
- Add comprehensive tests

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Links

- [Python ADK](https://github.com/google/adk-python) - Original Python implementation
- [ADK Documentation](https://google.github.io/adk-docs) - Comprehensive documentation
- [ADK Samples](https://github.com/google/adk-samples) - Example projects

---

*Building the future of AI agents with Go*
