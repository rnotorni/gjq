package gjq

/*
#cgo LDFLAGS: -ljq
#cgo linux LDFLAGS: -lm
#include <stdlib.h>
#include <jq.h>
#include <jv.h>

typedef struct {
  jv msg;
} gjq_error_data;

static void gjq_error_cb(void *data, jv msg) {
   ((gjq_error_data *)data)->msg = jq_format_error(msg);
}

static void gjq_set_error_cb(jq_state *jq, void *data) {
   jq_set_error_cb(jq, gjq_error_cb, data);
}
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

// GJQ .
type GJQ struct {
	jqState     *C.jq_state
	script      string
	errorData   C.gjq_error_data
	defferFuncs []func(g *GJQ) error
}

// NewGJQ .
func NewGJQ(script string) (*GJQ, error) {
	ret := &GJQ{
		jqState: C.jq_init(),
		defferFuncs: []func(g *GJQ) error{
			freeErrorData,
			treadown,
		},
	}

	C.gjq_set_error_cb(ret.jqState, unsafe.Pointer(&ret.errorData))
	err := ret.compile(script)
	if err != nil {
		ret.Close()
		return nil, err
	}
	return ret, nil
}

func (g *GJQ) Error() error {
	return fmt.Errorf("in jq error: %v", C.GoString(C.jv_string_value(g.errorData.msg)))
}

func (g *GJQ) cString(str string) *C.char {
	ret := C.CString(str)
	g.defferFuncs = append(g.defferFuncs, func(*GJQ) error {
		C.free(unsafe.Pointer(ret))
		return nil
	})
	return ret
}

func (g *GJQ) compile(script string) error {
	if result := C.jq_compile(g.jqState, g.cString(script)); result == 0 {
		return g.Error()
	}
	return nil
}

func freeErrorData(g *GJQ) error {
	C.jv_free(g.errorData.msg)
	return nil
}

func treadown(g *GJQ) error {
	C.jq_teardown(&g.jqState)
	return nil
}

// Close with free
func (g *GJQ) Close() {
	for _, f := range g.defferFuncs {
		f(g)
	}
}

// Run .
// TODO: input to any and use buffer
func (g *GJQ) Run(json string) (string, error) {
	inputJv := C.jv_parse(g.cString(json))
	if C.jv_is_valid(inputJv) == 0 {
		// TODO:
		return "", fmt.Errorf("can'input: parse error")
	}
	C.jq_start(g.jqState, inputJv, C.int(0))
	defer C.jq_start(g.jqState, C.jv_null(), C.int(0))

	out := make([]string, 0)
	for tmp := C.jq_next(g.jqState); C.jv_is_valid(tmp) == 1; tmp = C.jq_next(g.jqState) {
		out = append(out, JvDumpString(tmp))
		C.jv_free(tmp)
	}
	return strings.Join(out, "\n"), nil
}

// JvDumpString .
func JvDumpString(str C.jv) string {
	dumpedjv := C.jv_dump_string(C.jv_copy(str), C.int(0))

	defer C.jv_free(dumpedjv)
	return C.GoString(C.jv_string_value(dumpedjv))
}
