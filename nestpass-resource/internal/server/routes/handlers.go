package routes

import (
	"nestpass/internal/dependencies"
	"nestpass/internal/users"
	"nestpass/internal/users/categories"
	"nestpass/internal/users/passwords"
)

type APIHandler struct {
	User     *users.Handler
	Category *categories.Handler
	Password *passwords.Handler
}

func NewAPIHandler(deps *dependencies.Dependencies) *APIHandler {
	return &APIHandler{
		User:     users.NewHandler(deps),
		Category: categories.NewHandler(deps),
		Password: passwords.NewHandler(deps),
	}
}
