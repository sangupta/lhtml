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
type HtmlDocument struct {
	Nodes []*HtmlNode // list of nodes at the top level
}

//
// Function that returns a new empty `HtmlDocument` object.
// This has no nodes defined and is totally empty. It is
// used to initialize the internal structure.
//
func NewHtmlDocument() *HtmlDocument {
	return &HtmlDocument{
		Nodes: make([]*HtmlNode, 0),
	}
}

//----- basic property accessors

func (document *HtmlDocument) NumNodes() int {
	if document.Nodes == nil {
		return 0
	}

	return len(document.Nodes)
}

//
// Check if this document is empty or not. A document is considered
// empty if it has no child node.
//
func (document *HtmlDocument) IsEmpty() bool {
	if len(document.Nodes) == 0 {
		return true
	}

	return false
}

//----- Get various dedicated nodes

func (document *HtmlDocument) Head() *HtmlNode {
	elements := document.GetElementsByName("head")
	if elements == nil || elements.NumNodes() == 0 {
		return nil
	}

	return elements.Nodes[0]
}

func (document *HtmlDocument) Body() *HtmlNode {
	elements := document.GetElementsByName("body")
	if elements == nil || elements.NumNodes() == 0 {
		return nil
	}

	return elements.Nodes[0]
}

func (doc *HtmlDocument) GetDocType() *HtmlNode {
	if doc.IsEmpty() {
		return nil
	}

	for _, node := range doc.Nodes {
		if node.NodeType == DoctypeNode {
			return node
		}
	}

	return nil
}

func (document *HtmlDocument) FirstNode() *HtmlNode {
	if document.NumNodes() == 0 {
		return nil
	}

	return document.Nodes[0]
}

func (document *HtmlDocument) LastNode() *HtmlNode {
	num := document.NumNodes()
	if num == 0 {
		return nil
	}

	return document.Nodes[num-1]
}

func (document *HtmlDocument) GetNode(index int) *HtmlNode {
	num := document.NumNodes()
	if index < 0 || index >= num {
		return nil
	}

	return document.Nodes[num]
}

//----- FIND methods

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

//----- Manipulation methods

func (document *HtmlDocument) InsertFirst(node *HtmlNode) {
	document.Nodes = append([]*HtmlNode{node}, document.Nodes...)
}

//
// Append the given node to the list of nodes in this document.
//
func (document *HtmlDocument) appendNode(node *HtmlNode) {
	node.document = document
	document.Nodes = append(document.Nodes, node)
}

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
			replacement._parent = nil
			replacement.document = doc

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
func (document *HtmlDocument) addNodeToStack(node *HtmlNode, stack *nodeStack) {
	// push this node on to the stack (and optionally document) itself
	if !stack.isEmpty() {
		// stack already has something so let's set parent-child relationship
		parent := stack.peek()
		parent.addChild(node)
		node._parent = parent
	} else {
		// stack is empty, as this is top-most node
		// we will add it to document
		document.appendNode(node)
	}

	// push this node to stack
	stack.push(node)
}
