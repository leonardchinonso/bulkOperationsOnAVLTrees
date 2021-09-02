package main

import (
	"bulkOperations/avl"
	"fmt"
)

func main() {
	b := []byte("\x99]kd\xech\xec\xec\xec\xec\xef\x05n\x90\x0e\x0e\x0e\f\x13\v\x12\x03\x10d\x06\x0e\x0e\x0e\x11\x0e\x0e\xef,\xb46\x1c\x1c\x10\x1a%R\x88\xf2X\x14")
	temp1, temp2 := avl.SplitByteArray(&b)
	fmt.Println(temp1, temp2)

	set := make(map[string]bool)

	b1 := make([][]byte, 0)
	if temp1 != nil {
		b1 = *(avl.EmbedByteArray(*temp1, &set))
	}

	b2 := make([][]byte, 0)
	if temp1 != nil {
		b2 = *(avl.EmbedByteArray(*temp2, &set))
	}

	t1 := avl.CreateTree(&b1)
	t2 := avl.CreateTree(&b2)

	avl.Visualize("t1", t1)
	avl.Visualize("t2", t2)

	tU := avl.Union(t1, t2)
	avl.Visualize("tu", tU)

	// Check that all nodes in tU are either in t1 or t2
	for _, key := range *(avl.GetInorderTraversal(tU)) {
		fmt.Println(key)
		if !avl.IsInTree(t1, &key) {
			if !avl.IsInTree(t2, &key) {
				fmt.Printf("Key: %v in tU not in t1 nor t2", key)
			}
		}
	}

	// Check that all nodes in t1 are in tU
	for _, key := range *(avl.GetInorderTraversal(t1)) {
		fmt.Println(key)
		if !avl.IsInTree(tU, &key) {
			fmt.Printf("Key: %v in t1 not in tU", key)
		}
	}

	// Check that all nodes in t1 are in tU
	for _, key := range *(avl.GetInorderTraversal(t2)) {
		fmt.Println(key)
		if !avl.IsInTree(tU, &key) {
			fmt.Printf("Key: %v in t2 not in tU", key)
		}
	}
}
