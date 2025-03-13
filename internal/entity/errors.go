package entity

import "errors"

var (
	ErrUnmarshal               = errors.New("error unmarshalling request body")
	ErrEmptyTitle              = errors.New("title cannot be empty")
	ErrEmptyID                 = errors.New("id cannot be empty")
	ErrMissingRepeatParams     = errors.New("missing repeat parameters")
	ErrNoInterval              = errors.New("no interval specified")
	ErrMaxIntervalExceeded     = errors.New("maximum interval exceeded")
	ErrDateCalculation         = errors.New("error in date calculation")
	ErrUnsupportedRepeatFormat = errors.New("unsupported repeat format")
	ErrInvalidDate             = errors.New("invalid date")
	ErrTaskNotFound            = errors.New("task not found")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrAuthDisabled            = errors.New("authentication is disabled")
	ErrAuthenticationRequired  = errors.New("authentication is required")
)
