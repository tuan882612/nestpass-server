package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"project/internal/auth"
	"project/internal/auth/jwt"
	"project/internal/config"
)

const userInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"

type Service struct {
	oauthCfg   *oauth2.Config
	repo       *auth.Repository
	jwtManager *jwt.Manager
}

func NewService(oauthCfg *config.OAuthConfig, deps *auth.Dependencies) *Service {
	cfg := &oauth2.Config{
		ClientID:     oauthCfg.ClientID,
		ClientSecret: oauthCfg.ClientSecret,
		RedirectURL:  "http://localhost:2001/api/v1/oauth/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	return &Service{
		oauthCfg:   cfg,
		repo:       deps.Repository,
		jwtManager: deps.JWTManager,
	}
}

// Generates a CSRF state token
func (s *Service) generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (s *Service) getOAuthData(ctx context.Context, token *oauth2.Token) (*OAuthData, error) {
	// build new request
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		log.Error().Str("location", "getOAuthData").Msgf("Failed to create request: %s", err.Error())
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	// initiate the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Str("location", "getOAuthData").Msgf("Failed to get user data: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// decode the response body for email and name
	data := &OAuthData{}
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		log.Error().Err(err).Msgf("Failed to decode response body: %s", resp.Body)
		return nil, err
	}

	return data, nil
}

func (s *Service) StartOAuth(ctx context.Context, state string) string {
	return s.oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *Service) CallbackOAuth(ctx context.Context, code string, state string, expectedState string) (*oauth2.Token, error) {
	if state != expectedState {
		return nil, errors.New("invalid state token")
	}

	token, err := s.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) UserLoginSignup(ctx context.Context, token *oauth2.Token) (string, error) {
	data, err := s.getOAuthData(ctx, token)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetUserCredentials(ctx, data.Email)
	if err != nil {
		switch err.(type) {
		case apiutils.ErrNotFound:
			newUser := data.NewUser()

			// new context with a timeou
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()

				tx, err := s.repo.StartTx(ctx)
				if err != nil {
					log.Error().Str("location", "UserLoginSignup").Msgf("%v: failed to start transaction: %v", newUser.UserID, err)
					return
				}
				defer tx.Rollback(ctx)

				if err := s.repo.AddUser(ctx, tx, newUser); err != nil {
					return
				}

				if err := tx.Commit(ctx); err != nil {
					log.Error().Str("location", "UserLoginSignup").Msgf("%v: failed to commit transaction: %v", newUser.UserID, err)
					return
				}

				log.Info().Msgf("%v: successfully registered user", newUser.UserID)
			}()

			authToken, err := s.jwtManager.GenerateToken(newUser.UserID)
			if err != nil {
				return "", err
			}

			return authToken, nil
		default:
			// some other error occurred
			return "", err
		}
	}

	authToken, err := s.jwtManager.GenerateToken(user.UserID)
	if err != nil {
		return "", err
	}

	return authToken, nil
}
