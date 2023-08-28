// This code was generated .
package main

import (
	"errors"
    "fmt"
	"sync"
	"sync/atomic"
    "unsafe"

	orig "flutter-go-bridge/example"
	"flutter-go-bridge/runtime"
)

/*
#include <stdlib.h>
#include <stdint.h>

typedef struct {
    int a;
    void* d;
} fgb_vt_inner;

typedef struct {
    int v_1;
    int v_2;
    fgb_vt_inner i_1;
} fgb_vt_some_val;

typedef struct {
    void* err;
} fgb_ret_example;

typedef struct {
    fgb_vt_some_val res;
    void* err;
} fgb_ret_other;

typedef struct {
    int res;
    void* err;
} fgb_ret_call_me;
*/
import "C"

var (
	handles   = sync.Map{}
	handleIdx uint64
	ErrDart   = errors.New("dart")
)

// Required by cgo
func main() {}

//export fgb_internal_init
func fgb_internal_init(p unsafe.Pointer) unsafe.Pointer {
	err := runtime.InitializeApi(p)

    var cerr unsafe.Pointer
    if err != nil {
        cerr = unsafe.Pointer(C.CString(err.Error()))
    }

	return cerr
}

func mapToString(from unsafe.Pointer) string {
	res := C.GoString((*C.char)(from))
	C.free(from)
	return res
}

func mapFromString(from string) unsafe.Pointer {
	return unsafe.Pointer(C.CString(from))
}

func mapToError(from unsafe.Pointer) error {
	res := C.GoString((*C.char)(from))
	C.free(from)
	return fmt.Errorf("%w: %v", ErrDart, res)
}

func mapFromError(from error) unsafe.Pointer {
	return unsafe.Pointer(C.CString(from.Error()))
}

//export fgbempty_inner
func fgbempty_inner() (res C.fgb_vt_inner) {
    return
}

func mapToInner(from C.fgb_vt_inner) (res orig.Inner) {
	res.A = (int)(from.a)
	res.D = mapToString(from.d)
	return
}

func mapFromInner(from orig.Inner) (res C.fgb_vt_inner) {
	res.a = (C.int)(from.A)
	res.d = mapFromString(from.D)
	return
}

//export fgbempty_some_val
func fgbempty_some_val() (res C.fgb_vt_some_val) {
    return
}

func mapToSomeVal(from C.fgb_vt_some_val) (res orig.SomeVal) {
	res.V1 = (int)(from.v_1)
	res.V2 = (int)(from.v_2)
	res.I1 = mapToInner(from.i_1)
	return
}

func mapFromSomeVal(from orig.SomeVal) (res C.fgb_vt_some_val) {
	res.v_1 = (C.int)(from.V1)
	res.v_2 = (C.int)(from.V2)
	res.i_1 = mapFromInner(from.I1)
	return
}

//export fgb_example
func fgb_example(v C.fgb_vt_some_val) (resw C.fgb_ret_example) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		resw = C.fgb_ret_example {
			err: unsafe.Pointer(C.CString(fmt.Sprintf("panic: %v", r))),
		}
	}()
	
	vGo := mapToSomeVal(v)
	gerr := orig.Example(vGo)
    if gerr != nil {
		return C.fgb_ret_example {
			err: unsafe.Pointer(C.CString(gerr.Error())),
		}
    }
    

    return C.fgb_ret_example {
    }
}

//export fgbasync_example
func fgbasync_example(v C.fgb_vt_some_val, fgbPort int64) {
	go func() {
		h := atomic.AddUint64(&handleIdx, 1)
		if h == 0 {
			panic("ran out of handle space")
		}

		handles.Store(h, fgb_example(v))

		sent := runtime.Send(fgbPort, []uint64{h}, func() {
			handles.LoadAndDelete(h)
		})
        if !sent {
            handles.LoadAndDelete(h)
        }
	}()
}

//export fgbasyncres_example
func fgbasyncres_example(h uint64) C.fgb_ret_example {
	v, ok := handles.LoadAndDelete(h)
	if !ok {
		return C.fgb_ret_example{
			err: unsafe.Pointer(C.CString("result handle is not valid")),
		}
	}

	return (v).(C.fgb_ret_example)
}

//export fgb_other
func fgb_other() (resw C.fgb_ret_other) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		resw = C.fgb_ret_other {
			err: unsafe.Pointer(C.CString(fmt.Sprintf("panic: %v", r))),
		}
	}()
	
	gres := orig.Other()
    
	cres := mapFromSomeVal(gres)

    return C.fgb_ret_other {
        res: cres,
    }
}

//export fgbasync_other
func fgbasync_other(fgbPort int64) {
	go func() {
		h := atomic.AddUint64(&handleIdx, 1)
		if h == 0 {
			panic("ran out of handle space")
		}

		handles.Store(h, fgb_other())

		sent := runtime.Send(fgbPort, []uint64{h}, func() {
			handles.LoadAndDelete(h)
		})
        if !sent {
            handles.LoadAndDelete(h)
        }
	}()
}

//export fgbasyncres_other
func fgbasyncres_other(h uint64) C.fgb_ret_other {
	v, ok := handles.LoadAndDelete(h)
	if !ok {
		return C.fgb_ret_other{
			err: unsafe.Pointer(C.CString("result handle is not valid")),
		}
	}

	return (v).(C.fgb_ret_other)
}

//export fgb_call_me
func fgb_call_me() (resw C.fgb_ret_call_me) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		resw = C.fgb_ret_call_me {
			err: unsafe.Pointer(C.CString(fmt.Sprintf("panic: %v", r))),
		}
	}()
	
	gres, gerr := orig.CallMe()
    if gerr != nil {
		return C.fgb_ret_call_me {
			err: unsafe.Pointer(C.CString(gerr.Error())),
		}
    }
    
	cres := (C.int)(gres)

    return C.fgb_ret_call_me {
        res: cres,
    }
}

//export fgbasync_call_me
func fgbasync_call_me(fgbPort int64) {
	go func() {
		h := atomic.AddUint64(&handleIdx, 1)
		if h == 0 {
			panic("ran out of handle space")
		}

		handles.Store(h, fgb_call_me())

		sent := runtime.Send(fgbPort, []uint64{h}, func() {
			handles.LoadAndDelete(h)
		})
        if !sent {
            handles.LoadAndDelete(h)
        }
	}()
}

//export fgbasyncres_call_me
func fgbasyncres_call_me(h uint64) C.fgb_ret_call_me {
	v, ok := handles.LoadAndDelete(h)
	if !ok {
		return C.fgb_ret_call_me{
			err: unsafe.Pointer(C.CString("result handle is not valid")),
		}
	}

	return (v).(C.fgb_ret_call_me)
}