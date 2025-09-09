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

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adrienveepee/adk-go/internal/models"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", HealthHandler)
	router.POST("/agents", CreateAgent)
	router.GET("/agents/:id", GetAgent)
	router.POST("/agents/:id/run", RunAgent)
	return router
}

func TestHealthHandler(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
	}
}

func TestCreateAgent(t *testing.T) {
	router := setupTestRouter()

	config := models.AgentConfig{
		Name:        "Test Agent",
		Model:       "gemini-2.0-flash",
		Instruction: "You are a helpful assistant",
		Description: "A test agent",
	}

	jsonBody, _ := json.Marshal(config)
	req, _ := http.NewRequest("POST", "/agents", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var agent models.Agent
	err := json.Unmarshal(w.Body.Bytes(), &agent)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if agent.Name != config.Name {
		t.Errorf("Expected agent name %s, got %s", config.Name, agent.Name)
	}
}

func TestCreateAgentInvalidPayload(t *testing.T) {
	router := setupTestRouter()

	// Send invalid JSON
	req, _ := http.NewRequest("POST", "/agents", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}