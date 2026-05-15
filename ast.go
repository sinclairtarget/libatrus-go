package atrus

/*
#include <stdlib.h>
#include <atrus.h>
*/
import "C"

import (
	"fmt"
)

type ASTNode struct {
	cNode *C.struct_atrus_node
}

// Wraps C AST node in a Go struct.
func newASTNode(cNode *C.struct_atrus_node) *ASTNode {
	return &ASTNode{
		cNode: cNode,
	}
}

// Returns a type name for the node.
//
// These names match the type names given in the MyST spec node index:
// https://mystmd.org/spec/myst-schema
func (n ASTNode) Type() string {
	name := C.atrus_node_type_name(n.cNode)
	return C.GoString(name)
}

// Returns a slice containing the node's children.
//
// If the node has no children, or is a leaf node, just returns an empty slice.
func (n ASTNode) Children() []*ASTNode {
	outSlice := []*ASTNode{}

	nChildren := C.atrus_node_num_children(n.cNode)
	for i := range nChildren {
		child := C.atrus_node_child(n.cNode, i)
		outSlice = append(outSlice, newASTNode(child))
	}

	return outSlice
}

// ----------------------------------------------------------------------------
// Remaining types and functions all handle recovering the various union
// payloads from a node.
// ----------------------------------------------------------------------------

type Heading struct {
	Depth C.ushort
}

func (n ASTNode) Heading() Heading {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_HEADING {
		msg := formatTypePanicMsg("Heading()", n)
		panic(msg) // called Heading() on an Atrus AST node of type X
	}

	depth := C.atrus_node_heading_depth(n.cNode)
	return Heading{
		Depth: depth,
	}
}

type Text struct {
	Value string
}

func (n ASTNode) Text() Text {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_TEXT {
		msg := formatTypePanicMsg("Text()", n)
		panic(msg) // called Text() on an Atrus AST node of type X
	}

	value := C.atrus_node_text_value(n.cNode)
	return Text{
		Value: C.GoString(value),
	}
}

type Code struct {
	Value string
	Lang  string
	ShowLineNumbers bool
}

func (n ASTNode) Code() Code {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_CODE {
		msg := formatTypePanicMsg("Code()", n)
		panic(msg) // called Code() on an Atrus AST node of type X
	}

	value := C.atrus_node_code_value(n.cNode)
	lang := C.atrus_node_code_lang(n.cNode)
	showLineNumbers := C.atrus_node_code_show_line_numbers(n.cNode)
	return Code{
		Value: C.GoString(value),
		Lang:  C.GoString(lang),
		ShowLineNumbers: bool(showLineNumbers),
	}
}

type Link struct {
	URL   string
	Title string
}

func (n ASTNode) Link() Link {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_LINK {
		msg := formatTypePanicMsg("Link()", n)
		panic(msg) // called Link() on an Atrus AST node of type X
	}

	url := C.atrus_node_link_url(n.cNode)
	title := C.atrus_node_link_title(n.cNode)
	return Link{
		URL:   C.GoString(url),
		Title: C.GoString(title),
	}
}

type LinkDefinition struct {
	URL string
	Title string
	Label string
}

func (n ASTNode) LinkDefinition() LinkDefinition {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_DEFINITION {
		msg := formatTypePanicMsg("LinkDefinition()", n)
		panic(msg) // called LinkDefinition() on an Atrus AST node of type X
	}

	url := C.atrus_node_definition_url(n.cNode)
	title := C.atrus_node_definition_title(n.cNode)
	label := C.atrus_node_definition_label(n.cNode)
	return LinkDefinition{
		URL:   C.GoString(url),
		Title: C.GoString(title),
		Label: C.GoString(label),
	}
}

type Image struct {
	URL   string
	Title string
	Alt   string
}

func (n ASTNode) Image() Image {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_IMAGE {
		msg := formatTypePanicMsg("Image()", n)
		panic(msg) // called Image() on an Atrus AST node of type X
	}

	url := C.atrus_node_image_url(n.cNode)
	title := C.atrus_node_image_title(n.cNode)
	alt := C.atrus_node_image_alt(n.cNode)
	return Image{
		URL:   C.GoString(url),
		Title: C.GoString(title),
		Alt:   C.GoString(alt),
	}
}

type Container struct {
	Kind string
}

func (n ASTNode) Container() Container {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_CONTAINER {
		msg := formatTypePanicMsg("Container()", n)
		panic(msg) // called Container() on an Atrus AST node of type X
	}

	kind := C.atrus_node_container_kind(n.cNode)
	return Container{
		Kind: C.GoString(kind),
	}
}

type MySTRole struct {
	Name string
	Value string
}

func (n ASTNode) MySTRole() MySTRole {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_MYST_ROLE {
		msg := formatTypePanicMsg("MySTRole()", n)
		panic(msg) // called MySTRole() on an Atrus AST node of type X
	}

	name := C.atrus_node_myst_role_name(n.cNode)
	value := C.atrus_node_myst_role_value(n.cNode)
	return MySTRole{
		Name:  C.GoString(name),
		Value: C.GoString(value),
	}
}

type MySTRoleError struct {
	Value string
}

func (n ASTNode) MySTRoleError() MySTRoleError {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_MYST_ROLE_ERROR {
		msg := formatTypePanicMsg("MySTRoleError()", n)
		panic(msg) // called MySTRoleError() on an Atrus AST node of type X
	}

	value := C.atrus_node_myst_role_error_value(n.cNode)
	return MySTRoleError{
		Value: C.GoString(value),
	}
}

type Abbreviation struct {
	Title string
}

func (n ASTNode) Abbreviation() Abbreviation {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_ABBREVIATION {
		msg := formatTypePanicMsg("Abbreviation()", n)
		panic(msg) // called Abbreviation() on an Atrus AST node of type X
	}

	title := C.atrus_node_abbreviation_title(n.cNode)
	return Abbreviation{
		Title: C.GoString(title),
	}
}

type MySTDirective struct {
	Name string
	Args string
	Value string
}

func (n ASTNode) MySTDirective() MySTDirective {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_MYST_DIRECTIVE {
		msg := formatTypePanicMsg("MySTDirective()", n)
		panic(msg) // called MySTDirective() on an Atrus AST node of type X
	}

	name := C.atrus_node_myst_directive_name(n.cNode)
	args := C.atrus_node_myst_directive_args(n.cNode)
	value := C.atrus_node_myst_directive_value(n.cNode)
	return MySTDirective{
		Name:  C.GoString(name),
		Args:  C.GoString(args),
		Value: C.GoString(value),
	}
}

type MySTDirectiveError struct {
	Message string
}

func (n ASTNode) MySTDirectiveError() MySTDirectiveError {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_MYST_DIRECTIVE_ERROR {
		msg := formatTypePanicMsg("MySTDirectiveError()", n)
		panic(msg) // called MySTDirectiveError() on an Atrus AST node of type X
	}

	message := C.atrus_node_myst_directive_error_message(n.cNode)
	return MySTDirectiveError{
		Message: C.GoString(message),
	}
}

type Admonition struct {
	Kind string
}

func (n ASTNode) Admonition() Admonition {
	if C.atrus_node_type(n.cNode) != C.ATRUS_NODE_TYPE_ADMONITION {
		msg := formatTypePanicMsg("Admonition()", n)
		panic(msg) // called Admonition() on an Atrus AST node of type X
	}

	kind := C.atrus_node_admonition_kind(n.cNode)
	return Admonition{
		Kind: C.GoString(kind),
	}
}

func formatTypePanicMsg(fnName string, node ASTNode) string {
	t := node.Type()
	return fmt.Sprintf(
		"called %s on an Atrus AST node of type \"%s\"",
		fnName,
		t,
	)
}
