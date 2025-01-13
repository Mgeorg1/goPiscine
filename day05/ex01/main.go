package main

import (
	"day05"
	"fmt"
)

type stack struct {
	nodes []*day05.TreeNode
}

func (s *stack) push(node *day05.TreeNode) {
	s.nodes = append(s.nodes, node)
}

func (s *stack) pop() *day05.TreeNode {
	if len(s.nodes) == 0 {
		return nil
	}
	node := s.nodes[len(s.nodes)-1]
	s.nodes = s.nodes[:len(s.nodes)-1]
	return node
}

func makeStack() *stack {
	return &stack{nodes: make([]*day05.TreeNode, 0)}
}

func traverseAndUnroll(node *day05.TreeNode) []bool {
	if node == nil {
		return []bool{}
	}
	ret := make([]bool, 0)
	s := makeStack()
	s.push(node)
	leftToRight := true
	for len(s.nodes) > 0 {
		size := len(s.nodes)
		for i := 0; i < size; i++ {
			node := s.pop()
			ret = append(ret, node.HasToy)
			if leftToRight {
				if node.Left != nil {
					s.push(node.Left)
				}
				if node.Right != nil {
					s.push(node.Right)
				}
			} else {
				if node.Right != nil {
					s.push(node.Right)
				}
				if node.Left != nil {
					s.push(node.Left)
				}
			}
		}
		leftToRight = !leftToRight
	}
	return ret
}

func unrollGarland(tree *day05.BinaryTree) []bool {
	return traverseAndUnroll(tree.Root)
}

func main() {
	tree := day05.NewBinaryTree(true)
	tree.Root.Left = day05.NewTreeNode(true)
	tree.Root.Right = day05.NewTreeNode(false)
	node := tree.Root.Left
	node.Left = day05.NewTreeNode(true)
	node.Right = day05.NewTreeNode(false)
	node = tree.Root.Right
	node.Left = day05.NewTreeNode(true)
	node.Right = day05.NewTreeNode(true)
	slice := unrollGarland(tree)
	fmt.Println(slice)
}
