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

func newNode(name string) *HtmlNode {
	return &HtmlNode{
		TagName: name,
	}
}
