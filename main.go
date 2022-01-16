package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	avl2 "github.com/leonardchinonso/bulkOperations/cairo-avl"
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

	b1 := *(avl2.EmbedByteArray(convertToBytes(os.Args[1])))
	b2 := *(avl2.EmbedByteArray(convertToBytes(os.Args[2])))

	if len(b1) == 0 || len(b2) == 0 {
		panic("One of the input files does not have enough unique data!")
	}

	fmt.Println("# UNION:")
	t1 := avl2.BuildTreeFromInorder(&b1)
	t2 := avl2.BuildDictTreeFromInorder(&b2)

	t2NodeType := (*t2).ConvertToNode()

	numOfExposedNodesInUnion := 0
	numOfHeightTakenNodesInUnion := 0
	numOfRevisitedNodesInUnion := 0
	newNodesCountInUnion := 0
	tU := avl2.Union(t1, t2, &numOfExposedNodesInUnion, &numOfHeightTakenNodesInUnion, &numOfRevisitedNodesInUnion)
	avl2.CountNumberOfNewHashes(tU, &newNodesCountInUnion)
	fmt.Println("Number of nodes in the original tree: ", len(b1))
	fmt.Println("Number of nodes in the update tree: ", len(b2))
	fmt.Println("number of re-hashes to be made: ", newNodesCountInUnion)
	fmt.Println("number of nodes exposed in union: ", numOfExposedNodesInUnion)
	fmt.Println("number of hashes required to expose nodes in union: ", numOfExposedNodesInUnion*3)
	fmt.Println("number of nodes with height taken in union: ", numOfHeightTakenNodesInUnion)
	fmt.Println("number of hashes required for taking heights in union: ", numOfHeightTakenNodesInUnion*2)
	fmt.Println("number of nodes visited more than once: ", numOfRevisitedNodesInUnion)

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
	fmt.Println("# DIFFERENCE:")

	tt1 := avl2.BuildTreeFromInorder(&b1)
	tt2 := avl2.BuildDictTreeFromInorder(&b2)

	tt2NodeType := (*tt2).ConvertToNode()

	numOfExposedNodesInDifference := 0
	numOfHeightTakenNodesInDifference := 0
	numOfRevisitedNodesInDifference := 0
	newNodesCountInDifference := 0
	tD := avl2.Difference(tt1, tt2, &numOfExposedNodesInDifference, &numOfHeightTakenNodesInDifference, &numOfRevisitedNodesInDifference)
	avl2.CountNumberOfNewHashes(tD, &newNodesCountInDifference)
	fmt.Println("Number of nodes in the original tree: ", len(b1))
	fmt.Println("Number of nodes in the update tree: ", len(b2))
	fmt.Println("number of re-hashes to be made: ", newNodesCountInDifference)
	fmt.Println("number of nodes exposed in difference: ", numOfExposedNodesInDifference)
	fmt.Println("number of hashes required to expose nodes in difference: ", numOfExposedNodesInDifference*3)
	fmt.Println("number of nodes with height taken in difference: ", numOfHeightTakenNodesInDifference)
	fmt.Println("number of hashes required for taking heights in difference: ", numOfHeightTakenNodesInDifference*2)
	fmt.Println("number of nodes visited more than once: ", numOfRevisitedNodesInDifference)

	// Check that all nodes in t1 are either in tD or t2 but not both
	for _, key := range *(avl2.GetInorderTraversal(tt1)) {
		in_tD := avl2.IsInTree(tD, &key)
		in_tt2 := avl2.IsInTree(tt2NodeType, &key)

		if !in_tD && !in_tt2 {
			fmt.Printf("Key: %v not in tD and not in tt2", key)
		}

		if in_tD && in_tt2 {
			fmt.Printf("Key: %v in tD and tt2", key)
		}
	}

	// Check that all nodes in t2 are either in tD or t1 but not both
	for _, key := range *(avl2.GetInorderTraversal(tt2NodeType)) {
		in_tD := avl2.IsInTree(tD, &key)
		in_tt1 := avl2.IsInTree(tt2NodeType, &key)

		if !in_tD && !in_tt1 {
			fmt.Printf("Key: %v not in tD and not in tt1", key)
		}

		if in_tD && in_tt1 {
			fmt.Printf("Key: %v in tD and tt1", key)
		}
	}

	if !avl2.IsValidBST(tD) {
		fmt.Println("tD is an invalid BST")
	}

	if !avl2.IsBalanced(tD) {
		fmt.Println("tD is an unbalanced BST")
	}

}
