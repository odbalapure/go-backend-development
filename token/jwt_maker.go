package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	minSecretKeySize = 32
)

type JWTMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	maker := &JWTMaker{secretKey}

	return maker, nil
}

// Creates token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	// NOTE: The Payload type must implement the `Valid()` method
	// `NewWithClaims` creates a new token with the payload and the signing method
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// `jwtToken.SignedString` signs the token with the secret key in the `JWTMaker` struct
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, nil
}

// Checks if the token is valid
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// `keyFunc` is a function that returns the secret key for the token
	// It is used to verify the token signature
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Verify the token signature
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		// Return the secret key for the token
		return []byte(maker.secretKey), nil
	}

	// `ParseWithClaims` parses the token and returns the payload
	// It is used to verify the token signature
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		// The ok variable is used to check if the error is the same as the one we want to check for
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// jwtToken.Claims returns the payload of the token
	// It is used to verify the token signature
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
