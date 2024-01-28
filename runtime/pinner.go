package runtime

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
#include <stdlib.h>
*/
import "C"

var (
	handles     = sync.Map{}
	handleIdx   = &atomic.Uintptr{}
	handleCount = &atomic.Int64{}
)

func ActiveHandles() int64 {
	return handleCount.Load()
}

func Pin(value any) uintptr {
	// TODO: On 32bit, recycle addresses
	h := handleIdx.Add(1)
	if h == 0 {
		panic("ran out of handle space")
	}

	handles.Store(h, value)
	handleCount.Add(1)

	return h
}

func GetPin[T any](ptr uintptr) T {
	v, ok := handles.Load(ptr)
	if !ok {
		panic(fmt.Sprintf("invalid handle: %v", ptr))
	}

	return v.(T)
}

func FreePin(ptr uintptr) {
	_, ok := handles.LoadAndDelete(ptr)
	if ok {
		handleCount.Add(-1)
	}
}

//export fgbinternal_freepingo
func fgbinternal_freepingo(ptr uintptr) {
	FreePin(ptr)
}
