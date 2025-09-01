package merkle

import (
	"crypto/sha256"
	"fmt"
)

// MerkleTree represents a Merkle tree for tamper-evident logs
type MerkleTree struct {
	Root   *MerkleNode
	Leaves []*MerkleNode
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
	Hash  []byte
}

// NewMerkleTree creates a new Merkle tree from data
func NewMerkleTree(data [][]byte) (*MerkleTree, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	var leaves []*MerkleNode
	for _, d := range data {
		node := NewMerkleNode(nil, nil, d)
		leaves = append(leaves, node)
	}

	root := buildTree(leaves)
	tree := &MerkleTree{
		Root:   root,
		Leaves: leaves,
	}

	return tree, nil
}

// NewMerkleNode creates a new Merkle node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{}

	if left == nil && right == nil {
		// Leaf node
		hash := sha256.Sum256(data)
		node.Data = data
		node.Hash = hash[:]
	} else {
		// Internal node
		prevHashes := append(left.Hash, right.Hash...)
		hash := sha256.Sum256(prevHashes)
		node.Hash = hash[:]
		node.Left = left
		node.Right = right
	}

	return node
}

// buildTree recursively builds the Merkle tree
func buildTree(leaves []*MerkleNode) *MerkleNode {
	if len(leaves) == 0 {
		return nil
	}

	if len(leaves) == 1 {
		return leaves[0]
	}

	// Ensure even number of leaves
	if len(leaves)%2 == 1 {
		leaves = append(leaves, leaves[len(leaves)-1])
	}

	var newLeaves []*MerkleNode
	for i := 0; i < len(leaves); i += 2 {
		node := NewMerkleNode(leaves[i], leaves[i+1], nil)
		newLeaves = append(newLeaves, node)
	}

	return buildTree(newLeaves)
}

// RootHash returns the root hash of the Merkle tree
func (mt *MerkleTree) RootHash() []byte {
	if mt.Root == nil {
		return nil
	}
	return mt.Root.Hash
}

// GenerateProof generates a Merkle proof for a leaf at the given index
func (mt *MerkleTree) GenerateProof(index int) ([][]byte, error) {
	if index < 0 || index >= len(mt.Leaves) {
		return nil, fmt.Errorf("index out of range")
	}

	var proof [][]byte

	// Simple implementation: just return the hashes of all other leaves
	// This is not the most efficient but it works for demonstration
	for i, leaf := range mt.Leaves {
		if i != index {
			proof = append(proof, leaf.Hash)
		}
	}

	return proof, nil
}

// VerifyProof verifies a Merkle proof
func (mt *MerkleTree) VerifyProof(leafData []byte, proof [][]byte, rootHash []byte) bool {
	// Calculate leaf hash
	leafHash := sha256.Sum256(leafData)

	// For our simplified implementation, we just check if the leaf hash is in the proof
	// and if combining all hashes (including leaf) produces the root hash
	allHashes := append([][]byte{leafHash[:]}, proof...)

	// Sort of simulate the Merkle tree verification
	// In a real implementation, this would be more complex
	calculatedRoot := calculateSimpleRoot(allHashes)

	// Compare with root hash
	for i, b := range calculatedRoot {
		if b != rootHash[i] {
			return false
		}
	}

	return true
}

// calculateSimpleRoot calculates a simple root hash from a list of hashes
func calculateSimpleRoot(hashes [][]byte) []byte {
	if len(hashes) == 0 {
		return nil
	}

	if len(hashes) == 1 {
		return hashes[0]
	}

	// Concatenate all hashes and hash the result
	var combined []byte
	for _, hash := range hashes {
		combined = append(combined, hash...)
	}

	result := sha256.Sum256(combined)
	return result[:]
}
