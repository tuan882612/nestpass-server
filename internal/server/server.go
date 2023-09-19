package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"project/internal/auth"
	"project/internal/auth/singlefa"
	"project/internal/config"
	"project/internal/server/routes"
)

type Server struct {
	Router     *chi.Mux
	ApiUrl     string
	ApiVersion string
	Cfg        *config.Configuration
}

func NewServer() (*Server, error) {
	log.Info().Msg("initializing server...")

	// initialize and validate new configuration instance
	cfg := config.NewConfiguration()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Server{
		Router:     chi.NewRouter(),
		ApiUrl:     cfg.Host + ":" + cfg.Port,
		ApiVersion: cfg.ApiVersion,
		Cfg:        cfg,
	}, nil
}

func (s *Server) SetupRouter() error {
	log.Info().Msg("initializing " + s.ApiVersion + " api routes...")

	s.setupMiddleware()

	authRepo, err := auth.NewRepository(s.Cfg)
	if err != nil {
		return err
	}

	authService, err := auth.NewService(authRepo)
	if err != nil {
		return err
	}

	jwtHandler, err := auth.NewJWTManager(s.Cfg.SignKey, s.Cfg.Duration)
	if err != nil {
		return err
	}

	sfaHandler, err := singlefa.NewHandler(authService, jwtHandler)
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

func (s *Server) setupMiddleware() {
	s.Router.Use(middleware.Logger)

}

func (s *Server) Run() {
	log.Info().Msg("server is running on " + s.ApiUrl + "/api/" + s.ApiVersion)
	http.ListenAndServe(s.ApiUrl, s.Router)
}
