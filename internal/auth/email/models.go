package email

import (
	"encoding/json"
	"io"

	"github.com/rs/zerolog/log"
)

// Body for twofa data from redis.
type TwofaBody struct {
	Code       string
	Retries    int
	UserStatus string
}

// Deserialize the json data into the struct.
func (tfa *TwofaBody) Deserialize(data string) error {
	if err := json.Unmarshal([]byte(data), tfa); err != nil {
		log.Err(err).Str("location", "TwofaBody.Unmarshal()").Msg("failed to unmarshal json")
		return err
	}

	return nil
}

// Serialize the struct into json data and return it as a string.
func (tfa *TwofaBody) Serialize() (string, error) {
	data, err := json.Marshal(tfa)

	if err != nil {
		log.Err(err).Str("location", "TwofaBody.Serialize()").Msg("failed to marshal json")
		return "", err
	}

	return string(data), nil
}

// Body for twofa request data.
type TokenInput struct {
	Token string `json:"token"`
}

// Deserialize the json data into the struct.
func (ti *TokenInput) Deserialize(data io.ReadCloser) error {
	if err := json.NewDecoder(data).Decode(&ti); err != nil {
		log.Err(err).Str("location", "TokenInput.Unmarshal()").Msg("failed to unmarshal json")
		return err
	}

	return nil
}
