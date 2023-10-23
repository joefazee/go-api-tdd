package domain

import "time"

type JWTPayload struct {
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"iat"`
	IssuedAt  time.Time `json:"exp"`
	Token     string    `json:"token"`
}

func (p *JWTPayload) Valid() error {

	if time.Now().After(p.ExpiresAt) {
		return ErrExpiredToken
	}

	if p.UserID == 0 {
		return ErrUserNotFound
	}

	return nil
}
