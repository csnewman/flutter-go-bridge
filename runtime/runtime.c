#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <pthread.h>
#include "dart_api_dl.h"

// GC related callbacks (pins, ports) need to spawn a separate thread before calling Go to prevent crashing on darwin.

extern void fgbinternal_freepingo(void* ptr);

void fgbinternal_freepinthread(void* ptr) {
    fgbinternal_freepingo(ptr);

    pthread_exit(0);
}

void fgbinternal_freepin(void* token) {
    pthread_t id;
    pthread_create(&id, NULL, fgbinternal_freepinthread, token);
    pthread_detach(id);
}

extern void ifgb_callback(void* ptr);

void ifgb_callbackthread(void* ptr) {
	ifgb_callback(ptr);

    pthread_exit(0);
}

void InternalBlobCallback(void* isolate_callback_data, void* peer) {
	pthread_t id;
	pthread_create(&id, NULL, ifgb_callbackthread, peer);
	pthread_detach(id);
}

bool InternalPostBlob(Dart_Port port_id, intptr_t len, void* data, uintptr_t peer) {
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
