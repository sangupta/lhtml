package core

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getAttributeDoc() (*HtmlDocument, error) {
	html := "<html class='a1' class='b1' class='c1'>Hello World</html>"
	reader := strings.NewReader(html)
	return Parse(reader)
}

func TestAttributes(t *testing.T) {
	doc, err := getAttributeDoc()
	assert.NoError(t, err)

	assert.Equal(t, 1, doc.NumNodes())
	assert.Equal(t, 1, doc.Nodes[0].NumChildren())

	node := doc.Nodes[0]
	assert.True(t, node.HasAttribute("class"))
	assert.False(t, node.HasAttribute("id"))

	assert.Equal(t, "a1", node.GetAttribute("class").Value)
	assert.Equal(t, 3, len(node.GetAttributes("class")))

	assert.NotNil(t, node.GetAttributeWithValue("class", "b1"))
	assert.NotNil(t, node.GetAttributeWithValue("class", "c1"))
	assert.NotNil(t, node.GetAttributeWithValue("class", "a1"))

	assert.Nil(t, node.GetAttributeWithValue("class", "d1"))
}
