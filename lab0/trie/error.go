package trie

import (
	"errors"
)

var (
	ErrEmptyTrie   = errors.New("empty trie")
	ErrEmptyKey    = errors.New("empty key") // empty key is not allowed
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExists   = errors.New("key exists")
)
