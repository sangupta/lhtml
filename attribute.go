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

import "strings"

//
// Holds the values for an attribute pair.
//
type HtmlAttribute struct {
	Name  string // the name of this attribute
	Value string // the value of this attribute
}

func (node *HtmlNode) HasAttributes() bool {
	if node.Attributes == nil || len(node.Attributes) == 0 {
		return false
	}

	return true
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
	if !node.HasAttributes() {
		return nil
	}

	for _, attr := range node.Attributes {
		if strings.EqualFold(attr.Name, key) {
			return attr
		}
	}

	return nil
}

func (node *HtmlNode) RemoveAttribute(key string) bool {
	if !node.HasAttributes() {
		return false
	}

	modified := false
	newAttributes := make([]*HtmlAttribute, 0)
	for _, attr := range node.Attributes {
		if strings.EqualFold(attr.Name, key) {
			modified = true
			continue
		}
		newAttributes = append(newAttributes, attr)
	}

	if modified {
		node.Attributes = newAttributes
	}

	return modified
}

func (node *HtmlNode) SetAttribute(key string, value string) bool {
	if node.HasAttributes() {
		// do we have the attribute with matching name?
		for _, attr := range node.Attributes {
			if strings.EqualFold(attr.Name, key) {
				attr.Value = value
				return true
			}
		}
	}

	// no attributes in node, add one
	if node.Attributes == nil {
		node.Attributes = make([]*HtmlAttribute, 0)
	}

	node.Attributes = append(node.Attributes, &HtmlAttribute{
		Name:  key,
		Value: value,
	})

	return true
}

// //
// // This method remove all duplicate attributes from the node.
// // If `merge` is set to `true`, the value from each attribute is
// // combined together a space. If it's set to `false` the attributes
// // occuring later are removed.
// //
// func (node *HtmlNode) RemoveDuplicateAttributes(merge bool) {

// }

//
// Find and return all attributes that have the given name.
//
// Returns a slice of `HtmlNode` if found, `nil` otherwise
//
func (node *HtmlNode) GetAttributes(key string) []*HtmlAttribute {
	if !node.HasAttributes() {
		return nil
	}

	result := make([]*HtmlAttribute, 0)
	for _, attr := range node.Attributes {
		if strings.EqualFold(attr.Name, key) {
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
	if !node.HasAttributes() {
		return nil
	}

	for _, attr := range node.Attributes {
		if attr.Name == key && attr.Value == value {
			return attr
		}
	}

	return nil
}
