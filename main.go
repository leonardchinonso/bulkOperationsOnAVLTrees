package main

import (
	"bufio"
	avl2 "bulkOperations/cairo-avl"
	"bytes"
	"fmt"
	"io"
	"os"
)

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func convertToBytes(fileName string) []byte {
	file, err := os.Open(fileName)
	handleError(err)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Can't close file")
		}
	}(file)

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))

	var chunk []byte
	var eol bool
	byteArr := make([]byte, 0)

	for {
		if chunk, eol, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(chunk)
		if !eol {
			byteArr = append(byteArr, buffer.Bytes()...)
			buffer.Reset()
		}
	}

	if err == io.EOF {
		err = nil
	}

	return byteArr
}

func main() {
	if len(os.Args) != 3 {
		panic("Invalid number of arguments")
	}

	set := make(map[string]bool)

	b1 := *(avl2.EmbedByteArray(convertToBytes(os.Args[1]), &set))
	b2 := *(avl2.EmbedByteArray(convertToBytes(os.Args[2]), &set))

	if len(b1) == 0 || len(b2) == 0 {
		panic("One of the input files does not have enough unique data!")
	}

	t1 := avl2.BuildTreeFromInorder(&b1)
	t2 := avl2.BuildDictTreeFromInorder(&b2)

	t2NodeType := (*t2).ConvertToNode()

	tU, numOfExposedNodesInUnion := avl2.Union(t1, t2)

	newNodesCount := 0
	avl2.CountNumberOfNewHashes(tU, &newNodesCount)
	fmt.Println("number of new nodes created: ", newNodesCount)
	fmt.Println("number of nodes exposed in union: ", numOfExposedNodesInUnion)

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

	if !avl2.IsValidBST(tU) {
		fmt.Println("tU is an invalid BST")
	}

	if !avl2.IsBalanced(tU) {
		fmt.Println("tU is an unbalanced BST")
	}

	// DIFFERENCE: Yet to count the hashes for the difference operation
	fmt.Println()
	tD := avl2.Difference(t1, t2)

	// Check that all nodes in t1 are either in tD or t2 but not both
	for _, key := range *(avl2.GetInorderTraversal(t1)) {
		in_tD := avl2.IsInTree(tD, &key)
		in_t2 := avl2.IsInTree(t2NodeType, &key)

		if !in_tD && !in_t2 {
			fmt.Printf("Key: %v not in tD and not in t2", key)
		}

		if in_tD && in_t2 {
			fmt.Printf("Key: %v in tD and t2", key)
		}
	}

	// Check that all nodes in t2 are either in tD or t1 but not both
	for _, key := range *(avl2.GetInorderTraversal(t2NodeType)) {
		in_tD := avl2.IsInTree(tD, &key)
		in_t1 := avl2.IsInTree(t2NodeType, &key)

		if !in_tD && !in_t1 {
			fmt.Printf("Key: %v not in tD and not in t1", key)
		}

		if in_tD && in_t1 {
			fmt.Printf("Key: %v in tD and t1", key)
		}
	}

	if !avl2.IsValidBST(tD) {
		fmt.Println("tD is an invalid BST")
	}

	if !avl2.IsBalanced(tD) {
		fmt.Println("tD is an unbalanced BST")
	}

}
