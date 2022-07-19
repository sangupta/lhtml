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

* [Details](#details)
* [API](#api)
* [Usage](#usage)
* [Changelog](#changelog)
* [Related Projects](#related-projects)
* [License](#license)

# Details

`lhtml` diffes from standard `html` parser in the following ways:

* Tag can have multiple attributes with same name
* Tags can occur in any order, under any parent
* No sanitization of the resulting DOM
* Allows using custom tag/attributes
* Provides node discovery functions, such as, `GetElementById`
* Provides function to walk entire DOM tree easily
* Provides functions to `remove`, `replace` nodes

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
$ go get sangupta.com/lhtml
```

And then, use it to parse your HTML markup:

```go
import (
    "sangupta.com/lhtml"
    "sangupta.com/lhtml/core"
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

# Changelog

* Version 0.1.0
  - Initial release

# Related projects

* [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)

# License

MIT License. Copyright (C) 2022, Sandeep Gupta.
