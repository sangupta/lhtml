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

package lhtml

//
// A wrapper representing a HTML document which is nothing
// but an array of HTML elements.
//
type HtmlDocument struct {
	HtmlElements // the child elements
}

//
// Returns the `head` element in the document if any.
// Only the top level elements are searched for the desired
// element. `nil` is returned if the document is empty
// or the element is not found.
//
func (document *HtmlDocument) Head() *HtmlNode {
	if document.IsEmpty() {
		return nil
	}

	// do we have an html node?
	html := document.GetChildrenByName("html")
	if html.Length() == 0 {
		return nil
	}

	return html.First().GetElementsByName("head").First()
}

//
// Returns the `body` element in the document if any.
// Only the top level elements are searched for the desired
// element. `nil` is returned if the document is empty
// or the element is not found.
//
func (document *HtmlDocument) Body() *HtmlNode {
	if document.IsEmpty() {
		return nil
	}

	// do we have an html node?
	html := document.GetChildrenByName("html")
	if html.Length() == 0 {
		return nil
	}

	return html.First().GetElementsByName("body").First()
}

//
// Return the DocType element associated if any.
// Only the top level elements are searched for the desired
// element. `nil` is returned if the document is empty
// or the element is not found.
//
func (document *HtmlDocument) GetDocType() *HtmlNode {
	if document.IsEmpty() {
		return nil
	}

	for _, node := range document.nodes {
		if node.NodeType == DoctypeNode {
			return node
		}
	}

	return nil
}
