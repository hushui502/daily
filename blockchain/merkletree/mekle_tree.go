package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
)

type Content interface {
	CalculateHash() ([]byte, error)
	Equals(other Content) (bool, error)
}

type MerkleTree struct {
	Root *Node
	merkleRoot []byte
	Leafs []*Node
	hashStrategy func() hash.Hash
}

type Node struct {
	Tree *MerkleTree
	Parent *Node
	Left *Node
	Right *Node
	leaf bool
	dup bool
	Hash []byte
	C Content
}

func (n *Node) calculateNodeHashFromLeaf() ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash()
	}

	leftBytes, err := n.Left.calculateNodeHashFromLeaf()
	if err != nil {
		return nil, err
	}

	rightBytes, err := n.Right.calculateNodeHashFromLeaf()
	if err != nil {
		return nil, err
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (n *Node) calculateNodeHash() ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash()
	}
	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(n.Left.Hash, n.Right.Hash...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func NewMerkleTree(cs []Content) (*MerkleTree, error) {
	var defaultHashStrategy = sha256.New
	t := &MerkleTree{
		hashStrategy: defaultHashStrategy,
	}
	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash

	return t, nil
}

func NewMerkleTreeWithHashStrategy(cs []Content, hashStrategy func() hash.Hash) (*MerkleTree, error) {
	t := &MerkleTree{
		hashStrategy: hashStrategy,
	}

	root, leafs, err := buildWithContent(cs, t)
	if err != nil {
		return nil, err
	}
	t.Root = root
	t.Leafs = leafs
	t.merkleRoot = root.Hash
	return t, nil
}

func (m *MerkleTree) GetMerklePath(content Content) ([][]byte, []int64, error) {
	for _, current := range m.Leafs {
		ok, err := current.C.Equals(content)
		if err != nil {
			return nil, nil, err
		}

		if ok {
			currentParent := current.Parent
			var marklePath [][]byte
			var index []int64
			for currentParent != nil {
				if bytes.Equal(currentParent.Left.Hash, current.Hash) {
					marklePath = append(marklePath, currentParent.Right.Hash)
					index = append(index, 1)
				} else {
					marklePath = append(marklePath, currentParent.Left.Hash)
					index = append(index, 0)
				}
				current = currentParent
				currentParent = currentParent.Parent
			}
			return marklePath, index, nil
		}
	}
	return nil, nil, nil
}

func buildWithContent(cs []Content, t *MerkleTree) (*Node, []*Node, error) {
	if len(cs) == 0 {
		return nil, nil, errors.New("error: can not construct tree with no content")
	}
	var leafs []*Node
	for _, c := range cs {
		hash, err := c.CalculateHash()
		if err != nil {
			return nil, nil, err
		}

		leafs = append(leafs, &Node{
			Hash:hash,
			C:c,
			leaf:true,
			Tree:t,
		})
	}
	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash:leafs[len(leafs)-1].Hash,
			C:leafs[len(leafs)-1].C,
			leaf:true,
			dup:true,
			Tree:t,
		}
		leafs = append(leafs, duplicate)
	}
	root, err := buildIntermediate(leafs, t)
	if err != nil {
		return nil, nil, err
	}

	return root, leafs, nil
}

func buildIntermediate(n1 []*Node, t *MerkleTree) (*Node, error) {
	var nodes []*Node
	for i := 0; i < len(n1); i += 2 {
		h := t.hashStrategy()
		var left, right int = i, i + 1
		if i+1 == len(n1) {
			right = i
		}
		chash := append(n1[left].Hash, n1[right].Hash...)
		if _, err := h.Write(chash); err != nil {
			return nil, err
		}
		n := &Node{
			Left:n1[left],
			Right:n1[right],
			Hash:h.Sum(nil),
			Tree:t,
		}
		nodes = append(nodes, n)
		n1[left].Parent = n
		n1[right].Parent = n
		if len(n1) == 2 {
			return n, nil
		}
	}
	return buildIntermediate(nodes, t)
}

func (m *MerkleTree) MerkleRoot() []byte {
	return m.merkleRoot
}

func (m *MerkleTree) RebuildTree() error {
	var cs []Content
	for _, c := range m.Leafs {
		cs = append(cs, c.C)
	}
	root, leafs, err := buildWithContent(cs, m)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

func (m *MerkleTree) RebuildTreeWith(cs []Content) error {
	root, leafs, err := buildWithContent(cs, m)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

func (m *MerkleTree) VerifyTree() (bool, error) {
	calculateMerkleRoot, err := m.Root.calculateNodeHashFromLeaf()
	if err != nil {
		return false, err
	}

	if bytes.Compare(m.merkleRoot, calculateMerkleRoot) == 0 {
		return true, nil
	}

	return false, nil
}

func (m *MerkleTree) VerifyContent(content Content) (bool, error) {
	for _, l := range m.Leafs {
		ok, err := l.C.Equals(content)
		if err != nil {
			return false, nil
		}

		if ok {
			currentParent := l.Parent
			for currentParent != nil {
				h := m.hashStrategy()
				rightBytes, err := currentParent.Right.calculateNodeHash()
				if err != nil {
					return false, err
				}

				leftBytes, err := currentParent.Left.calculateNodeHash()
				if err != nil {
					return false, err
				}

				if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
					return false, err
				}
				if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
					return false, nil
				}

				currentParent = currentParent.Parent
			}
			return true, nil
		}
	}
	return false, nil
}

func (n *Node) String() string {
	if n.leaf {
		return fmt.Sprintf("%s %#v", BufferToHexString(n.Hash, false), n.C)
	}
	return BufferToHexString(n.Hash, false)
}

func (m *MerkleTree) String() string {
	str := ""
	m.Root.CurString(&str)
	return str
}


func (n *Node) CurString(str *string) {
	if n.leaf {
		*str += fmt.Sprintln(n)
		return
	}

	*str += fmt.Sprintln(n)
	n.Left.CurString(str)
	n.Right.CurString(str)
}

func BufferToHexString(data []byte, prefix bool) string {
	hex := fmt.Sprintf(`%x`, data)
	if prefix {
		return `0x` + hex
	}

	return hex
}
























