package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// create a new payload with a username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, err
}

func (payload *Payload) Valid() error {
	fmt.Printf("DEBUG: Valid() method called\n")

	now := time.Now()
	fmt.Printf("DEBUG: Current time: %v\n", now)
	fmt.Printf("DEBUG: Token expires: %v\n", payload.ExpiredAt)
	fmt.Printf("DEBUG: Is current time after expiration? %v\n", now.After(payload.ExpiredAt))
	fmt.Printf("DEBUG: Time difference: %v\n", payload.ExpiredAt.Sub(now))

	if now.After(payload.ExpiredAt) {
		fmt.Printf("DEBUG: Token has expired!\n")
		return ErrExpiredToken
	}

	fmt.Printf("DEBUG: Token is valid!\n")
	return nil
}
