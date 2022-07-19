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

func TestStack(t *testing.T) {
	a := newNode("a")
	b := newNode("b")
	c := newNode("c")

	stack := newNodeStack()
	assert.Equal(t, 0, stack.NumNodes())
	stack.push(a)
	stack.push(b)
	assert.Equal(t, 2, stack.NumNodes())
	assert.Equal(t, b, stack.pop())
	assert.Equal(t, 1, stack.NumNodes())
	stack.push(c)
	assert.Equal(t, 2, stack.NumNodes())
	assert.False(t, stack.isEmpty())
	assert.Equal(t, c, stack.peek())
	assert.Equal(t, c, stack.pop())
	assert.Equal(t, a, stack.pop())
	assert.Nil(t, stack.peek())
	assert.Nil(t, stack.pop())
	assert.Equal(t, 0, stack.NumNodes())
	assert.True(t, stack.isEmpty())
}

func TestUninitializedStack(t *testing.T) {
	stack := nodeStack{}

	assert.Equal(t, 0, stack.NumNodes())
}
