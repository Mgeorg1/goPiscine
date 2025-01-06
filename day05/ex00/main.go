package main

import (
	"day05"
	"fmt"
)

func countToys(node *day05.TreeNode) int {
	if node == nil {
		return 0
	}
	count := 0
	if node.HasToy {
		count++
	}
	count += countToys(node.Left)
	count += countToys(node.Right)
	return count
}

func areToysBalanced(tree *day05.BinaryTree) bool {
	left := countToys(tree.Root.Left)
	right := countToys(tree.Root.Right)
	return left == right
}

func main() {
	tree := day05.NewBinaryTree(true)
	tree.Root.Left = day05.NewTreeNode(false)
	tree.Root.Right = day05.NewTreeNode(true)
	left := tree.Root.Left
	left.Left = day05.NewTreeNode(false)
	left.Right = day05.NewTreeNode(true)
	fmt.Println(areToysBalanced(tree))
}
