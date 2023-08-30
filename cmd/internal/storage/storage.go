package storage

import "errors"

// TODO: Refactor error handling
var (
	ErrSegmentNotFound = errors.New("url not found")
	ErrSegmentExists   = errors.New("url exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserHasSegment  = errors.New("user already has such segment")
)
