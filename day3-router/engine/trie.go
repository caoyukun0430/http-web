package engine

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string // URL pattern being matched, pattern only be set at leaf node so that we can use it to decide if route is matched
	segment  string // URL segment being split
	children []*node
	isWild   bool // segment starts with : or *
}

// Print for debug
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, segment=%s, isWild=%t}", n.pattern, n.segment, n.isWild)
}

// implement two basic method recursively, insert and search

// insert node into correct trie tree
// pattern is generic like /:name, segments is specific like /tom, /alice
func (n *node) insert(pattern string, segments []string, height int) {
	// base case
	if height == len(segments) {
		// this is criteria for routing match
		n.pattern = pattern
		// fmt.Println(n)
		return
	}
	seg := segments[height]
	child := n.matchChild(seg)
	// fmt.Println(child)
	// if there is no match, we construct the child node
	if child == nil {
		child = &node{
			segment: seg,
			isWild:  seg[0] == ':' || seg[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, segments, height+1)
}

// search top to botton segments of the node, if a path is found, i.e.
// leaf node pattern in not nil, return node
func (n *node) search(segments []string, height int) *node {
	// base case
	if height == len(segments) || strings.HasPrefix(n.segment, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	seg := segments[height]
	children := n.matchChildren(seg)
	for _, child := range children {
		result := child.search(segments, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

// HELPER FUNCS
// matchChildren is use to find all children with same pattern part during routing search
// static routes have priority
func (n *node) matchChildren(segment string) []*node {
	nodes := make([]*node, 0)
	wildNodes := make([]*node, 0)
	for _, child := range n.children {
		if child.segment == segment {
			nodes = append(nodes, child)
		}
		if child.isWild {
			wildNodes = append(wildNodes, child)
		}
	}
	// variadic expansion and used to append all elements of a slice to another slice
	nodes = append(nodes, wildNodes...)
	return nodes
}

// matchChild is used for insert to build up the trie tree, since it's being build,
// there is no multi match, we simply return the first match
// verify explictly if segment is dynamic routing to avoid static overwritten
// e.g. /:age and /18, both will not overwrite each other, but wildcard will overwrite itself
func (n *node) matchChild(segment string) *node {
	for _, child := range n.children {
		if child.segment == segment || ((segment[0] == ':' || segment[0] == '*') && child.isWild) {
			return child
		}
	}
	return nil
}

func (n *node) getPatternNodes(nodes *([]*node)) {
	if n.pattern != "" {
		*nodes = append(*nodes, n)
	}
	for _, e := range n.children {
		e.getPatternNodes(nodes)
	}
}
