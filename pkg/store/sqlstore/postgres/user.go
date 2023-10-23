package postgres

import (
	"database/sql"
	"errors"
	"github.com/joefazee/go-api-tdd/pkg/common"
	"github.com/joefazee/go-api-tdd/pkg/domain"
)

var (
	sqlCreateUser = `INSERT INTO users 
    			(name, email, password) VALUES ($1, $2, $3) 
    		   RETURNING id, name, email, password`

	sqlDeleteUserByID = `DELETE FROM users WHERE  id = $1`

	sqlFindUserByEmail = `SELECT id, name, email, password FROM users WHERE email = $1`

	sqlFindUserByID = `SELECT id, name, email, password FROM users WHERE id = $1`

	sqlDeleteAllUsers = `DELETE FROM users`
)

func (q *postgresStore) CreateUser(user *domain.User) (*domain.User, error) {

	user.Password, _ = common.PasswordHash(user.Password)

	err := q.db.QueryRow(sqlCreateUser, user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (q *postgresStore) DeleteUserByID(ID int64) error {
	_, err := q.db.Exec(sqlDeleteUserByID, ID)
	if err != nil {
		return err
	}
	return nil
}

func (q *postgresStore) FindUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	err := q.db.QueryRow(sqlFindUserByEmail, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (q *postgresStore) FindUserByID(ID int64) (*domain.User, error) {
	user := &domain.User{}

	err := q.db.QueryRow(sqlFindUserByID, ID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (q *postgresStore) DeleteAllUsers() error {
	_, err := q.db.Exec(sqlDeleteAllUsers)
	if err != nil {
		return err
	}
	return nil
}
