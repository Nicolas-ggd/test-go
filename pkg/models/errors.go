package models

import "errors"

var (
	errNoRecord = errors.New("models: no matching record found")

	errInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")
)
