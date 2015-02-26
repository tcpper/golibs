package rbtree

import (
	//"fmt"
	"testing"
)

type MyInt int

func (a MyInt) LessEqual(b Comparable) bool {
	t := b.(MyInt)
	return a <= t
}

func TestInsert(t *testing.T) {
	tree := NewRBTree(false)
	for i := 0; i < 13; i++ {
		tree.Insert(MyInt(i))
		if err := tree.Verify(); err != nil {
			t.Errorf("verify error %s", err.Error())
		}
	}

	if tree.size != 13 {
		t.Errorf("size error, expect %d v.s. %d", 13, tree.size)
	}
}

func TestInsertDupable(t *testing.T) {
	tree := NewRBTree(true)
	for i := 0; i < 13; i++ {
		tree.Insert(MyInt(13))
		if err := tree.Verify(); err != nil {
			t.Errorf("verify error %s", err.Error())
		}
	}

	if tree.size != 13 {
		t.Errorf("size error, expect %d v.s. %d", 13, tree.size)
	}
}

func TestFind(t *testing.T) {
	tree := NewRBTree(false)
	for i := 0; i < 13; i++ {
		tree.Insert(MyInt(i))
	}

	for i := -10; i < 100; i++ {
		if i < 0 || i >= 13 {
			if vs := tree.Find(MyInt(i)); len(vs) > 0 {
				t.Errorf("found nonexist %d: %d", i, vs[0])
			}
		}

		if 0 <= i && i < 13 {
			if vs := tree.Find(MyInt(i)); len(vs) != 1 {
				t.Errorf("not found %d", i)
			} else if vs[0] != MyInt(i) {
				t.Errorf("not equal %d v.s. %d", i, vs[0])
			}
		}
	}
}

func TestFindDupable(t *testing.T) {
	tree := NewRBTree(true)
	for i := 0; i < 13; i++ {
		for j := 0; j < 130000; j++ {
			tree.Insert(MyInt(j))
		}
	}

	for i := -10; i < 140000; i++ {
		if i < 0 || i >= 130000 {
			if vs := tree.Find(MyInt(i)); len(vs) > 0 {
				t.Errorf("found nonexist %d: %d", i, vs[0])
			}
		}

		if 0 <= i && i < 130000 {
			if vs := tree.Find(MyInt(i)); len(vs) != 13 {
				t.Errorf("count of %d not match 13, really %d", i, len(vs))
			} else if vs[0] != MyInt(i) {
				t.Errorf("not equal %d v.s. %d", i, vs[0])
			}
		}
	}
}

func TestMinMax(t *testing.T) {
	tree := NewRBTree(false)
	for i := 0; i < 13; i++ {
		tree.Insert(MyInt(i))
	}
	if v := tree.Max(); v != MyInt(12) {
		t.Errorf("max expected to be %d v.s. %d", 12, v)
	}
	if v := tree.Min(); v != MyInt(0) {
		t.Errorf("min expected to be %d v.s. %d", 0, v)
	}
}

func TestInteration(t *testing.T) {
	tree := NewRBTree(false)
	for i := 0; i < 13; i++ {
		tree.Insert(MyInt(i))
	}

	cur := tree.MinRaw()

	var i MyInt
	for cur != tree.Nil {
		if i > 12 {
			t.Error("too many, should be at most 12")
		}

		if Compare(i, cur.Bag) != Equal {
			t.Errorf("expect %d, not %d", i, cur.Bag)
		}
		i++
		cur = tree.NextRaw(cur)
	}
}

func TestInterationBackward(t *testing.T) {
	tree := NewRBTree(false)
	for i := 0; i < 13; i++ {
		tree.Insert(MyInt(i))
	}

	cur := tree.MaxRaw()

	var i MyInt
	i = MyInt(12)

	for cur != tree.Nil {
		if i < 0 {
			t.Error("too many, should be at less 0")
		}

		if Compare(i, cur.Bag) != Equal {
			t.Errorf("expect %d, not %d", i, cur.Bag)
		}
		i--
		cur = tree.PrevRaw(cur)
	}
}

func TestNonDupable(t *testing.T) {
	tree := NewRBTree(false)
	for i := 0; i < 13; i++ {
		tree.Insert(MyInt(i))
	}

	for i := -10; i < 100; i++ {
		err := tree.Insert(MyInt(i))
		if 0 <= i && i < 13 {
			if err == nil {
				t.Errorf("error expaected for nondupable tree, and key %d", i)
			}
		} else {
			if err != nil {
				t.Errorf("error unexpaected for nondupable tree, and key %d", i)
			}
		}
	}
}

func TestDupableTreeFind(t *testing.T) {
	tree := NewRBTree(true)

	items := []int{0, 0, 1, 2, 3, 4, 4, 4, 5, 6, 7, 7, 8, 9, 10}
	for _, item := range items {
		err := tree.Insert(MyInt(item))
		if err != nil {
			t.Errorf("unexpected error for dupable tree and key %d", item)
		}
	}

	bs := tree.Find(MyInt(4))
	if len(bs) == 3 {
		if bs[0] != MyInt(4) ||
			bs[1] != MyInt(4) ||
			bs[2] != MyInt(4) {
			t.Error("not equal")
		}
	} else {
		t.Error("expect to find 3 items")
	}
}

func TestDupableTreeInteration(t *testing.T) {
	tree := NewRBTree(true)

	items := []int{0, 0, 1, 2, 3, 4, 4, 4, 5, 6, 7, 7, 8, 9, 10}
	for _, item := range items {
		err := tree.Insert(MyInt(item))
		if err != nil {
			t.Errorf("unexpected error for dupable tree and key %d", item)
		}
	}

	i := 0
	cur := tree.MinRaw()
	for cur != tree.Nil {
		if Compare(MyInt(items[i]), cur.Bag) != Equal {
			t.Errorf("expect %d, not %d", i, cur.Bag)
		}
		i++
		cur = tree.NextRaw(cur)
	}
	if i != len(items) {
		t.Errorf("too many tree item %d", i)
	}
}

func TestDupableTreeInterationBackward(t *testing.T) {
	tree := NewRBTree(true)

	items := []int{0, 0, 1, 2, 3, 4, 4, 4, 5, 6, 7, 7, 8, 9, 10}
	for _, item := range items {
		err := tree.Insert(MyInt(item))
		if err != nil {
			t.Errorf("unexpected error for dupable tree and key %d", item)
		}
	}

	i := len(items) - 1
	cur := tree.MaxRaw()
	for cur != tree.Nil {
		if Compare(MyInt(items[i]), cur.Bag) != Equal {
			t.Errorf("expect %d, not %d", i, cur.Bag)
		}
		i--
		cur = tree.PrevRaw(cur)
	}
	if i != -1 {
		t.Errorf("too many tree item %d", i)
	}
}

func filter(items []int, value int) (left []int) {
	for _, item := range items {
		if item != value {
			left = append(left, item)
		}
	}
	return
}

func count(items []int, value int) (c int) {
	for _, item := range items {
		if item == value {
			c += 1
		}
	}
	return
}

func plainDeleteHelper(tree *RBTree, items []int, value int, count int, t *testing.T) {
	nodes := tree.FindRaw(MyInt(value))
	if len(nodes) != count {
		t.Errorf("expect %d four v.s. %d", count, len(nodes))
	} else {
		for i := 0; i < count; i++ {
			tree.PlainDeleteRaw(nodes[0])
			nodes = tree.FindRaw(MyInt(value))
			if len(nodes) != count-1-i {
				t.Errorf("should be %d four left", count-1-i)
			}
		}
	}

	i := 0
	cur := tree.MinRaw()
	for cur != tree.Nil {
		if Compare(MyInt(items[i]), cur.Bag) != Equal {
			t.Errorf("expect %d, not %d", i, cur.Bag)
		}
		i++
		cur = tree.NextRaw(cur)
	}
	if i != len(items) {
		t.Errorf("too many tree item %d", i)
	}
}

func deleteHelper(tree *RBTree, items []int, value int, count int, t *testing.T) {
	nodes := tree.FindRaw(MyInt(value))
	if len(nodes) != count {
		t.Errorf("expect %d four v.s. %d", count, len(nodes))
	} else {
		for i := 0; i < count; i++ {
			tree.DeleteRaw(nodes[0])
			if err := tree.Verify(); err != nil {
				t.Errorf("verify failed after delete, %s", err.Error())
			}
			nodes = tree.FindRaw(MyInt(value))
			if len(nodes) != count-1-i {
				t.Errorf("should be %d four left", count-1-i)
			}
		}
	}

	i := 0
	cur := tree.MinRaw()
	for cur != tree.Nil {
		if Compare(MyInt(items[i]), cur.Bag) != Equal {
			t.Errorf("expect %d, not %d", i, cur.Bag)
		}
		i++
		cur = tree.NextRaw(cur)
	}
	if i != len(items) {
		t.Errorf("too many tree item %d", i)
	}
}

func TestDupableTreePlainDelete(t *testing.T) {
	tree := NewRBTree(true)

	items := []int{0, 0, 1, 2, 3, 4, 4, 4, 5, 6, 7, 7, 8, 9, 10}
	for _, item := range items {
		err := tree.Insert(MyInt(item))
		if err != nil {
			t.Errorf("unexpected error for dupable tree and key %d", item)
		}
	}

	if err := tree.PlainDelete(MyInt(40), true); err == nil {
		t.Error("delete not exist node should give out error")
	}

	if err := tree.PlainDelete(MyInt(4), false); err == nil {
		if l := len(tree.Find(MyInt(4))); l != 2 {
			t.Error("should left 2, really %d", l)
		}

		if err := tree.PlainDelete(MyInt(4), true); err == nil {
			if l := len(tree.Find(MyInt(4))); l != 0 {
				t.Error("delete all failed")
			}
		} else {
			t.Errorf("unexpected error, %s", err.Error())
		}
	} else {
		t.Errorf("delete unexpected error, %s", err.Error())
	}

}

func TestDupableTreeDelete(t *testing.T) {
	tree := NewRBTree(true)

	items := []int{0, 0, 1, 2, 3, 4, 4, 4, 5, 6, 7, 7, 8, 9, 10}
	for _, item := range items {
		err := tree.Insert(MyInt(item))
		if err != nil {
			t.Errorf("unexpected error for dupable tree and key %d", item)
		}
	}

	if err := tree.Delete(MyInt(40), true); err == nil {
		t.Error("delete not exist node should give out error")
	}

	if err := tree.Delete(MyInt(4), false); err == nil {
		if err := tree.Verify(); err != nil {
			t.Error("verify failed after delete, %s", err.Error())
		}
		if l := len(tree.Find(MyInt(4))); l != 2 {
			t.Error("should left 2, really %d", l)
		}

		if err := tree.Delete(MyInt(4), true); err == nil {
			if err := tree.Verify(); err != nil {
				t.Error("verify failed after delete, %s", err.Error())
			}
			if l := len(tree.Find(MyInt(4))); l != 0 {
				t.Error("delete all failed")
			}
		} else {
			t.Errorf("unexpected error, %s", err.Error())
		}
	} else {
		t.Errorf("delete unexpected error, %s", err.Error())
	}

}

func TestDupableTreePlainDeleteRaw(t *testing.T) {
	tree := NewRBTree(true)

	items := []int{0, 0, 1, 2, 3, 4, 4, 4, 5, 6, 7, 7, 8, 9, 10}
	for _, item := range items {
		err := tree.Insert(MyInt(item))
		if err != nil {
			t.Errorf("unexpected error for dupable tree and key %d", item)
		}
	}

	if err := tree.PlainDeleteRaw(tree.NewRBNode(MyInt(40), Red)); err == nil {
		t.Error("delete non match root node should be error")
	}

	for _, v := range []int{4, 0, 1, 7, 2, 3, 5, 6, 8, 9, 10} {
		c := count(items, v)
		items = filter(items, v)
		//fmt.Println(c, items)
		plainDeleteHelper(tree, items, v, c, t)
	}

	if tree.size != 0 {
		t.Error("size of tree should be zero")
	}
}

func TestDupableTreeDeleteRaw(t *testing.T) {
	tree := NewRBTree(true)

	items := []int{0, 0, 1, 2, 3, 4, 4, 4, 5, 6, 7, 7, 8, 9, 10}
	for _, item := range items {
		err := tree.Insert(MyInt(item))
		if err != nil {
			t.Errorf("unexpected error for dupable tree and key %d", item)
		}
	}

	for _, v := range []int{4, 0, 1, 7, 2, 3, 5, 6, 8, 9, 10} {
		c := count(items, v)
		items = filter(items, v)
		//fmt.Println(c, items)
		deleteHelper(tree, items, v, c, t)
	}

	if tree.size != 0 {
		t.Error("size of tree should be zero")
	}
}

func TestDupableTreeRotate(t *testing.T) {
	tree := NewRBTree(true)

	items := []int{0, 0, 1, 2, 3, 4, 4, 4, 5, 6, 7, 7, 8, 9, 10}
	for _, item := range items {
		err := tree.Insert(MyInt(item))
		if err != nil {
			t.Errorf("unexpected error for dupable tree and key %d", item)
		}
	}

	min := tree.MinRaw()
	if min == tree.Nil {
		t.Error("unexpected nil")
	}

	next := min
	for next = tree.NextRaw(next); next != tree.Nil; next = tree.NextRaw(next) {
		if next.right == tree.Nil {
			if err := tree.rotateLeft(next); err == nil {
				t.Error("should give out error when do left-rotate with right child be nil")
			}
		} else {
			if err := tree.rotateLeft(next); err != nil {
				t.Error("should not give out error when do left-rotate with right child not be nil")
			}
		}
	}

	next = tree.MinRaw()
	i := 0
	for ; next != tree.Nil; next = tree.NextRaw(next) {
		if Compare(next.Bag, MyInt(items[i])) != Equal {
			t.Error("not equal %d v.s. %d", next.Bag, items[i])
		}
		i++
	}

	if i != len(items) {
		t.Error("size not match %d, expect %d", i, len(items))
	}
}

func TestDupableTreeFindAfterRotate(t *testing.T) {
	tree := NewRBTree(true)

	items := []int{0, 0, 1, 2, 3, 4, 4, 4, 5, 6, 7, 7, 8, 9, 10}
	for _, item := range items {
		err := tree.Insert(MyInt(item))
		if err != nil {
			t.Errorf("unexpected error for dupable tree and key %d", item)
		}
	}

	min := tree.MinRaw()
	if min == tree.Nil {
		t.Error("unexpected nil")
	}

	next := min
	for next = tree.NextRaw(next); next != tree.Nil; next = tree.NextRaw(next) {
		if next.left == tree.Nil {
			if err := tree.rotateRight(next); err == nil {
				t.Error("should give out error when do right-rotate with left child be nil")
			}
		} else {
			if err := tree.rotateRight(next); err != nil {
				t.Error("should not give out error when do right-rotate with left child not be nil")
			}
		}
	}

	if vs := tree.Find(MyInt(4)); len(vs) != 3 {
		t.Error("expect 3 items(4), really %d", len(vs))
	}

	// if err := tree.Verify(); err != nil {
	// 	t.Errorf("verify fail, %s", err.Error())
	// }
}
