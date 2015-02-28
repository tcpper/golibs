package tries

import (
	"errors"
	//"flag"
	"fmt"
	"math"
	//"sort"
	//"strings"
)

const RunWidth = 26

type Node struct {
	children [RunWidth]*Node
	exists   bool
}

type Tries struct {
	Node
	hmin int
	hmax int
}

func (t *Tries) Insert(word string) (err error) {
	if word == "" {
		// never set root.exists = true
		return errors.New("empty word")
	}

	h := 0
	cur := &t.Node
	for _, c := range word {
		if c < 'a' || c > 'z' {
			return fmt.Errorf("%s has chars not allowed", word)
		}

		if cur.children[c-'a'] == nil {
			cur.children[c-'a'] = &Node{}
		}
		cur = cur.children[c-'a']

		h += 1
	}
	cur.exists = true

	if h < t.hmin {
		t.hmin = h
	}

	if h > t.hmax {
		t.hmax = h
	}

	return
}

func (n *Node) dump(prefix string) (words []string) {
	if n.exists {
		words = append(words, prefix)
	}

	for i, n := range n.children {
		if n != nil {
			words = append(words, n.dump(prefix+string(i+'a'))...)
		}
	}

	return
}

func (t *Tries) Match(str string) bool {
	if str == "" {
		return false
	}

	n := &t.Node
	for _, c := range str {
		if n.children[c-'a'] != nil {
			n = n.children[c-'a']
		} else {
			return false
		}
	}

	if n.exists {
		return true
	}

	return false
}

func (t *Tries) MatchPartial(str string) (res []string) {
	n := &t.Node
	for _, c := range str {
		if n.children[c-'a'] != nil {
			n = n.children[c-'a']
		} else {
			return
		}
	}

	return n.dump(str)
}

func NewTries() *Tries {
	return &Tries{hmin: math.MaxInt32, hmax: math.MinInt32}
}
