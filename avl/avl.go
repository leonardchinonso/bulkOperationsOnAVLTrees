/*
Implementation of Adelson-Velsky and Landis (AVL) self-balancing binary search tree in Go.
References:
	- https://www.cs.cmu.edu/~blelloch/papers/BFS16.pdf
	- https://www.cise.ufl.edu/~nemo/cop3530/AVL-Tree-Rotations.pdf
	- https://github.com/cmuparlay/PAM
	- https://github.com/canepat/balanced-search-tree
*/

package avl

import (
	"bytes"
	"math"
)

type Node struct {
	Key    []byte
	Left   *Node
	Right  *Node
	Height int
	Path   string
}

// NewNode is a custom constructor method to initialise height of the node
func NewNode(key []byte, left, right *Node) *Node {
	node := new(Node)
	node.Key = key
	node.Left = left
	node.Right = right
	node.Height = getHeight(node)
	return node
}

// populatePaths attaches node paths from the root of a node down to the node
func (n *Node) populatePaths(path string) {
	n.Path = path
	if n.Left != nil {
		n.Left.populatePaths(path + "L")
	}
	if n.Right != nil {
		n.Right.populatePaths(path + "R")
	}
}

// Get height of a node
// This function exists because the join method can take a null Node pointer
// and accessing the height property of a null Node pointer will fail
func heightOf(node *Node) int {
	if node == nil {
		return -1
	}
	return node.Height
}

// getHeight sets the height of a node with respect to its children, returns the current node's height
func getHeight(tree *Node) int {
	var height int
	if tree.Left != nil && tree.Right != nil {
		height = int(math.Max(float64(getHeight(tree.Left)), float64(getHeight(tree.Right)))) + 1
	} else if tree.Left != nil {
		height = getHeight(tree.Left) + 1
	} else if tree.Right != nil {
		height = getHeight(tree.Right) + 1
	}
	tree.Height = height
	return height
}

// updateHeight sets the height of a node with respect to its children, returns nothing
func updateHeight(tree *Node) {
	var leftHeight int
	var rightHeight int
	if tree.Left != nil {
		leftHeight = tree.Left.Height
	}
	if tree.Right != nil {
		rightHeight = tree.Right.Height
	}
	tree.Height = int(math.Max(float64(leftHeight), float64(rightHeight))) + 1
}

// expose returns key of a node and it's left and right children
func expose(tree *Node) (tree1 *Node, k []byte, tree2 *Node) {
	if tree != nil {
		return tree.Left, tree.Key, tree.Right
	}
	return nil, []byte{}, nil
}

// rotateRight rotates a node to the right to maintain the AVL balance criteria
func rotateRight(tree *Node) *Node {
	newRoot := tree.Left
	tree.Left = newRoot.Right
	newRoot.Right = tree
	updateHeight(newRoot.Right)
	updateHeight(newRoot)
	return newRoot
}

// rotateLeft rotates a node to the left to maintain the AVL balance criteria
func rotateLeft(tree *Node) *Node {
	newRoot := tree.Right
	tree.Right = newRoot.Left
	newRoot.Left = tree
	updateHeight(newRoot.Left)
	updateHeight(newRoot)
	return newRoot
}

// joinRight concatenates a left tree, k and a right tree
func joinRight(tree1 *Node, k []byte, tree2 *Node) *Node {
	l, kPrime, c := expose(tree1)
	if heightOf(c) <= heightOf(tree2)+1 {
		treePrime := NewNode(k, c, tree2)
		if heightOf(treePrime) <= heightOf(l)+1 {
			return NewNode(kPrime, l, treePrime)
		}
		return rotateLeft(NewNode(kPrime, l, rotateRight(treePrime)))
	}
	treePrime := joinRight(c, k, tree2)
	treePrimePrime := NewNode(kPrime, l, treePrime)
	if heightOf(treePrime) <= heightOf(l)+1 {
		return treePrimePrime
	}
	return rotateLeft(treePrimePrime)
}

// joinLeft concatenates a left tree, k and a right tree
func joinLeft(tree1 *Node, k []byte, tree2 *Node) *Node {
	c, kPrime, r := expose(tree2)
	if heightOf(c) <= heightOf(tree1)+1 {
		treePrime := NewNode(k, tree1, c)
		if heightOf(treePrime) <= heightOf(r)+1 {
			return NewNode(kPrime, treePrime, r)
		}
		return rotateRight(NewNode(kPrime, rotateLeft(treePrime), r))
	}
	treePrime := joinLeft(tree1, k, c)
	treePrimePrime := NewNode(kPrime, treePrime, r)
	if heightOf(treePrime) <= heightOf(r)+1 {
		return treePrimePrime
	}
	return rotateRight(treePrimePrime)
}

// join concatenates a left tree, k and a right tree
func join(tree1 *Node, k []byte, tree2 *Node) *Node {
	if heightOf(tree1) > heightOf(tree2)+1 {
		return joinRight(tree1, k, tree2)
	} else if heightOf(tree2) > heightOf(tree1)+1 {
		return joinLeft(tree1, k, tree2)
	}
	return NewNode(k, tree1, tree2)
}

// split separates a tree into two distinct trees at value k
func split(tree *Node, k []byte) (*Node, bool, *Node) {
	if tree == nil {
		return nil, false, nil
	}
	l, m, r := expose(tree)
	if bytes.Compare(k, m) == 0 {
		return l, true, r
	}
	if bytes.Compare(k, m) == -1 {
		ll, b, lr := split(l, k)
		return ll, b, join(lr, m, r)
	}
	rl, b, rr := split(r, k)
	return join(l, m, rl), b, rr
}

// splitLast separates a tree into two distinct trees at the rightmost node
func splitLast(tree *Node) (*Node, []byte) {
	l, k, r := expose(tree)
	if r == nil {
		return l, k
	}
	treePrime, kPrime := splitLast(r)
	return join(l, k, treePrime), kPrime
}

// join2 concatenates a left tree and a right tree
func join2(tree1 *Node, tree2 *Node) *Node {
	if tree1 == nil {
		return tree2
	}
	tree1Prime, k := splitLast(tree1)
	return join(tree1Prime, k, tree2)
}

// Insert inserts a node into a tree
func Insert(tree *Node, k []byte) *Node {
	tree1, _, tree2 := split(tree, k)
	return join(tree1, k, tree2)
}

// deleteNode deletes a node from a tree
func deleteNode(tree *Node, k []byte) *Node {
	tree1, _, tree2 := split(tree, k)
	return join2(tree1, tree2)
}

// Union carries out the union operation on two trees
func Union(tree1 *Node, tree2 *Node) *Node {
	if tree1 == nil {
		return tree2
	}
	if tree2 == nil {
		return tree1
	}
	l2, k2, r2 := expose(tree2)
	l1, _, r1 := split(tree1, k2)
	treeLeft := Union(l1, l2)
	treeRight := Union(r1, r2)
	return join(treeLeft, k2, treeRight)
}

// Difference carries out the difference operation on two trees
func Difference(tree1 *Node, tree2 *Node) *Node {
	if tree1 == nil {
		return nil
	}
	if tree2 == nil {
		return tree1
	}
	l2, k2, r2 := expose(tree2)
	l1, _, r1 := split(tree1, k2)
	tree1 = Difference(l1, l2)
	tree2 = Difference(r1, r2)
	return join2(tree1, tree2)
}
