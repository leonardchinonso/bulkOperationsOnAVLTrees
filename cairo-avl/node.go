package cairo_avl

import (
	"math"
)

// Node node representation of the data in a TreeNode object format
type Node struct {
	Key     []byte
	Value   []byte
	Left    *Node
	Right   *Node
	Nested  *Node
	Height  int
	Path    string
	Exposed bool
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

// NewNode is a custom constructor method to initialise height of the node
func NewNode(key []byte, value []byte, h int, left, right, nested *Node) *Node {
	node := new(Node)
	node.Key = key
	node.Value = value
	node.Left = left
	node.Right = right
	node.Nested = nested
	node.Height = h
	node.Exposed = true
	return node
}

func (n *Node) ConvertToDictNode() *DictNode {
	if n == nil {
		return nil
	}
	return NewDictNode(n.Key, n.Value, n.Height, n.Left.ConvertToDictNode(), n.Right.ConvertToDictNode())
}

// exposeNode opens up a node type
func exposeNode(tree *Node, numOfExposedNodes *int) (k []byte, v []byte, TL *Node, TR *Node, TN *Node) {
	if tree != nil {
		if !tree.Exposed {
			*numOfExposedNodes++
			tree.Exposed = true
		}
		return tree.Key, tree.Value, tree.Left, tree.Right, tree.Nested
	}
	return []byte{}, []byte{}, nil, nil, nil
}

// insertNode inserts a node into the tree
func insertNode(T *Node, k []byte) *Node {
	numOfExposedNodes := 0
	TL, TR, TN := split(T, k, &numOfExposedNodes)
	return join(k, k, nil, nil, TL, TR, TN, &numOfExposedNodes)
}

func _createTree(arr *[][]byte) *Node {
	var root *Node
	for _, b := range *arr {
		root = insertNode(root, b)
	}
	return root
}

func BuildTreeFromInorder(arr *[][]byte) *Node {
	sortArray(arr)
	n := len(*arr)
	height := int(math.Floor(math.Log2(float64(n))))
	root := _buildTreeFromInorder(arr, height+1)
	setExposure(root, false)
	return root
}

func _buildTreeFromInorder(arr *[][]byte, height int) *Node {
	if len(*arr) == 0 {
		return nil
	}
	mid := len(*arr) / 2

	root := NewNode((*arr)[mid], (*arr)[mid], height, nil, nil, nil)
	leftArr := (*arr)[:mid]
	rightArr := (*arr)[mid+1:]
	root.Left = _buildTreeFromInorder(&leftArr, height-1)
	root.Right = _buildTreeFromInorder(&rightArr, height-1)

	return root
}

func setExposure(root *Node, b bool) {
	if root == nil {
		return
	}
	setExposure(root.Left, b)
	root.Exposed = false
	setExposure(root.Right, b)
}

func CreateTree(arr *[][]byte) *Node {
	root := _createTree(arr)
	setExposure(root, false)
	return root
}

func CountNumberOfNewHashes(root *Node, newNodesCount *int) {
	if root == nil {
		return
	}
	if root.Exposed {
		*newNodesCount++
	}
	CountNumberOfNewHashes(root.Left, newNodesCount)
	CountNumberOfNewHashes(root.Right, newNodesCount)
}
