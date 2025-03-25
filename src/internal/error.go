package internal

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

	ErrEmptyTrie = errors.New("empty trie")

	// ErrEmptyKey empty key is not allowed
	ErrEmptyKey    = errors.New("empty key")
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExists   = errors.New("key exists")
)
