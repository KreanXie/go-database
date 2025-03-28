package internal

// Trie class
type Trie struct {
	root *TrieNode
}

// TrieNode is a Trie Node
type TrieNode struct {
	key byte

	// if val is nil, we consider it's not the end of the key
	// in someways, this can replace isEnd attribute
	val    any
	childs map[byte]*TrieNode
}

// NewTrie returns an initialized trie
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			key:    ' ',
			val:    nil,
			childs: make(map[byte]*TrieNode),
		},
	}
}

// Put does insert a kv into trie
func (t *Trie) Put(key string, val any) error {
	switch {
	case len(key) == 0:
		return ErrEmptyKey
	case t.root == nil:
		return ErrEmptyTrie
	}

	cur := t.root
	i := 0
	for i < len(key) {
		if _, ok := cur.childs[key[i]]; !ok {
			cur.childs[key[i]] = &TrieNode{
				key:    key[i],
				val:    nil,
				childs: make(map[byte]*TrieNode),
			}
		}
		cur = cur.childs[key[i]]
		i++
	}

	if cur.val != nil {
		return ErrKeyExists
	}

	cur.val = val
	return nil
}

func (t *Trie) Get(key string) (any, error) {
	switch {
	case len(key) == 0:
		return "", ErrEmptyKey
	case t.root == nil:
		return "", ErrEmptyTrie
	}

	cur := t.root
	i := 0
	for i < len(key) {
		if _, ok := cur.childs[key[i]]; ok {
			cur = cur.childs[key[i]]
		} else {
			return "", ErrKeyNotFound
		}
		i++
	}

	if cur.val == nil {
		return "", ErrKeyNotFound
	}

	return cur.val, nil
}

func (t *Trie) Delete(key string) error {
	switch {
	case len(key) == 0:
		return ErrEmptyKey
	case t.root == nil:
		return ErrEmptyTrie
	}

	cur := t.root
	i := 0
	for i < len(key) {
		if _, ok := cur.childs[key[i]]; ok {
			cur = cur.childs[key[i]]
		} else {
			return ErrKeyNotFound
		}
		i++
	}
	if cur.val == nil {
		return ErrKeyNotFound
	}
	cur.val = nil
	return nil
}
