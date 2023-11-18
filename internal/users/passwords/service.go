package passwords

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/pbkdf2"

	"nestpass/pkg/httputils"
)

type service struct {
	repo *repository
}

func NewService(repo *repository) *service {
	return &service{repo: repo}
}

func (s *service) getKDFKey(ctx context.Context, userID uuid.UUID, kdf kdfType) ([]byte, error) {
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
			log.Error().Str("location", "getKDFKey").Msgf("%v: %v", userID, err)
			return nil, err
		}

		key = string(key64)
	}

	combinedInput := key + ":" + userID.String()
	kdfKey := pbkdf2.Key([]byte(combinedInput), kdfData.Salt, 4096, 32, sha256.New)
	return kdfKey, nil
}

func (s *service) GetAllPasswords(ctx context.Context, userID uuid.UUID, pageParams *httputils.Pagination) ([]*Password, error) {
	// retrieve kdf key
	key, err := s.getKDFKey(ctx, userID, CurrKDF)
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
		psw, err := password.Decrypt(userID, key)
		if err != nil {
			return nil, err
		}

		decryptedPsw = append(decryptedPsw, psw)
	}

	return decryptedPsw, nil
}

func (s *service) GetAllPasswordsByCategory(ctx context.Context, userID, categoryID uuid.UUID, pageParams *httputils.Pagination) ([]*Password, error) {
	// retrieve kdf key
	key, err := s.getKDFKey(ctx, userID, CurrKDF)
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
		psw, err := password.Decrypt(userID, key)
		if err != nil {
			return nil, err
		}

		decryptedPsw = append(decryptedPsw, psw)
	}

	return decryptedPsw, nil
}

func (s *service) GetPassword(ctx context.Context, passwordID, categoryID, userID uuid.UUID) (*Password, error) {
	// retrieve kdf key
	key, err := s.getKDFKey(ctx, userID, CurrKDF)
	if err != nil {
		return nil, err
	}

	// retrieve encrypted password
	password, err := s.repo.GetPassword(ctx, passwordID, categoryID, userID)
	if err != nil {
		return nil, err
	}

	// decrypt password
	psw, err := password.Decrypt(userID, key)
	if err != nil {
		return nil, err
	}

	return psw, nil
}

func (s *service) CreatePassword(ctx context.Context, psw *Password) error {
	key, err := s.getKDFKey(ctx, psw.UserID, CurrKDF)
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
	key, err := s.getKDFKey(ctx, psw.UserID, CurrKDF)
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
	currKeyCh, prevKeyCh := make(chan kdfKeyRes, 1), make(chan kdfKeyRes, 1)

	// Fetch currKey and prevKey concurrently
	go func() {
		key, err := s.getKDFKey(ctx, userID, CurrKDF)
		currKeyCh <- kdfKeyRes{Key: key, Err: err}
	}()
	go func() {
		key, err := s.getKDFKey(ctx, userID, PrevKDF)
		prevKeyCh <- kdfKeyRes{Key: key, Err: err}
	}()

	currKeyRes := <-currKeyCh
	if currKeyRes.Err != nil {
		return currKeyRes.Err
	}

	prevKeyRes := <-prevKeyCh
	if prevKeyRes.Err != nil {
		return prevKeyRes.Err
	}

	// retrieve encrypted passwords
	passwords, err := s.repo.GetAllPasswordsNonPaged(ctx, userID)
	if err != nil {
		return err
	}

	n := len(passwords)
	wg := new(sync.WaitGroup)
	chunkSize := 4
	if n < chunkSize {
		chunkSize = n
	}

	log.Info().Str("location", "ReUpdateAllPasswords").Msgf("%v: %v passwords to rehash now starting...", userID, n)
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)

		go func(chunk []*PasswordEncrypt) {
			defer wg.Done()

			// create transaction
			tx, err := s.repo.postgres.Begin(ctx)
			if err != nil {
				log.Error().Str("location", "ReUpdateAllPasswords").Msgf("%v: %v", userID, err)
				return
			}
			defer tx.Rollback(ctx)

			for _, password := range chunk {
				data, err := password.Decrypt(userID, prevKeyRes.Key)
				if err != nil {
					log.Error().Str("location", "ReUpdateAllPasswords").Msgf("%v: %v", userID, err)
					return
				}

				newData, err := NewPasswordEncrypt(data, currKeyRes.Key)
				if err != nil {
					log.Error().Str("location", "ReUpdateAllPasswords").Msgf("%v: %v", userID, err)
					return
				}

				if err := s.repo.UpdatePassword(ctx, tx, newData); err != nil {
					log.Error().Str("location", "ReUpdateAllPasswords").Msgf("%v: %v", userID, err)
					return
				}
			}

			// commit transaction
			if err := tx.Commit(ctx); err != nil {
				log.Error().Str("location", "ReUpdateAllPasswords").Msgf("%v: %v", userID, err)
				return
			}
		}(passwords[i:end])
	}

	wg.Wait()
	log.Info().Str("location", "ReUpdateAllPasswords").Msgf("%v: %v passwords rehashed", userID, n)

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
