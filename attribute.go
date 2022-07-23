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

func (node *HtmlNode) NumAttributes() int {
	if node.Attributes == nil {
		return 0
	}

	return len(node.Attributes)
}

func (node *HtmlNode) ContainsAttributes() bool {
	if node.NumAttributes() <= 0 {
		return false
	}

	return true
}

//
// Add a new attribute to this node. By design, we allow a single
// tag to hold multiple values for the same attribute name. This is
// to ensure that we can parse JSX-like syntax to allow templates
// to hold individual values, and then let the template engines to
// merge them into a single value.
//
func (node *HtmlNode) AddAttribute(key string, value string) {
	if len(node.Attributes) == 0 {
		node.Attributes = make([]*HtmlAttribute, 0)
	}

	node.Attributes = append(node.Attributes, &HtmlAttribute{
		Name:  key,
		Value: value,
	})
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
	if !node.ContainsAttributes() {
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
	if !node.ContainsAttributes() {
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
	if node.ContainsAttributes() {
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

//
// Find and return all attributes that have the given name.
//
// Returns a slice of `HtmlNode` if found, `nil` otherwise
//
func (node *HtmlNode) GetAttributes(key string) []*HtmlAttribute {
	if !node.ContainsAttributes() {
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
// Function removes all duplicate attributes from the node.
// The first available value is kept and other values are
// dropped. The function returns `true` if the attributes
// were modified (duplicates were removed), `false` otherwise.
//
func (node *HtmlNode) RemoveDuplicateAttributes() bool {
	if !node.ContainsAttributes() {
		return false
	}

	attributeMap := make(map[string]*HtmlAttribute, 0)
	for _, attr := range node.Attributes {
		attributeMap[attr.Name] = attr
	}

	if len(attributeMap) == len(node.Attributes) {
		return false
	}

	result := make([]*HtmlAttribute, 0)
	for _, value := range attributeMap {
		result = append(result, value)
	}

	node.Attributes = result
	return true
}

//
// Find and return the `HtmlAttribute` which has the given name and value
//
// Returns either a `HtmlAttribute` instance, `nil` otherwise
//
func (node *HtmlNode) GetAttributeWithValue(key string, value string) *HtmlAttribute {
	if !node.ContainsAttributes() {
		return nil
	}

	for _, attr := range node.Attributes {
		if attr.Name == key && attr.Value == value {
			return attr
		}
	}

	return nil
}
