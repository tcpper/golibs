package rbtree

import (
	"errors"
)

// caller should make sure node is in tree
func (t *RBTree) PlainDeleteNode(n *RBNode) error {
	if n == nil || n == t.Nil {
		return errors.New("rbtree: can not delete nil node")
	}

	if n.p == t.Nil && n != t.root {
		return errors.New("rbtree: node to be deleted is expected to belongs to the tree")
	}

	t.size -= 1

	var nextc, prevc bool
	var candidate *RBNode
	if candidate = t.nextChild(n); candidate != t.Nil {
		nextc = true
	} else if candidate = t.prevChild(n); candidate != t.Nil {
		prevc = true
	}

	if !prevc && !nextc {
		if n.p != t.Nil {
			if n.p.left == n {
				n.p.left = t.Nil
			} else {
				n.p.right = t.Nil
			}
			n.left = t.Nil
			n.right = t.Nil
			n.p = t.Nil
			return nil
		} else {
			// delete single root
			t.root = t.Nil
			n.left = t.Nil
			n.right = t.Nil
			// n.p = nil // not needed for root
			return nil
		}
	}

	if nextc {
		if candidate.p == n { // n.right == candidate
			if n.p != t.Nil {
				if n.p.left == n {
					n.p.left = candidate

				} else {
					n.p.right = candidate
				}
			} else {
				t.root = candidate
			}
			candidate.p = n.p
			// candidate.left == nil
			candidate.left = n.left
			if candidate.left != t.Nil {
				candidate.left.p = candidate
			}
			// keep right leaf of candidate untouched
			n.left = t.Nil
			n.right = t.Nil
			n.p = t.Nil
			return nil
		}

		// candidate should be left leaf of parent
		if n.p != t.Nil {
			if n.p.left == n {
				n.p.left = candidate
			} else {
				n.p.right = candidate
			}
		} else {
			t.root = candidate
		}
		// candidate.left == nil
		candidate.p.left = candidate.right
		if candidate.right != t.Nil {
			candidate.right.p = candidate.p
		}

		candidate.p = n.p
		candidate.left = n.left
		if candidate.left != t.Nil {
			candidate.left.p = candidate
		}
		candidate.right = n.right
		candidate.right.p = candidate
		return nil
	}

	// prevc == true
	if candidate.p == n { // n.left == candidate
		if n.p != t.Nil {
			if n.p.left == n {
				n.p.left = candidate
			} else {
				n.p.right = candidate
			}
		} else {
			t.root = candidate
		}
		candidate.p = n.p
		// candidate.right == nil
		candidate.right = n.right
		if candidate.right != t.Nil {
			candidate.right.p = candidate
		}
		// keep left leaf of candidate untouched
		n.left = t.Nil
		n.right = t.Nil
		n.p = t.Nil
		return nil
	}

	// candidate shoudl be right leaf of parent
	if n.p != t.Nil {
		if n.p.left == n {
			n.p.left = candidate
		} else {
			n.p.right = candidate
		}
	} else {
		t.root = candidate
	}
	// candidate.right == nil
	candidate.p.right = candidate.left
	if candidate.left != t.Nil {
		candidate.left.p = candidate.p
	}

	candidate.p = n.p
	candidate.left = n.left
	candidate.left.p = candidate
	candidate.right = n.right
	if candidate.right != t.Nil {
		candidate.right.p = candidate
	}
	return nil
}

// if not delete all, delete leftmost match
func (t *RBTree) PlainDelete(comp Comparable, all bool) error {
	if !t.dupable && all {
		return errors.New("no need to delete all for nondupable tree")
	}

	if nodes := t.FindNode(comp); len(nodes) > 0 {
		if !all {
			node := nodes[len(nodes)-1]
			return t.PlainDeleteNode(node)
		} else {
			for _, node := range nodes {
				if err := t.PlainDeleteNode(node); err != nil {
					return err
				}
			}
			return nil
		}
	} else {
		return errors.New("key not found")
	}
}
