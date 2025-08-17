package token

import "time"

// Interface for managing tokens
type Maker interface {
	// Creates token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	// Checks if the token is valid
	VerifyToken(token string) (*Payload, error)
}
