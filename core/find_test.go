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

func TestGetElementsByName(t *testing.T) {
	doc, err := getDoc("")
	assert.NoError(t, err)

	assert.Nil(t, doc.GetElementsByName("html"))
}

func TestGetElementsByNameNode(t *testing.T) {
	doc, err := getDoc("<html><head><title>hello world</title></HEAD><body>Hello world</body></html>")
	assert.NoError(t, err)

	assert.NotNil(t, doc.GetElementsByName("html"))
	assert.NotNil(t, doc.Nodes[0].GetElementsByName("head"))
	assert.NotNil(t, doc.Nodes[0].GetElementsByName("HEAD"))
}

func TestGetElementById(t *testing.T) {
	doc, err := getDoc("<html><head /></html>")
	assert.NoError(t, err)

	assert.Nil(t, doc.GetElementById(""))      // empty id
	assert.Nil(t, doc.GetElementById("hello")) // valid id

	// id but different case
	doc, err = getDoc("<html><head id='HELLO' /></html>")
	assert.NoError(t, err)
	assert.Nil(t, doc.GetElementById("hello"))

	// id same case
	doc, err = getDoc("<html><head id='HELLO' /></html>")
	assert.NoError(t, err)
	assert.NotNil(t, doc.GetElementById("HELLO"))

	// empty doc
	doc, err = getDoc("")
	assert.NoError(t, err)
	assert.Nil(t, doc.GetElementById("hello"))
}
