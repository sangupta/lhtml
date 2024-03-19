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

//
// Return the node name, also known as tag name for
// this element.
//
func (node *HtmlNode) NodeName() string {
	return strings.TrimSpace(node._tagName)
}

//
// Return the total number of children this node has.
//
func (node *HtmlNode) NumChildren() int {
	if node._children == nil {
		return 0
	}

	return len(node._children)
}

//
// Return the parent of this node, if any. A node at the root
// level (such as <html />) does not have a parent, but may
// have an internal `wrappingElement`. This allows us to provide
// functions to replace/remove node directly.
//
func (node *HtmlNode) Parent() *HtmlNode {
	return node._parent
}

//
// Get a list of all children of this `HtmlNode`.
//
func (node *HtmlNode) Children() []*HtmlNode {
	return node._children
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

//
// Return the node at a given index. If the index is out
// of bounds, `nil` is returned.
//
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

//
// Return the first child node, if any. Returns `nil` if
// the node has no children.
//
func (node *HtmlNode) First() *HtmlNode {
	return node.Get(0)
}

//
// Return the last child node, if any. Returns `nil` if
// the node has no children.
//
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

//
// Return the elements/nodes that match the given tag name, including this element.
// The node hierarchy is not maintained in results.
//
func (node *HtmlNode) GetElementsByName(name string) *HtmlElements {
	elements := NewHtmlElements()
	node.getElementsByNameInternal(name, elements)
	return elements
}

//
// Internal method to help with collection of nodes that match the tag name.
//
func (node *HtmlNode) getElementsByNameInternal(name string, elements *HtmlElements) {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)

	if node._tagName == name {
		elements.appendNode(node)
	}

	if !node.HasChildren() {
		return
	}

	for _, child := range node._children {
		child.getElementsByNameInternal(name, elements)
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

//
// Return the child at a given index. If the index is out of bounds'
// `nil` is returned.
//
func (node *HtmlNode) GetChild(index int) *HtmlNode {
	num := node.NumChildren()
	if index < 0 || index >= num {
		return nil
	}

	return node._children[index]
}

//
// Return the node before the given child node. Returns `nil`
// if the child is `nil`, or is not a direct child of this node,
// or if this is the first node in list.
//
func (node *HtmlNode) GetChildBefore(child *HtmlNode) *HtmlNode {
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

//
// Return the node after the given child node. Returns `nil`
// if the child is `nil`, or is not a direct child of this node,
// or if this is the last node in list.
//
func (node *HtmlNode) GetChildAfter(child *HtmlNode) *HtmlNode {
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

//
// Return the node before this node in the list. Returns
// `nil` if this node is detached, or has no previous sibling.
//
func (node *HtmlNode) PrevSibling() *HtmlNode {
	if node._parent != nil {
		return node._parent.GetChildAfter(node)
	}

	return nil
}

//
// Return the node after this node in the list. Returns
// `nil` if this node is detached, or has no next sibling.
//
func (node *HtmlNode) NextSibling() *HtmlNode {
	if node._parent != nil {
		return node._parent.GetChildAfter(node)
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

//
// Insert the child at the given index. If the index is less than zero
// the node is inserted as the first node. If the index is greater than
// the last node index it is inserted as the last node.
//
func (node *HtmlNode) InsertChildAt(index int, additional *HtmlNode) {
	// first addition
	if index <= 0 {
		node._children = append([]*HtmlNode{additional}, node._children...)
		return
	}

	// falls at the end
	num := len(node._children)
	if index >= num {
		node._children = append(node._children, additional)
		return
	}

	// falls in between
	prefix := node._children[:index]
	suffix := node._children[index:]
	node._children = append(prefix, additional)
	node._children = append(node._children, suffix...)
	return
}

//
// Insert a node before given child. Returns `true` if the node
// is inserted. Returns `false` if the node has no children, or
// the given child does not belong to this node.
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

			node.InsertChildAt(newIndex, additional)
			return true
		}
	}

	return false
}

//
// Insert a node after given child. Returns `true` if the node
// is inserted. Returns `false` if the node has no children, or
// the given child does not belong to this node.
//
func (node *HtmlNode) InsertAfterChild(child *HtmlNode, additional *HtmlNode) bool {
	if !node.HasChildren() {
		return false
	}

	for index, kid := range node._children {
		if kid == child {
			node.InsertChildAt(index+1, additional)
			return true
		}
	}

	return false
}

//
// Insert a node before this node in its parent's child nodes.
// Returns `true` if the node was inserted. Returns `false` if
// this node has no parent.
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
// Insert a node after this node in its parent's child nodes.
// Returns `true` if the node was inserted. Returns `false` if
// this node has no parent.
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

//----- tostring()

func (node *HtmlNode) String() string {
	builder := strings.Builder{}
	node.WriteToBuilder(&builder)
	return builder.String()
}

func (node *HtmlNode) WriteToBuilder(builder *strings.Builder) {
	if node.NodeType == DoctypeNode || node.NodeType == TextNode || node.NodeType == CommentNode {
		builder.WriteString(node.Data)
		return
	}

	builder.WriteString("<")
	builder.WriteString(node.NodeName())

	// attributes
	if node.ContainsAttributes() {
		for _, attr := range node.Attributes {
			builder.WriteString(" ")
			builder.WriteString(attr.Name)
			builder.WriteString("=\"")
			builder.WriteString(attr.Value)
			builder.WriteString("\"")
		}
	}

	// self-closing?
	if !node.HasChildren() {
		builder.WriteString(" />")
	} else {
		builder.WriteString(">")
		for _, child := range node._children {
			child.WriteToBuilder(builder)
		}

		// close
		builder.WriteString("</")
		builder.WriteString(node.NodeName())
		builder.WriteString(">")
	}

	builder.WriteString(" ")
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

//
// Detach the given node. Remove its associated with its
// parent or its wrapping element.
//
func (node *HtmlNode) detach() {
	node._parent = nil
	node._wrappingElements = nil
}
