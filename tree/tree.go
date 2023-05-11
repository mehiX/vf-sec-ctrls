package tree

import (
	"strings"
)

// Node represents a node in the tree
type Node struct {
	Value    string           // 1.2
	Children map[string]*Node // {1.2.1: n1, 1.2.2: n2, 1.2.3: n3}
}

// New takes in a slice of strings (each in the form x.x.x.x...),
// Returns a pointer to the root node of the tree
func New(val []string) *Node {
	root := &Node{Value: "", Children: make(map[string]*Node)}
	for _, v := range val {
		add(root, 0, strings.Split(v, "."))
	}

	return root
}

func add(n *Node, level int, val []string) {
	if len(val) == level {
		return
	}

	key := strings.Join(val[0:level+1], ".")
	child, ok := n.Children[key]
	if !ok {
		child = &Node{Value: key, Children: make(map[string]*Node)}
		n.Children[key] = child
	}

	add(child, level+1, val)
}

// FindNode looks at all the children and descendants of `root`.
// Returns the first occurrence of `id` as value of a Node
func FindNode(root *Node, id string) *Node {
	return findNode(root, 0, strings.Split(id, "."))
}

func findNode(n *Node, level int, parts []string) *Node {
	if n == nil {
		return n
	}

	if level == len(parts) {
		return n
	}

	key := strings.Join(parts[0:level+1], ".")
	next, ok := n.Children[key]
	if !ok {
		return nil
	}
	return findNode(next, level+1, parts)

}

func EdgesFrom(n *Node) []string {
	return allEdgesFrom(n, "", []string{})
}

func allEdgesFrom(n *Node, s string, results []string) []string {
	if n == nil {
		return results
	}

	if len(n.Children) == 0 {
		return append(results, n.Value)
	}

	for _, c := range n.Children {
		results = allEdgesFrom(c, n.Value, results)
	}

	return results
}
