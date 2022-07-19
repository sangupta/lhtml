package core

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplestHtml(t *testing.T) {
	html := "<html>Hello World</html>"
	reader := strings.NewReader(html)
	doc, err := Parse(reader)

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 1, doc.Nodes[0].NumChildren())
}

func TestOnlyString(t *testing.T) {
	html := "Hello World"
	reader := strings.NewReader(html)
	doc, err := Parse(reader)

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, TextNode, doc.Nodes[0].NodeType)
}

func TestHead(t *testing.T) {
	html := "<html><head>Hello World</head></html>"
	reader := strings.NewReader(html)
	doc, err := Parse(reader)

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 1, doc.Nodes[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[0].NumChildren())
	assert.Equal(t, doc.Nodes[0].Children[0], doc.Head())
}

func TestHeadWithError(t *testing.T) {
	html := "<html><head>Hello World</head><head>second head</head></html>"
	reader := strings.NewReader(html)
	doc, err := Parse(reader)

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 2, doc.Nodes[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[1].NumChildren())
	assert.Equal(t, doc.Nodes[0].Children[0], doc.Head())
}
