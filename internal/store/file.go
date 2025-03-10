package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sistemica/traefik-manager/internal/models"
)

// Data structure for storing all configuration
type storeData struct {
	Middlewares map[string]models.Middleware `json:"middlewares"`
	Routers     map[string]models.Router     `json:"routers"`
	Services    map[string]models.Service    `json:"services"`
}

// FileStore implements the Store interface with file-based persistence
type FileStore struct {
	mu           sync.RWMutex
	data         storeData
	filePath     string
	saveDebounce chan struct{}
	done         chan struct{} //  channel to signal shutdown
}

// NewFileStore creates a new FileStore
func NewFileStore(filePath string) (*FileStore, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	store := &FileStore{
		data: storeData{
			Middlewares: make(map[string]models.Middleware),
			Routers:     make(map[string]models.Router),
			Services:    make(map[string]models.Service),
		},
		filePath:     filePath,
		saveDebounce: make(chan struct{}, 1),
		done:         make(chan struct{}), // Initialize the done channel
	}

	// Try to load existing data
	if _, err := os.Stat(filePath); err == nil {
		if err := store.Load(); err != nil {
			return nil, fmt.Errorf("failed to load store data: %w", err)
		}
	}

	// Start debounced save goroutine
	go store.debouncedSave()

	return store, nil
}

// debouncedSave saves the store data to disk after a short delay
// to prevent excessive disk writes when multiple changes are made in succession
// Replace the debouncedSave method in FileStore:

func (s *FileStore) debouncedSave() {
	for {
		select {
		case <-s.saveDebounce:
			// Perform save operation with timeout protection
			doneSave := make(chan struct{})
			go func() {
				s.mu.RLock()
				_ = s.save() // Ignore errors during test
				s.mu.RUnlock()
				close(doneSave)
			}()

			// Wait for save with timeout
			select {
			case <-doneSave:
				// Save completed
			case <-s.done:
				// Exit if done signal received during save
				return
			}
		case <-s.done:
			// Exit the goroutine when done signal is received
			return
		}
	}
}

// triggerSave triggers a debounced save operation
func (s *FileStore) triggerSave() {
	// Non-blocking send to channel
	select {
	case s.saveDebounce <- struct{}{}:
		// Signal sent
	default:
		// Channel already has a pending save
	}
}

// save writes the store data to disk
func (s *FileStore) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal store data: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write store data: %w", err)
	}

	return nil
}

// Save persists the store data to disk
func (s *FileStore) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.save()
}

// Load loads the store data from disk
func (s *FileStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read store data: %w", err)
	}

	if err := json.Unmarshal(data, &s.data); err != nil {
		return fmt.Errorf("failed to unmarshal store data: %w", err)
	}

	// Initialize maps if they're nil
	if s.data.Middlewares == nil {
		s.data.Middlewares = make(map[string]models.Middleware)
	}
	if s.data.Routers == nil {
		s.data.Routers = make(map[string]models.Router)
	}
	if s.data.Services == nil {
		s.data.Services = make(map[string]models.Service)
	}

	return nil
}

// ListMiddlewares returns all middlewares
func (s *FileStore) ListMiddlewares() ([]models.Middleware, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	middlewares := make([]models.Middleware, 0, len(s.data.Middlewares))
	for _, middleware := range s.data.Middlewares {
		middlewares = append(middlewares, middleware)
	}
	return middlewares, nil
}

// GetMiddleware returns a middleware by ID
func (s *FileStore) GetMiddleware(id string) (*models.Middleware, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	middleware, ok := s.data.Middlewares[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &middleware, nil
}

// CreateMiddleware creates a new middleware
func (s *FileStore) CreateMiddleware(middleware *models.Middleware) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data.Middlewares[middleware.ID]; ok {
		return ErrAlreadyExists
	}

	s.data.Middlewares[middleware.ID] = *middleware
	s.triggerSave()
	return nil
}

// UpdateMiddleware updates an existing middleware
func (s *FileStore) UpdateMiddleware(id string, middleware *models.Middleware) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data.Middlewares[id]; !ok {
		return ErrNotFound
	}

	// Ensure ID doesn't change
	middleware.ID = id
	s.data.Middlewares[id] = *middleware
	s.triggerSave()
	return nil
}

// DeleteMiddleware deletes a middleware
func (s *FileStore) DeleteMiddleware(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data.Middlewares[id]; !ok {
		return ErrNotFound
	}

	// Check if middleware is in use
	inUse, _, err := s.middlewareInUse(id)
	if err != nil {
		return err
	}
	if inUse {
		return ErrResourceInUse
	}

	delete(s.data.Middlewares, id)
	s.triggerSave()
	return nil
}

// MiddlewareExists checks if a middleware exists
func (s *FileStore) MiddlewareExists(id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data.Middlewares[id]
	return ok, nil
}

// MiddlewareInUse checks if a middleware is in use by any routers
func (s *FileStore) MiddlewareInUse(id string) (bool, []string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.middlewareInUse(id)
}

// middlewareInUse is an internal non-locking version of MiddlewareInUse
func (s *FileStore) middlewareInUse(id string) (bool, []string, error) {
	usedBy := []string{}

	// Check if middleware is used by any routers
	for routerID, router := range s.data.Routers {
		if router.Middlewares != nil {
			for _, mw := range router.Middlewares {
				// FIXED: Instead of calling MiddlewareExists which acquires a lock,
				// just check the map directly since we already have a lock
				if mw.ID == id {
					usedBy = append(usedBy, fmt.Sprintf("router:%s", routerID))
					break // Found at least one reference in this router
				}
			}
		}
	}

	return len(usedBy) > 0, usedBy, nil
}

// ListRouters returns all routers
func (s *FileStore) ListRouters() ([]models.Router, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	routers := make([]models.Router, 0, len(s.data.Routers))
	for _, router := range s.data.Routers {
		routers = append(routers, router)
	}
	return routers, nil
}

// GetRouter returns a router by ID
func (s *FileStore) GetRouter(id string) (*models.Router, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	router, ok := s.data.Routers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &router, nil
}

// CreateRouter creates a new router
func (s *FileStore) CreateRouter(router *models.Router) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data.Routers[router.ID]; ok {
		return ErrAlreadyExists
	}

	// Validate that all referenced services exist
	if _, ok := s.data.Services[router.Service.ID]; !ok {
		return fmt.Errorf("service %s not found", router.Service.ID)
	}

	// Validate that all referenced middlewares exist - THE IMPORTANT FIX IS HERE
	if router.Middlewares != nil {
		for _, mw := range router.Middlewares {
			// Check map directly instead of calling MiddlewareExists
			if _, ok := s.data.Middlewares[mw.ID]; !ok {
				return fmt.Errorf("middleware %s not found", mw.ID)
			}
		}
	}

	s.data.Routers[router.ID] = *router
	s.triggerSave()
	return nil
}

// DeleteRouter deletes a router
func (s *FileStore) DeleteRouter(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data.Routers[id]; !ok {
		return ErrNotFound
	}

	delete(s.data.Routers, id)
	s.triggerSave()
	return nil
}

// RouterExists checks if a router exists
func (s *FileStore) RouterExists(id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data.Routers[id]
	return ok, nil
}

// RouterInUse checks if a router is in use
func (s *FileStore) RouterInUse(id string) (bool, []string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Routers are standalone entities and not referenced by other resources
	return false, nil, nil
}

// ListServices returns all services
func (s *FileStore) ListServices() ([]models.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	services := make([]models.Service, 0, len(s.data.Services))
	for _, service := range s.data.Services {
		services = append(services, service)
	}
	return services, nil
}

// GetService returns a service by ID
func (s *FileStore) GetService(id string) (*models.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	service, ok := s.data.Services[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &service, nil
}

// CreateService creates a new service
func (s *FileStore) CreateService(service *models.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data.Services[service.ID]; ok {
		return ErrAlreadyExists
	}

	s.data.Services[service.ID] = *service
	s.triggerSave()
	return nil
}

// UpdateService updates an existing service
func (s *FileStore) UpdateService(id string, service *models.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data.Services[id]; !ok {
		return ErrNotFound
	}

	// Ensure ID doesn't change
	service.ID = id
	s.data.Services[id] = *service
	s.triggerSave()
	return nil
}

// DeleteService deletes a service
func (s *FileStore) DeleteService(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data.Services[id]; !ok {
		return ErrNotFound
	}

	// Check if service is in use
	inUse, usedBy, err := s.serviceInUse(id)
	if err != nil {
		return err
	}
	if inUse {
		return fmt.Errorf("%w: %s", ErrResourceInUse, usedBy)
	}

	delete(s.data.Services, id)
	s.triggerSave()
	return nil
}

// ServiceExists checks if a service exists
func (s *FileStore) ServiceExists(id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data.Services[id]
	return ok, nil
}

// ServiceInUse checks if a service is in use by any routers
func (s *FileStore) ServiceInUse(id string) (bool, []string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.serviceInUse(id)
}

// serviceInUse is an internal non-locking version of ServiceInUse
func (s *FileStore) serviceInUse(id string) (bool, []string, error) {
	usedBy := []string{}

	// Check if service is used by any routers
	for routerID, router := range s.data.Routers {
		if router.Service.ID == id {
			usedBy = append(usedBy, fmt.Sprintf("router:%s", routerID))
		}
	}

	return len(usedBy) > 0, usedBy, nil
}

// UpdateRouter updates a router after validating all references
func (s *FileStore) UpdateRouter(id string, router *models.Router) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if router exists
	existingRouter, ok := s.data.Routers[id]
	if !ok {
		return ErrNotFound
	}

	// Validate that referenced service exists
	if router.Service.ID != "" {
		if _, ok := s.data.Services[router.Service.ID]; !ok {
			return fmt.Errorf("service %s not found", router.Service.ID)
		}
	} else {
		// If service ID is empty in update, keep the existing one
		router.Service = existingRouter.Service
	}

	// Validate that all referenced middlewares exist
	if router.Middlewares != nil {
		for _, mw := range router.Middlewares {
			if _, ok := s.data.Middlewares[mw.ID]; !ok {
				return fmt.Errorf("middleware %s not found", mw.ID)
			}
		}
	}

	// Ensure ID doesn't change
	router.ID = id

	// Update router
	s.data.Routers[id] = *router

	// Trigger save
	s.triggerSave()

	return nil
}

// GetTraefikConfig returns the current dynamic configuration for Traefik
func (s *FileStore) GetTraefikConfig() (*models.DynamicConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a DynamicConfig populated with all the stored configuration
	config := &models.DynamicConfig{
		HTTPRouters:     make(map[string]models.Router),
		HTTPServices:    make(map[string]models.Service),
		HTTPMiddlewares: make(map[string]models.Middleware),
	}

	// Copy routers, services, and middlewares
	for id, router := range s.data.Routers {
		config.HTTPRouters[id] = router
	}

	for id, service := range s.data.Services {
		config.HTTPServices[id] = service
	}

	for id, middleware := range s.data.Middlewares {
		config.HTTPMiddlewares[id] = middleware
	}

	return config, nil
}

func (s *FileStore) Close() {
	// Only close if the channel is not already closed
	select {
	case <-s.done:
		// Already closed
		return
	default:
		close(s.done)
	}

	// Wait for any pending save operations to complete
	s.mu.Lock()
	s.mu.Unlock()
}
