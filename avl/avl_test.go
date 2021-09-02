package avl

import (
	"testing"
)

func FuzzUnion(f *testing.F) {

	f.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})

	f.Fuzz(func(t *testing.T, b []byte) {
		set := make(map[string]bool)

		temp1, temp2 := SplitByteArray(&b)

		b1 := make([][]byte, 0)
		if temp1 != nil {
			b1 = *(EmbedByteArray(*temp1, &set))
		}

		b2 := make([][]byte, 0)
		if temp1 != nil {
			b2 = *(EmbedByteArray(*temp2, &set))
		}

		t1 := CreateTree(&b1)
		t2 := CreateTree(&b2)

		tU := Union(t1, t2)

		// Check that all nodes in tU are either in t1 or t2
		for _, key := range *(GetInorderTraversal(tU)) {
			if !IsInTree(t1, &key) {
				if !IsInTree(t2, &key) {
					t.Fatalf("Key: %v in tU not in t1 nor t2", key)
				}
			}
		}

		// Check that all nodes in t1 are in tU
		for _, key := range *(GetInorderTraversal(t1)) {
			if !IsInTree(tU, &key) {
				t.Fatalf("Key: %v in t1 not in tU", key)
			}
		}

		// Check that all nodes in t1 are in tU
		for _, key := range *(GetInorderTraversal(t2)) {
			if !IsInTree(tU, &key) {
				t.Fatalf("Key: %v in t2 not in tU", key)
			}
		}

		// Check that tU is balanced
		if !IsBalanced(tU) {
			t.Fatalf("tU with root: %v is height unbalanced", tU)
		}

		// Check that tU is a valid BST
		if !IsValidBST(tU) {
			t.Fatalf("tU with root: %v is not a valid BST", tU)
		}
	})
}

func FuzzDifference(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})

	f.Fuzz(func(t *testing.T, b []byte) {
		set := make(map[string]bool)

		temp1, temp2 := SplitByteArray(&b)

		b1 := make([][]byte, 0)
		if temp1 != nil {
			b1 = *(EmbedByteArray(*temp1, &set))
		}

		b2 := make([][]byte, 0)
		if temp1 != nil {
			b2 = *(EmbedByteArray(*temp2, &set))
		}

		t1 := CreateTree(&b1)
		t2 := CreateTree(&b2)

		tD := Difference(t1, t2)

		// Check that all nodes in t1 are either in tD or t2 but not both
		for _, key := range *(GetInorderTraversal(t1)) {
			in_tD := IsInTree(tD, &key)
			in_t2 := IsInTree(t2, &key)

			if !in_tD && !in_t2 {
				t.Fatalf("Key: %v not in tD and not in t2", key)
			}

			if in_tD && in_t2 {
				t.Fatalf("Key: %v in tD and t2", key)
			}
		}

		// Check that all nodes in t2 are either in tD or t1 but not both
		for _, key := range *(GetInorderTraversal(t2)) {
			in_tD := IsInTree(tD, &key)
			in_t1 := IsInTree(t2, &key)

			if !in_tD && !in_t1 {
				t.Fatalf("Key: %v not in tD and not in t1", key)
			}

			if in_tD && in_t1 {
				t.Fatalf("Key: %v in tD and t1", key)
			}
		}
	})
}
