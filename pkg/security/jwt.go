package security

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/joefazee/go-api-tdd/pkg/domain"
	"time"
)

type jwtToken struct {
	key string
}

// NewJWT create new instance of jwtToken with a key.
func NewJWT(key string) (domain.JWT, error) {
	if len(key) < 32 {
		return nil, errors.New("key too short. must be more than 32 characters")
	}

	return &jwtToken{key: key}, nil
}

func (j *jwtToken) CreateToken(user domain.User, duration time.Duration) (*domain.JWTPayload, error) {

	now := time.Now()

	payload := &domain.JWTPayload{
		UserID:    user.ID,
		ExpiresAt: now.Add(duration),
		IssuedAt:  now,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(j.key))
	if err != nil {
		return nil, err
	}
	payload.Token = tokenString

	return payload, nil

}

func (j *jwtToken) VerifyToken(tokenString string) (*domain.JWTPayload, error) {

	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTPayload{}, j.keyFunc)
	if err != nil {
		var terr *jwt.ValidationError
		ok := errors.As(err, &terr)
		if ok && errors.Is(terr.Inner, domain.ErrExpiredToken) {
			return nil, domain.ErrExpiredToken
		}
		return nil, domain.ErrInvalidToken
	}

	payload, ok := token.Claims.(*domain.JWTPayload)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	payload.Token = tokenString

	return payload, nil

}

func (j *jwtToken) keyFunc(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	return []byte(j.key), nil
}
