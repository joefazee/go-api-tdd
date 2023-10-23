package postgres

import (
	"errors"
	"github.com/joefazee/go-api-tdd/pkg/domain"
	"testing"
)

var (
	oldSqlCreateUser      = sqlCreateUser
	oldsqlDeleteUserByID  = sqlDeleteUserByID
	oldsqlFindUserByEmail = sqlFindUserByEmail
	oldsqlFindUserByID    = sqlFindUserByID
)

func TestCreateUser(t *testing.T) {

	pStore := NewPostgresStore(testDB)
	oldPassword := "password"

	user := &domain.User{
		Email:    "test@test.com",
		Password: oldPassword,
		Name:     "John Doe",
	}

	createdUser, err := pStore.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	if createdUser.ID == 0 {
		t.Errorf("want id not to be zero")
	}

	if user.Name != createdUser.Name {
		t.Errorf("expected %q; got %q", user.Name, createdUser.Name)
	}

	if createdUser.Password == oldPassword {
		t.Error("password was not hashed")
	}

	sqlCreateUser = "invalid"
	_, err = pStore.CreateUser(user)
	if err == nil {
		t.Errorf("expected error not to be nil for invalid CreateUser sql")
	}
	sqlCreateUser = oldSqlCreateUser

	err = pStore.DeleteUserByID(createdUser.ID)
	if err != nil {
		t.Errorf("expected nil error during DeleteUserByID; got %q", err)
	}

	sqlDeleteUserByID = "invalid"
	err = pStore.DeleteUserByID(createdUser.ID)
	if err == nil {
		t.Errorf("expected error not to be nil for invalid DeleteUserByID sql")
	}
	sqlDeleteUserByID = oldsqlDeleteUserByID

}

func TestFindUserByEmail(t *testing.T) {
	pStore := NewPostgresStore(testDB)

	user := &domain.User{
		Email:    "testwwww@test.com",
		Password: "password",
		Name:     "John Doe",
	}

	createdUser, err := pStore.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	uByEmail, err := pStore.FindUserByEmail(createdUser.Email)
	if err != nil {
		t.Errorf("expected no error; got %q", err)
	}

	if uByEmail.Email != createdUser.Email {
		t.Errorf("expect %q; got %q", createdUser.Email, uByEmail.Email)
	}

	_, err = pStore.FindUserByEmail("invalid@email")
	if err == nil {
		t.Errorf("want error; got nil for invalid email")
	}

	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("want domain.ErrUserNotFound error; got %q", err)
	}

	sqlFindUserByEmail = "invalid"
	_, err = pStore.FindUserByEmail(createdUser.Email)
	if err == nil {
		t.Errorf("want error; got nil for invalid FindUserByEmail sql")
	}
	sqlFindUserByEmail = oldsqlFindUserByEmail

	_ = pStore.DeleteUserByID(createdUser.ID)
}

func TestFindUserByID(t *testing.T) {
	pStore := NewPostgresStore(testDB)

	user := &domain.User{
		Email:    "testwwww@test.com",
		Password: "password",
		Name:     "John Doe",
	}

	createdUser, err := pStore.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	uByID, err := pStore.FindUserByID(createdUser.ID)
	if err != nil {
		t.Errorf("expected no error; got %q", err)
	}

	if uByID.ID != createdUser.ID {
		t.Errorf("expect %q; got %q", createdUser.ID, uByID.ID)
	}

	_, err = pStore.FindUserByID(-1)
	if err == nil {
		t.Errorf("want error; got nil for invalid id")
	}

	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("want domain.ErrUserNotFound error; got %q", err)
	}

	sqlFindUserByID = "invalid"
	_, err = pStore.FindUserByID(createdUser.ID)
	if err == nil {
		t.Errorf("want error; got nil for invalid FindUserByID sql")
	}
	sqlFindUserByID = oldsqlFindUserByID

	_ = pStore.DeleteUserByID(createdUser.ID)

	_ = pStore.DeleteAllUsers()
}
