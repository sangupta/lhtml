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
// The structure that holds the HTML document.
// This is different from the internal `html` package
// as the functions we provide are different than the
// standard ones.
//
type HtmlElements struct {
	nodes []*HtmlNode // list of nodes at the top level
}

//
// Function that returns a new empty `HtmlDocument` object.
// This has no nodes defined and is totally empty. It is
// used to initialize the internal structure.
//
func newHtmlElements() *HtmlElements {
	return &HtmlElements{
		nodes: make([]*HtmlNode, 0),
	}
}

func (elements *HtmlElements) AsHtmlDocument() *HtmlDocument {
	return &HtmlDocument{
		HtmlElements: *elements,
	}
}

//----- basic property accessors

func (elements *HtmlElements) NumNodes() int {
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

func (elements *HtmlElements) FirstNode() *HtmlNode {
	if elements.NumNodes() == 0 {
		return nil
	}

	return elements.nodes[0]
}

func (elements *HtmlElements) LastNode() *HtmlNode {
	num := elements.NumNodes()
	if num == 0 {
		return nil
	}

	return elements.nodes[num-1]
}

func (elements *HtmlElements) GetNode(index int) *HtmlNode {
	num := elements.NumNodes()
	if index < 0 || index >= num {
		return nil
	}

	return elements.nodes[num]
}

//----- FIND methods

//
// Find and return all elements in this document that match
// the given name/tag name/node name.
//
// Returns an instance of `HtmlDocument` which contains all the
// selected nodes.
//
func (elements *HtmlElements) GetElementsByName(name string) *HtmlElements {
	if elements.IsEmpty() {
		return nil
	}

	result := newHtmlElements()
	for _, child := range elements.nodes {
		child.getElementsByName(name, result)
	}

	return result
}

//
// Find a node within this document which has an
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

func (elements *HtmlElements) InsertFirst(node *HtmlNode) {
	elements.nodes = append([]*HtmlNode{node}, elements.nodes...)
}

//
// Append the given node to the list of nodes in this document.
//
func (elements *HtmlElements) appendNode(node *HtmlNode) {
	node.document = elements
	elements.nodes = append(elements.nodes, node)
}

//
// Remove all nodes from this document.
//
func (elements *HtmlElements) RemoveAllNodes() {
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
func (elements *HtmlElements) RemoveNode(node *HtmlNode) bool {
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
func (elements *HtmlElements) ReplaceNode(original *HtmlNode, replacement *HtmlNode) bool {
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
			replacement.document = elements

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
