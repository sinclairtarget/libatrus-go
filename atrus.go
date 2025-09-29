package atrus

/*
#cgo pkg-config: libatrus

#include <stdlib.h>
#include <atrus.h>
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type ASTNode struct {
	ptr *C.atrus_ast_node
}

func (n *ASTNode) Free() {
	if n.ptr != nil {
		C.atrus_ast_free(n.ptr)
		n.ptr = nil
	}
}

func ParseAST(md string) (*ASTNode, error) {
	cMd := C.CString(md)
	defer C.free(unsafe.Pointer(cMd))

	var ptr *C.atrus_ast_node
	retcode := C.atrus_ast_parse(cMd, &ptr)
	if retcode != 0 {
		return nil, errors.New("parse failed")
	}

	node := &ASTNode{ ptr: ptr }
	runtime.SetFinalizer(node, (*ASTNode).Free)
	return node, nil
}

func RenderJSON(node *ASTNode) (string, error) {
	var out *C.char
	length := C.atrus_render_json(node.ptr, &out)
	if length < 0 {
		return "", errors.New("render json failed")
	}
	defer C.free(unsafe.Pointer(out))

	return C.GoStringN(out, length), nil
}

func RenderHTML(node *ASTNode) (string, error) {
	var out *C.char
	length := C.atrus_render_html(node.ptr, &out)
	if length < 0 {
		return "", errors.New("render html failed")
	}
	defer C.free(unsafe.Pointer(out))

	return C.GoStringN(out, length), nil
}
