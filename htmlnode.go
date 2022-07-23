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

import (
	"strings"
)

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
	_tagName          string
	_parent           *HtmlNode `json:"-"`
	Attributes        []*HtmlAttribute
	_children         []*HtmlNode
	IsSelfClosing     bool
	NodeType          HtmlNodeType
	Data              string
	_wrappingElements *HtmlElements // the document node that this node belongs to
}

func newNode(name string) *HtmlNode {
	return &HtmlNode{
		_tagName: name,
	}
}

//----- basic property accessors

func (node *HtmlNode) NodeName() string {
	return strings.TrimSpace(node._tagName)
}

func (node *HtmlNode) NumChildren() int {
	if node._children == nil {
		return 0
	}

	return len(node._children)
}

func (node *HtmlNode) Parent() *HtmlNode {
	return node._parent
}

//
// Quick check to see if this node has any children
// or not.
//
func (node *HtmlNode) HasChildren() bool {
	if node._children == nil || len(node._children) == 0 {
		return false
	}

	return true
}

//----- FIND methods

func (node *HtmlNode) Get(index int) *HtmlNode {
	if index < 0 {
		return nil
	}

	num := len(node._children)
	if index >= num {
		return nil
	}

	return node._children[index]
}

func (node *HtmlNode) First() *HtmlNode {
	return node.Get(0)
}

func (node *HtmlNode) Last() *HtmlNode {
	return node.Get(node.NumChildren() - 1)
}

func (node *HtmlNode) GetChildByName(name string) *HtmlNode {
	if node.NumChildren() == 0 {
		return nil
	}

	for _, child := range node._children {
		if strings.EqualFold(child.NodeName(), name) {
			return child
		}
	}

	return nil
}

func (node *HtmlNode) GetElementsByName(name string) *HtmlElements {
	elements := NewHtmlElements()
	node.getElementsByName(name, elements)
	return elements
}

func (node *HtmlNode) getElementsByName(name string, elements *HtmlElements) {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)

	if node._tagName == name {
		elements.appendNode(node)
	}

	if !node.HasChildren() {
		return
	}

	for _, child := range node._children {
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
	for _, child := range node._children {
		found := child.GetElementById(id)
		if found != nil {
			return found
		}
	}

	return nil
}

func (node *HtmlNode) GetChild(index int) *HtmlNode {
	num := node.NumChildren()
	if index < 0 || index >= num {
		return nil
	}

	return node._children[index]
}

func (node *HtmlNode) GetBefore(child *HtmlNode) *HtmlNode {
	if child == nil {
		return nil
	}

	if !node.HasChildren() {
		return nil
	}

	for index, kid := range node._children {
		if kid == child {
			return node.GetChild(index - 1)
		}
	}

	if node._wrappingElements != nil {
		node._wrappingElements.GetBefore(node)
	}

	return nil
}

func (node *HtmlNode) GetAfter(child *HtmlNode) *HtmlNode {
	if child == nil {
		return nil
	}

	if !node.HasChildren() {
		return nil
	}

	for index, kid := range node._children {
		if kid == child {
			return node.GetChild(index + 1)
		}
	}

	if node._wrappingElements != nil {
		node._wrappingElements.GetAfter(node)
	}

	return nil
}

func (node *HtmlNode) PrevSibling() *HtmlNode {
	if node._parent != nil {
		return node._parent.GetAfter(node)
	}

	return nil
}

func (node *HtmlNode) NextSibling() *HtmlNode {
	if node._parent != nil {
		return node._parent.GetAfter(node)
	}

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

	node._children = make([]*HtmlNode, 0)
}

//
// Remove this node from its parent node, or from the document.
//
// Returns `true` if the node was actually removed, `false`
// otherwise
//
func (node *HtmlNode) RemoveMe() bool {
	if node._parent == nil {
		if node._wrappingElements == nil {
			return false
		}

		return node._wrappingElements.Remove(node)
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

	for index, c := range node._children {
		if c == child {
			child.detach()
			node._children = append(node._children[:index], node._children[index+1:]...)
			return true
		}
	}

	return false
}

//
// ReplaceMe the given node with provided replacement by ensuring
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
		return node._wrappingElements.Replace(node, replacement)
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

	for index, child := range node._children {
		if child == original {
			original.detach()
			node._children[index] = replacement
			return true
		}
	}

	return false
}

func (node *HtmlNode) InsertChildAt(index int, additional *HtmlNode) bool {
	// first addition
	if index <= 0 {
		node._children = append([]*HtmlNode{additional}, node._children...)
		return true
	}

	// falls at the end
	num := len(node._children)
	if index >= num {
		node._children = append(node._children, additional)
		return true
	}

	// falls in between
	prefix := node._children[:index]
	suffix := node._children[index:]
	node._children = append(prefix, additional)
	node._children = append(node._children, suffix...)
	return true
}

//
// Insert a node before given child
//
func (node *HtmlNode) InsertBeforeChild(child *HtmlNode, additional *HtmlNode) bool {
	if !node.HasChildren() {
		return false
	}

	for index, kid := range node._children {
		if kid == child {
			newIndex := index - 1
			if newIndex == -1 {
				newIndex = 0
			}
			return node.InsertChildAt(newIndex, additional)
		}
	}

	return false
}

//
// Insert a node after given child.
//
func (node *HtmlNode) InsertAfterChild(child *HtmlNode, additional *HtmlNode) bool {
	if !node.HasChildren() {
		return false
	}

	for index, kid := range node._children {
		if kid == child {
			return node.InsertChildAt(index+1, additional)
		}
	}

	return false
}

//
// Insert a node before me.
//
func (node *HtmlNode) InsertBeforeMe(additional *HtmlNode) bool {
	if node._parent != nil {
		return node._parent.InsertBeforeChild(node, additional)
	}

	if node._wrappingElements != nil {
		return node._wrappingElements.InsertBefore(node, additional)
	}

	return false
}

//
// Insert a node after me.
//
func (node *HtmlNode) InsertAfterMe(additional *HtmlNode) bool {
	if node._parent != nil {
		return node._parent.InsertAfterChild(node, additional)
	}

	if node._wrappingElements != nil {
		return node._wrappingElements.InsertAfter(node, additional)
	}

	return false
}

//----- Internal methods

//
// Add a child node to this node.
//
func (node *HtmlNode) addChild(child *HtmlNode) {
	if len(node._children) == 0 {
		node._children = make([]*HtmlNode, 0)
	}

	node._children = append(node._children, child)
}

func (node *HtmlNode) detach() {
	node._parent = nil
	node._wrappingElements = nil
}
