package treap

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"
)

// Tree contains a reference to the root of the tree
type Tree struct {
	Root *Node
	rnd  *rand.Rand
}

// Randomized priorities are in the range of [0 - 2^31)
const maxPriority = math.MaxInt32

// NewTree returns an empty treap Tree
func NewTree() *Tree {
	return &Tree{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Get searches the Tree for a target, returns node ptr and boolean indicating if found
func (tree *Tree) Get(target interface{}) (*Node, bool) {
	root := tree.Root

	for root != nil {
		if compare(target, root.Value) == 0 {
			return root, true
		}
		if compare(target, root.Value) < 0 {
			root = root.Left
		} else {
			root = root.Right
		}
	}

	return nil, false
}

// Insert will add a new node to the tree with the given value
func (tree *Tree) Insert(value interface{}) {
	current := tree.naiveInsert(value)

	// Bubble up while the current node's priority is lower than its parent's
	for current.Parent != nil && compare(current.Priority, current.Parent.Priority) < 0 {
		if current == current.Parent.Left {
			current.Parent.rightRotate()
		} else {
			current.Parent.leftRotate()
		}
	}
	if current.Parent == nil {
		tree.Root = current
	}
}

// Naive BST insertion for a given value
func (tree *Tree) naiveInsert(value interface{}) *Node {
	var inserted *Node

	root := tree.Root
	if root == nil {
		inserted = &Node{Value: value, Priority: tree.rnd.Intn(maxPriority)}
		tree.Root = inserted
	}

	for inserted == nil {
		if compare(value, root.Value) < 0 {
			if root.Left == nil {
				root.Left = &Node{Value: value, Priority: tree.rnd.Intn(maxPriority), Parent: root}
				inserted = root.Left
			} else {
				root = root.Left
			}

		} else {
			// Duplicate values placed on the right
			if root.Right == nil {
				root.Right = &Node{Value: value, Priority: tree.rnd.Intn(maxPriority), Parent: root}
				inserted = root.Right
			} else {
				root = root.Right
			}
		}
	}
	return inserted
}

func (tree *Tree) toSlice() []*Node {
	arr := make([]*Node, 0)
	tree.Root.flatten(&arr)
	return arr
}

// Node is a sub-tree
type Node struct {
	Value    interface{}
	Priority int
	Left     *Node
	Right    *Node
	Parent   *Node
}

func (node *Node) rightRotate() {
	child := node.Left
	parent := node.Parent

	// Promote node to be its grandparent's child
	// If the current node is the parent's left child, make node.Left the parent's left child
	if parent != nil && compare(node.Value, parent.Value) < 0 {
		parent.Left = child

	} else if parent != nil && compare(node.Value, parent.Value) >= 0 {
		parent.Right = child

	}
	child.Parent = parent

	// Hand over the Right child of the Left node
	node.Left = child.Right
	if child.Right != nil {
		child.Right.Parent = node
	}

	// Swap parent/child relationship
	child.Right = node
	node.Parent = child
}

func (node *Node) leftRotate() {
	child := node.Right
	parent := node.Parent

	// Promote node to be its grandparent's child
	if parent != nil && compare(node.Value, parent.Value) < 0 {
		parent.Left = child

	} else if parent != nil && compare(node.Value, parent.Value) >= 0 {
		parent.Right = child

	}
	child.Parent = parent

	// Hand over the Left child of the Right node
	node.Right = child.Left
	if child.Left != nil {
		child.Left.Parent = node
	}

	// Swap parent/child relationship
	child.Left = node
	node.Parent = child
}

// In order traversal to flatten tree into slice
func (node *Node) flatten(arr *[]*Node) {
	if node == nil {
		return
	}

	node.Left.flatten(arr)

	*arr = append(*arr, node)

	node.Right.flatten(arr)
}

func (node *Node) String() string {
	return fmt.Sprintf("(%v, %v)", node.Value, node.Priority)
}

// TODO: Don't panic
func compare(a, b interface{}) int {
	intA, okA := a.(int)
	intB, okB := b.(int)
	if !okA || !okB {
		err := fmt.Errorf("compare expected: (int, int), got: (%v, %v)", reflect.TypeOf(a), reflect.TypeOf(b))
		panic(err)
	}
	return intA - intB
}
