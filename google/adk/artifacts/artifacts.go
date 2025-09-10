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

package artifacts

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ArtifactVersion represents a version of an artifact
type ArtifactVersion struct {
	Version   int       `json:"version"`
	Data      []byte    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ArtifactService interface for managing artifacts
type ArtifactService interface {
	// SaveArtifact saves an artifact with the given key and data
	SaveArtifact(ctx context.Context, key string, data []byte, metadata map[string]interface{}) error
	
	// LoadArtifact loads an artifact by key
	LoadArtifact(ctx context.Context, key string) ([]byte, error)
	
	// DeleteArtifact deletes an artifact by key
	DeleteArtifact(ctx context.Context, key string) error
	
	// ListArtifactKeys lists all artifact keys
	ListArtifactKeys(ctx context.Context) ([]string, error)
	
	// ListVersions lists all versions of an artifact
	ListVersions(ctx context.Context, key string) ([]*ArtifactVersion, error)
}

// InMemoryArtifactService provides an in-memory implementation of ArtifactService
type InMemoryArtifactService struct {
	mu        sync.RWMutex
	artifacts map[string][]*ArtifactVersion // Key: artifact key, Value: versions
}

// NewInMemoryArtifactService creates a new in-memory artifact service
func NewInMemoryArtifactService() *InMemoryArtifactService {
	return &InMemoryArtifactService{
		artifacts: make(map[string][]*ArtifactVersion),
	}
}

// SaveArtifact saves an artifact with the given key and data
func (a *InMemoryArtifactService) SaveArtifact(ctx context.Context, key string, data []byte, metadata map[string]interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	versions, exists := a.artifacts[key]
	if !exists {
		versions = make([]*ArtifactVersion, 0)
	}
	
	newVersion := &ArtifactVersion{
		Version:   len(versions) + 1,
		Data:      make([]byte, len(data)),
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
	
	copy(newVersion.Data, data)
	for k, v := range metadata {
		newVersion.Metadata[k] = v
	}
	
	versions = append(versions, newVersion)
	a.artifacts[key] = versions
	
	return nil
}

// LoadArtifact loads an artifact by key (returns the latest version)
func (a *InMemoryArtifactService) LoadArtifact(ctx context.Context, key string) ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	versions, exists := a.artifacts[key]
	if !exists || len(versions) == 0 {
		return nil, fmt.Errorf("artifact not found: %s", key)
	}
	
	// Return the latest version
	latestVersion := versions[len(versions)-1]
	result := make([]byte, len(latestVersion.Data))
	copy(result, latestVersion.Data)
	
	return result, nil
}

// DeleteArtifact deletes an artifact by key
func (a *InMemoryArtifactService) DeleteArtifact(ctx context.Context, key string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	delete(a.artifacts, key)
	return nil
}

// ListArtifactKeys lists all artifact keys
func (a *InMemoryArtifactService) ListArtifactKeys(ctx context.Context) ([]string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	keys := make([]string, 0, len(a.artifacts))
	for key := range a.artifacts {
		keys = append(keys, key)
	}
	
	return keys, nil
}

// ListVersions lists all versions of an artifact
func (a *InMemoryArtifactService) ListVersions(ctx context.Context, key string) ([]*ArtifactVersion, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	versions, exists := a.artifacts[key]
	if !exists {
		return nil, fmt.Errorf("artifact not found: %s", key)
	}
	
	// Return a copy of the versions
	result := make([]*ArtifactVersion, len(versions))
	copy(result, versions)
	
	return result, nil
}

// GcsArtifactService provides a Google Cloud Storage implementation of ArtifactService
type GcsArtifactService struct {
	// TODO: Add GCS specific fields (bucket name, client, etc.)
}

// NewGcsArtifactService creates a new GCS artifact service
func NewGcsArtifactService() *GcsArtifactService {
	return &GcsArtifactService{}
}

// SaveArtifact saves an artifact to GCS
func (g *GcsArtifactService) SaveArtifact(ctx context.Context, key string, data []byte, metadata map[string]interface{}) error {
	// TODO: Implement GCS save
	return nil
}

// LoadArtifact loads an artifact from GCS
func (g *GcsArtifactService) LoadArtifact(ctx context.Context, key string) ([]byte, error) {
	// TODO: Implement GCS load
	return nil, nil
}

// DeleteArtifact deletes an artifact from GCS
func (g *GcsArtifactService) DeleteArtifact(ctx context.Context, key string) error {
	// TODO: Implement GCS delete
	return nil
}

// ListArtifactKeys lists all artifact keys in GCS
func (g *GcsArtifactService) ListArtifactKeys(ctx context.Context) ([]string, error) {
	// TODO: Implement GCS list
	return nil, nil
}

// ListVersions lists all versions of an artifact in GCS
func (g *GcsArtifactService) ListVersions(ctx context.Context, key string) ([]*ArtifactVersion, error) {
	// TODO: Implement GCS version listing
	return nil, nil
}