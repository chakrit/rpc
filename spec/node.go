package spec

import "sort"

// Node marks a type allowing it to be used as the spec's AST node.
type Node interface {
	name() string
	node() // marker method
}

type merger interface {
	Merge(Node) Node
}

type byNodeName []Node

var _ sort.Interface = byNodeName{}

func (n byNodeName) Len() int           { return len(n) }
func (n byNodeName) Less(i, j int) bool { return n[i].name() < n[j].name() }
func (n byNodeName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
