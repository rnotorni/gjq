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
