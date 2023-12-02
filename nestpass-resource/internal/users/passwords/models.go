package passwords

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type kdfType string

const (
	currKDF kdfType = "curr"
	prevKDF kdfType = "prev"
)

// New 128 bit Galois Counter Mode wrapped block cipher
func newGCMBlock(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error().Str("location", "encrypt").Msg(err.Error())
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error().Str("location", "encrypt").Msg(err.Error())
		return nil, err
	}

	return aesgcm, nil
}

type kdfKeyRes struct {
	Key []byte
	Err error
}

type kdfData struct {
	PswHash string `json:"password"`
	Salt    []byte `json:"salt"`
}

func (k *kdfData) Scan(row pgx.Row) error {
	return row.Scan(&k.PswHash, &k.Salt)
}

type pswData struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type PasswordEncrypt struct {
	PasswordID uuid.UUID `json:"password_id"`
	UserID     uuid.UUID `json:"user_id"`
	CategoryID uuid.UUID `json:"category_id"`
	Website    string    `json:"website"`
	Nonce      []byte    `json:"nonce"`
	Encrypted  []byte    `json:"encrypted"`
}

func (p *PasswordEncrypt) Scan(row pgx.Row) error {
	return row.Scan(
		&p.PasswordID,
		&p.UserID,
		&p.CategoryID,
		&p.Website,
		&p.Nonce,
		&p.Encrypted,
	)
}

func (p *PasswordEncrypt) Decrypt(userID uuid.UUID, key []byte) (*Password, error) {
	// create aesgcm block
	aesgcm, err := newGCMBlock(key)
	if err != nil {
		return nil, err
	}

	// decrypt data
	decrypted, err := aesgcm.Open(nil, p.Nonce, p.Encrypted, nil)
	if err != nil {
		log.Error().Str("location", "decrypt").Msg(err.Error())
		return nil, err
	}

	// deserialize decrypted data
	data := &pswData{}
	if err := json.Unmarshal(decrypted, data); err != nil {
		log.Error().Str("location", "decrypt").Msg(err.Error())
		return nil, err
	}

	// create Password
	return &Password{
		PasswordID:  p.PasswordID,
		UserID:      p.UserID,
		CategoryID:  p.CategoryID,
		Website:     p.Website,
		Username:    data.Username,
		Password:    data.Password,
		Description: data.Description,
	}, nil
}

func NewPasswordEncrypt(psw *Password, dKey []byte) (*PasswordEncrypt, error) {
	// pull out data from Password
	data := &pswData{
		Username:    psw.Username,
		Password:    psw.Password,
		Description: psw.Description,
	}

	// serialize the data
	rawData, err := json.Marshal(data)
	if err != nil {
		log.Error().Str("location", "NewPasswordEncrypt").Msg(err.Error())
		return nil, err
	}

	// encrypt the data
	aesgcm, err := newGCMBlock(dKey)
	if err != nil {
		return nil, err
	}

	// create nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Error().Str("location", "encrypt").Msg(err.Error())
		return nil, err
	}

	// encrypt data
	encrypted := aesgcm.Seal(nil, nonce, rawData, nil)

	var passwordID uuid.UUID
	if psw.PasswordID != uuid.Nil {
		passwordID = psw.PasswordID
	} else {
		passwordID = uuid.New()
	}

	return &PasswordEncrypt{
		PasswordID: passwordID,
		UserID:     psw.UserID,
		CategoryID: psw.CategoryID,
		Website:    psw.Website,
		Nonce:      nonce,
		Encrypted:  encrypted,
	}, nil
}

type Password struct {
	PasswordID  uuid.UUID `json:"password_id,omitempty"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Website     string    `json:"website" validate:"required"`
	Username    string    `json:"username" validate:"required"`
	Password    string    `json:"password" validate:"required"`
	Description string    `json:"description"`
}

func (p *Password) Deserialize(data io.ReadCloser) error {
	if err := json.NewDecoder(data).Decode(p); err != nil {
		log.Error().Str("location", "deserialize").Msg(err.Error())
		return err
	}

	if err := validator.New().Struct(p); err != nil {
		log.Error().Str("location", "deserialize").Msg(err.Error())
		return err
	}

	return nil
}

func NewPassword(data *pswData, website string, userID, categoryID uuid.UUID) *Password {
	return &Password{
		UserID:      userID,
		CategoryID:  categoryID,
		Website:     website,
		Username:    data.Username,
		Password:    data.Password,
		Description: data.Description,
	}
}
