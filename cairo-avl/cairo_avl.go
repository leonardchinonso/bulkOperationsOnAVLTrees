package cairo_avl

import (
	"bytes"
)

var numOfExposedNodes int
var numOfExposedNodesInDifference int

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
func rotateLeft(k []byte, v []byte, TL *Node, TR *Node, TN *Node) (int, *Node) {
	kR, vR, TRL, TRR, TRN := exposeNode(TR)
	hL := HeightOf(TL)
	hRL := HeightOf(TRL)
	hRR := HeightOf(TRR)
	hP := balancedHeight(hL, hRL)
	TP := NewNode(k, v, hP, TL, TRL, TN)
	h := balancedHeight(hP, hRR)
	return h, NewNode(kR, vR, h, TP, TRR, TRN)
}

// rotateLeft rotates a node to the left to maintain the AVL balance criteria
func rotateRight(k []byte, v []byte, TL *Node, TR *Node, TN *Node) (int, *Node) {
	kL, vL, TLL, TLR, TLN := exposeNode(TL)
	hR := HeightOf(TR)
	hLL := HeightOf(TLL)
	hLR := HeightOf(TLR)
	hP := balancedHeight(hR, hLR)
	TP := NewNode(k, v, hP, TLR, TR, TN)
	h := balancedHeight(hP, hLL)
	return h, NewNode(kL, vL, h, TLL, TP, TLN)
}

// joinRight concatenates a left tree, k and a right tree
func joinRight(k []byte, v []byte, TL *Node, TR *Node, TN *Node) (int, *Node) {
	kL, vL, TLL, TLR, TLN := exposeNode(TL)
	hLR := HeightOf(TLR)
	hR := HeightOf(TR)
	hLL := HeightOf(TLL)
	if hLR <= hR+1 {
		hP := balancedHeight(hLR, hR)
		if hP <= hLL+1 {
			h := balancedHeight(hLL, balancedHeight(hLR, hR))
			return h, NewNode(kL, vL, h, TLL, NewNode(k, v, hP, TLR, TR, TN), TLN)
		}
		_, TP := rotateRight(k, v, TLR, TR, TN)
		return rotateLeft(kL, vL, TLL, TP, TLN)
	}
	hP, TP := joinRight(k, v, TLR, TR, TN)
	if hP <= hLL+1 {
		h := balancedHeight(hP, hLL)
		return h, NewNode(kL, vL, h, TLL, TP, TLN)
	}
	return rotateLeft(kL, vL, TLL, TP, TLN)
}

func joinLeft(k []byte, v []byte, TL *Node, TR *Node, TN *Node) (int, *Node) {
	kR, vR, TRL, TRR, TRN := exposeNode(TR)
	hRL := HeightOf(TRL)
	hL := HeightOf(TL)
	hRR := HeightOf(TRR)
	if hRL <= hL+1 {
		hP := balancedHeight(hL, hRL)
		if hP <= hRR+1 {
			h := balancedHeight(hRR, balancedHeight(hL, hRL))
			return h, NewNode(kR, vR, h, NewNode(k, v, hP, TL, TRL, TN), TRR, TRN)
		}
		_, TP := rotateLeft(k, v, TL, TRL, TN)
		return rotateRight(kR, vR, TP, TRR, TRN)
	}
	hP, TP := joinLeft(k, v, TL, TRL, TN)
	if hP <= hRR+1 {
		h := balancedHeight(hP, hRR)
		return h, NewNode(kR, vR, h, TP, TRR, TRN)
	}
	return rotateRight(kR, vR, TP, TRR, TRN)
}

func join(k []byte, v []byte, DU *DictNode, DD *DictNode, TL *Node, TR *Node, TN *Node) *Node {
	hL := HeightOf(TL)
	hR := HeightOf(TR)
	if hL > hR+1 {
		_, T := joinRight(k, v, TL, TR, TN)
		return T
	}
	if hR > hL+1 {
		_, T := joinLeft(k, v, TL, TR, TN)
		return T
	}
	N, _ := Union(Difference(TN, DU), DD)
	h := balancedHeight(hL, hR)
	return NewNode(k, v, h, TL, TR, N)
}

func splitLast(T *Node) (*Node, []byte, []byte, *Node) {
	m, v, L, R, N := exposeNode(T)
	if R == nil {
		return L, m, v, N
	}

	TP, kP, vP, NP := splitLast(R)
	return join(m, v, nil, nil, L, TP, N), kP, vP, NP
}

func join2(TL *Node, TR *Node) *Node {
	if TL == nil {
		return TR
	}
	TLP, k, v, N := splitLast(TL)
	return join(k, v, nil, nil, TLP, TR, N)
}

func split(t *Node, k []byte) (*Node, *Node, *Node) {
	if t == nil {
		return nil, nil, nil
	}

	m, v, L, R, N := exposeNode(t)
	if bytes.Compare(k, m) == 0 {
		return L, R, N
	}

	if bytes.Compare(k, m) == -1 {
		LL, LR, LN := split(L, k)
		return LL, join(m, v, nil, nil, LR, R, N), LN
	}

	RL, RR, RN := split(R, k)
	return join(m, v, nil, nil, L, RL, N), RR, RN
}

func Union(T0 *Node, D *DictNode) (*Node, int) {
	if T0 == nil {
		return D.ConvertToNode(), numOfExposedNodes
	}
	if D == nil {
		return T0, numOfExposedNodes
	}

	k, v, DL, DR, DU, DD := exposeDict(D)
	TL, TR, TN := split(T0, k)
	L, _ := Union(TL, DL)
	R, _ := Union(TR, DR)
	joined := join(k, v, DU, DD, L, R, TN)
	return joined, numOfExposedNodes
}

func Difference(T0 *Node, D *DictNode) *Node {
	if T0 == nil {
		return nil
	}
	if D == nil {
		return T0
	}

	k, _, DL, DR, _, _ := exposeDict(D)
	TL, TR, _ := split(T0, k)
	L := Difference(TL, DL)
	R := Difference(TR, DR)
	return join2(L, R)
}
