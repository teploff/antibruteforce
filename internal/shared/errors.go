package shared

import "errors"

var (
	ErrNotFound     = errors.New("bucket doesn't exist")
	ErrAlreadyExist = errors.New("bucket already exists")
)
