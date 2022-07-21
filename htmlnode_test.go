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

func TestNumChildren(t *testing.T) {
	node := HtmlNode{}

	// must check for `nil` children slice
	assert.Equal(t, 0, node.NumChildren())
}

func getRemoveDoc() (*HtmlElements, error) {
	html := "<html><head>Hello World</head><head>second head</head></html>"
	reader := strings.NewReader(html)
	return ParseHtml(reader)
}

func TestNodeRemoveAllChildren(t *testing.T) {
	doc, err := getDoc("<html />")
	assert.NoError(t, err)

	doc.nodes[0].RemoveAllChildren()

	// if you have kids
	doc, err = getDoc("<html><head /><body /></html>")
	assert.NoError(t, err)

	doc.nodes[0].RemoveAllChildren()
}

func TestNodeRemoveMe(t *testing.T) {
	doc, err := getDoc("<html><head /><body /></html>")
	assert.NoError(t, err)

	head := doc.AsHtmlDocument().Head()
	// try removing head
	assert.Equal(t, 2, doc.nodes[0].NumChildren())
	assert.True(t, head.RemoveMe())
	assert.Equal(t, 1, doc.nodes[0].NumChildren())

	// removing again?
	assert.False(t, head.RemoveMe())
	assert.Equal(t, 1, doc.nodes[0].NumChildren())

	// remove html?
	assert.Equal(t, 1, doc.Length())
	assert.True(t, doc.nodes[0].RemoveMe())
	assert.Equal(t, 0, doc.Length())
}

func TestNodeRemoveChild(t *testing.T) {
	doc, err := getDoc("<html><head><title></title></head><body /></html>")
	assert.NoError(t, err)

	head := doc.AsHtmlDocument().Head()
	node := newNode("a1")

	// node which is not a child
	assert.Equal(t, 1, head.NumChildren())
	assert.False(t, head.RemoveChild(node))

	// direct child
	assert.Equal(t, 1, head.NumChildren())
	assert.True(t, head.RemoveChild(head.Children[0]))
	assert.Equal(t, 0, head.NumChildren())

	// another one
	doc, err = getDoc("<html />")
	assert.False(t, doc.nodes[0].RemoveChild(node))
}

func getReplaceDoc() (*HtmlElements, error) {
	html := "<html><head></head></html>"
	reader := strings.NewReader(html)
	return ParseHtml(reader)
}

func TestNodeReplaceMe(t *testing.T) {
	doc, err := getRemoveDoc()
	assert.NoError(t, err)

	assert.False(t, doc.nodes[0].ReplaceMe(nil))

	node := newNode("a1")
	assert.Equal(t, "html", doc.nodes[0].NodeName())
	assert.True(t, doc.nodes[0].ReplaceMe(node))
	assert.Equal(t, "a1", doc.nodes[0].NodeName())
}

func TestNodeReplaceChild(t *testing.T) {
	doc, err := getRemoveDoc()
	assert.NoError(t, err)

	node := newNode("a1")
	assert.False(t, doc.nodes[0].ReplaceChild(nil, node))
	assert.False(t, doc.nodes[0].ReplaceChild(doc.nodes[0].Children[0], nil))

	assert.Equal(t, "head", doc.nodes[0].Children[0].NodeName())
	assert.True(t, doc.nodes[0].ReplaceChild(doc.nodes[0].Children[0], node))
	assert.Equal(t, "a1", doc.nodes[0].Children[0].NodeName())

	// when node has no child
	doc, err = getDoc("<html></html>")
	assert.NoError(t, err)
	assert.False(t, doc.nodes[0].ReplaceChild(node, node))
}

func TestNodeGetElementsByName(t *testing.T) {
	doc, err := getDoc("<html><head><title>hello world</title></HEAD><body>Hello world</body></html>")
	assert.NoError(t, err)

	assert.NotNil(t, doc.GetElementsByName("html"))
	assert.NotNil(t, doc.nodes[0].GetElementsByName("head"))
	assert.NotNil(t, doc.nodes[0].GetElementsByName("HEAD"))
}
