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
// Find and return all elements in this document that match
// the given name/tag name/node name.
//
// Returns an instance of `HtmlElements` which contains all the
// selected nodes.
//
func (elements *HtmlElements) GetElementsByName(name string) *HtmlElements {
	if elements.IsEmpty() {
		return nil
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
// Insert the given node as the first node in the list
// of elements.
//
func (elements *HtmlElements) InsertFirst(node *HtmlNode) {
	elements.nodes = append([]*HtmlNode{node}, elements.nodes...)
}

//
// Insert the given node as the last node in the list
// of elements.
//
func (elements *HtmlElements) InsertLast(node *HtmlNode) {
	elements.nodes = append(elements.nodes, node)
}

//
// Remove all nodes from this document.
//
func (elements *HtmlElements) Empty() {
	if elements.IsEmpty() {
		return
	}

	elements.nodes = make([]*HtmlNode, 0)
}

//
// Remove given node from document if it is a direct child.
//
// Returns `true` if the node was actually removed, `false`
// otherwise
//
func (elements *HtmlElements) Remove(node *HtmlNode) bool {
	if elements.IsEmpty() {
		return false
	}

	for index, child := range elements.nodes {
		if child == node {
			node.detach()
			elements.nodes = append(elements.nodes[:index], elements.nodes[index+1:]...)
			return true
		}
	}

	return false
}

//
// Replace the given node with provided replacement
// it it exists in the list of nodes
//
// Returns `true` if the node was actually replaced, `false`
// otherwise
//
func (elements *HtmlElements) Replace(original *HtmlNode, replacement *HtmlNode) bool {
	if original == nil {
		return false
	}

	if replacement == nil {
		return false
	}

	if elements.IsEmpty() {
		return false
	}

	for index, child := range elements.nodes {
		if child == original {
			elements.nodes[index] = replacement

			// detach & attach
			original.detach()
			replacement._parent = nil
			replacement._wrappingElements = elements

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
