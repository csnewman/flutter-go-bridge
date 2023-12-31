// This code was generated by flutter-go-bridge. Do not manually edit.
package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	orig "github.com/csnewman/flutter-go-bridge/exampleapp/go"
	"github.com/csnewman/flutter-go-bridge/runtime"
)

/*
#include <stdlib.h>
#include <stdint.h>

typedef struct {
	int x;
	int y;
	void* name;
} fgb_vt_point;

typedef struct {
	int res;
	void* err;
} fgb_ret_add;

typedef struct {
	fgb_vt_point res;
	void* err;
} fgb_ret_add_points;

typedef struct {
	int res;
	void* err;
} fgb_ret_add_error;

typedef struct {
	void* res;
	void* err;
} fgb_ret_new_obj;

typedef struct {
	void* err;
} fgb_ret_modify_obj;

typedef struct {
	void* res;
	void* err;
} fgb_ret_format_obj;
*/
import "C"

var (
	handles   = sync.Map{}
	handleIdx uint64
	ErrDart   = errors.New("dart")
)

// Required by cgo
func main() {}

//export fgbinternal_init
func fgbinternal_init(p unsafe.Pointer) unsafe.Pointer {
	err := runtime.InitializeApi(p)

	var cerr unsafe.Pointer
	if err != nil {
		cerr = unsafe.Pointer(C.CString(err.Error()))
	}

	return cerr
}

//export fgbinternal_alloc
func fgbinternal_alloc(size C.intptr_t) unsafe.Pointer {
	return C.malloc((C.uintptr_t)(size))
}

//export fgbinternal_free
func fgbinternal_free(ptr unsafe.Pointer) {
	C.free(ptr)
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

//export fgbempty_point
func fgbempty_point() (res C.fgb_vt_point) {
	return
}

func mapToPoint(from C.fgb_vt_point) (res orig.Point) {
	res.X = (int)(from.x)
	res.Y = (int)(from.y)
	res.Name = mapToString(from.name)
	return
}

func mapFromPoint(from orig.Point) (res C.fgb_vt_point) {
	res.x = (C.int)(from.X)
	res.y = (C.int)(from.Y)
	res.name = mapFromString(from.Name)
	return
}

func mapToObj(from unsafe.Pointer) *orig.Obj {
	h := uint64(uintptr(from))

	v, ok := handles.Load(h)
	if !ok {
		panic(fmt.Sprintf("invalid handle: %v", h))
	}

	return v.(*orig.Obj)
}

func mapFromObj(from *orig.Obj) unsafe.Pointer {
	h := atomic.AddUint64(&handleIdx, 1)
	if h == 0 {
		panic("ran out of handle space")
	}

	handles.Store(h, from)

	return unsafe.Pointer(uintptr(h))
}

//export fgbfree_obj
func fgbfree_obj(from unsafe.Pointer) {
	h := uint64(uintptr(from))

	handles.Delete(h)
}

//export fgb_add
func fgb_add(arg_a C.int, arg_b C.int) (resw C.fgb_ret_add) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		resw = C.fgb_ret_add{
			err: unsafe.Pointer(C.CString(fmt.Sprintf("panic: %v", r))),
		}
	}()
	
	arggo_a := (int)(arg_a)
	arggo_b := (int)(arg_b)
	gres := orig.Add(arggo_a, arggo_b)
	
	cres := (C.int)(gres)

	return C.fgb_ret_add{
		res: cres,
	}
}

//export fgbasync_add
func fgbasync_add(arg_a C.int, arg_b C.int, fgbPort int64) {
	go func() {
		h := atomic.AddUint64(&handleIdx, 1)
		if h == 0 {
			panic("ran out of handle space")
		}

		handles.Store(h, fgb_add(arg_a, arg_b))

		sent := runtime.Send(fgbPort, []uint64{h}, func() {
			handles.LoadAndDelete(h)
		})
		if !sent {
			handles.LoadAndDelete(h)
		}
	}()
}

//export fgbasyncres_add
func fgbasyncres_add(h uint64) C.fgb_ret_add {
	v, ok := handles.LoadAndDelete(h)
	if !ok {
		return C.fgb_ret_add{
			err: unsafe.Pointer(C.CString("result handle is not valid")),
		}
	}

	return (v).(C.fgb_ret_add)
}

//export fgb_add_points
func fgb_add_points(arg_a C.fgb_vt_point, arg_b C.fgb_vt_point) (resw C.fgb_ret_add_points) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		resw = C.fgb_ret_add_points{
			err: unsafe.Pointer(C.CString(fmt.Sprintf("panic: %v", r))),
		}
	}()
	
	arggo_a := mapToPoint(arg_a)
	arggo_b := mapToPoint(arg_b)
	gres := orig.AddPoints(arggo_a, arggo_b)
	
	cres := mapFromPoint(gres)

	return C.fgb_ret_add_points{
		res: cres,
	}
}

//export fgbasync_add_points
func fgbasync_add_points(arg_a C.fgb_vt_point, arg_b C.fgb_vt_point, fgbPort int64) {
	go func() {
		h := atomic.AddUint64(&handleIdx, 1)
		if h == 0 {
			panic("ran out of handle space")
		}

		handles.Store(h, fgb_add_points(arg_a, arg_b))

		sent := runtime.Send(fgbPort, []uint64{h}, func() {
			handles.LoadAndDelete(h)
		})
		if !sent {
			handles.LoadAndDelete(h)
		}
	}()
}

//export fgbasyncres_add_points
func fgbasyncres_add_points(h uint64) C.fgb_ret_add_points {
	v, ok := handles.LoadAndDelete(h)
	if !ok {
		return C.fgb_ret_add_points{
			err: unsafe.Pointer(C.CString("result handle is not valid")),
		}
	}

	return (v).(C.fgb_ret_add_points)
}

//export fgb_add_error
func fgb_add_error(arg_a C.int, arg_b C.int) (resw C.fgb_ret_add_error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		resw = C.fgb_ret_add_error{
			err: unsafe.Pointer(C.CString(fmt.Sprintf("panic: %v", r))),
		}
	}()
	
	arggo_a := (int)(arg_a)
	arggo_b := (int)(arg_b)
	gres, gerr := orig.AddError(arggo_a, arggo_b)
	if gerr != nil {
		return C.fgb_ret_add_error{
			err: unsafe.Pointer(C.CString(gerr.Error())),
		}
	}
	
	cres := (C.int)(gres)

	return C.fgb_ret_add_error{
		res: cres,
	}
}

//export fgbasync_add_error
func fgbasync_add_error(arg_a C.int, arg_b C.int, fgbPort int64) {
	go func() {
		h := atomic.AddUint64(&handleIdx, 1)
		if h == 0 {
			panic("ran out of handle space")
		}

		handles.Store(h, fgb_add_error(arg_a, arg_b))

		sent := runtime.Send(fgbPort, []uint64{h}, func() {
			handles.LoadAndDelete(h)
		})
		if !sent {
			handles.LoadAndDelete(h)
		}
	}()
}

//export fgbasyncres_add_error
func fgbasyncres_add_error(h uint64) C.fgb_ret_add_error {
	v, ok := handles.LoadAndDelete(h)
	if !ok {
		return C.fgb_ret_add_error{
			err: unsafe.Pointer(C.CString("result handle is not valid")),
		}
	}

	return (v).(C.fgb_ret_add_error)
}

//export fgb_new_obj
func fgb_new_obj(arg_name unsafe.Pointer, arg_other C.int) (resw C.fgb_ret_new_obj) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		resw = C.fgb_ret_new_obj{
			err: unsafe.Pointer(C.CString(fmt.Sprintf("panic: %v", r))),
		}
	}()
	
	arggo_name := mapToString(arg_name)
	arggo_other := (int)(arg_other)
	gres := orig.NewObj(arggo_name, arggo_other)
	
	cres := mapFromObj(gres)

	return C.fgb_ret_new_obj{
		res: cres,
	}
}

//export fgbasync_new_obj
func fgbasync_new_obj(arg_name unsafe.Pointer, arg_other C.int, fgbPort int64) {
	go func() {
		h := atomic.AddUint64(&handleIdx, 1)
		if h == 0 {
			panic("ran out of handle space")
		}

		handles.Store(h, fgb_new_obj(arg_name, arg_other))

		sent := runtime.Send(fgbPort, []uint64{h}, func() {
			handles.LoadAndDelete(h)
		})
		if !sent {
			handles.LoadAndDelete(h)
		}
	}()
}

//export fgbasyncres_new_obj
func fgbasyncres_new_obj(h uint64) C.fgb_ret_new_obj {
	v, ok := handles.LoadAndDelete(h)
	if !ok {
		return C.fgb_ret_new_obj{
			err: unsafe.Pointer(C.CString("result handle is not valid")),
		}
	}

	return (v).(C.fgb_ret_new_obj)
}

//export fgb_modify_obj
func fgb_modify_obj(arg_o unsafe.Pointer) (resw C.fgb_ret_modify_obj) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		resw = C.fgb_ret_modify_obj{
			err: unsafe.Pointer(C.CString(fmt.Sprintf("panic: %v", r))),
		}
	}()
	
	arggo_o := mapToObj(arg_o)
	orig.ModifyObj(arggo_o)
	

	return C.fgb_ret_modify_obj{
	}
}

//export fgbasync_modify_obj
func fgbasync_modify_obj(arg_o unsafe.Pointer, fgbPort int64) {
	go func() {
		h := atomic.AddUint64(&handleIdx, 1)
		if h == 0 {
			panic("ran out of handle space")
		}

		handles.Store(h, fgb_modify_obj(arg_o))

		sent := runtime.Send(fgbPort, []uint64{h}, func() {
			handles.LoadAndDelete(h)
		})
		if !sent {
			handles.LoadAndDelete(h)
		}
	}()
}

//export fgbasyncres_modify_obj
func fgbasyncres_modify_obj(h uint64) C.fgb_ret_modify_obj {
	v, ok := handles.LoadAndDelete(h)
	if !ok {
		return C.fgb_ret_modify_obj{
			err: unsafe.Pointer(C.CString("result handle is not valid")),
		}
	}

	return (v).(C.fgb_ret_modify_obj)
}

//export fgb_format_obj
func fgb_format_obj(arg_o unsafe.Pointer) (resw C.fgb_ret_format_obj) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		resw = C.fgb_ret_format_obj{
			err: unsafe.Pointer(C.CString(fmt.Sprintf("panic: %v", r))),
		}
	}()
	
	arggo_o := mapToObj(arg_o)
	gres := orig.FormatObj(arggo_o)
	
	cres := mapFromString(gres)

	return C.fgb_ret_format_obj{
		res: cres,
	}
}

//export fgbasync_format_obj
func fgbasync_format_obj(arg_o unsafe.Pointer, fgbPort int64) {
	go func() {
		h := atomic.AddUint64(&handleIdx, 1)
		if h == 0 {
			panic("ran out of handle space")
		}

		handles.Store(h, fgb_format_obj(arg_o))

		sent := runtime.Send(fgbPort, []uint64{h}, func() {
			handles.LoadAndDelete(h)
		})
		if !sent {
			handles.LoadAndDelete(h)
		}
	}()
}

//export fgbasyncres_format_obj
func fgbasyncres_format_obj(h uint64) C.fgb_ret_format_obj {
	v, ok := handles.LoadAndDelete(h)
	if !ok {
		return C.fgb_ret_format_obj{
			err: unsafe.Pointer(C.CString("result handle is not valid")),
		}
	}

	return (v).(C.fgb_ret_format_obj)
}
