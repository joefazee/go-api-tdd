package domain

import "time"

type Store interface {
	CreateUser(user *User) (*User, error)
	DeleteUserByID(ID int64) error
	DeleteAllUsers() error
	FindUserByID(ID int64) (*User, error)
	FindUserByEmail(email string) (*User, error)

	CreatePost(post *Post) (*Post, error)
	GetAllUserPosts(userID int64) ([]Post, error)
	DeleteAllPosts() error
	DeletePostByID(ID int64) error
}

type JWT interface {
	CreateToken(user User, duration time.Duration) (*JWTPayload, error)
	VerifyToken(tokenString string) (*JWTPayload, error)
}
