package spec

import "sort"

// TODO: Convert to a tree which should speedup both name-based lookup and iteration
type Mappings map[string]Node

func (c *Mappings) AddIfNew(n Node) (Node, bool) {
	if *c == nil {
		*c = Mappings{}
	}

	name := n.name()
	if existing, exists := (*c)[name]; exists {
		return existing, false
	} else {
		(*c)[name] = n
		return n, true
	}
}

func (c *Mappings) Add(n Node) Node {
	name := n.name()
	existing, isNew := c.AddIfNew(n)
	if isNew {
		return n
	} else if merger, ok := existing.(merger); ok {
		n = merger.Merge(n)
	}

	// if n doesn't implement (merger), this'll result in a replace.
	//
	// TODO: Warn about implicit replacements.
	(*c)[name] = n // merged
	return n
}

func (c Mappings) SortedByName() []Node {
	var nodes []Node
	for _, item := range c {
		nodes = append(nodes, item)
	}

	sort.Sort(byNodeName(nodes))
	return nodes
}
