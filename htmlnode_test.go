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

// test removing all children on a node
func TestNodeRemoveAllChildren(t *testing.T) {
	doc, err := getDoc("<html />")
	assert.NoError(t, err)

	doc.nodes[0].RemoveAllChildren()

	// if you have kids
	doc, err = getDoc("<html><head /><body /></html>")
	assert.NoError(t, err)

	doc.Get(0).RemoveAllChildren()
}

// test removing a particular node using `RemoveMe`
func TestNodeRemoveMe(t *testing.T) {
	doc, err := getDoc("<html><head /><body /></html>")
	assert.NoError(t, err)

	head := doc.AsHtmlDocument().Head()
	// try removing head
	assert.Equal(t, 2, doc.First().NumChildren())
	assert.True(t, head.RemoveMe())
	assert.Equal(t, 1, doc.First().NumChildren())

	// removing again?
	assert.False(t, head.RemoveMe())
	assert.Equal(t, 1, doc.First().NumChildren())

	// remove html?
	assert.Equal(t, 1, doc.Length())
	assert.True(t, doc.First().RemoveMe())
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
	assert.True(t, head.RemoveChild(head.First()))
	assert.Equal(t, 0, head.NumChildren())

	// another one
	doc, err = getDoc("<html />")
	assert.False(t, doc.First().RemoveChild(node))
}

func getReplaceDoc() (*HtmlElements, error) {
	html := "<html><head></head></html>"
	reader := strings.NewReader(html)
	return ParseHtml(reader)
}

func TestNodeReplaceMe(t *testing.T) {
	doc, err := getRemoveDoc()
	assert.NoError(t, err)

	assert.False(t, doc.First().ReplaceMe(nil))

	node := newNode("a1")
	assert.Equal(t, "html", doc.First().NodeName())
	assert.True(t, doc.First().ReplaceMe(node))
	assert.Equal(t, "a1", doc.First().NodeName())
}

func TestNodeReplaceChild(t *testing.T) {
	doc, err := getRemoveDoc()
	assert.NoError(t, err)

	node := newNode("a1")
	assert.False(t, doc.First().ReplaceChild(nil, node))
	assert.False(t, doc.First().ReplaceChild(doc.First().First(), nil))

	assert.Equal(t, "head", doc.First().First().NodeName())
	assert.True(t, doc.First().ReplaceChild(doc.First().First(), node))
	assert.Equal(t, "a1", doc.First().First().NodeName())

	// when node has no child
	doc, err = getDoc("<html></html>")
	assert.NoError(t, err)
	assert.False(t, doc.First().ReplaceChild(node, node))
}

func TestNodeGetElementsByName(t *testing.T) {
	doc, err := getDoc("<html><head><title>hello world</title></HEAD><body>Hello world</body></html>")
	assert.NoError(t, err)

	assert.NotNil(t, doc.GetElementsByName("html"))
	assert.NotNil(t, doc.First().GetElementsByName("head"))
	assert.NotNil(t, doc.First().GetElementsByName("HEAD"))
}
