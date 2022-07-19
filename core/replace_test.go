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

func getReplaceDoc() (*HtmlDocument, error) {
	html := "<html><head></head></html>"
	reader := strings.NewReader(html)
	return Parse(reader)
}

func TestReplaceNode(t *testing.T) {
	doc, err := getRemoveDoc()
	assert.NoError(t, err)

	node := newNode("a1")

	assert.False(t, doc.ReplaceNode(nil, node))
	assert.False(t, doc.ReplaceNode(node, nil))
	assert.False(t, doc.ReplaceNode(node, node))

	assert.Equal(t, "html", doc.Nodes[0].NodeName())
	assert.True(t, doc.ReplaceNode(doc.Nodes[0], node))
	assert.Equal(t, "a1", doc.Nodes[0].NodeName())
}

func TestReplaceMe(t *testing.T) {
	doc, err := getRemoveDoc()
	assert.NoError(t, err)

	assert.False(t, doc.Nodes[0].ReplaceMe(nil))

	node := newNode("a1")
	assert.Equal(t, "html", doc.Nodes[0].NodeName())
	assert.True(t, doc.Nodes[0].ReplaceMe(node))
	assert.Equal(t, "a1", doc.Nodes[0].NodeName())
}

func TestReplaceChild(t *testing.T) {
	doc, err := getRemoveDoc()
	assert.NoError(t, err)

	node := newNode("a1")
	assert.False(t, doc.Nodes[0].ReplaceChild(nil, node))
	assert.False(t, doc.Nodes[0].ReplaceChild(doc.Nodes[0].Children[0], nil))

	assert.Equal(t, "head", doc.Nodes[0].Children[0].NodeName())
	assert.True(t, doc.Nodes[0].ReplaceChild(doc.Nodes[0].Children[0], node))
	assert.Equal(t, "a1", doc.Nodes[0].Children[0].NodeName())
}

func TestReplaceEmptyDoc(t *testing.T) {
	doc, err := getDoc("")
	assert.NoError(t, err)

	node1 := newNode("a1")
	node2 := newNode("b1")
	assert.False(t, doc.ReplaceNode(node1, node2))
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
