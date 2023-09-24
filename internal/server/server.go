package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"project/internal/auth"
	"project/internal/auth/singlefa"
	"project/internal/auth/twofa"
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
	AuthDeps *auth.Dependencies
}

// Creates a new HTTP server along with initializing all needed dependencies.
func New(cfg *config.Configuration) (*Server, error) {
	log.Info().Msg("initializing server...")

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

// Setups all routes and middlewares.
func (s *Server) SetupRouter() error {
	log.Info().Msg("initializing " + s.ApiVersion + " api routes...")
	s.setupMiddleware()

	// setting up handlers
	sfaHandler := singlefa.NewHandler(s.AuthDeps)
	twofaHandler := twofa.NewHandler(s.AuthDeps)

	// routing all api endpoints
	s.Router.NotFound(NotFoundHandler)
	s.Router.Route("/api/"+s.ApiVersion, func(r chi.Router) {
		r.Get("/health", HealthHandler)
		r.Route("/sfa", routes.SingleFA(sfaHandler))
		r.Route("/twofa", routes.TwoFA(twofaHandler, r))
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
