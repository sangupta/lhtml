# lhtml

[![Build Status](https://github.com/sangupta/lhtml/actions/workflows/unittest.yml/badge.svg?branch=main)](https://github.com/sangupta/lhtml/actions)
[![Code Coverage](https://codecov.io/gh/sangupta/lhtml/branch/main/graphs/badge.svg?branch=main)](https://codecov.io/gh/sangupta/lhtml)
[![go.mod version](https://img.shields.io/github/go-mod/go-version/sangupta/lhtml.svg)](https://github.com/sangupta/lhtml)
![GitHub](https://img.shields.io/github/license/sangupta/lhtml)

`lhtml` is a lenient HTML parser for Go. 

It differs from the standard `html` package because it will not re-order 
any of the encountered elements, nor will it try to sanitize your HTML 
file. This package is intended to be used for HTML-template based systems 
which want to process their own custom tags and attributes.

# Table of contents

* [Features](#features)
* [API](#api)
* [Usage](#usage)
* [Examples](#examples)
* [Changelog](#changelog)
* [Related Projects](#related-projects)
* [License](#license)

# Features

`lhtml` diffes from standard `html` parser in the following ways:

* Single parsing funtion that handles both documents as well as fragments
  - `ParseHtml`
* You may allow tags to have multiple attributes with same name
  - `ParseOption#AllowMultipleAttributesWithSameName`
* No sanitization of the resulting DOM
  - [example](#no-dom-sanitization)
* Provides node discovery functions
  - `GetElementById`
  - `GetElementsByName`
  - `GetBefore`
  - `GetAfter`
  - `Get` (at index)
  - `First`
  - `Last`
* Manipulation functions
  - `InsertFirst`
  - `InsertLast`
  - `EmptyChildren`
  - `Remove`
  - `Replace`
* Visitor functions when building tree, or to walk tree
  - `Traverse(visitor)` ([example](#traversing-the-dom))

# API

`lhtml` only has a single API that works both on the HTML document as
well as HTML fragments. 

```go
func ParseHtml(reader io.Reader) (*core.HtmlDocument, error)
```

We also expose a convenience method in case you would like to use `strings`
instead of a `io.Reader`:

```go
func ParseHtmlString(html string) (*core.HtmlDocument, error)
```

# Usage

Simply add the library to your project:

```sh
$ go get github.com/sangupta/lhtml@v0.1.0
```

And then, use it to parse your HTML markup:

```go
import (
    "github.com/sangupta/lhtml"
    "github.com/sangupta/lhtml/core"
)

func test() {
    html := "<html class='test1' class='test2' custom:title='hello'>Hello World <custom:PageBody /></html>"
    doc, err := lhtml.ParseHtmlString(htmlString)
    if err != nil {
        panic(err)
    }

    visitor := func(node *core.HtmlNode) bool {
        if node.NodeType == core.ElementNode {
            fmt.Println(node.TagName)
        }

        return true
    }
}
```

# Examples

## No DOM sanitization

For example, the HTML `title` tag cannot contain another tag. Given the
following html:

```html
<html>
    <head>
        <title>
            <custom:PageTitle />
        </title>
    </head>
</html>
```

The standard Go implementation will parse it to:

```html
<html>
    <head>
        <title>
            &lt; custom:PageTitle /&gt;
        </title>
    </head>
</html>
```

However, when using `lhtml` you will get the exact markup as defined
above. It is left to the callee code on how it wants to interpret and
use the parsed DOM nodes.

## Traversing the DOM

```go
func test() {
    doc, err := lhtml.ParseString("<html><head><title>Example</title></head><body><h1>Hello World</h1></body></html>")
    if err != nil {
        panic(err)
    }

    s := ""
	called := 0
	visitor := func(node *HtmlNode) bool {
		called++
		if node.NodeType != ElementNode {
			return true
		}
		s = s + " " + node.NodeName()
		return true
	}

    doc.Traverse(visitor)
	fmt.Println(s)          // " html head title body h1"
    fmt.Println(called)     // 7 (5 element nodes, 2 text nodes)
}
```

# Changelog

* **Version 0.1.0**
  - Initial release `$ go get github.com/sangupta/lhtml@v0.1.0`

# Related projects

* [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)

# License

MIT License. Copyright (C) 2022, Sandeep Gupta.
