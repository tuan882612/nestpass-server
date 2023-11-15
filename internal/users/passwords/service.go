package passwords

import (
	"context"
	"encoding/base64"

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

func (s *service) GetKDFKey(ctx context.Context, userID uuid.UUID, kdf KDFType) ([]byte, error) {
	kdfData, err := s.repo.GetKDFData(ctx, userID)
	if err != nil {
		return nil, err
	}

	var key string

	switch kdf {
	case CurrKDF:
		key = kdfData.PswHash
	case PrevKDF:
		data, err := s.repo.GetResetHash(ctx, userID)
		if err != nil {
			return nil, err
		}

		key64, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			log.Error().Str("location", "GetKDFKey").Msgf("%v: %v", userID, err)
			return nil, err
		}

		key = string(key64)
	}

	return KeyDerivation(key, userID.String(), kdfData.Salt), nil
}

func (s *service) DecryptAndGetPsw(password *PasswordEncrypt, userID uuid.UUID, key []byte) (*Password, error) {
	// decrypt rawData from encrypted
	data, err := Decrypt(password.Nonce, password.Encrypted, key)
	if err != nil {
		log.Error().Str("location", "GetAllPasswords").Msgf("%v: decrypt err %v", userID, err)
		return nil, err
	}

	// create Password
	return &Password{
		PasswordID:  password.PasswordID,
		UserID:      password.UserID,
		CategoryID:  password.CategoryID,
		Website:     password.Website,
		Username:    data.Username,
		Password:    data.Password,
		Description: data.Description,
	}, nil
}

func (s *service) GetAllPasswords(ctx context.Context, userID uuid.UUID, pageParams *httputils.Pagination) ([]*Password, error) {
	// retrieve kdf key
	key, err := s.GetKDFKey(ctx, userID, CurrKDF)
	if err != nil {
		return nil, err
	}

	// retrieve encrypted passwords
	passwords, err := s.repo.GetAllPasswords(ctx, userID, pageParams)
	if err != nil {
		return nil, err
	}

	// decrypt passwords
	decryptedPsw := []*Password{}
	for _, password := range passwords {
		psw, err := s.DecryptAndGetPsw(password, userID, key)
		if err != nil {
			return nil, err
		}

		decryptedPsw = append(decryptedPsw, psw)
	}

	return decryptedPsw, nil
}

func (s *service) GetAllPasswordsByCategory(ctx context.Context, userID, categoryID uuid.UUID, pageParams *httputils.Pagination) ([]*Password, error) {
	// retrieve kdf key
	key, err := s.GetKDFKey(ctx, userID, CurrKDF)
	if err != nil {
		return nil, err
	}

	// retrieve encrypted passwords
	passwords, err := s.repo.GetAllPasswordsByCategory(ctx, userID, categoryID, pageParams)
	if err != nil {
		return nil, err
	}

	// decrypt passwords
	decryptedPsw := []*Password{}
	for _, password := range passwords {
		psw, err := s.DecryptAndGetPsw(password, userID, key)
		if err != nil {
			return nil, err
		}

		decryptedPsw = append(decryptedPsw, psw)
	}

	return decryptedPsw, nil
}

func (s *service) GetPassword(ctx context.Context, passwordID, categoryID, userID uuid.UUID) (*Password, error) {
	// retrieve kdf key
	key, err := s.GetKDFKey(ctx, userID, CurrKDF)
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
	key, err := s.GetKDFKey(ctx, psw.UserID, CurrKDF)
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

func (s *service) UpdatePassword(ctx context.Context, psw *Password) error {
	key, err := s.GetKDFKey(ctx, psw.UserID, CurrKDF)
	if err != nil {
		return err
	}

	data, err := NewPasswordEncrypt(psw, key)
	if err != nil {
		return err
	}

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

func (s *service) ReUpdateAllPasswords(ctx context.Context, userID uuid.UUID) error {
    // Get both KDF keys
    currKey, err := s.GetKDFKey(ctx, userID, CurrKDF)
    if err != nil {
        return err
    }
	log.Print(currKey)

    prevKey, err := s.GetKDFKey(ctx, userID, PrevKDF)
    if err != nil {
        return err
    }
	log.Print(prevKey)

	// empty pageParams
	pageParams := &httputils.Pagination{
		Index: "",
		// max int value
		Limit: "2147483647",
	}

	// retrieve encrypted passwords
	passwords, err := s.repo.GetAllPasswords(ctx, userID, pageParams)
	if err != nil {
		return err
	}

	// create transaction
	tx, err := s.repo.postgres.Begin(ctx)
	if err != nil {
		log.Error().Str("location", "ReUpdateAllPasswords").Msgf("%v: %v", userID, err)
		return err
	}
	
	// process passwords
	for _, password := range passwords {
		// decrypt rawData from encrypted	
		data, err := s.DecryptAndGetPsw(password, userID, prevKey)
		if err != nil {
			return err
		}
		
		// re-encrypt password
		newData, err := NewPasswordEncrypt(data, currKey)
		if err != nil {
			return err
		}

		// update password
		if err := s.repo.UpdatePassword(ctx, tx, newData); err != nil {
			return err
		}
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Error().Str("location", "ReUpdateAllPasswords").Msgf("%v: %v", userID, err)
		return err
	}

	// delete reset hash
	go func() {
		if err := s.repo.DeleteResetHash(ctx, userID); err != nil {
			log.Error().Str("location", "ReUpdateAllPasswords").Msgf("%v: %v", userID, err)
			return
		}

		log.Info().Str("location", "ReUpdateAllPasswords").Msgf("%v: reset hash deleted", userID)
	}()

	return nil
}

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
