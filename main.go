package main

import (
	avl2 "bulkOperations/cairo-avl"
	"fmt"
)

func main() {
	set := make(map[string]bool)

	//temp1, temp2 := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, []byte{14, 15, 16, 17, 18, 19, 20}
	//temp1, temp2 := []byte{1, 2, 3, 4, 5, 6}, []byte{14, 15, 16, 17, 18, 19}
	temp1, temp2 := []byte("\x1d\xff\xbc \x88"), []byte("\x00\x00\x00d\xbc")

	b1 := make([][]byte, 0)
	if temp1 != nil {
		b1 = *(avl2.EmbedByteArray(temp1, &set))
	}

	b2 := make([][]byte, 0)
	if temp1 != nil {
		b2 = *(avl2.EmbedByteArray(temp2, &set))
	}

	t1 := avl2.BuildTreeFromInorder(&b1)
	t2 := avl2.BuildDictTreeFromInorder(&b2)
	t2NodeType := (*t2).ConvertToNode()

	avl2.VisualizeNodeTree("t1", t1)
	avl2.VisualizeNodeTree("t2", t2NodeType)

	tU := avl2.Union(t1, t2)
	avl2.VisualizeNodeTree("tu", tU)

	// Check that all nodes in tU are either in t1 or t2
	for _, key := range *(avl2.GetInorderTraversal(tU)) {
		if !avl2.IsInTree(t1, &key) {
			if !avl2.IsInTree(t2NodeType, &key) {
				fmt.Printf("Key: %v in tU not in t1 nor t2", key)
				fmt.Println()
			}
		}
	}

	// Check that all nodes in t1 are in tU
	for _, key := range *(avl2.GetInorderTraversal(t1)) {
		if !avl2.IsInTree(tU, &key) {
			fmt.Printf("Key: %v in t1 not in tU", key)
			fmt.Println()
		}
	}

	// Check that all nodes in t1 are in tU
	for _, key := range *(avl2.GetInorderTraversal(t2NodeType)) {
		if !avl2.IsInTree(tU, &key) {
			fmt.Printf("Key: %v in t2 not in tU", key)
			fmt.Println()
		}
	}

	TL := avl2.Node{
		Key:     []byte{0, 0},
		Value:   []byte{0, 0},
		Left:    nil,
		Right:   nil,
		Nested:  nil,
		Height:  1,
		Path:    "",
		Exposed: false,
	}

	LF := avl2.Node{
		Key:     []byte{29, 255},
		Value:   []byte{29, 255},
		Left:    nil,
		Right:   nil,
		Nested:  nil,
		Height:  1,
		Path:    "",
		Exposed: false,
	}

	TRL := avl2.Node{
		Key:     []byte{136},
		Value:   []byte{136},
		Left:    &LF,
		Right:   nil,
		Nested:  nil,
		Height:  2,
		Path:    "",
		Exposed: false,
	}

	TRR := avl2.Node{
		Key:     []byte{188, 32},
		Value:   []byte{188, 32},
		Left:    nil,
		Right:   nil,
		Nested:  nil,
		Height:  1,
		Path:    "",
		Exposed: false,
	}

	TR := avl2.Node{
		Key:     []byte{188},
		Value:   []byte{188},
		Left:    &TRL,
		Right:   &TRR,
		Nested:  nil,
		Height:  3,
		Path:    "",
		Exposed: false,
	}

	avl2.TestJoinLeft([]byte{0, 100}, []byte{0, 100}, &TL, &TR, nil)
}
