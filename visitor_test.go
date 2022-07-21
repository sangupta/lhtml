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
	"testing"

	"github.com/stretchr/testify/assert"
)

func getVisitorDoc() (*HtmlDocument, error) {
	html := "<html><head><title>Hello world</title></head><body><div>Hello world</div></body></html>"
	reader := strings.NewReader(html)
	return ParseHtml(reader)
}

func TestEmptyDoc(t *testing.T) {
	html := ""
	reader := strings.NewReader(html)
	doc, err := ParseHtml(reader)
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

func TestNilVisitorDoc(t *testing.T) {
	doc, err := getDoc("<html><head /></html>")
	assert.NoError(t, err)

	doc.Traverse(nil)
}

func TestNilVisitorNode(t *testing.T) {
	doc, err := getVisitorDoc()
	assert.NoError(t, err)

	assert.False(t, doc.Nodes[0].Children[0].Traverse(nil))
}

func TestTraverseDoc(t *testing.T) {
	doc, err := getVisitorDoc()
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

func TestTraverseNode(t *testing.T) {
	doc, err := getVisitorDoc()
	assert.NoError(t, err)

	s := ""
	visitor := func(node *HtmlNode) bool {
		if node.NodeType != ElementNode {
			return true
		}
		s = s + " " + node.NodeName()
		return true
	}

	doc.Nodes[0].Children[0].Traverse(visitor)

	assert.Equal(t, " head title", s)
}

func TestTraverseDocBreak(t *testing.T) {
	doc, err := getVisitorDoc()
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

func TestTraverseNodeBreak(t *testing.T) {
	doc, err := getVisitorDoc()
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
