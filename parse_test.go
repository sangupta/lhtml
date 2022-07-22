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
