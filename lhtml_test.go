package lhtml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test parsing an empty string.
func TestEmptyString(t *testing.T) {
	elements, err := ParseHtmlString("")
	assert.NoError(t, err)
	assert.Equal(t, 0, elements.Length(), "Empty string should return zero elements")
}

// Test parsing a blank string.
func TestBlankString(t *testing.T) {
	actual, err := ParseHtmlString("            ")
	assert.NoError(t, err)
	assert.Equal(t, 0, actual.Length(), "Blank string should zero elements")
}

// test reader mode with empty string.
func TestReaderWithEmptyString(t *testing.T) {
	html := ""
	reader := strings.NewReader(html)
	actual, err := ParseHtml(reader)
	assert.NoError(t, err)
	assert.Equal(t, 0, actual.Length(), "Empty string should return zero elements")
}

// test reader mode with blank string
func TestReaderWithBlankString(t *testing.T) {
	html := "          "
	reader := strings.NewReader(html)
	actual, err := ParseHtml(reader)
	assert.NoError(t, err)
	assert.Equal(t, 0, actual.Length(), "Blank string should return zero elements")
}

// Check when nil reader is passed.
func TestNilReader(t *testing.T) {
	actual, err := ParseHtml(nil)
	assert.NotNil(t, err)
	assert.Nil(t, actual)
}
