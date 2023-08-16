package runtime

import (
	"runtime/cgo"
	"unsafe"
)

/*
#include <stdlib.h>
*/
import "C"

//export ifgb_callback
func ifgb_callback(p unsafe.Pointer) {
	h := cgo.Handle(p)

	data := (h.Value()).(*sendData)

	if data.callback != nil {
		data.callback()
	}

	C.free(data.ptr)

	h.Delete()
}
