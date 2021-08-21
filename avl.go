/*
Implementation of Adelson-Velsky and Landis (AVL) self-balancing binary search tree in Go.
References:
	- https://www.cs.cmu.edu/~blelloch/papers/BFS16.pdf
	- https://www.cise.ufl.edu/~nemo/cop3530/AVL-Tree-Rotations.pdf
	- https://github.com/cmuparlay/PAM
	- https://github.com/canepat/balanced-search-tree
*/

package main

import (
	"bufio"
	"fmt"
	"gopkg.in/eapache/queue.v1"
	"log"
	"math"
	"os"
	"os/exec"
	"reflect"
	"sort"
)

type Node struct {
	Key    int64
	Left   *Node
	Right  *Node
	Height int
	Path   string
}

// NewNode is a custom constructor method to initialise height of the node
func NewNode(key int64, left, right *Node) *Node {
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
func expose(tree *Node) (tree1 *Node, k int64, tree2 *Node) {
	if tree != nil {
		return tree.Left, tree.Key, tree.Right
	}
	return nil, 0, nil
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
func joinRight(tree1 *Node, k int64, tree2 *Node) *Node {
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
func joinLeft(tree1 *Node, k int64, tree2 *Node) *Node {
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
func join(tree1 *Node, k int64, tree2 *Node) *Node {
	if heightOf(tree1) > heightOf(tree2)+1 {
		return joinRight(tree1, k, tree2)
	} else if heightOf(tree2) > heightOf(tree1)+1 {
		return joinLeft(tree1, k, tree2)
	}
	return NewNode(k, tree1, tree2)
}

// split separates a tree into two distinct trees at value k
func split(tree *Node, k int64) (*Node, bool, *Node) {
	if tree == nil {
		return nil, false, nil
	}
	l, m, r := expose(tree)
	if k == m {
		return l, true, r
	}
	if k < m {
		ll, b, lr := split(l, k)
		return ll, b, join(lr, m, r)
	}
	rl, b, rr := split(r, k)
	return join(l, m, rl), b, rr
}

// splitLast separates a tree into two distinct trees at the rightmost node
func splitLast(tree *Node) (*Node, int64) {
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

// insert inserts a node into a tree
func insert(tree *Node, k int64) *Node {
	tree1, _, tree2 := split(tree, k)
	return join(tree1, k, tree2)
}

// deleteNode deletes a node from a tree
func deleteNode(tree *Node, k int64) *Node {
	tree1, _, tree2 := split(tree, k)
	return join2(tree1, tree2)
}

// union carries out the union operation on two trees
func union(tree1 *Node, tree2 *Node) *Node {
	if tree1 == nil {
		return tree2
	}
	if tree2 == nil {
		return tree1
	}
	l2, k2, r2 := expose(tree2)
	l1, _, r1 := split(tree1, k2)
	tree1 = union(l1, l2)
	tree2 = union(r1, r2)
	return join(tree1, k2, tree2)
}

// intersect carries out the intersect operation on two trees
func intersect(tree1 *Node, tree2 *Node) *Node {
	if tree1 == nil {
		return nil
	}
	if tree2 == nil {
		return nil
	}
	l2, k2, r2 := expose(tree2)
	l1, b, r1 := split(tree1, k2)
	tree1 = intersect(l1, l2)
	tree2 = intersect(r1, r2)
	if b {
		return join(tree1, k2, tree2)
	}
	return join2(tree1, tree2)
}

// difference carries out the difference operation on two trees
func difference(tree1 *Node, tree2 *Node) *Node {
	if tree1 == nil {
		return nil
	}
	if tree2 == nil {
		return tree1
	}
	l2, k2, r2 := expose(tree2)
	l1, _, r1 := split(tree1, k2)
	tree1 = difference(l1, l2)
	tree2 = difference(r1, r2)
	return join2(tree1, tree2)
}

// handleError takes care of generic errors
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// write writes data to a buffer and flushes the buffered data to disk
func write(data *string, w *bufio.Writer) {
	_, err := w.WriteString(*data)
	handleError(err)
	err2 := w.Flush()
	handleError(err2)
}

// queueNode is a struct to handle custom queue nodes
type queueNode struct {
	Parent *Node
	Node   *Node
	Dir    string
}

// getProperty gets the primary struct node property of the interface returned from a queue
func getProperty(node *interface{}, prop string) *Node {
	nodePtr := reflect.ValueOf(*node)
	item := reflect.Indirect(nodePtr)
	return item.FieldByName(prop).Interface().(*Node)
}

// getDirection gets the direction property of the struct node returned from a queue
func getDirection(node *interface{}) string {
	nodePtr := reflect.ValueOf(*node)
	item := reflect.Indirect(nodePtr)
	return item.FieldByName("Dir").Interface().(string)
}

// constructGraph creates a .dot file containing the tree structure of a binary tree
func constructGraph(filename string, root *Node) {
	root.populatePaths("N")

	var data string
	colors := make(map[string]string)
	colors["<RT>RT"] = "#FDF3D0"
	colors["<MD>MD"] = "#DCE8FA"
	colors["<LF>LF"] = "#F1CFCD"

	f, err := os.Create(filename + ".dot")
	handleError(err)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic("Could not close file")
		}
	}(f)

	// Create a buffer writer to batch write to disk
	w := bufio.NewWriter(f)

	data = "strict digraph {\n"
	write(&data, w)
	data = "node [shape=record];\n"
	write(&data, w)

	q := queue.New()
	q.Add(&queueNode{Parent: nil, Node: root, Dir: ""})
	for q.Length() > 0 {
		curr := q.Remove()

		parent := getProperty(&curr, "Parent")
		node := getProperty(&curr, "Node")
		dir := getDirection(&curr)

		var left string
		var right string
		if node.Left != nil {
			left = "<L>L"
		}
		if node.Right != nil {
			right = "<R>R"
		}

		nest := "<MD>MD"
		if parent == nil {
			nest = "<RT>RT"
		} else if node.Left == nil && node.Right == nil {
			nest = "<LF>LF"
		}

		data = fmt.Sprintf("%s [label=\"%s|{<C>%d|%s}|%s\" style=filled fillcolor=\"%s\"];\n", node.Path, left, node.Key, nest, right, colors[nest])
		write(&data, w)

		if parent != nil {
			data = fmt.Sprintf("%s:%s -> %s:C;\n", parent.Path, dir, node.Path)
			write(&data, w)
		}

		if node.Left != nil {
			q.Add(&queueNode{Parent: node, Node: node.Left, Dir: "L"})
		}
		if node.Right != nil {
			q.Add(&queueNode{Parent: node, Node: node.Right, Dir: "R"})
		}
	}

	data = "}\n"
	write(&data, w)
}

// visualize creates a .png file from a given .dot file
func visualize(filename string, root *Node) {
	constructGraph(filename, root)
	_, err := exec.Command("dot", "-Tpng", filename+".dot", "-o", filename+".png").Output()
	if err != nil {
		log.Fatal(err)
	}
}

// createTree creates a balanced binary search tree from an unsorted array of integers
func createTree(arr []int64) *Node {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	n := len(arr)
	height := int(math.Floor(math.Log2(float64(n))))

	return _createTree(arr, height)
}

// _createTree is a helper function to createTree
func _createTree(arr []int64, height int) *Node {
	if len(arr) == 0 {
		return nil
	}

	mid := len(arr) / 2

	root := NewNode(arr[mid], nil, nil)
	root.Height = height
	root.Left = _createTree(arr[:mid], height-1)
	root.Right = _createTree(arr[mid+1:], height-1)

	return root
}

func main() {
	// Test automatic tree creation
	root := createTree([]int64{4, 3, 2, 1})
	visualize("created_tree", root)

	// Test for the join operation
	tree1 := createTree([]int64{5, 6, 7, 8, 9, 10, 11, 12, 13})
	visualize("tree1", tree1)

	var k int64 = 50

	tree2 := createTree([]int64{34, 21, 78, 345, 1009, 23, 16, 30, 45, 66, 74, 32})
	visualize("tree2", tree2)

	joined := join(tree1, k, tree2)
	visualize("joined", joined)

	// Test for the union operation
	tree3 := createTree([]int64{157, 11, 19})
	visualize("tree3", tree3)

	union1 := union(tree3, joined)
	visualize("union1", union1)

	// Test for the difference operation
	tree4 := createTree([]int64{1, 4, 5})
	visualize("tree4", tree4)

	tree5 := createTree([]int64{3, 2, 7})
	visualize("tree5", tree5)

	union2 := union(tree4, tree5)
	visualize("union2", union2)

	difference1 := difference(union2, tree4)
	visualize("difference1", difference1)

	difference2 := difference(union2, tree5)
	visualize("difference2", difference2)
}
