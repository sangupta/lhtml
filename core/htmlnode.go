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
	Parent        *HtmlNode `json:"-"`
	HasAttributes bool
	Attributes    []*HtmlAttribute
	Children      []*HtmlNode
	IsSelfClosing bool
	NodeType      HtmlNodeType
	Data          string
	document      *HtmlDocument // the document node that this node belongs to
}

func (node *HtmlNode) NodeName() string {
	return strings.TrimSpace(node.TagName)
}

func (node *HtmlNode) NumChildren() int {
	if node.Children == nil {
		return 0
	}

	return len(node.Children)
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

func (node *HtmlNode) detach() {
	node.Parent = nil
	node.document = nil
}

func newNode(name string) *HtmlNode {
	return &HtmlNode{
		TagName: name,
	}
}
