package runtime

/*

#include "dart_api_dl.h"
#include <string.h>
#include <stdlib.h>

extern void ifgb_callback(void* ptr);

void InternalBlobCallback(void* isolate_callback_data, void* peer) {
    ifgb_callback(peer);
}

bool InternalPostBlob(Dart_Port port_id, intptr_t len, void* data, void* peer) {
    Dart_CObject obj;
    obj.type = Dart_CObject_kUnmodifiableExternalTypedData;
    obj.value.as_external_typed_data.type = Dart_TypedData_kUint64;
    obj.value.as_external_typed_data.length = len;
    obj.value.as_external_typed_data.data = data;
    obj.value.as_external_typed_data.peer = peer;
    obj.value.as_external_typed_data.callback = &InternalBlobCallback;
    return Dart_PostCObject_DL(port_id, &obj);
}

void* CloneUInt64Array(intptr_t len, void* data) {
    void* ptr = malloc(len * sizeof(uint64_t));
    memcpy(ptr, data, len * sizeof(uint64_t));
    return ptr;
}

*/
import "C"

import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

var ErrDartIncompatible = errors.New("DartVM major version does not match")

func InitializeApi(p unsafe.Pointer) error {
	v := C.Dart_InitializeApiDL(p)

	if v == -1 {
		return ErrDartIncompatible
	}

	return nil
}

type sendData struct {
	ptr      unsafe.Pointer
	callback func()
}

func Send(port int64, data []uint64, callback func()) bool {
	ptr := C.CloneUInt64Array(C.intptr_t(len(data)), unsafe.Pointer(&data[0]))
	h := cgo.NewHandle(&sendData{
		ptr:      ptr,
		callback: callback,
	})

	sent := bool(C.InternalPostBlob(C.Dart_Port_DL(port), C.intptr_t(len(data)), ptr, unsafe.Pointer(h)))
	if !sent {
		C.free(ptr)
		h.Delete()
	}

	return sent
}
