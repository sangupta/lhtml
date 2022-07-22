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
// The structure that holds a group of HTML elements.
// This is similar to `HtmlDocument` except that it can
// hold HTML fragments as well. Thus, it is different
// from the internal `html` package as the functions it
// provide are different than the standard ones.
//
type HtmlElements struct {
	nodes []*HtmlNode // list of nodes at the top level
}

//
// Function that returns a new empty `HtmlElements` object.
// This has no nodes defined and is totally empty. It is
// used to initialize the internal structure.
//
func NewHtmlElements() *HtmlElements {
	return &HtmlElements{
		nodes: make([]*HtmlNode, 0),
	}
}

//
// Convert this `HtmlElements` instance into a `HtmlDocument`
// instance. Note, in case of fragments, most of the `HtmlDocument`
// functions will return `nil` unless you add them (for example, for
// a simple fragment, `document.Head()` will return `nil`).
//
func (elements *HtmlElements) AsHtmlDocument() *HtmlDocument {
	return &HtmlDocument{
		HtmlElements: *elements,
	}
}

//----- basic property accessors

//
// Return the length of elements inside this instance.
//
func (elements *HtmlElements) Length() int {
	if elements.nodes == nil {
		return 0
	}

	return len(elements.nodes)
}

//
// Check if this document is empty or not. A document is considered
// empty if it has no child node.
//
func (elements *HtmlElements) IsEmpty() bool {
	if len(elements.nodes) == 0 {
		return true
	}

	return false
}

//----- Get various dedicated nodes

//
// Return the first node in this list of nodes.
//
func (elements *HtmlElements) First() *HtmlNode {
	if elements.Length() == 0 {
		return nil
	}

	return elements.nodes[0]
}

//
// Return the last node in this list of nodes.
//
func (elements *HtmlElements) Last() *HtmlNode {
	num := elements.Length()
	if num == 0 {
		return nil
	}

	return elements.nodes[num-1]
}

//
// Return the node at the given index. If index is out
// of bounds, this function shall return `nil`.
//
func (elements *HtmlElements) Get(index int) *HtmlNode {
	num := elements.Length()
	if index < 0 || index >= num {
		return nil
	}

	return elements.nodes[num]
}

//----- FIND methods

//
// Get the node occuring before this node in the list.
//
func (elements *HtmlElements) GetBefore(child *HtmlNode) *HtmlNode {
	if child == nil {
		return nil
	}

	if elements.Length() == 0 {
		return nil
	}

	for index, node := range elements.nodes {
		if node == child {
			return elements.Get(index - 1)
		}
	}

	return nil
}

//
// Get the node occuring after this node in the list.
//
func (elements *HtmlElements) GetAfter(child *HtmlNode) *HtmlNode {
	if child == nil {
		return nil
	}

	if elements.Length() == 0 {
		return nil
	}

	for index, node := range elements.nodes {
		if node == child {
			return elements.Get(index + 1)
		}
	}

	return nil
}

//
// Find and return all elements in this list's direct children
// that match the given name/tag name/node name.
// Returns an instance of `HtmlElements` which contains all the
// selected nodes. If no match is found, an empty list is returned.
// This method never returns a `nil`.
//
func (elements *HtmlElements) GetChildrenByName(name string) *HtmlElements {
	if elements.IsEmpty() {
		return NewHtmlElements()
	}

	result := NewHtmlElements()
	for _, child := range elements.nodes {
		if strings.EqualFold(child.NodeName(), name) {
			result.nodes = append(result.nodes, child)
		}
	}

	return result
}

//
// Find and return all elements in this list of elements and its
// children that match the given name/tag name/node name. This function
// searches the entire tree for a match.
//
// Returns an instance of `HtmlElements` which contains all the
// selected nodes. If no match is found, an empty list is returned.
// This method never returns a `nil`.
//
func (elements *HtmlElements) GetElementsByName(name string) *HtmlElements {
	if elements.IsEmpty() {
		return NewHtmlElements()
	}

	result := NewHtmlElements()
	for _, child := range elements.nodes {
		child.getElementsByName(name, result)
	}

	return result
}

//
// Find a node within this list of elements which has an
// ID value as the given value.
//
// Returns `HtmlNode` instance if found, `nil` otherwise
//
func (elements *HtmlElements) GetElementById(id string) *HtmlNode {
	if id == "" {
		return nil
	}

	if elements.IsEmpty() {
		return nil
	}

	for _, child := range elements.nodes {
		found := child.GetElementById(id)
		if found != nil {
			return found
		}
	}

	return nil
}

//----- Manipulation methods

//
// Insert a node at given index.
// If index is less than or equal to zero, the node is inserted as first element.
// If index is equal or greater than length, the node is inserted as last element.
//
func (elements *HtmlElements) InsertAt(index int, newNode *HtmlNode) {
	// attach the node
	newNode._parent = nil
	newNode._wrappingElements = elements

	// first addition
	if index <= 0 {
		elements.nodes = append([]*HtmlNode{newNode}, elements.nodes...)
		return
	}

	// falls at the end
	num := len(elements.nodes)
	if index >= num {
		elements.nodes = append(elements.nodes, newNode)
		return
	}

	// falls in between
	prefix := elements.nodes[:index]
	suffix := elements.nodes[index:]
	elements.nodes = append(prefix, newNode)
	elements.nodes = append(elements.nodes, suffix...)
	return
}

//
// Insert a newNode before another childNode.
// Returns `true` if the newNode was added successfully.
// Returns `false` if there are no elements in this instance
// or the child instance cannot be found.
//
func (elements *HtmlElements) InsertBefore(childNode *HtmlNode, newNode *HtmlNode) bool {
	if !childNode.HasChildren() {
		return false
	}

	for index, child := range elements.nodes {
		if child == childNode {
			newIndex := index - 1
			if newIndex == -1 {
				newIndex = 0
			}

			elements.InsertAt(newIndex, newNode)
			return true
		}
	}

	return false
}

//
// Insert a newNode after given childNode.
// Returns `true` if the newNode was added successfully.
// Returns `false` if there are no elements in this instance
// or the child instance cannot be found.
//
func (elements *HtmlElements) InsertAfter(childNode *HtmlNode, newNode *HtmlNode) bool {
	if !childNode.HasChildren() {
		return false
	}

	for index, kid := range elements.nodes {
		if kid == childNode {
			elements.InsertAt(index+1, newNode)
			return true
		}
	}

	return false
}

//
// Insert the given newNode as the first node in the list
// of elements.
//
func (elements *HtmlElements) InsertFirst(newNode *HtmlNode) {
	elements.InsertAt(-1, newNode)
}

//
// Insert the given node as the last node in the list
// of elements.
//
func (elements *HtmlElements) InsertLast(newNode *HtmlNode) {
	elements.InsertAt(elements.Length()+1, newNode)
}

//
// Remove all nodes from this list of elements. All removed
// nodes are detached.
//
func (elements *HtmlElements) Empty() {
	if elements.IsEmpty() {
		return
	}

	// detach nodes
	for _, node := range elements.nodes {
		node._parent = nil
		node._wrappingElements = nil
	}

	// create new slice
	elements.nodes = make([]*HtmlNode, 0)
}

//
// Remove given childNode from document if it is a direct child.
//
// Returns `true` if the childNode was actually removed, `false`
// otherwise. The removed childNode is detached.
//
func (elements *HtmlElements) Remove(childNode *HtmlNode) bool {
	if elements.IsEmpty() {
		return false
	}

	for index, child := range elements.nodes {
		if child == childNode {
			childNode.detach()
			elements.nodes = append(elements.nodes[:index], elements.nodes[index+1:]...)
			return true
		}
	}

	return false
}

//
// Replace the given childNode with provided newNode replacement
// if it exists in the list of nodes within this element.
// Returns `true` if the node was actually replaced, `false`
// otherwise. The removed childNode is detached.
//
func (elements *HtmlElements) Replace(childNode *HtmlNode, newNode *HtmlNode) bool {
	if childNode == nil {
		return false
	}

	if newNode == nil {
		return false
	}

	if elements.IsEmpty() {
		return false
	}

	for index, child := range elements.nodes {
		if child == childNode {
			elements.nodes[index] = newNode

			// detach & attach
			childNode.detach()
			newNode._parent = nil
			newNode._wrappingElements = elements

			// all done
			return true
		}
	}

	return false
}

//----- Internal methods

//
// Add the given node to the stack ensuring that we also add the parent
// relationship, and add to list of direct child nodes of this document
// if there is nothing on the stack
//
func (elements *HtmlElements) addNodeToStack(node *HtmlNode, stack *nodeStack) {
	// push this node on to the stack (and optionally document) itself
	if !stack.isEmpty() {
		// stack already has something so let's set parent-child relationship
		parent := stack.peek()
		parent.addChild(node)
		node._parent = parent
	} else {
		// stack is empty, as this is top-most node
		// we will add it to document
		elements.appendNode(node)
	}

	// push this node to stack
	stack.push(node)
}

//
// Append the given node to the list of nodes in this document.
//
func (elements *HtmlElements) appendNode(node *HtmlNode) {
	node._wrappingElements = elements
	elements.nodes = append(elements.nodes, node)
}
