package users

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User struct {
	UserID     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Registered time.Time `json:"registered"`
	UserStatus string    `json:"user_status"`
}

func (u *User) Scan(row pgx.Row) error {
	return row.Scan(&u.UserID, &u.Email, &u.Name, &u.Registered, &u.UserStatus)
}
