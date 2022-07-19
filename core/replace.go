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
// Replace the given node with provided replacement
// it it exists in the list of nodes
//
// Returns `true` if the node was actually replaced, `false`
// otherwise
//
func (doc *HtmlDocument) ReplaceNode(original *HtmlNode, replacement *HtmlNode) bool {
	if original == nil {
		return false
	}

	if replacement == nil {
		return false
	}

	if doc.IsEmpty() {
		return false
	}

	for index, child := range doc.Nodes {
		if child == original {
			doc.Nodes[index] = replacement

			// detach & attach
			original.detach()
			replacement.Parent = nil
			replacement.document = doc

			// all done
			return true
		}
	}

	return false
}

//
// Replace the given node with provided replacement by ensuring
// whether it has a parent, or is directly attached to document.
//
// Returns `true` if the node was actually replaced, `false`
// otherwise
//
func (node *HtmlNode) ReplaceMe(replacement *HtmlNode) bool {
	if replacement == nil {
		return false
	}

	if node.Parent == nil {
		return node.document.ReplaceNode(node, replacement)
	}

	return node.Parent.ReplaceChild(node, replacement)
}

//
// Replace a child of this node with given replacement.
//
// Returns `true` if the node was actually replaced, `false`
// otherwise
//
func (node *HtmlNode) ReplaceChild(original *HtmlNode, replacement *HtmlNode) bool {
	if original == nil {
		return false
	}

	if replacement == nil {
		return false
	}

	if !node.HasChildren() {
		return false
	}

	for index, child := range node.Children {
		if child == original {
			original.detach()
			node.Children[index] = replacement
			return true
		}
	}

	return false
}
