package day05

type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

type BinaryTree struct {
	Root *TreeNode
}

func NewTreeNode(value bool) *TreeNode {
	return &TreeNode{HasToy: value, Left: nil, Right: nil}
}

func NewBinaryTree(rootValue bool) *BinaryTree {
	return &BinaryTree{Root: NewTreeNode(rootValue)}
}
