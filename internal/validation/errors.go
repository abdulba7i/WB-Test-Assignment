package validation

import "errors"

var (
	ErrNameTooLong = errors.New("name too long")
	ErrInvalidAge  = errors.New("age must be between 18 and 100")
)
