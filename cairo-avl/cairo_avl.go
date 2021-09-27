package cairo_avl

import (
	"bytes"
)

func balancedHeight(hL int, hR int) int {
	if hL == hR {
		return hL + 1
	} else if hL == hR+1 {
		return hL + 1
	}
	return hR + 1
}

// rotateRight rotates a node to the right to maintain the AVL balance criteria
func rotateLeft(k []byte, v []byte, TL *Node, TR *Node, TN *Node, numOfExposedNodes *int, numOfHeightTakenNodes *int) (int, *Node) {
	kR, vR, TRL, TRR, TRN := exposeNode(TR, numOfExposedNodes)
	hL := HeightOf(TL, numOfHeightTakenNodes)
	hRL := HeightOf(TRL, numOfHeightTakenNodes)
	hRR := HeightOf(TRR, numOfHeightTakenNodes)
	hP := balancedHeight(hL, hRL)
	TP := NewNode(k, v, hP, TL, TRL, TN)
	h := balancedHeight(hP, hRR)
	return h, NewNode(kR, vR, h, TP, TRR, TRN)
}

// rotateLeft rotates a node to the left to maintain the AVL balance criteria
func rotateRight(k []byte, v []byte, TL *Node, TR *Node, TN *Node, numOfExposedNodes *int, numOfHeightTakenNodes *int) (int, *Node) {
	kL, vL, TLL, TLR, TLN := exposeNode(TL, numOfExposedNodes)
	hR := HeightOf(TR, numOfHeightTakenNodes)
	hLL := HeightOf(TLL, numOfHeightTakenNodes)
	hLR := HeightOf(TLR, numOfHeightTakenNodes)
	hP := balancedHeight(hR, hLR)
	TP := NewNode(k, v, hP, TLR, TR, TN)
	h := balancedHeight(hP, hLL)
	return h, NewNode(kL, vL, h, TLL, TP, TLN)
}

// joinRight concatenates a left tree, k and a right tree
func joinRight(k []byte, v []byte, TL *Node, TR *Node, TN *Node, numOfExposedNodes *int, numOfHeightTakenNodes *int) (int, *Node) {
	kL, vL, TLL, TLR, TLN := exposeNode(TL, numOfExposedNodes)
	hLR := HeightOf(TLR, numOfHeightTakenNodes)
	hR := HeightOf(TR, numOfHeightTakenNodes)
	hLL := HeightOf(TLL, numOfHeightTakenNodes)
	if hLR <= hR+1 {
		hP := balancedHeight(hLR, hR)
		if hP <= hLL+1 {
			h := balancedHeight(hLL, balancedHeight(hLR, hR))
			return h, NewNode(kL, vL, h, TLL, NewNode(k, v, hP, TLR, TR, TN), TLN)
		}
		_, TP := rotateRight(k, v, TLR, TR, TN, numOfExposedNodes, numOfHeightTakenNodes)
		return rotateLeft(kL, vL, TLL, TP, TLN, numOfExposedNodes, numOfHeightTakenNodes)
	}
	hP, TP := joinRight(k, v, TLR, TR, TN, numOfExposedNodes, numOfHeightTakenNodes)
	if hP <= hLL+1 {
		h := balancedHeight(hP, hLL)
		return h, NewNode(kL, vL, h, TLL, TP, TLN)
	}
	return rotateLeft(kL, vL, TLL, TP, TLN, numOfExposedNodes, numOfHeightTakenNodes)
}

func joinLeft(k []byte, v []byte, TL *Node, TR *Node, TN *Node, numOfExposedNodes *int, numOfHeightTakenNodes *int) (int, *Node) {
	kR, vR, TRL, TRR, TRN := exposeNode(TR, numOfExposedNodes)
	hRL := HeightOf(TRL, numOfHeightTakenNodes)
	hL := HeightOf(TL, numOfHeightTakenNodes)
	hRR := HeightOf(TRR, numOfHeightTakenNodes)
	if hRL <= hL+1 {
		hP := balancedHeight(hL, hRL)
		if hP <= hRR+1 {
			h := balancedHeight(hRR, balancedHeight(hL, hRL))
			return h, NewNode(kR, vR, h, NewNode(k, v, hP, TL, TRL, TN), TRR, TRN)
		}
		_, TP := rotateLeft(k, v, TL, TRL, TN, numOfExposedNodes, numOfHeightTakenNodes)
		return rotateRight(kR, vR, TP, TRR, TRN, numOfExposedNodes, numOfHeightTakenNodes)
	}
	hP, TP := joinLeft(k, v, TL, TRL, TN, numOfExposedNodes, numOfHeightTakenNodes)
	if hP <= hRR+1 {
		h := balancedHeight(hP, hRR)
		return h, NewNode(kR, vR, h, TP, TRR, TRN)
	}
	return rotateRight(kR, vR, TP, TRR, TRN, numOfExposedNodes, numOfHeightTakenNodes)
}

func join(k []byte, v []byte, DU *DictNode, DD *DictNode, TL *Node, TR *Node, TN *Node, numOfExposedNodes *int, numOfHeightTakenNodes *int) *Node {
	hL := HeightOf(TL, numOfHeightTakenNodes)
	hR := HeightOf(TR, numOfHeightTakenNodes)
	if hL > hR+1 {
		_, T := joinRight(k, v, TL, TR, TN, numOfExposedNodes, numOfHeightTakenNodes)
		return T
	}
	if hR > hL+1 {
		_, T := joinLeft(k, v, TL, TR, TN, numOfExposedNodes, numOfHeightTakenNodes)
		return T
	}
	N := Union(Difference(TN, DU, numOfExposedNodes, numOfHeightTakenNodes), DD, numOfExposedNodes, numOfHeightTakenNodes)
	h := balancedHeight(hL, hR)
	return NewNode(k, v, h, TL, TR, N)
}

func splitLast(T *Node, numOfExposedNodes *int, numOfHeightTakenNodes *int) (*Node, []byte, []byte, *Node) {
	m, v, L, R, N := exposeNode(T, numOfExposedNodes)
	if R == nil {
		return L, m, v, N
	}

	TP, kP, vP, NP := splitLast(R, numOfExposedNodes, numOfHeightTakenNodes)
	return join(m, v, nil, nil, L, TP, N, numOfExposedNodes, numOfHeightTakenNodes), kP, vP, NP
}

func join2(TL *Node, TR *Node, numOfExposedNodes *int, numOfHeightTakenNodes *int) *Node {
	if TL == nil {
		return TR
	}
	TLP, k, v, N := splitLast(TL, numOfExposedNodes, numOfHeightTakenNodes)
	return join(k, v, nil, nil, TLP, TR, N, numOfExposedNodes, numOfHeightTakenNodes)
}

func split(t *Node, k []byte, numOfExposedNodes *int, numOfHeightTakenNodes *int) (*Node, *Node, *Node) {
	if t == nil {
		return nil, nil, nil
	}

	m, v, L, R, N := exposeNode(t, numOfExposedNodes)
	if bytes.Compare(k, m) == 0 {
		return L, R, N
	}

	if bytes.Compare(k, m) == -1 {
		LL, LR, LN := split(L, k, numOfExposedNodes, numOfHeightTakenNodes)
		return LL, join(m, v, nil, nil, LR, R, N, numOfExposedNodes, numOfHeightTakenNodes), LN
	}

	RL, RR, RN := split(R, k, numOfExposedNodes, numOfHeightTakenNodes)
	return join(m, v, nil, nil, L, RL, N, numOfExposedNodes, numOfHeightTakenNodes), RR, RN
}

func Union(T0 *Node, D *DictNode, numOfExposedNodes *int, numOfHeightTakenNodes *int) *Node {
	if T0 == nil {
		return D.ConvertToNode()
	}
	if D == nil {
		return T0
	}

	k, v, DL, DR, DU, DD := exposeDict(D)
	TL, TR, TN := split(T0, k, numOfExposedNodes, numOfHeightTakenNodes)
	L := Union(TL, DL, numOfExposedNodes, numOfHeightTakenNodes)
	R := Union(TR, DR, numOfExposedNodes, numOfHeightTakenNodes)
	joined := join(k, v, DU, DD, L, R, TN, numOfExposedNodes, numOfHeightTakenNodes)
	return joined
}

func Difference(T0 *Node, D *DictNode, numOfExposedNodes *int, numOfHeightTakenNodes *int) *Node {
	if T0 == nil {
		return nil
	}
	if D == nil {
		return T0
	}

	k, _, DL, DR, _, _ := exposeDict(D)
	TL, TR, _ := split(T0, k, numOfExposedNodes, numOfHeightTakenNodes)
	L := Difference(TL, DL, numOfExposedNodes, numOfHeightTakenNodes)
	R := Difference(TR, DR, numOfExposedNodes, numOfHeightTakenNodes)
	return join2(L, R, numOfExposedNodes, numOfHeightTakenNodes)
}
