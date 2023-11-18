package categories

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Category struct {
	CategoryID  uuid.UUID `json:"category_id,omitempty"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Scanable interface {
	Scan(row ...any) error
}

func (c *Category) Scan(row Scanable) error {
	return row.Scan(&c.CategoryID, &c.UserID, &c.Name, &c.Description)
}

func (c *Category) Deserialize(data io.ReadCloser) error {
	if err := json.NewDecoder(data).Decode(c); err != nil {
		log.Error().Str("location", "Deserialize").Msg(err.Error())
		return err
	}

	if err := validator.New().Struct(c); err != nil {
		log.Error().Str("location", "Deserialize").Msg(err.Error())
		return err
	}

	return nil
}

func New(name, description string, userID uuid.UUID) *Category {
	return &Category{
		CategoryID:  uuid.New(),
		UserID:      userID,
		Name:        name,
		Description: description,
	}
}
