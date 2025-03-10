// internal/api/handlers/base.go
package handlers

import (
	"github.com/sistemica/traefik-manager/internal/store"
)

// BaseHandler contains common dependencies for all handlers
type BaseHandler struct {
	Store store.Store
}

// NewBaseHandler creates a new BaseHandler with the given dependencies
func NewBaseHandler(store store.Store) BaseHandler {
	return BaseHandler{
		Store: store,
	}
}
