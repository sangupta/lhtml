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

import "strings"

//
// Find and return all elements in this document that match
// the given name/tag name/node name.
//
// Returns an instance of `HtmlDocument` which contains all the
// selected nodes.
//
func (doc *HtmlDocument) GetElementsByName(name string) *HtmlDocument {
	if doc.IsEmpty() {
		return nil
	}

	result := NewHtmlDocument()
	for _, child := range doc.Nodes {
		child.getElementsByName(name, result)
	}

	return result
}

func (node *HtmlNode) GetElementsByName(name string) *HtmlDocument {
	elements := NewHtmlDocument()
	node.getElementsByName(name, elements)
	return elements
}

func (node *HtmlNode) getElementsByName(name string, elements *HtmlDocument) {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)

	if node.TagName == name {
		elements.appendNode(node)
	}

	if !node.HasChildren() {
		return
	}

	for _, child := range node.Children {
		child.getElementsByName(name, elements)
	}
}

//
// Find a node within this document which has an
// ID value as the given value.
//
// Returns `HtmlNode` instance if found, `nil` otherwise
//
func (doc *HtmlDocument) GetElementById(id string) *HtmlNode {
	if id == "" {
		return nil
	}

	if doc.IsEmpty() {
		return nil
	}

	for _, child := range doc.Nodes {
		found := child.GetElementById(id)
		if found != nil {
			return found
		}
	}

	return nil
}

//
// Find a node within this node (including this one) which has an
// ID value as the given value.
//
// Returns `HtmlNode` instance if found, `nil` otherwise
//
func (node *HtmlNode) GetElementById(id string) *HtmlNode {
	if node.GetAttributeWithValue("id", id) != nil {
		return node
	}

	if !node.HasChildren() {
		return nil
	}
	for _, child := range node.Children {
		found := child.GetElementById(id)
		if found != nil {
			return found
		}
	}

	return nil
}
