package lhtml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyString(t *testing.T) {
	actual, err := ParseHtmlString("")
	assert.NoError(t, err)
	assert.Equal(t, 0, actual.NumNodes(), "Empty string should return a document with zero nodes")
}

func TestBlankString(t *testing.T) {
	actual, err := ParseHtmlString("            ")
	assert.NoError(t, err)
	assert.Equal(t, 0, actual.NumNodes(), "Blank string should return a document with zero nodes")
}

func TestReaderWithEmptyString(t *testing.T) {
	html := ""
	reader := strings.NewReader(html)
	actual, err := ParseHtml(reader)
	assert.NoError(t, err)
	assert.Equal(t, 0, actual.NumNodes(), "Empty string should return a document with zero nodes")
}

func TestReaderWithBlankString(t *testing.T) {
	html := "          "
	reader := strings.NewReader(html)
	actual, err := ParseHtml(reader)
	assert.NoError(t, err)
	assert.Equal(t, 0, actual.NumNodes(), "Blank string should return a document with zero nodes")
}

func TestNilReader(t *testing.T) {
	actual, err := ParseHtml(nil)
	assert.Error(t, err)
	assert.Nil(t, actual)
}
