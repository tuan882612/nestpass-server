package twofa

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/auth"
	"project/internal/auth/email"
	"project/internal/auth/jwt"
	"project/internal/proto/pb/twofapb"
)

// Service for handling two-factor authentication.
type Service struct {
	authRepo     *auth.Repository // base auth repository
	authService  *auth.Service    // base auth service
	cacheRepo    *auth.Cache      // cache repository
	jwtManager   *jwt.Manager
	emailManager *email.Manager
}

// Creates a new two-factor authentication service with the given dependencies.
func NewService(deps *auth.Dependencies) *Service {
	return &Service{
		authRepo:     deps.Repository,
		authService:  deps.Service,
		cacheRepo:    deps.Cache,
		jwtManager:   deps.JWTManager,
		emailManager: deps.EmailManager,
	}
}

func (s *Service) SendVerificationEmail(userID uuid.UUID, email string) {
	// create new context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Make new payload with input email and user id
	payload := &twofapb.TwoFAPayload{
		UserId: userID.String(),
		Email:  email,
	}

	// send the two-factor auth code to the user's email
	_, err := s.emailManager.Client.GenerateTwoFACode(ctx, payload)
	if err != nil {
		log.Error().Str("location", "SendVerificationEmail").Msg("failed to send verification email: " + err.Error())
		return
	}

	log.Info().Msg("successfully sent verification email")
}

// Verifies the two-factor auth code and returns a JWT token if the verification is successful.
func (s *Service) VerifyAuthToken(ctx context.Context, userID uuid.UUID, input *email.TokenInput) (string, error) {
	tfaBody, err := s.cacheRepo.GetTwofaCache(ctx, userID)
	if err != nil {
		return "", err
	}

	if tfaBody.Code != input.Token {
		// check if the user has any retries left
		if tfaBody.Retries -= 1; tfaBody.Retries == 0 {
			if err := s.cacheRepo.DeleteTwofaCache(ctx, userID); err != nil {
				return "", err
			}

			return "", apiutils.NewErrUnauthorized("too many retries")
		}

		// update the twofa retries
		if err := s.cacheRepo.UpdateTwofaCache(ctx, userID, tfaBody); err != nil {
			return "", err
		}

		return "", apiutils.NewErrUnauthorized("invalid code")
	}

	// delete the twofa data
	go func() {
		if err := s.cacheRepo.DeleteTwofaCache(ctx, userID); err != nil {
			log.Error().Str("location", "VerifyAuthToken").Msg(err.Error())
			return
		}
	}()

	// generate a JWT token
	tokenChan := make(chan string, 1)
	go func() {
		token, err := s.jwtManager.GenerateToken(userID)
		if err != nil {
			log.Error().Str("location", "VerifyAuthToken").Msg(err.Error())
			return
		}

		tokenChan <- token
	}()

	return <-tokenChan, nil
}

// Sends a two-factor authentication code to the user's email.
func (s *Service) LoginSend(ctx context.Context, input *auth.LoginInput) (string, error) {
	data, err := s.authService.VerifyUser(ctx, input.Email, input.Password)
	if err != nil {
		return "", err
	}

	// send the two-factor auth code to the user's email in the background
	go s.SendVerificationEmail(data, input.Email)

	return data.String(), nil
}

func (s *Service) RegisterSend(ctx context.Context, input *auth.RegisterInput) (string, error) {
	// convert RegisterInput to RegisterResp and validate the input
	regResp, err := auth.NewRegisterResp(input)
	if err != nil {
		return "", err
	}

	// register the user
	if err := s.authService.RegisterUser(ctx, regResp); err != nil {
		return "", err
	}

	// send the two-factor auth code to the user's email in the background
	go s.SendVerificationEmail(regResp.UserID, input.Email)

	return regResp.UserID.String(), nil
}

func (s *Service) RegisterVerify(ctx context.Context, userID uuid.UUID, input *email.TokenInput) (string, error) {
	token, err := s.VerifyAuthToken(ctx, userID, input)
	if err != nil {
		return "", err
	}

	// update the user's status in the background
	go func() {
		if err := s.authRepo.UpdateUserStatus(ctx, userID); err != nil {
			log.Error().Str("location", "RegisterVerify").Msg(err.Error())
		}
	}()

	return token, nil
}
