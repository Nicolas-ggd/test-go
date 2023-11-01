package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int
	Email           string
	Password        []byte
	ConfirmPassword string
}

type UserModel struct {
	DB *sql.DB
}

type UserModelInterface interface {
	Insert(email, password string) error
	Authentication(email, password string) (int, error)
	UserExists(id int) (bool, error)
}

func (us *UserModel) Insert(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO test_user (email, password)
	VALUES($1, $2)`
	_, err = us.DB.Exec(stmt, email, string(hashedPassword))
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) {
			if pgError.Code == "1062" && strings.Contains(pgError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (us *UserModel) Authentication(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM test_user WEHRE email = $1`

	err := us.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errInvalidCredentials
		} else {
			return 0, err
		}
	}

	// check if hashed password match plain text password or not
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, errInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (us *UserModel) UserExists(id int) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id = $1)`

	err := us.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
