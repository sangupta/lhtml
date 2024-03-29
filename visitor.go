/**
 * lhtml - Lenient HTML parser for Go.
 *
 * MIT License.
 * Copyright (c) 2022, Sandeep Gupta.
 * https://github.com/sangupta/lhtml
 *
 * Use of this source code is governed by a MIT style license
 * that can be found in LICENSE file in the code repository:
 */

package lhtml

//
// Defines a simple contract for a node visitor. The visitor
// receives a node, and then returns either `true` to continue
// traversing the html tree, or `false` to immediately stop
// walking the tree.
//
type HtmlNodeVisitor func(node *HtmlNode) bool

//
// Allow traversing over the `HtmlDocument`. If a `nil`
// visitor is supplied, no tree traversal happens.
//
func (doc *HtmlElements) Traverse(visitor HtmlNodeVisitor) {
	if visitor == nil {
		return
	}

	if doc.IsEmpty() {
		return
	}

	for _, node := range doc.nodes {
		shouldContinue := node.Traverse(visitor)
		if !shouldContinue {
			break
		}
	}
}

//
// Allow traversing over the `HtmlNode`.
//
func (node *HtmlNode) Traverse(visitor HtmlNodeVisitor) bool {
	if visitor == nil {
		return false
	}

	shouldContinue := visitor(node)
	if !shouldContinue {
		return false
	}

	if !node.HasChildren() {
		return true
	}

	for _, child := range node._children {
		shouldContinue := child.Traverse(visitor)
		if !shouldContinue {
			return false
		}
	}

	return true
}
