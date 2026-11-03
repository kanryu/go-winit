//go:build windows && amd64

package winit

/*
#include <stdint.h>
*/
import "C"

import "runtime/cgo"

import "unsafe"

//export goWinitEventCallback
func goWinitEventCallback(data C.uintptr_t, eventType C.uint32_t,
	high C.uint64_t, low C.uint64_t,
	resizedWidth C.uint32_t, resizedHeight C.uint32_t,
	cursorX C.double, cursorY C.double,
	scaleFactor C.double, scaleWidth C.uint32_t, scaleHeight C.uint32_t,
	receivedLength C.uint64_t, receivedByte1 C.uint8_t, receivedByte2 C.uint8_t,
	receivedByte3 C.uint8_t, receivedByte4 C.uint8_t,
	keyboardState C.uint32_t, hasPhysicalKey C.uint32_t,
	physicalKeyCode C.uint32_t, isSynthetic C.uint32_t,
	logicalKeyLength C.uint64_t, logicalKey *C.uint8_t,
	occluded C.uint32_t) (result C.uint32_t) {
	result = C.uint32_t(Exit)
	defer func() {
		if recover() != nil {
			result = C.uint32_t(Exit)
		}
	}()
	value, ok := callbackHandles.Load(uintptr(data))
	if !ok {
		return
	}
	handle, ok := value.(cgo.Handle)
	if !ok {
		return
	}
	raw := Event{
		Type:         EventType(eventType),
		WindowID:     [2]uint64{uint64(high), uint64(low)},
		ResizedWidth: uint32(resizedWidth), ResizedHeight: uint32(resizedHeight),
		CursorX: float64(cursorX), CursorY: float64(cursorY),
		ScaleFactor: float64(scaleFactor), ScaleWidth: uint32(scaleWidth), ScaleHeight: uint32(scaleHeight),
		ReceivedLength:  uint64(receivedLength),
		ReceivedBytes:   [4]byte{byte(receivedByte1), byte(receivedByte2), byte(receivedByte3), byte(receivedByte4)},
		KeyboardState:   InputElementState(keyboardState),
		HasPhysicalKey:  hasPhysicalKey != 0,
		PhysicalKeyCode: uint32(physicalKeyCode),
		IsSynthetic:     isSynthetic != 0,
		Occluded:        occluded != 0,
	}
	if logicalKeyLength > 64 {
		logicalKeyLength = 64
	}
	if logicalKeyLength != 0 && logicalKey != nil {
		raw.LogicalKey = string(C.GoBytes(unsafe.Pointer(logicalKey), C.int(logicalKeyLength)))
	}
	switch callback := handle.Value().(type) {
	case Callback:
		return C.uint32_t(callback(raw))
	case typedCallbackState:
		if raw.Type == UserEvent && callback.loop != nil && callback.loop.userEvents != nil {
			flow := Wait
			for _, value := range callback.loop.userEvents.drain() {
				flow = callback.callback(UserEventValue{Kind: UserEventWakeUp, Value: value})
				if flow == Exit {
					break
				}
			}
			return C.uint32_t(flow)
		}
		return C.uint32_t(callback.callback(eventValueFromRaw(raw)))
	default:
		return
	}
}
