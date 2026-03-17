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
	// Copied by value into Go struct, but contains pointers to externally
	// allocated memory regions that need to be freed manually.
	cNode C.struct_atrus_ast_node
}

// Wraps C AST node in a Go struct.
func NewASTNode(cNode *C.struct_atrus_ast_node) *ASTNode {
	return &ASTNode{
		cNode: *cNode,
	}
}

// Returns a type name for the node.
//
// These names match the type names given in the MyST spec node index:
// https://mystmd.org/spec/myst-schema
func (n ASTNode) Type() string {
	switch n.cNode.tag {
	case C.ATRUS_NODE_TYPE_ROOT:
		return "root"
	case C.ATRUS_NODE_TYPE_BLOCK:
		return "block"
	case C.ATRUS_NODE_TYPE_HEADING:
		return "heading"
	case C.ATRUS_NODE_TYPE_PARAGRAPH:
		return "paragraph"
	case C.ATRUS_NODE_TYPE_TEXT:
		return "text"
	case C.ATRUS_NODE_TYPE_CODE:
		return "code"
	case C.ATRUS_NODE_TYPE_THEMATIC_BREAK:
		return "thematicBreak"
	case C.ATRUS_NODE_TYPE_BREAK:
		return "break"
	case C.ATRUS_NODE_TYPE_EMPHASIS:
		return "emphasis"
	case C.ATRUS_NODE_TYPE_STRONG:
		return "strong"
	case C.ATRUS_NODE_TYPE_INLINE_CODE:
		return "inlineCode"
	case C.ATRUS_NODE_TYPE_LINK:
		return "link"
	case C.ATRUS_NODE_TYPE_DEFINITION:
		return "definition"
	case C.ATRUS_NODE_TYPE_IMAGE:
		return "image"
	case C.ATRUS_NODE_TYPE_BLOCKQUOTE:
		return "blockquote"
	default:
		panic("unknown Atrus node type")
	}
}

// Returns a slice containing the node's children.
//
// If the node has no children, or is a leaf node, just returns an empty slice.
func (n ASTNode) Children() []*ASTNode {
	outSlice := []*ASTNode{}

	payload := unsafe.Pointer(&n.cNode.payload)
	switch n.cNode.tag {
	case C.ATRUS_NODE_TYPE_ROOT:
		root := (*C.struct_atrus_ast_node_root)(payload)
		children := unsafe.Slice(root.children, root.n_children)

		for _, rawChild := range children {
			child := NewASTNode(rawChild)
			outSlice = append(outSlice, child)
		}
	case C.ATRUS_NODE_TYPE_BLOCK,
		C.ATRUS_NODE_TYPE_HEADING,
		C.ATRUS_NODE_TYPE_PARAGRAPH,
		C.ATRUS_NODE_TYPE_EMPHASIS,
		C.ATRUS_NODE_TYPE_STRONG,
		C.ATRUS_NODE_TYPE_BLOCKQUOTE:
		container := (*C.struct_atrus_ast_node_container)(payload)
		children := unsafe.Slice(container.children, container.n_children)

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
		C.ATRUS_NODE_TYPE_IMAGE:
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
	if n.cNode.tag != C.ATRUS_NODE_TYPE_HEADING {
		msg := formatTypePanicMsg("Heading()", n)
		panic(msg) // called Heading() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cHeading := (*C.struct_atrus_ast_node_heading)(payload)
	return Heading{
		Depth: cHeading.depth,
	}
}

type Text struct {
	Value string
}

func (n ASTNode) Text() Text {
	if n.cNode.tag != C.ATRUS_NODE_TYPE_TEXT {
		msg := formatTypePanicMsg("Text()", n)
		panic(msg) // called Text() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cText := (*C.struct_atrus_ast_node_text)(payload)
	return Text{
		Value: C.GoString(cText.value),
	}
}

type Code struct {
	Value string
	Lang  string
}

func (n ASTNode) Code() Code {
	if n.cNode.tag != C.ATRUS_NODE_TYPE_CODE {
		msg := formatTypePanicMsg("Code()", n)
		panic(msg) // called Code() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cCode := (*C.struct_atrus_ast_node_code)(payload)
	return Code{
		Value: C.GoString(cCode.value),
		Lang:  C.GoString(cCode.lang),
	}
}

type Link struct {
	URL   string
	Title string
}

func (n ASTNode) Link() Link {
	if n.cNode.tag != C.ATRUS_NODE_TYPE_LINK {
		msg := formatTypePanicMsg("Link()", n)
		panic(msg) // called Link() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cLink := (*C.struct_atrus_ast_node_link)(payload)
	return Link{
		URL:   C.GoString(cLink.url),
		Title: C.GoString(cLink.title),
	}
}

type Image struct {
	URL   string
	Title string
	Alt   string
}

func (n ASTNode) Image() Image {
	if n.cNode.tag != C.ATRUS_NODE_TYPE_IMAGE {
		msg := formatTypePanicMsg("Image()", n)
		panic(msg) // called Image() on an Atrus AST node of type X
	}

	payload := unsafe.Pointer(&n.cNode.payload)
	cImage := (*C.struct_atrus_ast_node_image)(payload)
	return Image{
		URL:   C.GoString(cImage.url),
		Title: C.GoString(cImage.title),
		Alt:   C.GoString(cImage.alt),
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
