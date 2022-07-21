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

import "strings"

//
// Enum to define the NodeType
//
type HtmlNodeType uint32

// Enumeration
const (
	ErrorNode HtmlNodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

//
// Defines the structure for a `node` in the HTML. Before working
// with a node, do check the `NodeType` value to ensure that the
// property you are reading will contain a value or not.
//
type HtmlNode struct {
	TagName       string
	_parent       *HtmlNode `json:"-"`
	Attributes    []*HtmlAttribute
	Children      []*HtmlNode
	IsSelfClosing bool
	NodeType      HtmlNodeType
	Data          string
	document      *HtmlDocument // the document node that this node belongs to
}

func newNode(name string) *HtmlNode {
	return &HtmlNode{
		TagName: name,
	}
}

//----- basic property accessors

func (node *HtmlNode) NodeName() string {
	return strings.TrimSpace(node.TagName)
}

func (node *HtmlNode) NumChildren() int {
	if node.Children == nil {
		return 0
	}

	return len(node.Children)
}

func (node *HtmlNode) Parent() *HtmlNode {
	return node._parent
}

//
// Quick check to see if this node has any children
// or not.
//
func (node *HtmlNode) HasChildren() bool {
	if node.Children == nil || len(node.Children) == 0 {
		return false
	}

	return true
}

func (node *HtmlNode) HasAttributes() bool {
	if node.Attributes == nil || len(node.Attributes) == 0 {
		return false
	}

	return true
}

//----- FIND methods

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

func (node *HtmlNode) PrevSibling() *HtmlNode {
	return nil
}

func (node *HtmlNode) NextSibling() *HtmlNode {
	return nil
}

//----- Manipulation methods

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
// Remove this node from its parent node, or from the document.
//
// Returns `true` if the node was actually removed, `false`
// otherwise
//
func (node *HtmlNode) RemoveMe() bool {
	if node._parent == nil {
		if node.document == nil {
			return false
		}

		return node.document.RemoveNode(node)
	}

	return node._parent.RemoveChild(node)
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

	if node._parent == nil {
		return node.document.ReplaceNode(node, replacement)
	}

	return node._parent.ReplaceChild(node, replacement)
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

//----- Internal methods

//
// Add a new attribute to this node. By design, we allow a single
// tag to hold multiple values for the same attribute name. This is
// to ensure that we can parse JSX-like syntax to allow templates
// to hold individual values, and then let the template engines to
// merge them into a single value.
//
func (node *HtmlNode) addAttribute(key string, value string) {
	if len(node.Attributes) == 0 {
		node.Attributes = make([]*HtmlAttribute, 0)
	}

	node.Attributes = append(node.Attributes, &HtmlAttribute{
		Name:  key,
		Value: value,
	})
}

//
// Add a child node to this node.
//
func (node *HtmlNode) addChild(child *HtmlNode) {
	if len(node.Children) == 0 {
		node.Children = make([]*HtmlNode, 0)
	}

	node.Children = append(node.Children, child)
}

func (node *HtmlNode) detach() {
	node._parent = nil
	node.document = nil
}
