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

import (
	"errors"
	"io"
	"strings"
)

type ParseOptions struct {
	CaseSensitiveAttributes bool
	AllowMultipleAttributes bool
}

func getDefaultOptions() *ParseOptions {
	return &ParseOptions{
		CaseSensitiveAttributes: false,
		AllowMultipleAttributes: false,
	}
}

//
// A loose HTML parser that just returns the tags and their
// attributes in the order they appear. It makes no assumption
// on if the tag is permitted inside another or not. For example,
// you cannot include `iframe` within an `input` tag. But,
// using this package you still will be able to parse such syntax.
//
// This package should be useful for parsing and working with
// html-like template syntax where we define our own custom tags
// that can emit or alter the behavior of the final HTML code.
//
// Thus, this function returns a wrapper over the actual
// `[]*html.Node` nodes parsed via the `html` package. This wrapper
// provides convenience functions to achieve some of the
// templating work quickly.
//
func ParseHtml(reader io.Reader) (*HtmlElements, error) {
	if reader == nil {
		return nil, errors.New("Reader is required to parse html.")
	}

	return ParseWithOptions(reader, getDefaultOptions())
}

//
// Parse the given string as a HTML document or a fragment. It
// is a convenience method and calls `ParseHtml(reader io.Reader)`
// internally.
//
// Returns the `HtmlDocument` and any error if encountered. If
// error is `nil`, the `HtmlDocument` instance would be available.
// If the error is not `nil`, the `HtmlDocument` will be `nil`.
//
func ParseHtmlString(html string) (*HtmlElements, error) {
	if len(html) == 0 {
		return newHtmlElements(), nil
	}

	reader := strings.NewReader(html)
	return ParseHtml(reader)
}
