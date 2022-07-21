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

func TestSimplestHtml(t *testing.T) {
	doc, err := getDoc("<html>Hello World</html>")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.Length())
	assert.Equal(t, 1, doc.nodes[0].NumChildren())
}

func TestOnlyString(t *testing.T) {
	doc, err := getDoc("Hello World")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.Length())
	assert.Equal(t, TextNode, doc.nodes[0].NodeType)
	assert.Equal(t, "Hello World", doc.nodes[0].Data)
}

func TestHead(t *testing.T) {
	doc, err := getDoc("<html><head>Hello World</head></html>")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.Length())
	assert.Equal(t, 1, doc.nodes[0].NumChildren())
	assert.Equal(t, 1, doc.nodes[0].Children[0].NumChildren())
	assert.Equal(t, doc.nodes[0].Children[0], doc.AsHtmlDocument().Head())
}

func TestHeadWithError(t *testing.T) {
	doc, err := getDoc("<html><head>Hello World</head><head>second head</head></html>")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.Length())
	assert.Equal(t, 2, doc.nodes[0].NumChildren())
	assert.Equal(t, 1, doc.nodes[0].Children[0].NumChildren())
	assert.Equal(t, 1, doc.nodes[0].Children[1].NumChildren())
	assert.Equal(t, doc.nodes[0].Children[0], doc.AsHtmlDocument().Head())

	// no head
	doc, err = getDoc("<html></html>")
	assert.NoError(t, err)
	assert.Nil(t, doc.AsHtmlDocument().Head())
	assert.Equal(t, 0, doc.GetElementsByName("head").Length())
}

func TestBody(t *testing.T) {
	doc, err := getDoc("<html><body>Hello World</body></html>")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.Length())
	assert.Equal(t, 1, doc.nodes[0].NumChildren())
	assert.Equal(t, 1, doc.nodes[0].Children[0].NumChildren())
	assert.Equal(t, doc.nodes[0].Children[0], doc.AsHtmlDocument().Body())

	// no body
	doc, err = getDoc("<html></html>")
	assert.NoError(t, err)
	assert.Nil(t, doc.AsHtmlDocument().Body())
	assert.Equal(t, 0, doc.GetElementsByName("body").Length())
}

func TestDoctype(t *testing.T) {
	doc, err := getDoc("<!doctype html><html />")
	assert.NoError(t, err)

	assert.Equal(t, 2, doc.Length())
	assert.Equal(t, DoctypeNode, doc.nodes[0].NodeType)
}

func TestComment(t *testing.T) {
	doc, err := getDoc("<!doctype html><html><!-- this is a comment --></html>")
	assert.NoError(t, err)

	assert.Equal(t, 3, doc.Length())
	assert.Equal(t, DoctypeNode, doc.nodes[0].NodeType)
	assert.Equal(t, CommentNode, doc.nodes[2].NodeType)
	assert.Equal(t, " this is a comment ", doc.nodes[2].Data)
}

func TestEmptyText(t *testing.T) {
	doc, err := getDoc("<html> </html>")
	assert.NoError(t, err)

	assert.Equal(t, 1, doc.Length())
}
