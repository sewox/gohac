package domain

import "errors"

// Domain-specific errors
var (
	ErrPageNotFound      = errors.New("page not found")
	ErrPageAlreadyExists = errors.New("page with this slug already exists")
	ErrInvalidSlug       = errors.New("invalid slug format")
	ErrInvalidStatus     = errors.New("invalid page status")

	ErrBlockMissingID   = errors.New("block missing required id field")
	ErrBlockMissingType = errors.New("block missing required type field")
	ErrBlockMissingData = errors.New("block missing required data field")
	ErrInvalidBlockType = errors.New("invalid block type")
)
