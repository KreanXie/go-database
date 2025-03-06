package trie

import (
	"errors"
)

var (
	ErrEmptyTrie = errors.New("empty trie")

	// ErrEmptyKey empty key is not allowed
	ErrEmptyKey    = errors.New("empty key")
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExists   = errors.New("key exists")
)
