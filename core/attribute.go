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

//
// Holds the values for an attribute pair.
//
type HtmlAttribute struct {
	Name  string // the name of this attribute
	Value string // the value of this attribute
}

//
// Check if the node has an attribute with the given name.
//
// Returns `true` if the an attribute exists, `false`
// otherwise
//
func (node *HtmlNode) HasAttribute(key string) bool {
	attr := node.GetAttribute(key)
	return attr != nil
}

//
// Find and return the first attribute with the given name.
//
// Returns an `HtmlNode` if found, `nil` otherwise
//
func (node *HtmlNode) GetAttribute(key string) *HtmlAttribute {
	if !node.HasAttributes {
		return nil
	}

	for _, attr := range node.Attributes {
		if attr.Name == key {
			return attr
		}
	}

	return nil
}

//
// Find and return all attributes that have the given name.
//
// Returns a slice of `HtmlNode` if found, `nil` otherwise
//
func (node *HtmlNode) GetAttributes(key string) []*HtmlAttribute {
	if !node.HasAttributes {
		return nil
	}

	result := make([]*HtmlAttribute, 0)
	for _, attr := range node.Attributes {
		if attr.Name == key {
			result = append(result, attr)
		}
	}

	return result
}

//
// Find and return the `HtmlAttribute` which has the given name and value
//
// Returns either a `HtmlAttribute` instance, `nil` otherwise
//
func (node *HtmlNode) GetAttributeWithValue(key string, value string) *HtmlAttribute {
	if !node.HasAttributes {
		return nil
	}

	for _, attr := range node.Attributes {
		if attr.Name == key && attr.Value == value {
			return attr
		}
	}

	return nil
}
