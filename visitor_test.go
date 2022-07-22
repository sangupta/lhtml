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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyDoc(t *testing.T) {
	doc, err := ParseHtmlString("")
	assert.NoError(t, err)

	s := ""
	called := 0
	visitor := func(node *HtmlNode) bool {
		called++
		if node.NodeType != ElementNode {
			return true
		}
		s = s + " " + node.NodeName()
		return true
	}

	doc.Traverse(visitor)
	assert.Equal(t, "", s)
	assert.Equal(t, 0, called)
}

// check if a `nil` visitor is passed.
func TestNilVisitorDoc(t *testing.T) {
	doc, err := ParseHtmlString("<html><head /></html>")
	assert.NoError(t, err)

	doc.Traverse(nil)
}

func TestNilVisitorNode(t *testing.T) {
	doc, err := ParseHtmlString("<html><head><title>Hello world</title></head><body><div>Hello world</div></body></html>")
	assert.NoError(t, err)

	assert.False(t, doc.nodes[0]._children[0].Traverse(nil))
}

func TestTraverseDoc(t *testing.T) {
	doc, err := ParseHtmlString("<html><head><title>Hello world</title></head><body><div>Hello world</div></body></html>")
	assert.NoError(t, err)

	s := ""
	visitor := func(node *HtmlNode) bool {
		if node.NodeType != ElementNode {
			return true
		}
		s = s + " " + node.NodeName()
		return true
	}

	doc.Traverse(visitor)

	assert.Equal(t, " html head title body div", s)
}

// test traversal using a child node to start at.
func TestTraverseNode(t *testing.T) {
	doc, err := ParseHtmlString("<html><head><title>Hello world</title></head><body><div>Hello world</div></body></html>")
	assert.NoError(t, err)

	s := ""
	visitor := func(node *HtmlNode) bool {
		if node.NodeType != ElementNode {
			return true
		}
		s = s + " " + node.NodeName()
		return true
	}

	doc.First().First().Traverse(visitor)

	assert.Equal(t, " head title", s)
}

// test traversal breaking at the document level
func TestTraverseDocBreak(t *testing.T) {
	doc, err := ParseHtmlString("<html><head><title>Hello world</title></head><body><div>Hello world</div></body></html>")
	assert.NoError(t, err)

	s := ""
	visitor := func(node *HtmlNode) bool {
		if node.NodeType != ElementNode {
			return true
		}
		if node.NodeName() == "head" {
			return false
		}
		s = s + " " + node.NodeName()
		return true
	}

	doc.Traverse(visitor)

	assert.Equal(t, " html", s)
}

// test traversal breaking at a certain node
func TestTraverseNodeBreak(t *testing.T) {
	doc, err := ParseHtmlString("<html><head><title>Hello world</title></head><body><div>Hello world</div></body></html>")
	assert.NoError(t, err)

	s := ""
	visitor := func(node *HtmlNode) bool {
		if node.NodeType != ElementNode {
			return true
		}
		if node.NodeName() == "title" {
			return false
		}
		s = s + " " + node.NodeName()
		return true
	}

	doc.Traverse(visitor)
	assert.Equal(t, " html head", s)
}
