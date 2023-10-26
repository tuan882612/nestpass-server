package categories

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Category struct {
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
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

type CategoryResp struct {
	CategoryID  uuid.UUID `json:"category_id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Scanable interface {
	Scan(row ...any) error
}

func (c *CategoryResp) Scan(row Scanable) error {
	return row.Scan(&c.CategoryID, &c.UserID, &c.Name, &c.Description)
}

func (c *CategoryResp) Deserialize(data io.ReadCloser) error {
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

func New(name, description string, userID uuid.UUID) *CategoryResp {
	return &CategoryResp{
		CategoryID:  uuid.New(),
		UserID:      userID,
		Name:        name,
		Description: description,
	}
}
