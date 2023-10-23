package security

import (
	"github.com/joefazee/go-api-tdd/pkg/common"
	"github.com/joefazee/go-api-tdd/pkg/domain"
	"testing"
	"time"
)

func TestJwtToken(t *testing.T) {
	_, err := NewJWT("short-key")

	if err == nil {
		t.Error("NewJWT should return an error when key is too short")
	}
}

func TestJWTToken(t *testing.T) {

	key := common.RandomString(32)

	newJWT, err := NewJWT(key)
	if err != nil {
		t.Error("expected error to be nil")
	}

	user := domain.User{
		ID:    1,
		Email: "test@test.com",
	}

	payload, err := newJWT.CreateToken(user, 1*time.Minute)
	if err != nil {
		t.Error("expected error to be nil")
	}

	if payload.UserID != user.ID {
		t.Errorf("want %q; got %q", user.ID, payload.UserID)
	}

	if len(payload.Token) == 0 {
		t.Error("token string len is zero")
	}

	if payload.ExpiresAt.Before(time.Now()) {
		t.Error("CreateToken should return a token that expires in future")
	}

	_, err = newJWT.VerifyToken(payload.Token + "invalid")
	if err == nil {
		t.Error("VerifyToken should return error for invalid token")
	}

	expiredToken, err := newJWT.CreateToken(user, -1*time.Minute)
	if err != nil {
		t.Error("CreateToken should not return error")
	}

	_, err = newJWT.VerifyToken(expiredToken.Token)
	if err == nil {
		t.Error("VerifyToken should return error for expired token")
	}

}
