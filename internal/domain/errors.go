package domain

import "errors"

var (
	ErrTicketNotFound = errors.New("rating ticket not found")
	ErrTicketExpired  = errors.New("rating ticket expired")
	ErrTicketExists   = errors.New("rating for this direction already exists")
	ErrInvalidStar    = errors.New("stars must be between 1 and 5")
)
