package repository

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEndpointNotFound = errors.New("endpoint not found")
)
