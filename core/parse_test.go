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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplestHtml(t *testing.T) {
	doc, err := getDoc("<html>Hello World</html>")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 1, doc.Nodes[0].NumChildren())
}

func TestOnlyString(t *testing.T) {
	doc, err := getDoc("Hello World")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, TextNode, doc.Nodes[0].NodeType)
	assert.Equal(t, "Hello World", doc.Nodes[0].Data)
}

func TestHead(t *testing.T) {
	doc, err := getDoc("<html><head>Hello World</head></html>")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 1, doc.Nodes[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[0].NumChildren())
	assert.Equal(t, doc.Nodes[0].Children[0], doc.Head())
}

func TestHeadWithError(t *testing.T) {
	doc, err := getDoc("<html><head>Hello World</head><head>second head</head></html>")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 2, doc.Nodes[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[1].NumChildren())
	assert.Equal(t, doc.Nodes[0].Children[0], doc.Head())
}

func TestBody(t *testing.T) {
	doc, err := getDoc("<html><body>Hello World</body></html>")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 1, doc.Nodes[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[0].NumChildren())
	assert.Equal(t, doc.Nodes[0].Children[0], doc.Body())
}
