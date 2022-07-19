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
// The structure that holds the HTML document.
// This is different from the internal `html` package
// as the functions we provide are different than the
// standard ones.
//
type HtmlDocument struct {
	head  *HtmlNode
	body  *HtmlNode
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

//
// Append the given node to the list of nodes in this document.
//
func (document *HtmlDocument) appendNode(node *HtmlNode) {
	node.document = document

	if document.head == nil && node.NodeName() == "head" {
		document.head = node
	}

	if document.body == nil && node.NodeName() == "body" {
		document.body = node
	}

	document.Nodes = append(document.Nodes, node)
}

func (document *HtmlDocument) Head() *HtmlNode {
	return document.head
}

func (document *HtmlDocument) Body() *HtmlNode {
	return document.body
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
		node.Parent = parent
	} else {
		// stack is empty, as this is top-most node
		// we will add it to document
		document.appendNode(node)
	}

	// push this node to stack
	stack.push(node)
}

func (document *HtmlDocument) NumNodes() int {
	if document.Nodes == nil {
		return 0
	}

	return len(document.Nodes)
}

func (document *HtmlDocument) InsertFirrst(node *HtmlNode) {
	document.Nodes = append([]*HtmlNode{node}, document.Nodes...)
}
