package service

import "errors"

// Sentinel errors returned by UserService methods.
var (
	// ErrEmailAlreadyTaken is returned when registering with an email that is already in use.
	ErrEmailAlreadyTaken = errors.New("email already taken")
	// ErrInvalidCredentials is returned when login credentials do not match.
	ErrInvalidCredentials = errors.New("invalid email or password")
)
