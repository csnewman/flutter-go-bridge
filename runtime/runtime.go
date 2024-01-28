package runtime

/*
#include "dart_api_dl.h"
#include <stdlib.h>

extern bool InternalPostBlob(Dart_Port port_id, intptr_t len, void* data, uintptr_t peer);
extern void* CloneUInt64Array(intptr_t len, void* data);
*/
import "C"

import (
	"errors"
	"unsafe"
)

var ErrDartIncompatible = errors.New("DartVM major version does not match")

//export fgbinternal_init
func fgbinternal_init(p unsafe.Pointer) unsafe.Pointer {
	err := InitializeApi(p)

	var cerr unsafe.Pointer
	if err != nil {
		cerr = unsafe.Pointer(C.CString(err.Error()))
	}

	return cerr
}

func InitializeApi(p unsafe.Pointer) error {
	v := C.Dart_InitializeApiDL(p)

	if v == -1 {
		return ErrDartIncompatible
	}

	return nil
}

//export fgbinternal_alloc
func fgbinternal_alloc(size C.intptr_t) unsafe.Pointer {
	return C.malloc((C.uintptr_t)(size))
}

//export fgbinternal_free
func fgbinternal_free(ptr unsafe.Pointer) {
	C.free(ptr)
}

type sendData struct {
	ptr      unsafe.Pointer
	callback func()
}

func Send(port int64, data []uint64, callback func()) bool {
	ptr := C.CloneUInt64Array(C.intptr_t(len(data)), unsafe.Pointer(&data[0]))
	h := Pin(&sendData{
		ptr:      ptr,
		callback: callback,
	})

	sent := bool(C.InternalPostBlob(C.Dart_Port_DL(port), C.intptr_t(len(data)), ptr, (C.uintptr_t)(h)))
	if !sent {
		C.free(ptr)
		FreePin(h)
	}

	return sent
}

//export ifgb_callback
func ifgb_callback(p uintptr) {
	data := GetPin[*sendData](p)

	if data.callback != nil {
		data.callback()
	}

	C.free(data.ptr)
	FreePin(p)
}
