package tries

import (
	"testing"
)

func TestMatch(t *testing.T) {
	tree := NewTries()
	words := []string{"a", "ab", "abc", "abd", "abe", "hello", "world"}
	for _, word := range words {
		tree.Insert(word)
	}

	if !tree.Match("abd") {
		t.Error("expect match abd")
	}

	if tree.Match("worl") {
		t.Error("unexpected match worl")
	}

	if tree.Match("worlde") {
		t.Error("unexpected match worlde")
	}
}

func TestMatchPartial(t *testing.T) {
	tree := NewTries()
	words := []string{"a", "ab", "abc", "abd", "abe", "hello", "world"}
	for _, word := range words {
		tree.Insert(word)
	}

	partial := tree.MatchPartial("ab")
	if len(partial) != 4 {
		t.Error("expect 4 match")
	}
	if partial[0] != "ab" || partial[1] != "abc" || partial[2] != "abd" || partial[3] != "abe" {
		t.Errorf("not expected: %v", partial)
	}

	partial = tree.MatchPartial("")
	if len(partial) != len(words) {
		t.Errorf("expect %d match, found %d", len(words), len(partial))
	}
}
