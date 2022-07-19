package core

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getBasicDoc() (*HtmlDocument, error) {
	html := "<html><head>Hello World</head><head>second head</head></html>"
	reader := strings.NewReader(html)
	return Parse(reader)
}

func TestRemoveAll(t *testing.T) {
	doc, err := getBasicDoc()
	assert.NoError(t, err)

	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 2, doc.Nodes[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[1].NumChildren())
	assert.Equal(t, doc.Nodes[0].Children[0], doc.Head())

	doc.RemoveAllNodes()
	assert.Equal(t, 0, doc.NumNodes())
	assert.True(t, doc.IsEmpty())
}

func TestRemoveNode(t *testing.T) {
	doc, err := getBasicDoc()
	assert.NoError(t, err)

	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 2, doc.Nodes[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[0].NumChildren())
	assert.Equal(t, 1, doc.Nodes[0].Children[1].NumChildren())
	assert.Equal(t, doc.Nodes[0].Children[0], doc.Head())

	doc.RemoveNode(doc.Nodes[0])

	assert.Equal(t, 0, doc.NumNodes())
	assert.True(t, doc.IsEmpty())
}
