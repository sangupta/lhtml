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

type HtmlDocument struct {
	HtmlElements
}

func (document *HtmlDocument) Head() *HtmlNode {
	elements := document.GetElementsByName("head")
	if elements == nil || elements.Length() == 0 {
		return nil
	}

	return elements.nodes[0]
}

func (document *HtmlDocument) Body() *HtmlNode {
	elements := document.GetElementsByName("body")
	if elements == nil || elements.Length() == 0 {
		return nil
	}

	return elements.nodes[0]
}

func (doc *HtmlDocument) GetDocType() *HtmlNode {
	if doc.IsEmpty() {
		return nil
	}

	for _, node := range doc.nodes {
		if node.NodeType == DoctypeNode {
			return node
		}
	}

	return nil
}
