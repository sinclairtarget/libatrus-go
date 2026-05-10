package atrus

/*
#include <stdlib.h>
#include <atrus.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type ASTNode struct {
	cNode *C.struct_atrus_node
	cOpaqueNode *C.struct_atrus_node_opaque
}

// Wraps C AST node in a Go struct.
func NewASTNode(cNode *C.struct_atrus_node) *ASTNode {
	return &ASTNode{
		cNode: cNode,
		cOpaqueNode: nil,
	}
}

func (n ASTNode) Frozen() bool {
	return n.cNode == nil
}

func (n *ASTNode) setOpaque(opaqueNode *C.struct_atrus_node_opaque) {
	n.cOpaqueNode = opaqueNode
	n.cNode = nil
}

func (n ASTNode) checkFrozen() {
	if n.Frozen() {
		panic("invalid method called on frozen AST node")
	}
}

// Returns a type name for the node.
//
// These names match the type names given in the MyST spec node index:
// https://mystmd.org/spec/myst-schema
func (n ASTNode) Type() string {
	name := C.atrus_name(n.cNode.tag)
	return C.GoString(name)
}

// Returns a slice containing the node's children.
//
// If the node has no children, or is a leaf node, just returns an empty slice.
func (n ASTNode) Children() []*ASTNode {
	n.checkFrozen()

	outSlice := []*ASTNode{}

	payload := unsafe.Pointer(&n.cNode.payload)
	switch n.cNode.tag {
	case C.ATRUS_NODE_TYPE_ROOT:
		root := (*C.struct_atrus_node_root)(payload)
		children := unsafe.Slice(root.children, root.children_len)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_BLOCK,
		C.ATRUS_NODE_TYPE_PARAGRAPH,
		C.ATRUS_NODE_TYPE_EMPHASIS,
		C.ATRUS_NODE_TYPE_STRONG,
		C.ATRUS_NODE_TYPE_BLOCKQUOTE,
		C.ATRUS_NODE_TYPE_CAPTION,
		C.ATRUS_NODE_TYPE_SUBSCRIPT,
		C.ATRUS_NODE_TYPE_SUPERSCRIPT,
		C.ATRUS_NODE_TYPE_ADMONITION_TITLE:
		wrapper := (*C.struct_atrus_node_wrapper)(payload)
		children := unsafe.Slice(wrapper.children, wrapper.children_len)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_HEADING:
		heading := (*C.struct_atrus_node_heading)(payload)
		children := unsafe.Slice(heading.children, heading.children_len)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_LINK:
		link := (*C.struct_atrus_node_link)(payload)
		children := unsafe.Slice(link.children, link.children_len)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_CONTAINER:
		container := (*C.struct_atrus_node_container)(payload)
		children := unsafe.Slice(container.children, container.children_len)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_MYST_ROLE:
		myst_role := (*C.struct_atrus_node_myst_role)(payload)
		children := unsafe.Slice(myst_role.children, myst_role.children_len)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_ABBREVIATION:
		abbreviation := (*C.struct_atrus_node_abbreviation)(payload)
		children := unsafe.Slice(
			abbreviation.children,
			abbreviation.children_len,
		)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_MYST_DIRECTIVE:
		myst_directive := (*C.struct_atrus_node_myst_directive)(payload)
		children := unsafe.Slice(
			myst_directive.children,
			myst_directive.children_len,
		)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_MYST_DIRECTIVE_ERROR:
		myst_directive_error := (*C.struct_atrus_node_myst_directive_error)(
			payload,
		)
		children := unsafe.Slice(
			myst_directive_error.children,
			myst_directive_error.children_len,
		)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_ADMONITION:
		admonition := (*C.struct_atrus_node_admonition)(payload)
		children := unsafe.Slice(admonition.children, admonition.children_len)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_TEXT,
		C.ATRUS_NODE_TYPE_CODE,
		C.ATRUS_NODE_TYPE_THEMATIC_BREAK,
		C.ATRUS_NODE_TYPE_BREAK,
		C.ATRUS_NODE_TYPE_INLINE_CODE,
		C.ATRUS_NODE_TYPE_DEFINITION,
		C.ATRUS_NODE_TYPE_IMAGE,
		C.ATRUS_NODE_TYPE_HTML,
		C.ATRUS_NODE_TYPE_MYST_ROLE_ERROR:
		// Childless nodes
		break
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
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_HEADING {
		msg := formatTypePanicMsg("Heading()", n)
		panic(msg) // called Heading() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cHeading := (*C.struct_atrus_node_heading)(payload)
	return Heading{
		Depth: cHeading.depth,
	}
}

type Text struct {
	Value string
}

func (n ASTNode) Text() Text {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_TEXT {
		msg := formatTypePanicMsg("Text()", n)
		panic(msg) // called Text() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cText := (*C.struct_atrus_node_text)(payload)
	return Text{
		Value: C.GoString(cText.value),
	}
}

type Code struct {
	Value string
	Lang  string
	ShowLineNumbers bool
}

func (n ASTNode) Code() Code {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_CODE {
		msg := formatTypePanicMsg("Code()", n)
		panic(msg) // called Code() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cCode := (*C.struct_atrus_node_code)(payload)
	return Code{
		Value: C.GoString(cCode.value),
		Lang:  C.GoString(cCode.lang),
		ShowLineNumbers:  bool(cCode.show_line_numbers),
	}
}

type Link struct {
	URL   string
	Title string
}

func (n ASTNode) Link() Link {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_LINK {
		msg := formatTypePanicMsg("Link()", n)
		panic(msg) // called Link() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cLink := (*C.struct_atrus_node_link)(payload)
	return Link{
		URL:   C.GoString(cLink.url),
		Title: C.GoString(cLink.title),
	}
}

type LinkDefinition struct {
	URL string
	Title string
	Label string
}

func (n ASTNode) LinkDefinition() LinkDefinition {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_DEFINITION {
		msg := formatTypePanicMsg("LinkDefinition()", n)
		panic(msg) // called LinkDefinition() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cLinkDefinition := (*C.struct_atrus_node_link_definition)(payload)
	return LinkDefinition{
		URL:   C.GoString(cLinkDefinition.url),
		Title: C.GoString(cLinkDefinition.title),
		Label: C.GoString(cLinkDefinition.label),
	}
}

type Image struct {
	URL   string
	Title string
	Alt   string
}

func (n ASTNode) Image() Image {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_IMAGE {
		msg := formatTypePanicMsg("Image()", n)
		panic(msg) // called Image() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cImage := (*C.struct_atrus_node_image)(payload)
	return Image{
		URL:   C.GoString(cImage.url),
		Title: C.GoString(cImage.title),
		Alt:   C.GoString(cImage.alt),
	}
}

type Container struct {
	Kind string
}

func (n ASTNode) Container() Container {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_CONTAINER {
		msg := formatTypePanicMsg("Container()", n)
		panic(msg) // called Container() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cContainer := (*C.struct_atrus_node_container)(payload)
	return Container{
		Kind: C.GoString(cContainer.kind),
	}
}

type MySTRole struct {
	Name string
	Value string
}

func (n ASTNode) MySTRole() MySTRole {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_MYST_ROLE {
		msg := formatTypePanicMsg("MySTRole()", n)
		panic(msg) // called MySTRole() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cMySTRole := (*C.struct_atrus_node_myst_role)(payload)
	return MySTRole{
		Name:  C.GoString(cMySTRole.name),
		Value: C.GoString(cMySTRole.value),
	}
}

type MySTRoleError struct {
	Value string
}

func (n ASTNode) MySTRoleError() MySTRoleError {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_MYST_ROLE_ERROR {
		msg := formatTypePanicMsg("MySTRoleError()", n)
		panic(msg) // called MySTRoleError() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cMySTRoleError := (*C.struct_atrus_node_myst_role_error)(payload)
	return MySTRoleError{
		Value: C.GoString(cMySTRoleError.value),
	}
}

type Abbreviation struct {
	Title string
}

func (n ASTNode) Abbreviation() Abbreviation {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_ABBREVIATION {
		msg := formatTypePanicMsg("Abbreviation()", n)
		panic(msg) // called Abbreviation() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cAbbreviation := (*C.struct_atrus_node_abbreviation)(payload)
	return Abbreviation{
		Title: C.GoString(cAbbreviation.title),
	}
}

type MySTDirective struct {
	Name string
	Args string
	Value string
}

func (n ASTNode) MySTDirective() MySTDirective {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_MYST_DIRECTIVE {
		msg := formatTypePanicMsg("MySTDirective()", n)
		panic(msg) // called MySTDirective() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cMySTDirective := (*C.struct_atrus_node_myst_directive)(payload)
	return MySTDirective{
		Name:  C.GoString(cMySTDirective.name),
		Args:  C.GoString(cMySTDirective.args),
		Value: C.GoString(cMySTDirective.value),
	}
}

type MySTDirectiveError struct {
	Message string
}

func (n ASTNode) MySTDirectiveError() MySTDirectiveError {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_MYST_DIRECTIVE_ERROR {
		msg := formatTypePanicMsg("MySTDirectiveError()", n)
		panic(msg) // called MySTDirectiveError() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cMySTDirectiveError := (*C.struct_atrus_node_myst_directive_error)(payload)
	return MySTDirectiveError{
		Message: C.GoString(cMySTDirectiveError.message),
	}
}

type Admonition struct {
	Kind string
}

func (n ASTNode) Admonition() Admonition {
	n.checkFrozen()

	if n.cNode.tag != C.ATRUS_NODE_TYPE_ADMONITION {
		msg := formatTypePanicMsg("Admonition()", n)
		panic(msg) // called Admonition() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cAdmonition := (*C.struct_atrus_node_admonition)(payload)
	return Admonition{
		Kind: C.GoString(cAdmonition.kind),
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
