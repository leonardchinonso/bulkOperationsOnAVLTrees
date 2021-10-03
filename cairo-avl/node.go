package cairo_avl

import (
	"math"
)

// Node node representation of the data in a TreeNode object format
type Node struct {
	Key         []byte
	Value       []byte
	Left        *Node
	Right       *Node
	Nested      *Node
	Height      int
	Path        string
	Exposed     bool
	HeightTaken bool
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
	node.HeightTaken = true
	return node
}

func (n *Node) ConvertToDictNode() *DictNode {
	if n == nil {
		return nil
	}
	return NewDictNode(n.Key, n.Value, n.Height, n.Left.ConvertToDictNode(), n.Right.ConvertToDictNode())
}

// exposeNode opens up a node type
func exposeNode(tree *Node, numOfExposedNodes *int, numOfHeightTakenNodes *int) (k []byte, v []byte, TL *Node, TR *Node, TN *Node) {
	if tree != nil {
		if !tree.Exposed {
			if tree.HeightTaken {
				*numOfHeightTakenNodes--
			} else {
				tree.HeightTaken = true
			}
			*numOfExposedNodes++
			tree.Exposed = true
		}
		return tree.Key, tree.Value, tree.Left, tree.Right, tree.Nested
	}
	return []byte{}, []byte{}, nil, nil, nil
}

// HeightOf Get height of a node
// This function exists because the join method can take a null Node pointer
// and accessing the height property of a null Node pointer will fail
func HeightOf(node *Node, numOfHeightTakenNodes *int) int {
	if node == nil {
		return 0
	}
	if numOfHeightTakenNodes != nil && !node.Exposed && !node.HeightTaken {
		*numOfHeightTakenNodes++
	}
	node.HeightTaken = true
	return node.Height
}

// insertNode inserts a node into the tree
func insertNode(T *Node, k []byte) *Node {
	numOfExposedNodes := 0
	numOfHeightTakenNodes := 0
	TL, TR, TN := split(T, k, &numOfExposedNodes, &numOfHeightTakenNodes)
	return join(k, k, nil, nil, TL, TR, TN, &numOfExposedNodes, &numOfHeightTakenNodes)
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
	setExposureAndHeightTaken(root, false)
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

func setExposureAndHeightTaken(root *Node, b bool) {
	if root == nil {
		return
	}
	setExposureAndHeightTaken(root.Left, b)
	root.Exposed = false
	root.HeightTaken = false
	setExposureAndHeightTaken(root.Right, b)
}

func CreateTree(arr *[][]byte) *Node {
	root := _createTree(arr)
	setExposureAndHeightTaken(root, false)
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
