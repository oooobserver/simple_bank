package token

import "time"

// Manage the token
type Maker interface {
	// Creates a new token
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}
