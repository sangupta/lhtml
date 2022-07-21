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

func TestNumChildrenDoc(t *testing.T) {
	node := HtmlDocument{}

	// must check for `nil` children slice
	assert.Equal(t, 0, node.NumNodes())
}

func TestGetDocType(t *testing.T) {
	doc, err := getDoc("")
	assert.NoError(t, err)

	assert.Nil(t, doc.GetDocType()) // empty string

	doc, err = getDoc("<html />")
	assert.Nil(t, doc.GetDocType())

	doc, err = getDoc("<!doctype html><html />")
	assert.NotNil(t, doc.GetDocType())
}
