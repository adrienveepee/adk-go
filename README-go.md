# ADK Go Server

A Go implementation of the Agent Development Kit (ADK) server, providing RESTful APIs for managing AI agents and conversations.

## Overview

This Go server provides the core functionality of the ADK system, including:

- Agent creation and management
- Session management
- Message handling
- RESTful API endpoints

## Project Structure

```
cmd/
└── adk-server/          # Main application entry point
    └── main.go
internal/
├── agents/              # Agent service and business logic
│   ├── service.go
│   └── service_test.go
├── handlers/            # HTTP handlers for API endpoints
│   ├── handlers.go
│   └── handlers_test.go
└── models/              # Data models and structures
    └── models.go
```

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/adrienveepee/adk-go.git
cd adk-go
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build ./cmd/adk-server
```

4. Run the server:
```bash
./adk-server
```

The server will start on `http://localhost:8080`.

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## API Endpoints

### Health Check

- `GET /health` - Check server health

### Agents

- `POST /api/v1/agents` - Create a new agent
- `GET /api/v1/agents/:id` - Get agent by ID
- `POST /api/v1/agents/:id/run` - Run an agent with input

### Sessions

- `POST /api/v1/sessions` - Create a new session
- `GET /api/v1/sessions/:id` - Get session by ID
- `POST /api/v1/sessions/:id/messages` - Send a message to a session

## Example Usage

### Create an Agent

```bash
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Assistant",
    "model": "gemini-2.0-flash",
    "instruction": "You are a helpful assistant",
    "description": "A general-purpose AI assistant"
  }'
```

### Run an Agent

```bash
curl -X POST http://localhost:8080/api/v1/agents/{agent-id}/run \
  -H "Content-Type: application/json" \
  -d '{
    "input": "Hello, how can you help me?"
  }'
```

## Development

### Code Style

This project follows Go best practices and conventions:

- Use `gofmt` for code formatting
- Follow effective Go guidelines
- Write tests for all public functions
- Use proper error handling

### Adding New Features

1. Create feature branch from main
2. Implement functionality with tests
3. Ensure all tests pass
4. Submit pull request

## Configuration

The server can be configured using environment variables:

- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Gin mode (release, debug, test)

## License

Copyright 2025 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.