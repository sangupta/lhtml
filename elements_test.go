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

func TestNumChildrenDoc(t *testing.T) {
	node := HtmlElements{}

	// must check for `nil` children slice
	assert.Equal(t, 0, node.Length())
}

func TestDocReplaceEmpty(t *testing.T) {
	doc, err := getDoc("")
	assert.NoError(t, err)

	node1 := newNode("a1")
	node2 := newNode("b1")
	assert.False(t, doc.Replace(node1, node2))
}

func TestDocRemoveAll(t *testing.T) {
	html := "<html><head>Hello World</head><head>second head</head></html>"
	doc, err := getDoc(html)
	assert.NoError(t, err)

	assert.Equal(t, 1, doc.Length())
	assert.Equal(t, 2, doc.nodes[0].NumChildren())
	assert.Equal(t, 1, doc.nodes[0].Children[0].NumChildren())
	assert.Equal(t, 1, doc.nodes[0].Children[1].NumChildren())
	assert.Equal(t, doc.nodes[0].Children[0], doc.AsHtmlDocument().Head())

	// remove all on doc
	doc.Empty()
	assert.Equal(t, 0, doc.Length())
	assert.True(t, doc.IsEmpty())

	// empty doc
	doc, err = getDoc("")
	assert.NoError(t, err)
	assert.Equal(t, 0, doc.Length())
	doc.Empty()
	assert.Equal(t, 0, doc.Length())
}

func TestDocRemoveNode(t *testing.T) {
	html := "<html><head>Hello World</head><head>second head</head></html>"
	doc, err := getDoc(html)
	assert.NoError(t, err)

	assert.Equal(t, 1, doc.Length())
	assert.Equal(t, 2, doc.nodes[0].NumChildren())
	assert.Equal(t, 1, doc.nodes[0].Children[0].NumChildren())
	assert.Equal(t, 1, doc.nodes[0].Children[1].NumChildren())
	assert.Equal(t, doc.nodes[0].Children[0], doc.AsHtmlDocument().Head())

	doc.Remove(doc.nodes[0])

	assert.Equal(t, 0, doc.Length())
	assert.True(t, doc.IsEmpty())

	// empty doc
	doc, err = getDoc("")
	assert.NoError(t, err)
	assert.False(t, doc.Remove(newNode("a1")))

	// node not in doc
	doc, err = getDoc("<head id='hello' /><body />")
	assert.NoError(t, err)
	assert.Equal(t, 2, doc.Length())
	assert.False(t, doc.Remove(newNode("a1")))
	assert.False(t, doc.Remove(newNode("head")))
	assert.Equal(t, 2, doc.Length())
	assert.True(t, doc.Remove(doc.GetElementById("hello")))
	assert.Equal(t, 1, doc.Length())
}

func TestParsePlainText(t *testing.T) {
	doc, err := getDoc("hello world")
	assert.Nil(t, err)

	assert.Equal(t, 1, doc.Length())
	assert.Equal(t, TextNode, doc.nodes[0].NodeType)
}

func TestDocReplaceNode(t *testing.T) {
	html := "<html><head></head></html>"
	doc, err := getDoc(html)
	assert.NoError(t, err)

	node := newNode("a1")

	assert.False(t, doc.Replace(nil, node))
	assert.False(t, doc.Replace(node, nil))
	assert.False(t, doc.Replace(node, node))

	assert.Equal(t, "html", doc.nodes[0].NodeName())
	assert.True(t, doc.Replace(doc.nodes[0], node))
	assert.Equal(t, "a1", doc.nodes[0].NodeName())
}

func TestDocGetElementsByName(t *testing.T) {
	doc, err := getDoc("")
	assert.NoError(t, err)

	assert.Nil(t, doc.GetElementsByName("html"))
}

func TestDocGetElementById(t *testing.T) {
	doc, err := getDoc("<html><head /></html>")
	assert.NoError(t, err)

	assert.Nil(t, doc.GetElementById(""))      // empty id
	assert.Nil(t, doc.GetElementById("hello")) // valid id

	// id but different case
	doc, err = getDoc("<html><head id='HELLO' /></html>")
	assert.NoError(t, err)
	assert.Nil(t, doc.GetElementById("hello"))

	// id same case
	doc, err = getDoc("<html><head id='HELLO' /></html>")
	assert.NoError(t, err)
	assert.NotNil(t, doc.GetElementById("HELLO"))

	// empty doc
	doc, err = getDoc("")
	assert.NoError(t, err)
	assert.Nil(t, doc.GetElementById("hello"))
}
