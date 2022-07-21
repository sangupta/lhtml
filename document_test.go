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
	node := HtmlDocument{}

	// must check for `nil` children slice
	assert.Equal(t, 0, node.NumNodes())
}

func TestGetDocType(t *testing.T) {
	doc, err := getDoc("")
	assert.NoError(t, err)

	assert.Nil(t, doc.GetDocType()) // empty string

	doc, err = getDoc("<html />")
	assert.Nil(t, doc.GetDocType())

	doc, err = getDoc("<!doctype html><html />")
	assert.NotNil(t, doc.GetDocType())
}

func TestDocReplaceEmpty(t *testing.T) {
	doc, err := getDoc("")
	assert.NoError(t, err)

	node1 := newNode("a1")
	node2 := newNode("b1")
	assert.False(t, doc.ReplaceNode(node1, node2))
}

func TestDocRemoveAll(t *testing.T) {
	html := "<html><head>Hello World</head><head>second head</head></html>"
	doc, err := getDoc(html)
	assert.NoError(t, err)

	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 2, doc.Nodes[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[1].NumChildren())
	assert.Equal(t, doc.Nodes[0].Children[0], doc.Head())

	// remove all on doc
	doc.RemoveAllNodes()
	assert.Equal(t, 0, doc.NumNodes())
	assert.True(t, doc.IsEmpty())

	// empty doc
	doc, err = getDoc("")
	assert.NoError(t, err)
	assert.Equal(t, 0, doc.NumNodes())
	doc.RemoveAllNodes()
	assert.Equal(t, 0, doc.NumNodes())
}

func TestDocRemoveNode(t *testing.T) {
	html := "<html><head>Hello World</head><head>second head</head></html>"
	doc, err := getDoc(html)
	assert.NoError(t, err)

	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 2, doc.Nodes[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[1].NumChildren())
	assert.Equal(t, doc.Nodes[0].Children[0], doc.Head())

	doc.RemoveNode(doc.Nodes[0])

	assert.Equal(t, 0, doc.NumNodes())
	assert.True(t, doc.IsEmpty())

	// empty doc
	doc, err = getDoc("")
	assert.NoError(t, err)
	assert.False(t, doc.RemoveNode(newNode("a1")))

	// node not in doc
	doc, err = getDoc("<head id='hello' /><body />")
	assert.NoError(t, err)
	assert.Equal(t, 2, doc.NumNodes())
	assert.False(t, doc.RemoveNode(newNode("a1")))
	assert.False(t, doc.RemoveNode(newNode("head")))
	assert.Equal(t, 2, doc.NumNodes())
	assert.True(t, doc.RemoveNode(doc.GetElementById("hello")))
	assert.Equal(t, 1, doc.NumNodes())
}

func TestParsePlainText(t *testing.T) {
	doc, err := getDoc("hello world")
	assert.Nil(t, err)

	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, TextNode, doc.Nodes[0].NodeType)
}

func TestReplaceHead(t *testing.T) {
	doc, err := getDoc("<html><head /></html>")
	assert.NoError(t, err)

	node := newNode("a1")
	assert.NotNil(t, doc.Head())

	// head is not direct descendant of head
	assert.False(t, doc.ReplaceNode(doc.Head(), node))
	assert.NotNil(t, doc.Head())
	assert.NotEqual(t, node, doc.Head())
}

func TestDocReplaceNode(t *testing.T) {
	html := "<html><head></head></html>"
	doc, err := getDoc(html)
	assert.NoError(t, err)

	node := newNode("a1")

	assert.False(t, doc.ReplaceNode(nil, node))
	assert.False(t, doc.ReplaceNode(node, nil))
	assert.False(t, doc.ReplaceNode(node, node))

	assert.Equal(t, "html", doc.Nodes[0].NodeName())
	assert.True(t, doc.ReplaceNode(doc.Nodes[0], node))
	assert.Equal(t, "a1", doc.Nodes[0].NodeName())
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
