package lhtml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDocType(t *testing.T) {
	doc, err := getDoc("")
	assert.NoError(t, err)

	assert.Nil(t, doc.AsHtmlDocument().GetDocType()) // empty string

	doc, err = getDoc("<html />")
	assert.Nil(t, doc.AsHtmlDocument().GetDocType())

	doc, err = getDoc("<!doctype html><html />")
	assert.NotNil(t, doc.AsHtmlDocument().GetDocType())
}

func TestHead(t *testing.T) {
	// on an empty string
	elements, _ := ParseHtmlString("")
	assert.Nil(t, elements.AsHtmlDocument().Head())

	// on direct head
	elements, _ = ParseHtmlString("<head />")
	assert.Nil(t, elements.AsHtmlDocument().Head())

	// in a proper html string
	elements, _ = ParseHtmlString("<html><head>Hello World</head></html>")

	assert.Equal(t, 1, elements.Length())
	assert.Equal(t, 1, elements.First().NumChildren())
	assert.Equal(t, 1, elements.First().First().NumChildren())
	assert.Equal(t, elements.First().First(), elements.AsHtmlDocument().Head())
}

func TestHeadWithError(t *testing.T) {
	doc, err := getDoc("<html><head>Hello World</head><head>second head</head></html>")

	assert.NoError(t, err)
	assert.Equal(t, 1, doc.Length())
	assert.Equal(t, 2, doc.First().NumChildren())
	assert.Equal(t, 1, doc.First().First().NumChildren())
	assert.Equal(t, 1, doc.First().First().NumChildren())
	assert.Equal(t, doc.First().First(), doc.AsHtmlDocument().Head())

	// no head
	doc, err = getDoc("<html></html>")
	assert.NoError(t, err)
	assert.Nil(t, doc.AsHtmlDocument().Head())
	assert.Equal(t, 0, doc.GetElementsByName("head").Length())
}

func TestBody(t *testing.T) {
	// check if body exists or not
	elements, _ := ParseHtmlString("<html><body>Hello World</body></html>")
	assert.Equal(t, 1, elements.Length())                      // html node
	assert.Equal(t, 1, elements.First().NumChildren())         // body node
	assert.Equal(t, 1, elements.First().First().NumChildren()) // text node
	assert.Equal(t, elements.First().First(), elements.AsHtmlDocument().Body())

	// no body
	elements, _ = getDoc("<html></html>")
	assert.Nil(t, elements.AsHtmlDocument().Body())
	assert.Equal(t, 0, elements.GetElementsByName("body").Length())
}

func TestDoctype(t *testing.T) {
	doc, err := getDoc("<!doctype html><html />")
	assert.NoError(t, err)

	assert.Equal(t, 2, doc.Length())
	assert.Equal(t, DoctypeNode, doc.First().NodeType)
}

func TestReplaceHead(t *testing.T) {
	doc, err := getDoc("<html><head /></html>")
	assert.NoError(t, err)

	node := newNode("a1")
	assert.NotNil(t, doc.AsHtmlDocument().Head())

	// head is not direct descendant of head
	assert.False(t, doc.Replace(doc.AsHtmlDocument().Head(), node))
	assert.NotNil(t, doc.AsHtmlDocument().Head())
	assert.NotEqual(t, node, doc.AsHtmlDocument().Head())
}
