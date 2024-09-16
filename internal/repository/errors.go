package repository

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrTenderNotFound = errors.New("tender not found")
	ErrNoAccess       = errors.New("no access")
)
