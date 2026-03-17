// Go bindings for libatrus.
package atrus

/*
#include <stdlib.h>
#include <atrus.h>
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

func Version() string {
	return C.GoString(C.atrus_version)
}

type ParseLevel C.atrus_parse_option_parse_level_t

const (
	ParseLevelBlock ParseLevel = C.ATRUS_BLOCK_PARSE_LEVEL
	ParseLevelPre              = C.ATRUS_PRE_PARSE_LEVEL
	ParseLevelPost             = C.ATRUS_POST_PARSE_LEVEL
)

// Options for parsing.
type ParseOpts struct {
	ParseLevel C.atrus_parse_option_parse_level_t
}

func Parse(md string, opts ParseOpts) (*ASTNode, error) {
	cMd := C.CString(md)
	defer C.free(unsafe.Pointer(cMd))

	parseOpts := C.struct_atrus_parse_opts{
		parse_level: opts.ParseLevel,
	}

	var out *C.struct_atrus_ast_node
	retcode := C.atrus_parse(cMd, &out, &parseOpts)
	if retcode != 0 {
		return nil, errors.New("parse failed")
	}

	node := NewASTNode(out)

	// Set finalizer on root node.
	// Root node is responsible for freeing the whole tree when it gets GC-ed.
	runtime.SetFinalizer(node, func(n *ASTNode) {
		C.atrus_free(n.cNode)
	})

	return node, nil
}

func RenderHTML(node *ASTNode) (string, error) {
	var out *C.char
	length := C.atrus_render_html(node.cNode, &out)
	if length < 0 {
		return "", errors.New("render html failed")
	}
	defer C.free(unsafe.Pointer(out))

	return C.GoStringN(out, length), nil
}

type WhitespaceOption C.atrus_json_option_whitespace_t

const (
	JSONMinified WhitespaceOption = C.ATRUS_JSON_MINIFIED
	JSONIndent2                   = C.ATRUS_JSON_INDENT_2
	JSONIndent4                   = C.ATRUS_JSON_INDENT_4
)

type JSONOpts struct {
	Whitespace C.atrus_json_option_whitespace_t
}

func RenderJSON(node *ASTNode, opts JSONOpts) (string, error) {
	var out *C.char
	renderOpts := C.struct_atrus_json_opts{
		whitespace: opts.Whitespace,
	}

	length := C.atrus_render_json(node.cNode, &out, &renderOpts)
	if length < 0 {
		return "", errors.New("render json failed")
	}
	defer C.free(unsafe.Pointer(out))

	return C.GoStringN(out, length), nil
}

// func LoadJSON(s string) (*ASTNode, error) {}
