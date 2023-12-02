package auth

import (
	"project/internal/auth/email"
	"project/internal/auth/jwt"
	"project/internal/config"
	"project/internal/database"
	"project/internal/ping"
)

// Dependencies for the base authentication service.
type Dependencies struct {
	Repository   *Repository // base auth repository
	Cache        *Cache      // twofa cache repository
	JWTManager   *jwt.Manager
	EmailManager *email.Manager
	PingManager  *ping.PingManager
	ProdEnv      bool
}

// Constructor for creating all dependencies for the base authentication service.
func NewDependencies(cfg *config.Configuration) (*Dependencies, error) {
	// initialize data access
	databases, err := database.NewDataAccess(cfg.Database.PgURL, cfg.Database.RedisURL, cfg.Database.RedisPsw)
	if err != nil {
		return nil, err
	}
	repo := NewRepository(databases)
	cache := NewCache(databases)

	// initialize dependencies
	emailManager, err := email.NewManger(cfg)
	if err != nil {
		return nil, err
	}

	pingManager, err := ping.NewPingManager(cfg)
	if err != nil {
		return nil, err
	}

	jwtManager := jwt.NewManager(cfg)

	return &Dependencies{
		Repository:   repo,
		Cache:        cache,
		JWTManager:   jwtManager,
		EmailManager: emailManager,
		PingManager:  pingManager,
		ProdEnv:      cfg.Server.ProdEnv,
	}, nil
}
