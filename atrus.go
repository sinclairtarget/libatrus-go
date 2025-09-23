package atrus

/*
#cgo pkg-config: ${SRCDIR}/vendored/zig-out/share/pkgconfig/libatrus.pc

#include <atrus.h>
*/
import "C"

import (
	"errors"
)

type ASTNode struct {}

func ParseAST(md string) (*ASTNode, error) {
	var node *ASTNode
	retcode := C.atrus_ast_parse(md, &node)
	if retcode != 0 {
		return nil, errors.New("parse failed")
	}
	return node, nil
}

func FreeAST(node *ASTNode) {
}

func RenderJSON(node *ASTNode) (string, error) {
	return "", nil
}
