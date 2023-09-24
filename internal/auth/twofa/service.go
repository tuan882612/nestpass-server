package twofa

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/auth"
	"project/internal/auth/email"
	"project/internal/auth/jwt"
	"project/internal/proto/pb/twofapb"
)

// Service for handling two-factor authentication.
type service struct {
	authRepo     auth.Repository // base auth repository
	authService  auth.Service    // base auth service
	jwtManager   *jwt.Manager
	emailManager *email.Manager
}

// Creates a new two-factor authentication service with the given dependencies.
func NewService(deps *auth.Dependencies) Service {
	return &service{
		authRepo:     deps.Repository,
		authService:  deps.Service,
		jwtManager:   deps.JWTManager,
		emailManager: deps.EmailManager,
	}
}

func (s *service) SendVerificationEmail(ctx context.Context, userID uuid.UUID, email string) error {
	// Make new payload with input email and user id
	payload := &twofapb.TwoFAPayload{
		UserId: userID.String(),
		Email:  email,
	}

	// send the two-factor auth code to the user's email
	_, err := s.emailManager.Client.GenerateTwoFACode(ctx, payload)
	if err != nil {
		log.Error().Str("location", "SendVerificationEmail").Msg(err.Error())
		return err
	}

	return nil
}

// Verifies the two-factor auth code and returns a JWT token if the verification is successful.
func (s *service) VerifyAuthToken(ctx context.Context, userID uuid.UUID, input *email.TokenInput) (string, error) {
	tfaBody, err := s.authRepo.GetTwofaCache(ctx, userID)
	if err != nil {
		return "", err
	}

	if tfaBody.Code != input.Token {
		// check if the user has any retries left
		if tfaBody.Retries -= 1; tfaBody.Retries == 0 {
			if err := s.authRepo.DeleteTwofaCache(ctx, userID); err != nil {
				return "", err
			}

			return "", apiutils.NewErrUnauthorized("too many retries")
		}

		// update the twofa retries
		if err := s.authRepo.UpdateTwofaCache(ctx, userID, tfaBody); err != nil {
			return "", err
		}

		return "", apiutils.NewErrUnauthorized("invalid code")
	}

	// delete the twofa data
	go func() {
		if err := s.authRepo.DeleteTwofaCache(ctx, userID); err != nil {
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
func (s *service) LoginSend(ctx context.Context, input *auth.LoginInput) (string, error) {
	data, err := s.authService.VerifyUser(ctx, input.Email, input.Password)
	if err != nil {
		return "", err
	}

	// send the two-factor auth code to the user's email in the background
	go func() {
		if err := s.SendVerificationEmail(ctx, data, input.Email); err != nil {
			log.Error().Str("location", "LoginSend").Msg("failed to send email: " + err.Error())
		}
	}()

	return data.String(), nil
}

func (s *service) RegisterSend(ctx context.Context, input *auth.RegisterInput) (string, error) {
	// convert RegisterInput to RegisterResp and validate the input
	regResp, err := auth.NewRegisterResp(input)
	if err != nil {
		return "", err
	}

	// send the two-factor auth code to the user's email in the background
	go func() {
		if err := s.SendVerificationEmail(ctx, regResp.UserID, input.Email); err != nil {
			log.Error().Str("location", "RegisterSend").Msg("failed to send email: " + err.Error())
		}
	}()

	// register the user
	if err := s.authService.RegisterUser(ctx, regResp); err != nil {
		return "", err
	}

	return regResp.UserID.String(), nil
}

func (s *service) RegisterVerify(ctx context.Context, userID uuid.UUID, input *email.TokenInput) (string, error) {
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
