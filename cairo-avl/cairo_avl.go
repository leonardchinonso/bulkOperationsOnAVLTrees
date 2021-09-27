package cairo_avl

import (
	"bytes"
)

// HeightOf Get height of a node
// This function exists because the join method can take a null Node pointer
// and accessing the height property of a null Node pointer will fail
func HeightOf(node *Node) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func balancedHeight(hL int, hR int) int {
	if hL == hR {
		return hL + 1
	} else if hL == hR+1 {
		return hL + 1
	}
	return hR + 1
}

// rotateRight rotates a node to the right to maintain the AVL balance criteria
func rotateLeft(k []byte, v []byte, TL *Node, TR *Node, TN *Node, numOfExposedNodes *int) (int, *Node) {
	kR, vR, TRL, TRR, TRN := exposeNode(TR, numOfExposedNodes)
	hL := HeightOf(TL)
	hRL := HeightOf(TRL)
	hRR := HeightOf(TRR)
	hP := balancedHeight(hL, hRL)
	TP := NewNode(k, v, hP, TL, TRL, TN)
	h := balancedHeight(hP, hRR)
	return h, NewNode(kR, vR, h, TP, TRR, TRN)
}

// rotateLeft rotates a node to the left to maintain the AVL balance criteria
func rotateRight(k []byte, v []byte, TL *Node, TR *Node, TN *Node, numOfExposedNodes *int) (int, *Node) {
	kL, vL, TLL, TLR, TLN := exposeNode(TL, numOfExposedNodes)
	hR := HeightOf(TR)
	hLL := HeightOf(TLL)
	hLR := HeightOf(TLR)
	hP := balancedHeight(hR, hLR)
	TP := NewNode(k, v, hP, TLR, TR, TN)
	h := balancedHeight(hP, hLL)
	return h, NewNode(kL, vL, h, TLL, TP, TLN)
}

// joinRight concatenates a left tree, k and a right tree
func joinRight(k []byte, v []byte, TL *Node, TR *Node, TN *Node, numOfExposedNodes *int) (int, *Node) {
	kL, vL, TLL, TLR, TLN := exposeNode(TL, numOfExposedNodes)
	hLR := HeightOf(TLR)
	hR := HeightOf(TR)
	hLL := HeightOf(TLL)
	if hLR <= hR+1 {
		hP := balancedHeight(hLR, hR)
		if hP <= hLL+1 {
			h := balancedHeight(hLL, balancedHeight(hLR, hR))
			return h, NewNode(kL, vL, h, TLL, NewNode(k, v, hP, TLR, TR, TN), TLN)
		}
		_, TP := rotateRight(k, v, TLR, TR, TN, numOfExposedNodes)
		return rotateLeft(kL, vL, TLL, TP, TLN, numOfExposedNodes)
	}
	hP, TP := joinRight(k, v, TLR, TR, TN, numOfExposedNodes)
	if hP <= hLL+1 {
		h := balancedHeight(hP, hLL)
		return h, NewNode(kL, vL, h, TLL, TP, TLN)
	}
	return rotateLeft(kL, vL, TLL, TP, TLN, numOfExposedNodes)
}

func joinLeft(k []byte, v []byte, TL *Node, TR *Node, TN *Node, numOfExposedNodes *int) (int, *Node) {
	kR, vR, TRL, TRR, TRN := exposeNode(TR, numOfExposedNodes)
	hRL := HeightOf(TRL)
	hL := HeightOf(TL)
	hRR := HeightOf(TRR)
	if hRL <= hL+1 {
		hP := balancedHeight(hL, hRL)
		if hP <= hRR+1 {
			h := balancedHeight(hRR, balancedHeight(hL, hRL))
			return h, NewNode(kR, vR, h, NewNode(k, v, hP, TL, TRL, TN), TRR, TRN)
		}
		_, TP := rotateLeft(k, v, TL, TRL, TN, numOfExposedNodes)
		return rotateRight(kR, vR, TP, TRR, TRN, numOfExposedNodes)
	}
	hP, TP := joinLeft(k, v, TL, TRL, TN, numOfExposedNodes)
	if hP <= hRR+1 {
		h := balancedHeight(hP, hRR)
		return h, NewNode(kR, vR, h, TP, TRR, TRN)
	}
	return rotateRight(kR, vR, TP, TRR, TRN, numOfExposedNodes)
}

func join(k []byte, v []byte, DU *DictNode, DD *DictNode, TL *Node, TR *Node, TN *Node, numOfExposedNodes *int) *Node {
	hL := HeightOf(TL)
	hR := HeightOf(TR)
	if hL > hR+1 {
		_, T := joinRight(k, v, TL, TR, TN, numOfExposedNodes)
		return T
	}
	if hR > hL+1 {
		_, T := joinLeft(k, v, TL, TR, TN, numOfExposedNodes)
		return T
	}
	N := Union(Difference(TN, DU, numOfExposedNodes), DD, numOfExposedNodes)
	h := balancedHeight(hL, hR)
	return NewNode(k, v, h, TL, TR, N)
}

func splitLast(T *Node, numOfExposedNodes *int) (*Node, []byte, []byte, *Node) {
	m, v, L, R, N := exposeNode(T, numOfExposedNodes)
	if R == nil {
		return L, m, v, N
	}

	TP, kP, vP, NP := splitLast(R, numOfExposedNodes)
	return join(m, v, nil, nil, L, TP, N, numOfExposedNodes), kP, vP, NP
}

func join2(TL *Node, TR *Node, numOfExposedNodes *int) *Node {
	if TL == nil {
		return TR
	}
	TLP, k, v, N := splitLast(TL, numOfExposedNodes)
	return join(k, v, nil, nil, TLP, TR, N, numOfExposedNodes)
}

func split(t *Node, k []byte, numOfExposedNodes *int) (*Node, *Node, *Node) {
	if t == nil {
		return nil, nil, nil
	}

	m, v, L, R, N := exposeNode(t, numOfExposedNodes)
	if bytes.Compare(k, m) == 0 {
		return L, R, N
	}

	if bytes.Compare(k, m) == -1 {
		LL, LR, LN := split(L, k, numOfExposedNodes)
		return LL, join(m, v, nil, nil, LR, R, N, numOfExposedNodes), LN
	}

	RL, RR, RN := split(R, k, numOfExposedNodes)
	return join(m, v, nil, nil, L, RL, N, numOfExposedNodes), RR, RN
}

func Union(T0 *Node, D *DictNode, numOfExposedNodes *int) *Node {
	if T0 == nil {
		return D.ConvertToNode()
	}
	if D == nil {
		return T0
	}

	k, v, DL, DR, DU, DD := exposeDict(D)
	TL, TR, TN := split(T0, k, numOfExposedNodes)
	L := Union(TL, DL, numOfExposedNodes)
	R := Union(TR, DR, numOfExposedNodes)
	joined := join(k, v, DU, DD, L, R, TN, numOfExposedNodes)
	return joined
}

func Difference(T0 *Node, D *DictNode, numOfExposedNodes *int) *Node {
	if T0 == nil {
		return nil
	}
	if D == nil {
		return T0
	}

	k, _, DL, DR, _, _ := exposeDict(D)
	TL, TR, _ := split(T0, k, numOfExposedNodes)
	L := Difference(TL, DL, numOfExposedNodes)
	R := Difference(TR, DR, numOfExposedNodes)
	return join2(L, R, numOfExposedNodes)
}
