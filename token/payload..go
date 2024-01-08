package token

import (
	"time"

	"github.com/google/uuid"
)

// The message of one token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(name, role string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		id,
		name,
		role,
		time.Now(),
		time.Now().Add(duration),
	}

	return payload, nil
}

// Check if this token is valid
func (payload *Payload) Valid() bool {
	expired := payload.ExpiredAt
	return !time.Now().After(expired)
}
