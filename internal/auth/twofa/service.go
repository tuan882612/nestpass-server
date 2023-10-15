package twofa

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"
	"github.com/tuan882612/apiutils/securityutils"

	"project/internal/auth"
	"project/internal/auth/email"
	"project/internal/auth/jwt"
	"project/internal/proto/pb/twofapb"
)

// Service for handling two-factor authentication.
type Service struct {
	authRepo     *auth.Repository // base auth repository
	cacheRepo    *auth.Cache      // cache repository
	jwtManager   *jwt.Manager
	emailManager *email.Manager
}

// Creates a new two-factor authentication service with the given dependencies.
func NewService(deps *auth.Dependencies) *Service {
	return &Service{
		authRepo:     deps.Repository,
		cacheRepo:    deps.Cache,
		jwtManager:   deps.JWTManager,
		emailManager: deps.EmailManager,
	}
}

// Sends a two-factor authentication code to the user's email.
func (s *Service) SendVerificationEmail(ctx context.Context, userID uuid.UUID, email, status string) error {
	// check if data is restricted
	restricted, err := s.cacheRepo.IsRestricted(ctx, userID)
	if err != nil {
		return err
	}

	if restricted {
		return apiutils.NewErrForbidden("user is restricted")
	}

	// send request to email service to generate a two-factor auth code
	go func() {
		// create new context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Make new payload with input email and user id
		payload := &twofapb.TwoFAPayload{
			UserId:     userID.String(),
			Email:      email,
			UserStatus: status,
		}

		// send the two-factor auth code to the user's email
		_, err := s.emailManager.Client.GenerateTwoFACode(ctx, payload)
		if err != nil {
			log.Error().Str("location", "SendVerificationEmail").Msgf("%v: failed to send verification email: %v", userID, err)
			return
		}

		log.Info().Msgf("%v: successfully sent verification email", userID)
	}()

	return nil
}

// resend twofa code
func (s *Service) ResendCode(ctx context.Context, email string) error {
	user, err := s.authRepo.GetUserCredentials(ctx, email)
	if err != nil {
		return err
	}

	// check if the user is inactive
	if user.UserStatus == auth.InactiveUser {
		return apiutils.NewErrForbidden("user is inactive")
	}

	// send the two-factor auth code to the user's email in the background
	if err := s.SendVerificationEmail(ctx, user.UserID, email, user.UserStatus); err != nil {
		return err
	}

	return nil
}

// Verifies the two-factor auth code and returns a JWT token if the verification is successful.
func (s *Service) VerifyAuthToken(ctx context.Context, userID uuid.UUID, token string) (string, error) {
	tfaBody, err := s.cacheRepo.GetTwofa(ctx, userID)
	if err != nil {
		return "", err
	}

	if tfaBody.Code != token {
		// check if the user has any retries left
		if tfaBody.Retries -= 1; tfaBody.Retries == 0 {
			// add the user as restricted async
			go func() {
				if err := s.cacheRepo.AddRestricted(ctx, userID); err != nil {
					log.Error().Str("location", "VerifyAuthToken").Msgf("%v: failed to restrict user: %v", userID, err)
					return
				}

				log.Info().Msgf("%v: successfully restricted user", userID)
			}()

			return "", apiutils.NewErrUnauthorized("too many retries")
		}

		// update the twofa retries async
		go func() {
			if err := s.cacheRepo.UpdateTwofa(ctx, userID, tfaBody); err != nil {
				log.Error().Str("location", "VerifyAuthToken").Msgf("%v: failed to update twofa retries: %v", userID, err)
				return
			}

			log.Info().Msgf("%v: successfully updated twofa retries", userID)
		}()

		return "", apiutils.NewErrUnauthorized("invalid code")
	}

	// delete the twofa data async
	go func() {
		if err := s.cacheRepo.DeleteTwofa(ctx, userID); err != nil {
			log.Error().Str("location", "VerifyAuthToken").Msgf("%v: failed to delete twofa cache: %v", userID, err)
			return
		}

		log.Info().Msgf("%v: successfully deleted twofa cache", userID)
	}()

	fmt.Println(tfaBody)

	// update the user's status in the background if the user is a non-registered user
	if tfaBody.UserStatus == auth.NonRegUser {
		go func() {
			// new context with a timeout
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			if err := s.authRepo.UpdateUserStatus(ctx, userID); err != nil {
				return
			}

			log.Info().Msgf("%v: successfully updated user status", userID)
		}()
	}

	// generate a JWT token
	authToken, err := s.jwtManager.GenerateToken(userID)
	if err != nil {
		log.Error().Str("location", "VerifyAuthToken").Msgf("%v: failed to generate JWT token: %v", userID, err)
		return "", err
	}

	return authToken, nil
}

// Initial twofa login
func (s *Service) LoginSend(ctx context.Context, input *auth.LoginInput) (string, error) {
	// retrieve the user credentials from the database
	user, err := s.authRepo.GetUserCredentials(ctx, input.Email)
	if err != nil {
		return "", err
	}

	// validate the password
	if err := securityutils.ValidatePassword(user.Password, input.Password); err != nil {
		return "", apiutils.NewErrUnauthorized(err.Error())
	}

	// check if user is inactive
	if user.UserStatus == auth.InactiveUser {
		return "", apiutils.NewErrForbidden("user is inactive")
	}

	// send the two-factor auth code to the user's email in the background
	if err := s.SendVerificationEmail(ctx, user.UserID, input.Email, user.UserStatus); err != nil {
		return "", err
	}

	return user.UserID.String(), nil
}

// Initial twofa register
func (s *Service) RegisterSend(ctx context.Context, input *auth.RegisterInput) (string, error) {
	// convert RegisterInput to RegisterResp and validate the input
	regResp, err := auth.NewRegisterResp(input)
	if err != nil {
		return "", err
	}

	// register the user
	tx, err := s.authRepo.StartTx(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	// tries to add the user to the database
	if err := s.authRepo.AddUser(ctx, tx, regResp); err != nil {
		return "", err
	}

	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		log.Error().Str("location", "RegisterUser").Msgf("failed to commit transaction: %v", err)
		return "", err
	}

	// send the two-factor auth code to the user's email in the background
	if err := s.SendVerificationEmail(ctx, regResp.UserID, input.Email, regResp.UserStatus); err != nil {
		return "", err
	}

	return regResp.UserID.String(), nil
}

// Final twofa register
func (s *Service) RegisterVerify(ctx context.Context, userID uuid.UUID, input *email.TokenInput) (string, error) {
	token, err := s.VerifyAuthToken(ctx, userID, input.Token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Initial twofa reset password
func (s *Service) ResetPassword(ctx context.Context, email string) (string, error) {
	user, err := s.authRepo.GetUserCredentials(ctx, email)
	if err != nil {
		return "", err
	}

	// check if the user is inactive
	if user.UserStatus == auth.InactiveUser {
		return "", apiutils.NewErrForbidden("user is inactive")
	}

	// send the two-factor auth code to the user's email in the background
	if err := s.SendVerificationEmail(ctx, user.UserID, email, user.UserStatus); err != nil {
		return "", err
	}

	return user.UserID.String(), nil
}

// Final twofa reset password
func (s *Service) ResetPasswordFinal(ctx context.Context, userID uuid.UUID, password string) (string, error) {
	// verify the user's twofa code
	token, err := s.VerifyAuthToken(ctx, userID, password)
	if err != nil {
		return "", err
	}

	// hash the new password
	hashedPassword, err := securityutils.HashPassword(password)
	if err != nil {
		log.Error().Str("location", "ResetPasswordFinal").Msgf("%v: failed to hash password: %v", userID, err)
		return "", err
	}

	// update the user's password
	if err := s.authRepo.UpdateUserPassword(ctx, userID, hashedPassword); err != nil {
		return "", err
	}

	return token, nil
}
