package cairo_avl

// DictNode dict representation of the data in a TreeNode object format
type DictNode struct {
	Key    []byte
	Value  []byte
	Left   *DictNode
	Right  *DictNode
	Update *DictNode
	Delete *DictNode
	Height int
	Path   string
}

// NewDictNode is a custom constructor method to initialise the DictNode
func NewDictNode(key []byte, value []byte, h int, left, right *DictNode) *DictNode {
	node := new(DictNode)
	node.Key = key
	node.Value = value
	node.Left = left
	node.Right = right
	node.Update = nil
	node.Delete = nil
	node.Height = h
	return node
}

func (d *DictNode) ConvertToNode() *Node {
	if d == nil {
		return nil
	}
	return NewNode(d.Key, d.Value, d.Height, d.Left.ConvertToNode(), d.Right.ConvertToNode(), nil)
}

// exposeDict opens up a dict type
func exposeDict(d *DictNode) (k []byte, v []byte, DL *DictNode, DR *DictNode, DU *DictNode, DD *DictNode) {
	if d != nil {
		return d.Key, d.Value, d.Left, d.Right, d.Update, d.Delete
	}
	return []byte{}, []byte{}, nil, nil, nil, nil
}

// CreateDictTree creates a tree of the DictNode type from a list of bytes
func CreateDictTree(arr *[][]byte) *DictNode {
	var temp *Node
	for _, b := range *arr {
		temp = insertNode(temp, b)
	}

	var root *DictNode
	if temp != nil {
		root = temp.ConvertToDictNode()
	}
	return root
}

func BuildDictTreeFromInorder(arr *[][]byte) *DictNode {
	temp := BuildTreeFromInorder(arr)
	var root *DictNode
	if temp != nil {
		root = temp.ConvertToDictNode()
	}
	return root
}
