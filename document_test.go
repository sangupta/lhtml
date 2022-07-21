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

func TestReplaceHead(t *testing.T) {
	doc, err := getDoc("<html><head /></html>")
	assert.NoError(t, err)

	node := newNode("a1")
	assert.NotNil(t, doc.AsHtmlDocument().Head())

	// head is not direct descendant of head
	assert.False(t, doc.ReplaceNode(doc.AsHtmlDocument().Head(), node))
	assert.NotNil(t, doc.AsHtmlDocument().Head())
	assert.NotEqual(t, node, doc.AsHtmlDocument().Head())
}
