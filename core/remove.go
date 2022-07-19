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

package core

//
// Remove all nodes from this document.
//
func (doc *HtmlDocument) RemoveAllNodes() {
	if doc.IsEmpty() {
		return
	}

	doc.Nodes = make([]*HtmlNode, 0)
}

//
// Remove all children from this `HtmlNode`.
//
func (node *HtmlNode) RemoveAllChildren() {
	if !node.HasChildren() {
		return
	}

	node.Children = make([]*HtmlNode, 0)
}

//
// Remove given node from document if it is a direct child.
//
// Returns `true` if the node was actually removed, `false`
// otherwise
//
func (doc *HtmlDocument) RemoveNode(node *HtmlNode) bool {
	if doc.IsEmpty() {
		return false
	}

	for index, child := range doc.Nodes {
		if child == node {
			node.detach()
			doc.Nodes = append(doc.Nodes[:index], doc.Nodes[index+1:]...)
			return true
		}
	}

	return false
}

//
// Remove this node from its parent node, or from the document.
//
// Returns `true` if the node was actually removed, `false`
// otherwise
//
func (node *HtmlNode) RemoveMe() bool {
	if node.Parent == nil {
		if node.document == nil {
			return false
		}

		return node.document.RemoveNode(node)
	}

	return node.Parent.RemoveChild(node)
}

//
// Remove the given child from this node.
//
// Returns `true` if the node was actually removed, `false`
// otherwise
//
func (node *HtmlNode) RemoveChild(child *HtmlNode) bool {
	if !node.HasChildren() {
		return false
	}

	for index, c := range node.Children {
		if c == child {
			child.detach()
			node.Children = append(node.Children[:index], node.Children[index+1:]...)
			return true
		}
	}

	return false
}
