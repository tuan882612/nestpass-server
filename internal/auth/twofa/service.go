package twofa

import (
	"context"
	"errors"
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
	if _, err := s.cacheRepo.GetData(ctx, userID, auth.Restricted); err != nil {
		return err
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
func (s *Service) VerifyAuthToken(ctx context.Context, userID uuid.UUID, token, mode string) (string, error) {
	data, err := s.cacheRepo.GetData(ctx, userID, auth.TwoFA)
	if err != nil {
		return "", err
	}

	// check if the data type is correct
	tfaBody, ok := data.(*email.TwofaBody)
	if !ok {
		return "", errors.New("invalid twofa data")
	}

	var retriesErr error = nil
	// check if the code is correct
	if tfaBody.Code != token {
		tfaBody.Retries -= 1

		// check if the user has any retries left and update the twofa data async
		if tfaBody.Retries != 0 {
			go func() {
				if err := s.cacheRepo.UpdateTwofa(ctx, userID, tfaBody); err != nil {
					log.Error().Str("location", "VerifyAuthToken").Msgf("%v: failed to update twofa retries: %v", userID, err)
					return
				}
	
				log.Info().Msgf("%v: successfully updated twofa retries", userID)
			}()

			return "", apiutils.NewErrUnauthorized("invalid code")
		}

		// add the user as restricted async
		go func() {
			if err := s.cacheRepo.AddRestricted(ctx, userID); err != nil {
				log.Error().Str("location", "VerifyAuthToken").Msgf("%v: failed to restrict user: %v", userID, err)
				return
			}

			log.Info().Msgf("%v: successfully restricted user", userID)
		}()

		retriesErr = apiutils.NewErrUnauthorized("too many retries")
	}

	// delete the twofa data async
	go func() {
		if err := s.cacheRepo.DeleteData(ctx, userID, auth.TwoFA); err != nil {
			log.Error().Str("location", "VerifyAuthToken").Msgf("%v: failed to delete twofa cache: %v", userID, err)
			return
		}

		log.Info().Msgf("%v: successfully deleted twofa cache", userID)
	}()

	// check if there was a retries error	
	if retriesErr != nil {
		return "", retriesErr
	}

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

	if mode == "reset" {
		// creates 30 minute session on successful verification in the background
		go func() {
			if err := s.cacheRepo.AddSession(ctx, userID); err != nil {
				return
			}

			log.Info().Msgf("%v: successfully added 30 session", userID)
		}()

		return userID.String(), nil
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

	// check if user is oauth user
	if user.Password == "" {
		return "", apiutils.NewErrBadRequest("user is oauth user")
	}

	// check if user is inactive
	if user.UserStatus == auth.InactiveUser {
		return "", apiutils.NewErrForbidden("user is inactive")
	}

	// validate the password
	if err := securityutils.ValidatePassword(user.Password, input.Password); err != nil {
		return "", apiutils.NewErrUnauthorized(err.Error())
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

	// start a transaction for registering the user
	tx, err := s.authRepo.StartTx(ctx)
	if err != nil {
		log.Error().Str("location", "RegisterUser").Msgf("%v: failed to start transaction: %v", regResp.UserID, err)
		return "", err
	}
	defer tx.Rollback(ctx)

	// add the user to the database
	if err := s.authRepo.AddUser(ctx, tx, regResp); err != nil {
		return "", err
	}

	// commit the transaction
	if err := tx.Commit(ctx); err != nil {
		log.Error().Str("location", "RegisterUser").Msgf("%v: failed to commit transaction: %v", regResp.UserID, err)
		return "", err
	}

	// send the two-factor auth code to the user's email in the background
	if err := s.SendVerificationEmail(ctx, regResp.UserID, input.Email, regResp.UserStatus); err != nil {
		return "", err
	}

	return regResp.UserID.String(), nil
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
func (s *Service) ResetPasswordFinal(ctx context.Context, userID uuid.UUID, password string) error {
	// checks if the user has a 30 minute session
	if _, err := s.cacheRepo.GetData(ctx, userID, auth.Session); err != nil {
		return err
	}

	// check if password is duplicate
	psw, err := s.authRepo.GetUserPassword(ctx, userID)
	if err != nil {
		return err
	}

	if err := securityutils.ValidatePassword(psw, password); err == nil {
		return apiutils.NewErrBadRequest("password is duplicate")
	}

	// update the user's password in the background
	go func() {
		// new context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// hash the new password
		hashedPassword, err := securityutils.HashPassword(password)
		if err != nil {
			log.Error().Str("location", "ResetPasswordFinal").Msgf("%v: failed to hash password: %v", userID, err)
			return
		}

		if err := s.authRepo.UpdateUserPassword(ctx, userID, hashedPassword); err != nil {
			return
		}

		log.Info().Msgf("%v: successfully updated user password", userID)
	}()

	// delete the 30 minute session in the background
	go func() {
		// new context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := s.cacheRepo.DeleteData(ctx, userID, auth.Session); err != nil {
			log.Error().Str("location", "ResetPasswordFinal").Msgf("%v: failed to delete session: %v", userID, err)
			return
		}

		log.Info().Msgf("%v: successfully deleted session", userID)
	}()

	return nil
}
