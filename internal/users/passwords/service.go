package passwords

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"nestpass/pkg/httputils"
)

type service struct {
	repo *repository
}

func NewService(repo *repository) *service {
	return &service{repo: repo}
}

func (s *service) GetKDFKey(ctx context.Context, userID uuid.UUID) ([]byte, error) {
	kdfData, err := s.repo.GetKDFKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	return KeyDerivation(kdfData.PswHash, userID.String(), kdfData.Salt), nil
}

func (s *service) DecryptAndGetPsw(password *PasswordEncrypt, userID uuid.UUID, key []byte) (*PasswordResp, error) {
	// decrypt rawData from encrypted
	data, err := Decrypt(password.Nonce, password.Encrypted, key)
	if err != nil {
		log.Error().Str("location", "GetAllPasswords").Msgf("%v: decrypt err %v", userID, err)
		return nil, err
	}

	// create Password
	psw := NewPasswordResp(password.PasswordID, data, password.Website, userID, password.CategoryID)
	return psw, nil
}

func (s *service) GetAllPasswords(ctx context.Context, userID uuid.UUID, pageParams *httputils.Pagination) ([]*PasswordResp, error) {
	// retrieve kdf key
	key, err := s.GetKDFKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	// retrieve encrypted passwords
	passwords, err := s.repo.GetAllPasswords(ctx, userID, key, pageParams)
	if err != nil {
		return nil, err
	}

	// decrypt passwords
	decryptedPsw := []*PasswordResp{}
	for _, password := range passwords {
		psw, err := s.DecryptAndGetPsw(password, userID, key)
		if err != nil {
			return nil, err
		}

		decryptedPsw = append(decryptedPsw, psw)
	}

	return decryptedPsw, nil
}

func (s *service) GetAllPasswordsByCategory(ctx context.Context, userID, categoryID uuid.UUID, pageParams *httputils.Pagination) ([]*PasswordResp, error) {
	// retrieve kdf key
	key, err := s.GetKDFKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	// retrieve encrypted passwords
	passwords, err := s.repo.GetAllPasswordsByCategory(ctx, userID, categoryID, key, pageParams)
	if err != nil {
		return nil, err
	}

	// decrypt passwords
	decryptedPsw := []*PasswordResp{}
	for _, password := range passwords {
		psw, err := s.DecryptAndGetPsw(password, userID, key)
		if err != nil {
			return nil, err
		}

		decryptedPsw = append(decryptedPsw, psw)
	}

	return decryptedPsw, nil
}

func (s *service) GetPassword(ctx context.Context, passwordID, categoryID, userID uuid.UUID) (*PasswordResp, error) {
	// retrieve kdf key
	key, err := s.GetKDFKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	// retrieve encrypted password
	password, err := s.repo.GetPassword(ctx, passwordID, categoryID, userID)
	if err != nil {
		return nil, err
	}

	// decrypt password
	psw, err := s.DecryptAndGetPsw(password, userID, key)
	if err != nil {
		return nil, err
	}

	return psw, nil
}

func (s *service) CreatePassword(ctx context.Context, psw *Password) error {
	key, err := s.GetKDFKey(ctx, psw.UserID)
	if err != nil {
		return err
	}

	data, err := NewPasswordEncrypt(psw, key)
	if err != nil {
		return err
	}

	tx, err := s.repo.postgres.Begin(ctx)
	if err != nil {
		log.Error().Str("location", "CreatePassword").Msgf("%v: %v", psw.UserID, err)
		return err
	}

	if err := s.repo.CreatePassword(ctx, tx, data); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Str("location", "CreatePassword").Msgf("%v: %v", psw.UserID, err)
		return err
	}

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, psw *PasswordResp) error {
	key, err := s.GetKDFKey(ctx, psw.UserID)
	if err != nil {
		return err
	}

	data, err := NewPasswordEncrypt(&psw.Password, key)
	if err != nil {
		return err
	}
	// uuid is not being set properly

	tx, err := s.repo.postgres.Begin(ctx)
	if err != nil {
		log.Error().Str("location", "UpdatePassword").Msgf("%v: %v", psw.UserID, err)
		return err
	}

	if err := s.repo.UpdatePassword(ctx, tx, data); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Str("location", "UpdatePassword").Msgf("%v: %v", psw.UserID, err)
		return err
	}

	return nil
}

func (s *service) ReUpdateAllPasswords() {}

func (s *service) DeletePassword(ctx context.Context, userID, passwordID, categoryID uuid.UUID) error {
	tx, err := s.repo.postgres.Begin(ctx)
	if err != nil {
		log.Error().Str("location", "DeletePassword").Msgf("%v: %v", categoryID, err)
		return err
	}

	if err := s.repo.DeletePassword(ctx, tx, userID, passwordID, categoryID); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error().Str("location", "DeletePassword").Msgf("%v: %v", categoryID, err)
		return err
	}

	return nil
}
