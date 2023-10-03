package config

import "os"

type OAuthConfig struct {
	ClientID     string `validate:"required"`
	ProjectID    string `validate:"required"`
	AuthURI      string `validate:"required"`
	TokenURI     string `validate:"required"`
	CertURL      string `validate:"required"`
	ClientSecret string `validate:"required"`
}

func newOauthConfig() *OAuthConfig {
	return &OAuthConfig{
		ClientID:     os.Getenv("CLIENT_ID"),
		ProjectID:    os.Getenv("PROJECT_ID"),
		AuthURI:      os.Getenv("AUTH_URI"),
		TokenURI:     os.Getenv("TOKEN_URI"),
		CertURL:      os.Getenv("AUTH_PROVIDER_X509_CERT_URL"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}
}
