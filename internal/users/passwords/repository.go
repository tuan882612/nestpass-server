package passwords

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"nestpass/pkg/httputils"
)

type repository struct {
	postgres *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) *repository {
	return &repository{postgres: pg}
}

func (r *repository) GetKDFData(ctx context.Context, userID uuid.UUID) (*KDFData, error) {
	kdf := &KDFData{}

	// retrieves nonce and salt
	row := r.postgres.QueryRow(ctx, GetKDFDataQuery, userID)
	if err := kdf.Scan(row); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user id was not provided")
		}

		log.Error().Str("location", "GetKDFKey").Msgf("%v: %v", userID, err)
		return nil, err
	}

	return kdf, nil
}

func (r *repository) GetAllPasswords(
	ctx context.Context,
	userID uuid.UUID,
	key []byte,
	pageParams *httputils.Pagination) ([]*PasswordEncrypt, error) {

	rows, err := r.postgres.Query(ctx, GetAllPasswordsQuery, userID, pageParams.Index, pageParams.Limit)
	if err != nil {
		log.Error().Str("location", "GetAllPasswords").Msgf("%v: %v", userID, err)
		return nil, err
	}

	// retrieve passwords
	passwords := []*PasswordEncrypt{}
	for rows.Next() {
		password := &PasswordEncrypt{}
		if err := password.Scan(rows); err != nil {
			log.Error().Str("location", "GetAllPasswords").Msgf("%v: %v", userID, err)
			return nil, err
		}

		passwords = append(passwords, password)
	}

	return passwords, nil
}

func (r *repository) GetAllPasswordsByCategory(
	ctx context.Context,
	userID, categoryID uuid.UUID,
	key []byte,
	pageParms *httputils.Pagination) ([]*PasswordEncrypt, error) {

	rows, err := r.postgres.Query(ctx, GetAllPasswordsByCategoryQuery, userID, categoryID, pageParms.Index, pageParms.Limit)
	if err != nil {
		log.Error().Str("location", "GetAllPasswordsByCategory").Msgf("%v: %v", userID, err)
		return nil, err
	}

	// retrieve passwords
	passwords := []*PasswordEncrypt{}
	for rows.Next() {
		password := &PasswordEncrypt{}
		if err := password.Scan(rows); err != nil {
			log.Error().Str("location", "GetAllPasswordsByCategory").Msgf("%v: %v", userID, err)
			return nil, err
		}

		passwords = append(passwords, password)
	}

	return passwords, nil
}

func (r *repository) GetAllPasswordsNonPaged(ctx context.Context, userID uuid.UUID) ([]*PasswordEncrypt, error) {
	rows, err := r.postgres.Query(ctx, GetAllPasswordsNonPagedQuery, userID)
	if err != nil {
		log.Error().Str("location", "GetAllPasswordsNonPaged").Msgf("%v: %v", userID, err)
		return nil, err
	}

	// retrieve passwords
	passwords := []*PasswordEncrypt{}
	for rows.Next() {
		password := &PasswordEncrypt{}
		if err := password.Scan(rows); err != nil {
			log.Error().Str("location", "GetAllPasswordsNonPaged").Msgf("%v: %v", userID, err)
			return nil, err
		}

		passwords = append(passwords, password)
	}

	return passwords, nil
}

func (r *repository) GetPassword(ctx context.Context, passwordID, categoryID, userID uuid.UUID) (*PasswordEncrypt, error) {
	password := &PasswordEncrypt{}

	row := r.postgres.QueryRow(ctx, GetPasswordQuery, userID, passwordID, categoryID)
	if err := password.Scan(row); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apiutils.NewErrNotFound("password not found")
		}

		log.Error().Str("location", "GetPassword").Msgf("%v: %v", userID, err)
		return nil, err
	}

	return password, nil
}

func (r *repository) CreatePassword(ctx context.Context, tx pgx.Tx, data *PasswordEncrypt) error {
	_, err := tx.Exec(ctx, CreatePasswordQuery,
		&data.PasswordID,
		&data.UserID,
		&data.CategoryID,
		&data.Website,
		&data.Nonce,
		&data.Encrypted,
	)

	if err != nil {
		log.Error().Str("location", "CreatePassword").Msgf("%v: %v", &data.UserID, err)
		return err
	}

	return nil
}

func (r *repository) UpdatePassword(ctx context.Context, tx pgx.Tx, data *PasswordEncrypt) error {
	_, err := tx.Exec(ctx, UpdatePasswordQuery,
		&data.Website,
		&data.Nonce,
		&data.Encrypted,
		&data.PasswordID,
		&data.CategoryID,
		&data.UserID,
	)

	if err != nil {
		log.Error().Str("location", "UpdatePassword").Msgf("%v: %v", data.UserID, err)
		return err
	}

	return nil
}

func (r *repository) DeletePassword(ctx context.Context, tx pgx.Tx, userID, passwordID, categoryID uuid.UUID) error {
	_, err := tx.Exec(ctx, DeletePasswordQuery, passwordID, userID, categoryID)
	if err != nil {
		log.Error().Str("location", "DeletePassword").Msgf("%v: %v", userID, err)
		return err
	}

	return nil
}
