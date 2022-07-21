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

func getDoc(html string) (*HtmlElements, error) {
	reader := strings.NewReader(html)
	return ParseHtml(reader)
}

func getAttributeDoc() (*HtmlElements, error) {
	html := "<html class='a1' class='b1' class='c1'>Hello World</html>"
	reader := strings.NewReader(html)
	return ParseHtml(reader)
}

func TestAttributes(t *testing.T) {
	doc, err := getAttributeDoc()
	assert.NoError(t, err)

	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 1, doc.nodes[0].NumChildren())

	node := doc.nodes[0]
	assert.True(t, node.HasAttribute("class"))
	assert.False(t, node.HasAttribute("id"))

	assert.Equal(t, "a1", node.GetAttribute("class").Value)
	assert.Equal(t, 3, len(node.GetAttributes("class")))

	assert.NotNil(t, node.GetAttributeWithValue("class", "b1"))
	assert.NotNil(t, node.GetAttributeWithValue("class", "c1"))
	assert.NotNil(t, node.GetAttributeWithValue("class", "a1"))

	assert.Nil(t, node.GetAttributeWithValue("class", "d1"))
}

func TestEmptyAttributes(t *testing.T) {
	doc, err := getDoc("<html>Hello World</html>")
	assert.NoError(t, err)

	node := doc.nodes[0]
	assert.False(t, node.HasAttribute("class"))
	assert.Nil(t, node.GetAttribute("class"))
	assert.Nil(t, node.GetAttributes("class"))
	assert.Nil(t, node.GetAttributeWithValue("class", "b1"))
}
