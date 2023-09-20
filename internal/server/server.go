package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"project/internal/auth"
	"project/internal/auth/singlefa"
	"project/internal/config"
	"project/internal/server/routes"
)

// Server contains components and properties for the server.
type Server struct {
	// server components
	Router *chi.Mux
	Cfg    *config.Configuration
	// server properties
	ApiUrl     string
	ApiVersion string
	// dependencies
	AuthDeps *auth.Deps
}

// New is a constructor for the HTTP server which also initializes all needed dependencies.
func New(cfg *config.Configuration) (*Server, error) {
	log.Info().Msg("initializing server...")

	if cfg == nil {
		msg := "nil configuration"
		log.Error().Str("location", "New").Msg(msg)
		return nil, errors.New(msg)
	}

	// initializing auth dependencies
	authDeps, err := auth.NewDependencies(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		Router:     chi.NewRouter(),
		Cfg:        cfg,
		ApiUrl:     cfg.Host + ":" + cfg.Port,
		ApiVersion: cfg.ApiVersion,
		AuthDeps:   authDeps,
	}, nil
}

// This method setups all routes and middlewares.
func (s *Server) SetupRouter() error {
	log.Info().Msg("initializing " + s.ApiVersion + " api routes...")
	s.setupMiddleware()

	// setting up handlers
	sfaHandler, err := singlefa.NewHandler(s.AuthDeps.Service, s.AuthDeps.JWTManger)
	if err != nil {
		return err
	}

	// routing all api endpoints
	s.Router.Get("/health", HealthHandler)
	s.Router.NotFound(NotFoundHandler)
	s.Router.Route("/api/"+s.ApiVersion, func(r chi.Router) {
		r.Route("/sfa", routes.SingleFA(sfaHandler))
	})

	return nil
}

// helper: setupMiddleware setups all middlewares.
func (s *Server) setupMiddleware() {
	s.Router.Use(middleware.Logger)

}

// Starts the HTTP server.
func (s *Server) Run() {
	log.Info().Msg("server is running on " + s.ApiUrl + "/api/" + s.ApiVersion)
	http.ListenAndServe(s.ApiUrl, s.Router)
}
