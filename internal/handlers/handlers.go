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
	"context"
	"net/http"

	"github.com/adrienveepee/adk-go/internal/agents"
	"github.com/adrienveepee/adk-go/internal/models"
	"github.com/gin-gonic/gin"
)

var agentService = agents.NewService()

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HealthHandler handles health check requests
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "healthy",
		Message: "ADK Go server is running",
	})
}

// CreateAgent handles agent creation requests
func CreateAgent(c *gin.Context) {
	var config models.AgentConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent, err := agentService.CreateAgent(context.Background(), config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// GetAgent handles agent retrieval requests
func GetAgent(c *gin.Context) {
	id := c.Param("id")
	
	agent, err := agentService.GetAgent(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// RunAgent handles agent execution requests
func RunAgent(c *gin.Context) {
	id := c.Param("id")
	
	var request struct {
		Input string `json:"input" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := agentService.RunAgent(context.Background(), id, request.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

// CreateSession handles session creation requests
func CreateSession(c *gin.Context) {
	// TODO: Implement session creation logic
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Session creation not implemented yet",
	})
}

// GetSession handles session retrieval requests
func GetSession(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement session retrieval logic
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Session retrieval not implemented yet",
		"id":    id,
	})
}

// SendMessage handles message sending requests
func SendMessage(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement message sending logic
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Message sending not implemented yet",
		"id":    id,
	})
}