package common

import (
	"testing"
)

func TestPasswordHash(t *testing.T) {

	password := "password"
	hashed, err := PasswordHash(password)
	if err != nil {
		t.Fatal(err)
	}

	if len(hashed) == 0 {
		t.Errorf("want hash; got %q", hashed)
	}

	if hashed == password {
		t.Errorf("password was not hashed; got %q", hashed)
	}
}

func TestPasswordHash_Error(t *testing.T) {

	longPass := make([]byte, 73)

	_, err := PasswordHash(string(longPass))
	if err == nil {
		t.Errorf("expected an error;got nil")
	}

}

func TestCheckPassword(t *testing.T) {
	password := "password"

	hashed, err := PasswordHash(password)
	if err != nil {
		t.Fatal(err)
	}

	err = CheckPassword(password, hashed)
	if err != nil {
		t.Errorf("password verification failed; got %q", err)
	}
}
