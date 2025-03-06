package trie

import (
	"errors"
	"testing"
)

func TestNewTrie(t *testing.T) {
	_ = NewTrie()
}

func TestTrie_Put(t *testing.T) {
	t.Run("nil root", func(t *testing.T) {
		nilRootTrie := new(Trie)
		err := nilRootTrie.Put("key", "value")
		if !errors.Is(err, ErrEmptyTrie) {
			t.Error("expected ErrEmptyTrie")
		}
	})

	t.Run("empty key", func(t *testing.T) {
		trie := NewTrie()
		err := trie.Put("", "value")
		if !errors.Is(err, ErrEmptyKey) {
			t.Error("expected ErrEmptyKey")
		}
	})

	t.Run("existing key", func(t *testing.T) {
		trie := NewTrie()
		err := trie.Put("key", "value1")
		if err != nil {
			t.Error(err)
		}

		err = trie.Put("key", "value2")
		if !errors.Is(err, ErrKeyExists) {
			t.Error("expected ErrKeyExists")
		}
	})

	t.Run("normal case", func(t *testing.T) {
		trie := NewTrie()
		err := trie.Put("key1", "value1")
		if err != nil {
			t.Error(err)
		}

		err = trie.Put("key2", "value2")
		if err != nil {
			t.Error(err)
		}
	})
}

func TestTrie_Get(t *testing.T) {
	t.Run("nil root", func(t *testing.T) {
		nilRootTrie := new(Trie)
		_, err := nilRootTrie.Get("key")
		if !errors.Is(err, ErrEmptyTrie) {
			t.Error("expected ErrEmptyTrie")
		}
	})

	t.Run("empty key", func(t *testing.T) {
		trie := NewTrie()
		_, err := trie.Get("")
		if !errors.Is(err, ErrEmptyKey) {
			t.Error("expected ErrEmptyKey")
		}
	})

	t.Run("non existing key", func(t *testing.T) {
		trie := NewTrie()
		_, err := trie.Get("key")
		if !errors.Is(err, ErrKeyNotFound) {
			t.Error("expected ErrKeyNotFound")
		}

		err = trie.Put("key1", "value1")
		if err != nil {
			t.Error(err)
		}

		_, err = trie.Get("key")
		if !errors.Is(err, ErrKeyNotFound) {
			t.Error("expected ErrKeyNotFound")
		}
	})

	t.Run("normal case", func(t *testing.T) {
		trie := NewTrie()
		err := trie.Put("key", "value")
		if err != nil {
			t.Error(err)
		}

		val, err := trie.Get("key")
		if err != nil {
			t.Error(err)
		}
		if val != "value" {
			t.Error("expected value")
		}
	})
}

func TestTrie_Delete(t *testing.T) {
	t.Run("nil root", func(t *testing.T) {
		nilRootTrie := new(Trie)
		err := nilRootTrie.Delete("key")
		if !errors.Is(err, ErrEmptyTrie) {
			t.Error("expected ErrEmptyTrie")
		}
	})

	t.Run("empty key", func(t *testing.T) {
		trie := NewTrie()
		err := trie.Delete("")
		if !errors.Is(err, ErrEmptyKey) {
			t.Error("expected ErrEmptyKey")
		}
	})

	t.Run("non existing key", func(t *testing.T) {
		trie := NewTrie()
		err := trie.Delete("key")
		if !errors.Is(err, ErrKeyNotFound) {
			t.Error("expected ErrKeyNotFound")
		}

		err = trie.Put("key1", "value1")
		if err != nil {
			t.Error(err)
		}
		err = trie.Delete("key")
		if !errors.Is(err, ErrKeyNotFound) {
			t.Error("expected ErrKeyNotFound")
		}
	})

	t.Run("normal case", func(t *testing.T) {
		trie := NewTrie()
		err := trie.Put("key", "value")
		if err != nil {
			t.Error(err)
		}
		err = trie.Delete("key")
		if err != nil {
			t.Error(err)
		}
	})
}
