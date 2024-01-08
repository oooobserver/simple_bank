package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// Use this to ensure security
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// Customize error message
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto.NewV2(),
		[]byte(symmetricKey),
	}

	return maker, nil
}

// Implement the maker insterface
func (maker *PasetoMaker) CreateToken(name string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(name, role, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", payload, err
	}

	return token, payload, nil

}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !payload.Valid() {
		return nil, ErrExpiredToken
	}

	return payload, nil
}
