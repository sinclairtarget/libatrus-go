// Go bindings for libatrus.
package atrus

/*
#include <stdlib.h>
#include <atrus.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

func Version() string {
	return C.GoString(C.atrus_version())
}

// Returns an error if the linked version of libatrus is older than the version
// we were compiled against.
func CheckLinkedVersion() error {
	isOldEnough := C.atrus_version_at_least(
		C.ATRUS_MAJOR_VERSION,
		C.ATRUS_MINOR_VERSION,
		C.ATRUS_PATCH_VERSION,
	)

	if !isOldEnough {
		compileVersion := fmt.Sprintf(
			"%d.%d.%d",
			C.ATRUS_MAJOR_VERSION,
			C.ATRUS_MINOR_VERSION,
			C.ATRUS_PATCH_VERSION,
		)
		linkedVersion := Version()
		return fmt.Errorf(
			"linked version of libatrus (%s) is older than compile-time " +
			"version (%s)",
			linkedVersion,
			compileVersion,
		)
	}

	return nil
}

type ParseLevel C.atrus_parse_option_parse_level_t

const (
	ParseLevelBlock ParseLevel = C.ATRUS_PARSE_LEVEL_BLOCK
	ParseLevelRaw              = C.ATRUS_PARSE_LEVEL_RAW
	ParseLevelPre              = C.ATRUS_PARSE_LEVEL_PRE
	ParseLevelPost             = C.ATRUS_PARSE_LEVEL_POST
)

func Parse(md string, parseLevel ParseLevel) (*ASTNode, error) {
	cMd := C.CString(md)
	defer C.free(unsafe.Pointer(cMd))

	var cNode *C.struct_atrus_node
	retcode := C.atrus_parse(
		cMd,
		&cNode,
		// Not obvious to me why this cast is necessary. We've just said above
		// that ParseLevel is an alias for the underlying C type! But go will
		// be unhappy without this cast.
		(C.atrus_parse_option_parse_level_t)(parseLevel),
	)
	if retcode != 0 {
		return nil, errors.New("parse failed")
	}

	node := newASTNode(cNode, nil)
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

func RenderJSON(node *ASTNode, whitespace WhitespaceOption) (string, error) {
	var out *C.char
	length := C.atrus_render_json(
		node.cNode,
		&out,
		// Not obvious to me why this cast is necessary. We've just said above
		// that WhitespaceOption is an alias for the underlying C type! But go
		// will be unhappy without this cast.
		(C.atrus_json_option_whitespace_t)(whitespace),
	)
	if length < 0 {
		return "", errors.New("render json failed")
	}
	defer C.free(unsafe.Pointer(out))

	return C.GoStringN(out, length), nil
}

// func LoadJSON(s string) (*ASTNode, error) {}
