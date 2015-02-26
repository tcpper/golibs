package rbtree

import (
	"errors"
	"fmt"
)

type Comparable interface {
	LessEqual(o Comparable) (t bool)
}

type Relation int

const (
	Less    Relation = -1
	Equal   Relation = 0
	Greater Relation = 1
)

type Color bool

const (
	Red   Color = false
	Black Color = true
)

func Compare(a Comparable, b Comparable) Relation {
	if a.LessEqual(b) {
		if !b.LessEqual(a) {
			return Less
		}
		return Equal
	} else {
		return Greater
	}
}

type RBTree struct {
	root    *RBNode
	dupable bool
	size    int
	Nil     *RBNode
}

type RBNode struct {
	color Color
	left  *RBNode
	right *RBNode
	p     *RBNode
	Bag   Comparable
}

func NewRBNode(comp Comparable, color Color) *RBNode {
	return &RBNode{Bag: comp, color: color}
}

func NewRBTree(dupable bool) *RBTree {
	sentinel := NewRBNode(nil, Black)
	return &RBTree{dupable: dupable, Nil: sentinel, root: sentinel}
}

func (t *RBTree) NewRBNode(comp Comparable, color Color) *RBNode {
	return &RBNode{Bag: comp, color: color, left: t.Nil, right: t.Nil, p: t.Nil}
}

func (t *RBTree) rotateLeft(x *RBNode) error {
	if x.right == t.Nil {
		return errors.New("rotate left require right child not nil")
	}

	y := x.right
	y.p = x.p

	x.right = y.left
	if y.left != t.Nil {
		y.left.p = x
	}

	y.left = x

	if x.p == t.Nil {
		t.root = y
	} else {
		if x.p.left == x {
			x.p.left = y
		} else {
			x.p.right = y
		}
	}

	x.p = y
	return nil
}

func (t *RBTree) rotateRight(y *RBNode) error {
	if y.left == t.Nil {
		return errors.New("rotate right require left child not nil")
	}

	x := y.left
	x.p = y.p

	y.left = x.right
	if x.right != t.Nil {
		x.right.p = y
	}

	x.right = y

	if y.p == t.Nil {
		t.root = x
	} else {
		if y.p.left == y {
			y.p.left = x
		} else {
			y.p.right = x
		}
	}

	y.p = x
	return nil
}

func (t *RBTree) insertFix(z *RBNode) {
	for z.p.color == Red {
		// enter the for loop imply z is not root and z.p.p != t.Nil, because
		// if z is root (z.p == t.Nil), remember that sentinel(z.p) is black,
		// if grandparent is t.Nil (z.p is root), z.p.color can never be Red
		if z.p == z.p.p.left {
			if uncle := z.p.p.right; uncle.color == Red {
				// for loop repeat only if this case occures
				z.p.color = Black
				uncle.color = Black
				z.p.p.color = Red
				z = z.p.p
			} else { // uncle is black
				if z == z.p.right {
					z = z.p
					t.rotateLeft(z)
				}
				z.p.color = Black
				z.p.p.color = Red
				t.rotateRight(z.p.p)
				// this always just out of the loop (z.p.color == Black)
			}
		} else { // z.p == z.p.p.right
			if uncle := z.p.p.left; uncle.color == Red {
				// for loop repeat only if this case occures
				z.p.color = Black
				uncle.color = Black
				z.p.p.color = Red
				z = z.p.p
			} else { // uncle is black
				if z == z.p.left {
					z = z.p
					t.rotateRight(z)

				}
				z.p.color = Black
				z.p.p.color = Red
				t.rotateLeft(z.p.p)
				// this always just out of the loop (z.p.color == Black)
			}
		}
	}
	t.root.color = Black
}

func (t *RBTree) InsertNode(n *RBNode) error {
	parent := t.Nil
	t.size += 1
	// p is **RBNode
	p := &t.root
	for *p != t.Nil {
		parent = *p
		if t.dupable {
			if n.Bag.LessEqual((*p).Bag) {
				// we can ignore a "*()" for golang treat "." as "->"
				// just an illustration
				// p = &((*(*p)).left)
				p = &((*p).left)
			} else {
				// aha, different way to write
				p = &((*(*p)).right)
			}
		} else {
			switch Compare(n.Bag, (*p).Bag) {
			case Less:
				p = &((*p).left)
			case Greater:
				p = &((*(*p)).right)
			case Equal:
				return errors.New("duplicate key for nondupable tree")
			}
		}
	}
	n.p = parent
	*p = n
	// n.left = t.Nil
	// n.right = t.Nil
	t.insertFix(n)
	return nil
}

// it is user's responsibility to ensure key != nil
func (t *RBTree) Insert(comp Comparable) error {
	n := t.NewRBNode(comp, Red)
	return t.InsertNode(n)
}

// replace u with v, only update v and parent(maybe root) relationship
func (t *RBTree) transplantUp(u *RBNode, v *RBNode) {
	if u.p == t.Nil {
		t.root = v
	} else if u == u.p.left {
		u.p.left = v
	} else {
		u.p.right = v
	}
	v.p = u.p
}

func (t *RBTree) deleteFix(x *RBNode) {
	// x may be sentinel, but x.p have been set properly
	for x != t.root && x.color == Black {
		if x == x.p.left {
			// x points to a extra black node, so sibling of it can not be sentinel
			// otherwise, the tree violate bhHeight(x.p) equal rule
			w := x.p.right
			if w.color == Red {
				// if w is red, remember x is double black
				// due to reason alike bhHeight(x.p)
				// all children of w should exist and be black
				w.color = Black
				x.p.color = Red
				t.rotateLeft(x.p)
				w = x.p.right
			}
			// when arrived here, w(sibling of x) is black
			if w.left.color == Black && w.right.color == Black {
				// move x upward and make w red, no bhHeight change totally
				// cover the case that children of w both are sentinel
				w.color = Red
				x = x.p
			} else {
				// children of w can't be both sentinel otherwise they are all black
				// at least one of them are red
				if w.right.color == Black {
					// so w.left must be not be sentinel and must be red
					// do a rotation
					w.left.color = Black
					w.color = Red
					t.rotateRight(w)
					w = x.p.right
					// now w is black and right child of it is red
				}
				w.color = x.p.color
				x.p.color = Black
				w.right.color = Black
				t.rotateLeft(x.p)
				x = t.root
			}
		} else { // x == x.p.right
			w := x.p.left
			if w.color == Red {
				w.color = Black
				x.p.color = Red
				t.rotateRight(x.p)
				w = x.p.left
			}
			if w.left.color == Black && w.right.color == Black {
				w.color = Red
				x = x.p
			} else {
				if w.left.color == Black {
					w.right.color = Black
					w.color = Red
					t.rotateLeft(w)
					w = x.p.left
				}
				w.color = x.p.color
				x.p.color = Black
				w.left.color = Black
				t.rotateRight(x.p)
				x = t.root
			}
		}
	}
	// case 1: x points to a red-and-black node, make it black
	// case 2: x points to root, just drop the extra black
	// case 3: suitable rotations and recolorings done, exit loop
	x.color = Black
}

func (t *RBTree) DeleteNode(z *RBNode) {
	t.size -= 1
	var x *RBNode
	y := z
	yOrigColor := y.color
	if z.left == t.Nil {
		x = z.right
		t.transplantUp(z, z.right)
	} else if z.right == t.Nil {
		x = z.left
		t.transplantUp(z, z.left)
	} else {
		y = t.nextChild(z) // y could not be t.Nil
		yOrigColor = y.color
		x = y.right // x maybe t.Nil
		if y.p == z {
			// y is directly child of z, not need to pick and put
			// just replace z with y is OK
			//if x == t.Nil {
			x.p = y // TODO, only need when x is sentinel
			//}
		} else {
			// y need to be picked and put in the position of z
			t.transplantUp(y, y.right) // y.right maybe t.Nil
			// now y is detached from parent
			// deal with right child of z
			y.right = z.right
			y.right.p = y
		}
		// replace z with y for upward relationship
		t.transplantUp(z, y)
		// deal with left child of z
		y.left = z.left
		y.left.p = y
		y.color = z.color
	}
	if yOrigColor == Black {
		// we've removed a black node
		// x point to where the original black node reside
		t.deleteFix(x)
	}
}

func (t *RBTree) Delete(comp Comparable, all bool) error {
	if !t.dupable && all {
		return errors.New("no need to delete all for nondupable tree")
	}

	if nodes := t.FindNode(comp); len(nodes) > 0 {
		if !all {
			node := nodes[len(nodes)-1]
			t.DeleteNode(node)
			return nil
		} else {
			for _, node := range nodes {
				t.DeleteNode(node)
			}
			return nil
		}
	} else {
		return errors.New("key not found")
	}
}

func (t *RBTree) Find(key Comparable) (bags []Comparable) {
	nodes := t.FindNode(key)

	for _, n := range nodes {
		bags = append(bags, n.Bag)
	}
	return
}

// for dupable tree, random sequence
func (t *RBTree) FindNode(key Comparable) (nodes []*RBNode) {
	if t.size == 0 {
		return
	}
	n := t.root
LOOP:
	for n != t.Nil {
		switch Compare(key, n.Bag) {
		case Less:
			n = n.left
		case Greater:
			n = n.right
		case Equal:
			nodes = append(nodes, n)
			break LOOP
		}
	}
	if !t.dupable || n == t.Nil {
		return
	}

	next := n
	for next = t.NextNode(next); next != t.Nil; next = t.NextNode(next) {
		if Compare(key, next.Bag) == Equal {
			nodes = append(nodes, next)
		} else {
			break
		}
	}

	prev := n
	for prev = t.PrevNode(prev); prev != t.Nil; prev = t.PrevNode(prev) {
		if Compare(key, prev.Bag) == Equal {
			nodes = append(nodes, prev)
		} else {
			break
		}
	}
	return
}

func (t *RBTree) Min() Comparable {
	if t.size == 0 {
		return nil
	}
	n := t.root
	for n.left != t.Nil {
		n = n.left
	}
	return n.Bag
}

func (t *RBTree) MinNode() *RBNode {
	if t.size == 0 {
		return t.Nil
	}
	n := t.root
	for n.left != t.Nil {
		n = n.left
	}
	return n
}

func (t *RBTree) Max() Comparable {
	if t.size == 0 {
		return nil
	}
	n := t.root
	for n.right != t.Nil {
		n = n.right
	}
	return n.Bag
}

func (t *RBTree) MaxNode() *RBNode {
	if t.size == 0 {
		return t.Nil
	}
	n := t.root
	for n.right != t.Nil {
		n = n.right
	}
	return n
}

func (t *RBTree) prevChild(n *RBNode) *RBNode {
	if n.left != t.Nil {
		n = n.left
		for n.right != t.Nil {
			n = n.right
		}
		return n
	}
	return t.Nil
}

func (t *RBTree) prevParent(n *RBNode) *RBNode {
	for n.p != t.Nil {
		if n.p.right == n {
			return n.p
		}
		n = n.p
	}
	return t.Nil
}

func (t *RBTree) PrevNode(n *RBNode) *RBNode {
	if prev := t.prevChild(n); prev != t.Nil {
		return prev
	} else if prev := t.prevParent(n); prev != t.Nil {
		return prev
	}
	// we are left most node
	return t.Nil
}

func (t *RBTree) nextChild(n *RBNode) *RBNode {
	if n.right != t.Nil {
		n = n.right
		for n.left != t.Nil {
			n = n.left
		}
		return n
	}
	return t.Nil
}

func (t *RBTree) nextParent(n *RBNode) *RBNode {
	for n.p != t.Nil {
		if n.p.left == n {
			return n.p
		}
		n = n.p
	}
	return t.Nil
}

func (t *RBTree) NextNode(n *RBNode) *RBNode {
	if next := t.nextChild(n); next != t.Nil {
		return next
	} else if next := t.nextParent(n); next != t.Nil {
		return next
	}
	// we are right most node
	return t.Nil
}

func (t *RBTree) Verify() error {
	if t.size > 0 {
		if t.root.color != Black {
			return errors.New("root is red")
		}
		err, _ := t.VerifyNode(t.root)
		return err
	}
	return nil
}

func (t *RBTree) VerifyNode(n *RBNode) (err error, bh int) {
	if n.color == Red {
		if n.left.color == Red || n.right.color == Red {
			return errors.New("adjacent red node"), -1
		}
	}
	var bhLeft, bhRight int
	if n.left != t.Nil {
		err, bhLeft = t.VerifyNode(n.left)
		if err != nil {
			return err, -1
		}
	}

	if n.right != t.Nil {
		err, bhRight = t.VerifyNode(n.right)
		if err != nil {
			return err, -1
		}
	}
	if bhLeft == bhRight {
		bh = bhLeft
		if n.color == Black {
			bh++
		}
		return nil, bh
	}

	return fmt.Errorf("bh diff bhLeft: %d, bhRight: %d", bhLeft, bhRight), -1
}
