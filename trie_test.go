package trie

import (
	"testing"
)

// RuneTrie

func TestRuneTrie(t *testing.T) {
	trie := NewRuneTrie()
	testTrie(t, trie)
}

func TestRuneTrieNilBehavior(t *testing.T) {
	trie := NewRuneTrie()
	testNilBehavior(t, trie)
}

func TestRuneTrieRoot(t *testing.T) {
	trie := NewRuneTrie()
	testTrieRoot(t, trie)
}

func TestRuneTrieWalk(t *testing.T) {
	trie := NewRuneTrie()
	testTrieWalk(t, trie)
}

// PathTrie

func TestPathTrie(t *testing.T) {
	trie := NewPathTrie()
	testTrie(t, trie)
}

func TestPathTrieNilBehavior(t *testing.T) {
	trie := NewPathTrie()
	testNilBehavior(t, trie)
}

func TestPathTrieRoot(t *testing.T) {
	trie := NewPathTrie()
	testTrieRoot(t, trie)
}

func TestPathTrieWalk(t *testing.T) {
	trie := NewPathTrie()
	testTrieWalk(t, trie)
}

func testTrie(t *testing.T, trie Trier) {
	const firstPutValue = "first put"
	cases := []struct {
		key   string
		value interface{}
	}{
		{"fish", 0},
		{"/cat", 1},
		{"/dog", 2},
		{"/cats", 3},
		{"/caterpillar", 4},
		{"/cat/gideon", 5},
		{"/cat/giddy", 6},
	}

	// get missing keys
	for _, c := range cases {
		if value := trie.Get(c.key); value != nil {
			t.Errorf("expected key %s to be missing, found value %v", c.key, value)
		}
	}

	// initial put
	for _, c := range cases {
		if isNew := trie.Put(c.key, firstPutValue); !isNew {
			t.Errorf("expected key %s to be missing", c.key)
		}
	}

	// subsequent put
	for _, c := range cases {
		if isNew := trie.Put(c.key, c.value); isNew {
			t.Errorf("expected key %s to have a value already", c.key)
		}
	}

	// get
	for _, c := range cases {
		if value := trie.Get(c.key); value != c.value {
			t.Errorf("expected key %s to have value %v, got %v", c.key, c.value, value)
		}
	}

	// delete, expect Delete to return true indicating a node was nil'd
	for _, c := range cases {
		if deleted := trie.Delete(c.key); !deleted {
			t.Errorf("expected key %s to be deleted", c.key)
		}
	}

	// delete cleaned all the way to the first character
	// expect Delete to return false bc no node existed to nil
	for _, c := range cases {
		if deleted := trie.Delete(string(c.key[0])); deleted {
			t.Errorf("expected key %s to be cleaned by delete", string(c.key[0]))
		}
	}

	// get deleted keys
	for _, c := range cases {
		if value := trie.Get(c.key); value != nil {
			t.Errorf("expected key %s to be deleted, got value %v", c.key, value)
		}
	}
}

func testNilBehavior(t *testing.T, trie Trier) {
	cases := []struct {
		key   string
		value interface{}
	}{
		{"/cat", 1},
		{"/catamaran", 2},
		{"/caterpillar", nil},
	}
	expectNilValues := []string{"/", "/c", "/ca", "/caterpillar", "/other"}

	// initial put
	for _, c := range cases {
		if isNew := trie.Put(c.key, c.value); !isNew {
			t.Errorf("expected key %s to be missing", c.key)
		}
	}

	// get nil
	for _, key := range expectNilValues {
		if value := trie.Get(key); value != nil {
			t.Errorf("expected key %s to have value nil, got %v", key, value)
		}
	}
}

func testTrieRoot(t *testing.T, trie Trier) {
	const firstPutValue = "first put"
	const putValue = "value"

	if value := trie.Get(""); value != nil {
		t.Errorf("expected key '' to be missing, found value %v", value)
	}
	if !trie.Put("", firstPutValue) {
		t.Error("expected key '' to be missing")
	}
	if trie.Put("", putValue) {
		t.Error("expected key '' to have a value already")
	}
	if value := trie.Get(""); value != putValue {
		t.Errorf("expected key '' to have value %v, got %v", putValue, value)
	}
	if !trie.Delete("") {
		t.Error("expected key '' to be deleted")
	}
	if value := trie.Get(""); value != nil {
		t.Errorf("expected key '' to be deleted, got value %v", value)
	}
}

func testTrieWalk(t *testing.T, trie Trier) {
	table := map[string]interface{}{
		"fish":         0,
		"/cat":         1,
		"/dog":         2,
		"/cats":        3,
		"/caterpillar": 4,
		"/notes":       30,
		"/notes/new":   31,
		"/notes/:id":   32,
	}
	// key -> times walked
	walked := make(map[string]int)
	for key, _ := range table {
		walked[key] = 0
	}

	for key, value := range table {
		if isNew := trie.Put(key, value); !isNew {
			t.Errorf("expected key %s to be missing", key)
		}
	}

	walker := func(key string, value interface{}) error {
		// value for each walked key is correct
		if value != table[key] {
			t.Errorf("expected key %s to have value %v, got %v", key, table[key], value)
		}
		walked[key] += 1
		return nil
	}
	trie.Walk(walker)

	// each key/value walked exactly once
	for key, walkedCount := range walked {
		if walkedCount != 1 {
			t.Errorf("expected key %s to be walked exactly once, got %v", key, walkedCount)
		}
	}
}
