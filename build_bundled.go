//go:build bundled_libatrus
package atrus

/*
#cgo CFLAGS: -I${SRCDIR}/vendored/libatrus/zig-out/include
#cgo LDFLAGS: ${SRCDIR}/vendored/libatrus/zig-out/lib/libatrus.a
*/
import "C"
