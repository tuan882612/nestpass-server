package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"nestpass/internal/config"
	"nestpass/internal/dependencies"
	"nestpass/internal/server/routes"
)

// Server contains components and properties for the server.
type Server struct {
	// server components
	Router *chi.Mux
	Cfg    *config.Configuration
	// server properties
	ApiAddr    string
	ApiVersion string
	// dependencies
	Deps *dependencies.Dependencies
}

// Creates a new HTTP server along with initializing all needed dependencies.
func New(cfg *config.Configuration) (*Server, error) {
	log.Info().Msg("initializing server...")

	// initialize dependencies
	deps, err := dependencies.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		Router:     chi.NewRouter(),
		Cfg:        cfg,
		ApiAddr:    cfg.Host + ":" + cfg.Port,
		ApiVersion: "/api/" + cfg.ApiVersion,
		Deps:       deps,
	}, nil
}

// Setups all routes and middlewares.
func (s *Server) SetupRouter() error {
	log.Info().Msg("initializing " + s.ApiVersion + " api routes...")
	s.setupMiddleware()

	// initialize api handler
	apiHandler := routes.NewAPIHandler(s.Deps)

	// routing internal endpoints
	s.Router.Patch("/rehash", apiHandler.Password.RehashAllPasswords)

	// routing all api endpoints
	s.Router.NotFound(NotFoundHandler)
	s.Router.Route(s.ApiVersion, func(r chi.Router) {
		r.Get("/health", HealthHandler)
		r.Route("/user", routes.Users(apiHandler, s.Cfg))
	})

	return nil
}

// helper: setupMiddleware setups all middlewares.
func (s *Server) setupMiddleware() {
	s.Router.Use(middleware.Logger)

}

// Starts the HTTP server.
func (s *Server) Run() {
	log.Info().Msg("server is running on " + s.ApiAddr + s.ApiVersion)
	http.ListenAndServe(s.ApiAddr, s.Router)
}
