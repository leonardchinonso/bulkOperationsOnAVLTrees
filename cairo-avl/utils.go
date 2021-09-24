package cairo_avl

import (
	"bufio"
	"bytes"
	"fmt"
	"gopkg.in/eapache/queue.v1"
	"log"
	"os"
	"os/exec"
	"reflect"
	"sort"
)

func absInt(num int) int {
	if num < 0 {
		return num
	}
	return num
}

//// EmbedByteArray pairs bytes up as single values and removes duplicate byte pair entries
//func EmbedByteArray(b []byte, set *map[string]bool) *[][]byte {
//	var res [][]byte
//	for i := 0; i < len(b); i += 2 {
//		if i == len(b)-1 {
//			if _, ok := (*set)[string([]byte{b[i]})]; !ok {
//				res = append(res, []byte{b[i]})
//				(*set)[string([]byte{b[i]})] = true
//			}
//		} else {
//			if _, ok := (*set)[string([]byte{b[i], b[i+1]})]; !ok {
//				res = append(res, []byte{b[i], b[i+1]})
//				(*set)[string([]byte{b[i], b[i+1]})] = true
//			}
//		}
//	}
//
//	return &res
//}

// EmbedByteArray pairs bytes up as single values and removes duplicate byte pair entries
func EmbedByteArray(b []byte) *[][]byte {
	set := make(map[string]bool)

	var res [][]byte
	for i := 0; i < len(b); i += 4 {
		if len(b)-i <= 4 {
			if _, ok := set[string(b[i:])]; !ok {
				res = append(res, b[i:])
				set[string(b[i:])] = true
			}
		} else {
			if _, ok := set[string(b[i:i+4])]; !ok {
				res = append(res, b[i:i+4])
				set[string(b[i:i+4])] = true
			}
		}
	}

	return &res
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

// VisualizeNodeTree creates a .png file from a given .dot file
func VisualizeNodeTree(filename string, root *Node) {
	if root == nil {
		return
	}
	constructGraph(filename, root)
	_, err := exec.Command("dot", "-Tpng", filename+".dot", "-o", filename+".png").Output()
	if err != nil {
		log.Fatal(err)
	}
}

func IsInTree(root *Node, key *[]byte) bool {
	//if bytes.Compare(root.Key, []byte{236, 104}) == 0 {
	//	fmt.Println(root, key)
	//}
	if root == nil || key == nil {
		return false
	}
	if bytes.Compare(root.Key, *key) == 0 {
		return true
	}
	if bytes.Compare(root.Key, *key) == -1 {
		return IsInTree(root.Right, key)
	}
	return IsInTree(root.Left, key)
}

func GetInorderTraversal(root *Node) *[][]byte {
	arr := make([][]byte, 0)
	_getInorderTraversal(root, &arr)
	return &arr
}

func _getInorderTraversal(root *Node, arr *[][]byte) {
	if root == nil {
		return
	}
	if root.Left != nil {
		_getInorderTraversal(root.Left, arr)
	}
	*arr = append(*arr, root.Key)
	if root.Right != nil {
		_getInorderTraversal(root.Right, arr)
	}
}

func IsBalanced(root *Node) bool {
	if root == nil {
		return true
	}

	leftHeight := HeightOf(root.Left)
	rightHeight := HeightOf(root.Right)

	if absInt(leftHeight-rightHeight) <= 1 && IsBalanced(root.Left) && IsBalanced(root.Right) {
		return true
	}

	return false
}

func IsValidBST(root *Node) bool {
	if root == nil {
		return true
	}

	if root.Left != nil && bytes.Compare(root.Left.Key, root.Key) == 1 {
		return false
	}

	if root.Right != nil && bytes.Compare(root.Right.Key, root.Key) == -1 {
		return false
	}

	if !IsValidBST(root.Left) || !IsValidBST(root.Right) {
		return false
	}

	return true
}

func sortArray(arr *[][]byte) {
	sort.Slice(*arr, func(i, j int) bool {
		return bytes.Compare((*arr)[i], (*arr)[j]) == -1
	})
}

func PrintTree(root *Node, space string) {
	if root != nil {
		fmt.Println(space, root)
		PrintTree(root.Left, space+space)
		PrintTree(root.Right, space+space)
	}
}

func printDictTree(root *DictNode, space string) {
	if root != nil {
		fmt.Println(space, root)
		printDictTree(root.Left, space+space)
		printDictTree(root.Right, space+space)
	}
}
