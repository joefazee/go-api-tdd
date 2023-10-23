package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrExpiredToken = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)
