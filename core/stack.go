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
// A simple stack to hold our own `HtmlNode` objects.
//
type nodeStack struct {
	elements []*HtmlNode
}

//
// Create a new stack.
//
func newNodeStack() *nodeStack {
	return &nodeStack{
		elements: make([]*HtmlNode, 0),
	}
}

//
// Push given node to top of stack
//
func (stack *nodeStack) push(node *HtmlNode) {
	stack.elements = append(stack.elements, node)
}

//
// Pop an element from the stack. If there are no elements on
// the stack, this returns `nil`.
//
func (stack *nodeStack) pop() *HtmlNode {
	length := len(stack.elements)
	if length == 0 {
		return nil
	}

	element := stack.elements[length-1]
	stack.elements = stack.elements[0 : length-1]

	return element
}

//
// Peek the current element at top of stack. If there are no
// elements in stack, this shall return `nil`.
//
func (stack *nodeStack) peek() *HtmlNode {
	length := len(stack.elements)
	if length == 0 {
		return nil
	}

	element := stack.elements[length-1]
	return element
}

//
// Check if the stack is empty or not.
//
func (stack *nodeStack) isEmpty() bool {
	if len(stack.elements) == 0 {
		return true
	}

	return false
}

func (stack *nodeStack) NumNodes() int {
	if stack.elements == nil {
		return 0
	}

	return len(stack.elements)
}
