package passwords

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/pbkdf2"
)

type ctxKey string

const CtxKDFKey ctxKey = "kdf_key"

type KDFType string

const (
	CurrKDF KDFType = "curr"
	PrevKDF KDFType = "prev"
)

// Generates a key using PBKDF2 and returns the key
func KeyDerivation(passwordHash, userID string, salt []byte) []byte {
	combinedInput := passwordHash + ":" + userID
	return pbkdf2.Key([]byte(combinedInput), salt, 4096, 32, sha256.New)
}

// Encrypts the data using AES-256-GCM and returns the encrypted data
func Encrypt(data []byte, key []byte) (nonce, ciphertext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error().Str("location", "encrypt").Msg(err.Error())
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error().Str("location", "encrypt").Msg(err.Error())
		return nil, nil, err
	}

	nonce = make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Error().Str("location", "encrypt").Msg(err.Error())
		return nil, nil, err
	}

	ciphertext = aesgcm.Seal(nil, nonce, data, nil)
	return nonce, ciphertext, nil
}

// Decrypts the encrypted data using AES-256-GCM and returns the data
func Decrypt(nonce, ciphertext, key []byte) (data *PswData, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error().Str("location", "decrypt").Msg(err.Error())
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Error().Str("location", "decrypt").Msg(err.Error())
		return nil, err
	}

	decrypted, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Error().Str("location", "decrypt").Msg(err.Error())
		return nil, err
	}

	data = &PswData{}
	if err := json.Unmarshal(decrypted, data); err != nil {
		log.Error().Str("location", "decrypt").Msg(err.Error())
		return nil, err
	}

	return data, nil
}
