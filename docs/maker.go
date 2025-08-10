package docs

import "time"

// Interface for managing tokens
type Maker interface {
	// Creates token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)
	// Checks if the token is valid
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
