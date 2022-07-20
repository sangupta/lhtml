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
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

const whitespace = " \t\r\n\f"

//
// Generic parse function that takes the reader
// and tries to return the `HtmlDocument` on a best-effort
// basis.
//
func Parse(reader io.Reader) (*HtmlDocument, error) {
	// create a new stack
	stack := newNodeStack()

	// create a tokenizer
	tokenizer := html.NewTokenizer(reader)

	// create a document instance that we can use
	document := NewHtmlDocument()

	// let's start parsing
	for {
		token := tokenizer.Next()

		err := parseToken(document, tokenizer, &token, stack)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return document, nil
}

//
// Parse the given token and return an error, if any.
//
func parseToken(document *HtmlDocument, tokenizer *html.Tokenizer, token *html.TokenType, stack *nodeStack) error {
	switch *token {

	// handle the doctype token
	case html.DoctypeToken:
		return handleDocTypeToken(document, tokenizer)

	// handle error tokens
	case html.ErrorToken:
		return handleErrorToken(document, tokenizer)

	case html.TextToken:
		return handleTextToken(document, stack, tokenizer)

	// just add the comment as is
	case html.CommentToken:
		return handleCommentToken(document, tokenizer)

	// start of a token
	case html.StartTagToken:
		return handleStartTagToken(document, stack, tokenizer, false)

	// self-sufficient token
	case html.SelfClosingTagToken:
		return handleStartTagToken(document, stack, tokenizer, true)

	case html.EndTagToken:
		handleEndTagToken(document, stack, tokenizer)
	}

	// all processed
	return nil
}

//
// Read an element node from the tokenizer. An element node is
// basically a tag, such as `<br />`  or `<div ...>`.
// This method reads the tag as well as any attributes assigned
// to this tag.
//
func readElementNode(tokenizer *html.Tokenizer) *HtmlNode {
	// when a tag starts, we read the tag name
	tagName, hasAttributes := tokenizer.TagName()
	node := HtmlNode{
		TagName:       string(tagName),
		HasAttributes: hasAttributes,
		IsSelfClosing: false,
		NodeType:      ElementNode,
	}

	// copy attributes as needed
	if hasAttributes {
		for {
			key, value, more := tokenizer.TagAttr()
			if key != nil && value != nil {
				node.addAttribute(string(key), string(value))
			}
			if !more {
				break
			}
		}
	}

	return &node
}

func handleDocTypeToken(document *HtmlDocument, tokenizer *html.Tokenizer) error {
	docType := tokenizer.Token().Data

	// we currently do not parse doc type to reveal information
	// so add it to data attribute
	node := HtmlNode{
		Data:     docType,
		NodeType: DoctypeNode,
		Parent:   nil,
	}
	document.InsertFirrst(&node)
	return nil
}

func handleErrorToken(document *HtmlDocument, tokenizer *html.Tokenizer) error {
	err := tokenizer.Err()

	// if we ran into end-of-file we will return from
	// where ever we are
	if err != nil && err == io.EOF {
		return err
	}

	fmt.Println(err)

	// TODO: handle other error scenarios
	return errors.New("error not handled")
}

//
// we just got some text
// this may be an attribute
// or this may be some textnode as a child
// lets process
func handleTextToken(document *HtmlDocument, stack *nodeStack, tokenizer *html.Tokenizer) error {
	text := string(tokenizer.Text())
	text = strings.TrimLeft(text, whitespace)
	if len(text) == 0 {
		return nil
	}

	node := HtmlNode{
		NodeType: TextNode,
		Data:     text,
	}
	document.addNodeToStack(&node, stack)
	stack.pop()

	return nil
}

func handleCommentToken(document *HtmlDocument, tokenizer *html.Tokenizer) error {
	comment := tokenizer.Token().Data
	node := HtmlNode{
		Data:     comment,
		NodeType: CommentNode,
	}
	document.appendNode(&node)
	return nil
}

//
// for end token, we need to check if we have the right
// start tag at the top of the stack. Otherwise we may
// have to pop all the way down to see if we have a node
// with the end name
//
func handleEndTagToken(document *HtmlDocument, stack *nodeStack, tokenizer *html.Tokenizer) error {
	tagName, _ := tokenizer.TagName()

	// if stack is empty, this is a wrong tag
	if stack.isEmpty() {
		return errors.New("Encountered end tag when stack is empty: " + string(tagName))
	}

	// let's check what is at top
	element := stack.peek()
	if element.TagName == string(tagName) {
		// its the same tag, let's just pop and move ahead
		stack.pop()

		return nil
	}

	// this is not the same tag as the one at the top of stack
	// so we need to try and heal if we can, or just ignore this
	// TODO: fix this
	return errors.New("not implemented")
}

func handleStartTagToken(document *HtmlDocument, stack *nodeStack, tokenizer *html.Tokenizer, popElement bool) error {
	node := readElementNode(tokenizer)
	document.addNodeToStack(node, stack)

	if popElement {
		stack.pop()
	}

	return nil
}
