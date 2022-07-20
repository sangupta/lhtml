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

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getRemoveDoc() (*HtmlDocument, error) {
	html := "<html><head>Hello World</head><head>second head</head></html>"
	reader := strings.NewReader(html)
	return Parse(reader)
}

func TestRemoveAll(t *testing.T) {
	doc, err := getRemoveDoc()
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

func TestRemoveNode(t *testing.T) {
	doc, err := getRemoveDoc()
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

func TestRemoveAllChildren(t *testing.T) {
	doc, err := getDoc("<html />")
	assert.NoError(t, err)

	doc.Nodes[0].RemoveAllChildren()

	// if you have kids
	doc, err = getDoc("<html><head /><body /></html>")
	assert.NoError(t, err)

	doc.Nodes[0].RemoveAllChildren()
}

func TestRemoveMeNode(t *testing.T) {
	doc, err := getDoc("<html><head /><body /></html>")
	assert.NoError(t, err)

	head := doc.Head()
	// try removing head
	assert.Equal(t, 2, doc.Nodes[0].NumChildren())
	assert.True(t, head.RemoveMe())
	assert.Equal(t, 1, doc.Nodes[0].NumChildren())

	// removing again?
	assert.False(t, head.RemoveMe())
	assert.Equal(t, 1, doc.Nodes[0].NumChildren())

	// remove html?
	assert.Equal(t, 1, doc.NumNodes())
	assert.True(t, doc.Nodes[0].RemoveMe())
	assert.Equal(t, 0, doc.NumNodes())
}

func TestRemoveChildNode(t *testing.T) {
	doc, err := getDoc("<html><head><title></title></head><body /></html>")
	assert.NoError(t, err)

	head := doc.Head()
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
	assert.False(t, doc.Nodes[0].RemoveChild(node))
}

func TestPlainText(t *testing.T) {
	doc, err := getDoc("hello world")
	assert.Nil(t, err)

	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, TextNode, doc.Nodes[0].NodeType)
}
