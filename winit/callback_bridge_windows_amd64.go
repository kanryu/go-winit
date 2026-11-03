//go:build windows && amd64

package winit

/*
#cgo windows,amd64 CFLAGS: -I${SRCDIR}/../third_party/libwinit
#include "winit.h"
#include <stdint.h>
#include <stdlib.h>

extern uint32_t goWinitEventCallback(uintptr_t data, uint32_t eventType,
    uint64_t windowHigh, uint64_t windowLow,
    uint32_t resizedWidth, uint32_t resizedHeight,
    double cursorX, double cursorY,
    double scaleFactor, uint32_t scaleWidth, uint32_t scaleHeight,
    uint64_t receivedLength, uint8_t receivedByte1, uint8_t receivedByte2,
    uint8_t receivedByte3, uint8_t receivedByte4,
    uint32_t keyboardState, uint32_t hasPhysicalKey,
    uint32_t physicalKeyCode, uint32_t isSynthetic,
	uint64_t logicalKeyLength, const uint8_t *logicalKey,
	uint32_t occluded);

static WinitControlFlow go_winit_event_callback(void *data, WinitEvent *event) {
    return (WinitControlFlow)goWinitEventCallback(
        (uintptr_t)data,
        (uint32_t)event->event_type,
        event->window_id.high,
        event->window_id.low,
        event->window_resized.width,
        event->window_resized.height,
        event->cursor_moved.x,
        event->cursor_moved.y,
        event->scale_factor.scale_factor,
        event->scale_factor.width,
        event->scale_factor.height,
        (uint64_t)event->received_character.length,
        event->received_character.byte_1,
        event->received_character.byte_2,
        event->received_character.byte_3,
        event->received_character.byte_4,
        (uint32_t)event->keyboard_input.state,
        (uint32_t)event->keyboard_input.has_physical_keycode,
        event->keyboard_input.physical_keycode,
        (uint32_t)event->keyboard_input.is_synthetic,
        (uint64_t)event->keyboard_input.logical_key_length,
        event->keyboard_input.logical_key,
        (uint32_t)event->window_occluded.is_occluded);
}

static uintptr_t go_winit_event_callback_ptr(void) {
    return (uintptr_t)&go_winit_event_callback;
}

static void *go_winit_callback_token_new(void) {
    return malloc(1);
}

static void go_winit_callback_token_drop(void *token) {
    free(token);
}
*/
import "C"

import (
	"runtime/cgo"
	"sync"
)

var callbackHandles sync.Map // map[uintptr]cgo.Handle

type userEventQueue struct {
	mu     sync.Mutex
	events []any
}

func (q *userEventQueue) push(event any) {
	q.mu.Lock()
	q.events = append(q.events, event)
	q.mu.Unlock()
}

func (q *userEventQueue) drain() []any {
	q.mu.Lock()
	events := q.events
	q.events = nil
	q.mu.Unlock()
	return events
}

type typedCallbackState struct {
	loop     *EventLoop
	callback TypedCallback
}

func goCallbackPointer() uintptr {
	return uintptr(C.go_winit_event_callback_ptr())
}

func runGoCallback(loop *EventLoop, callback Callback) {
	runGoCallbackValue(loop, callback)
}

func runGoTypedCallback(loop *EventLoop, callback TypedCallback) {
	if callback == nil {
		return
	}
	runGoCallbackValue(loop, typedCallbackState{loop: loop, callback: callback})
}

func runGoCallbackValue(loop *EventLoop, callback any) {
	if loop == nil || loop.ptr == nil || callback == nil {
		return
	}
	handle := cgo.NewHandle(callback)
	token := C.go_winit_callback_token_new()
	if token == nil {
		handle.Delete()
		return
	}
	key := uintptr(token)
	callbackHandles.Store(key, handle)
	defer func() {
		callbackHandles.Delete(key)
		handle.Delete()
		C.go_winit_callback_token_drop(token)
	}()
	loop.RunNative(token, goCallbackPointer())
}
