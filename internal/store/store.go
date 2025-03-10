package store

import (
	"github.com/sistemica/traefik-manager/internal/models"
)

// Store is the interface for the storage layer
type Store interface {
	// Middlewares
	ListMiddlewares() ([]models.Middleware, error)
	GetMiddleware(id string) (*models.Middleware, error)
	CreateMiddleware(middleware *models.Middleware) error
	UpdateMiddleware(id string, middleware *models.Middleware) error
	DeleteMiddleware(id string) error
	MiddlewareExists(id string) (bool, error)
	MiddlewareInUse(id string) (bool, []string, error)

	// Routers
	ListRouters() ([]models.Router, error)
	GetRouter(id string) (*models.Router, error)
	CreateRouter(router *models.Router) error
	UpdateRouter(id string, router *models.Router) error
	DeleteRouter(id string) error
	RouterExists(id string) (bool, error)
	RouterInUse(id string) (bool, []string, error)

	// Services
	ListServices() ([]models.Service, error)
	GetService(id string) (*models.Service, error)
	CreateService(service *models.Service) error
	UpdateService(id string, service *models.Service) error
	DeleteService(id string) error
	ServiceExists(id string) (bool, error)
	ServiceInUse(id string) (bool, []string, error)

	// Persistence
	Save() error
	Load() error
	Close()
}
