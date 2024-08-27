package lruk

import (
	"errors"
)

var (
	ErrUnknownAccessType = errors.New("unknown access type")
	ErrInvalidFrameId    = errors.New("invalid frame id, should be greater than zero and less than replacer_size")
	ErrUnEvictableFrame  = errors.New("un evictable frame")
	ErrUnInitialized     = errors.New("un initialized")
	ErrNoEvictableFrame  = errors.New("no evictable frame")
	// ErrCapacityExceeded  = errors.New("capacity exceeded")

	ErrUnRemovableFrame = errors.New("un removable frame")
)
