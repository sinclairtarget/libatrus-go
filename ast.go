package atrus

/*
#include <stdlib.h>
#include <atrus.h>
*/
import "C"

import (
	"iter"
	"runtime"
	"unsafe"
)

type ASTNode struct {
	// Copied by value into Go struct, but contains pointers to externally
	// allocated memory regions that need to be freed manually.
	cNode C.struct_atrus_ast_node
}

// Wraps C AST node in a Go struct.
func NewASTNode(cNode *C.struct_atrus_ast_node) *ASTNode {
	n := &ASTNode{
		cNode: *cNode,
	}

	// Free manually allocated memory on cleanup
	runtime.SetFinalizer(n, func(n *ASTNode) {
		C.atrus_free(&n.cNode)
	})

	return n
}

// Returns a type name for the node.
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
		return "thematic_break"
	case C.ATRUS_NODE_TYPE_BREAK:
		return "break"
	case C.ATRUS_NODE_TYPE_EMPHASIS:
		return "emphasis"
	case C.ATRUS_NODE_TYPE_STRONG:
		return "strong"
	case C.ATRUS_NODE_TYPE_INLINE_CODE:
		return "inline_code"
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

// Returns an iterator through the node's children.
func (n ASTNode) Children() iter.Seq[*ASTNode] {
	seq := func(yield func(*ASTNode) bool) {
		payload := unsafe.Pointer(&n.cNode.payload)

		switch n.cNode.tag {
		case C.ATRUS_NODE_TYPE_ROOT:
			root := (*C.struct_atrus_ast_node_root)(payload)
			children := unsafe.Slice(root.children, root.n_children)

			for _, rawChild := range children {
				child := NewASTNode(rawChild)
				if !yield(child) {
					break
				}
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
				if !yield(child) {
					break
				}
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
	}

	return seq
}
