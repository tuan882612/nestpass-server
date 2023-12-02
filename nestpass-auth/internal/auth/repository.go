package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/tuan882612/apiutils"

	"project/internal/database"
)

// Base authentication repository.
type Repository struct {
	db *pgxpool.Pool
}

// Constructor for the base authentication repository.
func NewRepository(databases *database.DataAccess) *Repository {
	return &Repository{db: databases.Postgres}
}

// Retrieves the user's uuid and password from the database if the user exists.
func (r *Repository) GetUserCredentials(ctx context.Context, email string) (*User, error) {
	// initialize credential variables
	user := &User{}
	row := r.db.QueryRow(ctx, UserCredsQuery, email)

	// scan the row and check for errors
	if err := row.Scan(&user.UserID, &user.Password, &user.UserStatus); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apiutils.NewErrNotFound("user not found")
		}

		log.Error().Str("location", "GetUserCredentials").Msgf("failed to get user credentials: %v", err)
		return nil, err
	}

	return user, nil
}

// Retrieves the user's uuid and password from the database if the user exists.
func (r *Repository) GetUserPassword(ctx context.Context, userID uuid.UUID) (string, error) {
	// initialize credential variables
	var password string
	row := r.db.QueryRow(ctx, GetUserPasswordQuery, userID)

	// scan the row and check for errors
	if err := row.Scan(&password); err != nil {
		if err == pgx.ErrNoRows {
			return "", apiutils.NewErrNotFound("user not found")
		}

		log.Error().Str("location", "GetUserPassword").Msgf("failed to get user password: %v", err)
		return "", err
	}

	return password, nil
}

// Adds a new user to the database.
func (r *Repository) AddUser(ctx context.Context, tx pgx.Tx, input *Register) error {
	_, err := tx.Exec(ctx, AddUserQuery,
		&input.UserID,
		&input.Email,
		&input.Name,
		&input.Password,
		&input.Registered,
		&input.UserStatus,
		&input.Salt,
	)

	// checking for errors
	if err != nil {
		// initializing pgx error and checking for duplicate key error
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apiutils.NewErrConflict("user already exists")
		}

		log.Error().Str("location", "AddUser").Msgf("%v: failed to add user: %v", input.UserID, err)
		return err
	}

	return nil
}

// Updates the user's status to "active".
func (r *Repository) UpdateUserStatus(ctx context.Context, userID uuid.UUID) error {
	if _, err := r.db.Exec(ctx, UpdateUserStatusQuery, userID); err != nil {
		log.Error().Str("location", "UpdateUserStatus").Msgf("%v: failed to update user status: %v", userID, err)
		return err
	}

	return nil
}

// Updates the user's password.
func (r *Repository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, password string) error {
	if _, err := r.db.Exec(ctx, UpdateUserPasswordQuery, userID, password); err != nil {
		log.Error().Str("location", "UpdateUserPassword").Msgf("%v: failed to update user password: %v", userID, err)
		return err
	}

	return nil
}

// Starts a new postgres transaction.
func (r *Repository) StartTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error().Str("location", "startTx").Msgf("failed to start transaction: %v", err)
		return nil, err
	}

	return tx, nil
}
