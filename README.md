# Agent Development Kit (ADK) - Go SDK

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Go Tests](https://github.com/adrienveepee/adk-go/actions/workflows/go-tests.yml/badge.svg)](https://github.com/adrienveepee/adk-go/actions/workflows/go-tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/adrienveepee/adk-go)](https://goreportcard.com/report/github.com/adrienveepee/adk-go)

<html>
    <h2 align="center">
      <img src="https://raw.githubusercontent.com/google/adk-python/main/assets/agent-development-kit.png" width="256"/>
    </h2>
    <h3 align="center">
      An open-source, code-first Go toolkit for building, evaluating, and deploying sophisticated AI agents with flexibility and control.
    </h3>
    <h3 align="center">
      Important Links:
      <a href="https://google.github.io/adk-docs/">Docs</a>,
      <a href="https://github.com/google/adk-samples">Samples</a>,
      <a href="https://github.com/google/adk-python">Python ADK</a> &
      <a href="https://github.com/google/adk-web">ADK Web</a>.
    </h3>
</html>

Agent Development Kit (ADK) Go SDK is a flexible and modular framework for developing and deploying AI agents in Go. While optimized for Gemini and the Google ecosystem, ADK is model-agnostic, deployment-agnostic, and is built for compatibility with other frameworks. ADK was designed to make agent development feel more like software development, leveraging Go's performance, type safety, and concurrency to create, deploy, and orchestrate agentic architectures that range from simple tasks to complex workflows.

---

## ✨ What's new

- **Native Go Performance**: Built from the ground up in Go for superior performance and type safety
- **Full Feature Parity**: Complete implementation matching Python ADK capabilities
- **Concurrency First**: Leverages Go's goroutines and channels for high-performance agent orchestration

## ✨ Key Features

- **Rich Tool Ecosystem**: Utilize pre-built tools, custom functions, or integrate existing tools to give agents diverse capabilities, all with tight integration to the Google ecosystem.

- **Code-First Development**: Define agent logic, tools, and orchestration directly in Go for ultimate flexibility, testability, and versioning with compile-time safety.

- **Modular Multi-Agent Systems**: Design scalable applications by composing multiple specialized agents into flexible hierarchies.

- **Deploy Anywhere**: Easily containerize and deploy agents on Cloud Run or scale seamlessly with container orchestration.

- **High Performance**: Built with Go's native concurrency, goroutines, and channels for superior agent orchestration performance.

## 🚀 Installation

### Using Go Modules

You can install the ADK Go SDK using Go modules:

```bash
go mod init your-project
go get github.com/adrienveepee/adk-go
```

This version provides the latest stable release with full feature parity to the Python ADK.

### Requirements

- Go 1.21 or later
- Google Cloud Project (for Gemini models)
- API credentials configured

## 🏁 Feature Highlight

### Define a single agent:

```go
import (
    "github.com/adrienveepee/adk-go/google/adk/agents"
    "github.com/adrienveepee/adk-go/google/adk/tools"
)

rootAgent := agents.NewAgent(
    "search_assistant",
    "gemini-2.0-flash", // Or your preferred Gemini model
    "You are a helpful assistant. Answer user questions using search when needed.",
).SetDescription("An assistant that can search the web").
  SetTools([]tools.Tool{googleSearchTool})
```

### Define a multi-agent system:

Define a multi-agent system with coordinator agent, greeter agent, and task execution agent. Then ADK engine and the model will guide the agents to work together to accomplish the task.

```go
// Define individual agents
greeter := agents.NewAgent("greeter", "gemini-2.0-flash", 
    "You greet users warmly and help them feel welcome.")

taskExecutor := agents.NewAgent("task_executor", "gemini-2.0-flash",
    "You execute specific tasks efficiently and thoroughly.")

// Create parent agent and assign children via sub_agents
coordinator := agents.NewAgent(
    "Coordinator",
    "gemini-2.0-flash",
    "I coordinate greetings and tasks for optimal user experience.",
).SetDescription("I coordinate greetings and tasks.").
  AddSubAgent(greeter).
  AddSubAgent(taskExecutor)
```

## 🛒 Real-World Example: E-commerce Catalog Sale Preparation

Here's a comprehensive example of using ADK Go SDK for an e-commerce company preparing a catalog sale:

### E-commerce Sale Preparation Workflow

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/adrienveepee/adk-go/google/adk/agents"
    "github.com/adrienveepee/adk-go/google/adk/tools"
    "github.com/adrienveepee/adk-go/google/adk/runners"
    "github.com/adrienveepee/adk-go/google/adk/sessions"
    "github.com/adrienveepee/adk-go/google/adk/events"
)

func main() {
    ctx := context.Background()
    
    // 1. Product Analysis Agent - Analyzes product performance and inventory
    productAnalyzer := agents.NewAgent(
        "product_analyzer",
        "gemini-2.0-flash",
        `You are a product analyst. Analyze product data including:
        - Sales performance over the last 6 months
        - Current inventory levels
        - Seasonal trends
        - Customer reviews and ratings
        Provide recommendations for sale inclusion and pricing.`,
    ).SetTools([]tools.Tool{
        createInventoryTool(),
        createSalesDataTool(),
        createReviewAnalysisTool(),
    })
    
    // 2. Market Research Agent - Analyzes competitor pricing and market trends
    marketResearcher := agents.NewAgent(
        "market_researcher", 
        "gemini-2.0-flash",
        `You are a market research specialist. Research:
        - Competitor pricing for similar products
        - Market demand trends
        - Price elasticity analysis
        - Customer sentiment in the category
        Provide competitive pricing recommendations.`,
    ).SetTools([]tools.Tool{
        createCompetitorAnalysisTool(),
        createMarketTrendsTool(),
        createPriceComparisonTool(),
    })
    
    // 3. Content Creator Agent - Creates marketing content
    contentCreator := agents.NewAgent(
        "content_creator",
        "gemini-2.0-flash",
        `You are a creative marketing specialist. Create engaging content:
        - Product descriptions highlighting sale benefits
        - Email marketing copy
        - Social media posts
        - Banner ads and promotional materials
        Ensure all content is persuasive and brand-consistent.`,
    ).SetTools([]tools.Tool{
        createContentGenerationTool(),
        createImageGenerationTool(),
        createBrandGuidelinesTool(),
    })
    
    // 4. Campaign Coordinator - Orchestrates the entire workflow
    campaignCoordinator := agents.NewAgent(
        "campaign_coordinator",
        "gemini-2.0-flash",
        `You are a sale campaign coordinator. Orchestrate the entire sale preparation:
        1. Coordinate product analysis and market research
        2. Make final pricing and product selection decisions
        3. Approve marketing content and campaigns
        4. Set campaign timeline and launch strategy
        5. Monitor and adjust strategy based on performance data`,
    ).SetDescription("E-commerce Sale Campaign Coordinator").
      AddSubAgent(productAnalyzer).
      AddSubAgent(marketResearcher).
      AddSubAgent(contentCreator).
      SetTools([]tools.Tool{
        createCampaignPlanningTool(),
        createApprovalWorkflowTool(),
        createLaunchScheduleTool(),
    })
    
    // Create workflow orchestration
    salePreparationWorkflow := agents.NewSequentialAgent(
        "sale_preparation_workflow",
        []agents.Agent{
            productAnalyzer,
            marketResearcher, 
            contentCreator,
            campaignCoordinator,
        },
    )
    
    // Setup execution environment
    sessionService := sessions.NewInMemorySessionService()
    runner := runners.NewRunner(salePreparationWorkflow, "ecommerce_sale_app", sessionService)
    
    // Execute the workflow
    campaignBrief := &events.Content{
        Role: "user",
        Parts: []events.Part{{
            Text: `Prepare a summer sale campaign for our outdoor gear category. 
                   Target: 25% revenue increase, 40% inventory clearance.
                   Timeline: 2 weeks preparation, 3 weeks campaign duration.
                   Focus categories: camping equipment, hiking gear, water sports.`,
        }},
    }
    
    fmt.Println("🏪 Starting E-commerce Sale Preparation Workflow...")
    
    eventChan, err := runner.RunAsync(ctx, "campaign_manager", "summer_sale_2024", campaignBrief)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process workflow results
    for event := range eventChan {
        if event.Content != nil {
            switch event.Content.Role {
            case "product_analyzer":
                fmt.Printf("📊 Product Analysis: %s\n", event.Content.Parts[0].Text)
            case "market_researcher":
                fmt.Printf("📈 Market Research: %s\n", event.Content.Parts[0].Text)
            case "content_creator":
                fmt.Printf("✨ Content Created: %s\n", event.Content.Parts[0].Text)
            case "campaign_coordinator":
                fmt.Printf("🎯 Campaign Plan: %s\n", event.Content.Parts[0].Text)
            }
        }
        
        if event.ToolCalls != nil {
            fmt.Printf("🔧 Tool executed: %s\n", event.ToolCalls[0].Function.Name)
        }
    }
    
    fmt.Println("✅ Sale preparation workflow completed!")
}

// Helper functions to create domain-specific tools
func createInventoryTool() tools.Tool {
    inventoryCheck := func(category string, timeframe string) map[string]interface{} {
        // In real implementation, this would connect to your inventory system
        return map[string]interface{}{
            "category": category,
            "total_items": 1250,
            "low_stock_items": 45,
            "overstock_items": 23,
            "turnover_rate": "3.2x annually",
        }
    }
    
    tool, _ := tools.NewFunctionTool(inventoryCheck)
    tool.Name = "check_inventory"
    tool.Description = "Check current inventory levels and turnover rates for product categories"
    return tool
}

func createSalesDataTool() tools.Tool {
    salesAnalysis := func(category string, months int) map[string]interface{} {
        // Connect to your analytics platform
        return map[string]interface{}{
            "category": category,
            "revenue_trend": "increasing 15%",
            "top_performers": []string{"Ultra-light Tent", "Hiking Boots Pro", "Water Filter Plus"},
            "underperformers": []string{"Basic Compass", "Rain Poncho Standard"},
            "seasonal_peak": "June-August",
        }
    }
    
    tool, _ := tools.NewFunctionTool(salesAnalysis)
    tool.Name = "analyze_sales_data"
    tool.Description = "Analyze sales performance data for specific categories and timeframes"
    return tool
}

func createCompetitorAnalysisTool() tools.Tool {
    competitorPricing := func(product string) map[string]interface{} {
        // Connect to price monitoring services
        return map[string]interface{}{
            "product": product,
            "avg_competitor_price": "$89.99",
            "price_range": "$69.99 - $129.99",
            "our_current_price": "$94.99",
            "recommended_sale_price": "$74.99",
            "competitive_advantage": "superior quality, better warranty",
        }
    }
    
    tool, _ := tools.NewFunctionTool(competitorPricing)
    tool.Name = "analyze_competitor_pricing"
    tool.Description = "Analyze competitor pricing and market positioning for products"
    return tool
}

func createContentGenerationTool() tools.Tool {
    generateContent := func(contentType string, product string, saleDetails string) string {
        // Integrate with content generation tools or templates
        return fmt.Sprintf("Generated %s content for %s with sale details: %s", 
            contentType, product, saleDetails)
    }
    
    tool, _ := tools.NewFunctionTool(generateContent)
    tool.Name = "generate_marketing_content"
    tool.Description = "Generate marketing content for products and sales campaigns"
    return tool
}

// Additional tool creation functions...
func createReviewAnalysisTool() tools.Tool { /* implementation */ }
func createMarketTrendsTool() tools.Tool { /* implementation */ }
func createPriceComparisonTool() tools.Tool { /* implementation */ }
func createImageGenerationTool() tools.Tool { /* implementation */ }
func createBrandGuidelinesTool() tools.Tool { /* implementation */ }
func createCampaignPlanningTool() tools.Tool { /* implementation */ }
func createApprovalWorkflowTool() tools.Tool { /* implementation */ }
func createLaunchScheduleTool() tools.Tool { /* implementation */ }
```

### Advanced E-commerce Features

#### Parallel Processing for Performance

```go
// Process multiple product categories simultaneously
categoryAgents := []agents.Agent{
    createCategoryAgent("outdoor_gear"),
    createCategoryAgent("sports_equipment"), 
    createCategoryAgent("fitness_accessories"),
}

parallelProcessor := agents.NewParallelAgent("category_processor", categoryAgents)
```

#### Memory Integration for Campaign History

```go
import "github.com/adrienveepee/adk-go/google/adk/memory"

// Track previous campaign performance
memoryService := memory.NewInMemoryMemoryService()
campaignHistory, _ := memoryService.Search(ctx, "previous_summer_sales", "")

coordinator.SetMemoryService(memoryService)
```

#### Evaluation and A/B Testing

```go
import "github.com/adrienveepee/adk-go/google/adk/evaluation"

// Create evaluation set for campaign performance
evaluator := evaluation.NewAgentEvaluator(nil)
campaignTests := &evaluation.EvaluationSet{
    Name: "Campaign Performance Tests",
    TestCases: []*evaluation.EvaluationTestCase{
        {
            Name: "Revenue Target Test",
            Input: &events.Content{
                Role: "user", 
                Parts: []events.Part{{Text: "Evaluate campaign against 25% revenue increase target"}},
            },
            Expected: &events.Content{
                Role: "assistant",
                Parts: []events.Part{{Text: "Campaign projected to achieve 27% revenue increase"}},
            },
        },
    },
}

report, _ := evaluator.Evaluate(ctx, campaignCoordinator, campaignTests)
```

## 🔧 Advanced Features

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

### Workflow Agents

#### Sequential Execution
```go
sequential := agents.NewSequentialAgent("workflow", []agents.Agent{
    dataAgent, analysisAgent, reportAgent,
})
```

#### Parallel Execution  
```go
parallel := agents.NewParallelAgent("parallel_tasks", []agents.Agent{
    taskA, taskB, taskC,
})
```

#### Loop Execution
```go
loop := agents.NewLoopAgent("iterative_process", []agents.Agent{
    processAgent,
}, 5) // Max 5 iterations
```

### Session Management
Persistent conversation and state management:

```go
sessionService := sessions.NewInMemorySessionService()
session, _ := sessionService.CreateSession("app", "user123", "session456", nil)
```

### Memory Services
Long-term memory and retrieval:

```go
memoryService := memory.NewInMemoryMemoryService()
// Or use Vertex AI RAG
// memoryService := memory.NewVertexAiRagMemoryService()
```

### Code Execution
Safe code execution in multiple languages:

```go
executor := code_executors.NewUnsafeLocalCodeExecutor()
result, _ := executor.ExecuteCode(ctx, "print('Hello, World!')", "python", execCtx)

// Or use containerized execution
containerExec := code_executors.NewContainerCodeExecutor("python:3.9")
```

### Planning System
Multi-step reasoning and planning:

```go
planner := planners.NewBuiltInPlanner()
// Or use ReAct planner
// planner := planners.NewPlanReActPlanner(5)
```

### Evaluation Framework
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

## 📚 Getting Started

### Quick Start

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

### Development Setup

1. **Clone the repository:**
```bash
git clone https://github.com/adrienveepee/adk-go.git
cd adk-go
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Run tests:**
```bash
go test ./google/adk/...
```

4. **Build and run examples:**
```bash
go build ./examples/main.go
./main
```

## 🏗️ Architecture

The ADK Go SDK follows a modular architecture with clear separation of concerns:

```
google/adk/
├── agents/           # Agent types and interfaces (BaseAgent, LlmAgent, workflows)
├── tools/            # Tool system and implementations (Function, Agent, Example tools)
├── models/           # LLM model interfaces and implementations (Gemini integration)
├── sessions/         # Session and state management (InMemorySessionService)
├── runners/          # Agent execution orchestration (Runner, async execution)
├── events/           # Event system for agent communication (Content, ToolCalls)
├── memory/           # Memory services and storage (InMemory, VertexAI RAG)
├── artifacts/        # Artifact management and versioning
├── evaluation/       # Agent evaluation framework (test suites, metrics)
├── examples/         # Example providers and demonstrations
├── planners/         # Planning and reasoning systems (BuiltIn, ReAct)
└── code_executors/   # Code execution environments (Local, Container, VertexAI)
```

### Key Design Principles

- **Type Safety**: Full compile-time type checking throughout the SDK
- **Concurrency**: Built-in support for goroutines and channels for async agent execution
- **Modularity**: Clean interfaces allow easy extension and customization
- **Performance**: Optimized for high-throughput agent orchestration
- **Compatibility**: Feature parity with Python ADK while leveraging Go advantages

## 🤖 Model Support

The SDK currently supports Gemini models with extensible LLM interface:

- **Gemini 2.0 Flash** - Latest high-performance model with advanced capabilities
- **Gemini 1.5 Pro** - Advanced reasoning and long context understanding
- **Gemini 1.5 Flash** - Fast and efficient for most use cases
- **Gemini 1.0 Pro** - Stable baseline model

Additional model providers can be easily added through the `LLMConnection` interface.

## 🔄 Evaluation and Testing

### Running Tests

```bash
# Run all tests
go test ./google/adk/...

# Run tests with coverage
go test -cover ./google/adk/...

# Run specific package tests
go test ./google/adk/agents/
```

### Agent Evaluation

```bash
# Build evaluation tool
go build ./examples/evaluation/main.go

# Run evaluation suite
./main -agent ./examples/hello_world -evalset ./examples/hello_world/eval_set.json
```

## 🤝 Contributing

We welcome contributions from the community! Whether it's bug reports, feature requests, documentation improvements, or code contributions, please see our:

- [General contribution guideline and flow](https://google.github.io/adk-docs/contributing-guide/)
- [Code Contributing Guidelines](./CONTRIBUTING.md) to get started

### Development Guidelines

- Follow Go best practices and conventions
- Use `gofmt` for code formatting
- Write tests for all public functions
- Ensure proper error handling
- Add documentation for exported functions

### Adding New Features

1. Create feature branch from main
2. Implement functionality with comprehensive tests  
3. Ensure all tests pass and code is properly formatted
4. Submit pull request with clear description

## 📖 Documentation

Explore the full documentation for detailed guides on building, evaluating, and deploying agents:

* **[Documentation](https://google.github.io/adk-docs)** - Comprehensive guides and API reference
* **[Go SDK Examples](./examples/)** - Working examples and demonstrations
* **[Python ADK](https://github.com/google/adk-python)** - Original Python implementation

## 📄 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **[Python ADK](https://github.com/google/adk-python)** - Original Python implementation
- **[ADK Documentation](https://google.github.io/adk-docs)** - Comprehensive documentation  
- **[ADK Samples](https://github.com/google/adk-samples)** - Example projects and use cases
- **[ADK Web](https://github.com/google/adk-web)** - Web interface for agent development

---

*Building the future of AI agents with Go's performance and type safety*
