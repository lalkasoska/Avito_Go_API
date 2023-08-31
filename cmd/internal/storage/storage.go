package storage

import "errors"

var (
	ErrSegmentNotFound = errors.New("url not found")
	ErrSegmentExists   = errors.New("url exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserHasSegment  = errors.New("user already has such segment")
)
