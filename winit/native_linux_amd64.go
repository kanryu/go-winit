//go:build linux && amd64

package winit

/*
#cgo linux,amd64 LDFLAGS: -lX11 -lwayland-client -lxkbcommon -ldl -lpthread -lm
extern void* winit_event_loop_new(void);
extern void winit_event_loop_drop(void*);
extern void winit_event_loop_run_return(void*);
extern void* winit_window_attributes_new(void);
extern void winit_window_attributes_drop(void*);
extern void winit_window_attributes_with_title(void*, void*);
extern void winit_window_attributes_with_decorations(void*, _Bool);
extern void winit_window_attributes_with_transparency(void*, _Bool);
extern void winit_window_attributes_with_resizable(void*, _Bool);
extern void winit_window_attributes_with_dimensions(void*, double, double);
extern void winit_window_attributes_with_maximized(void*, _Bool);
extern void winit_window_attributes_with_visibility(void*, _Bool);
extern void winit_window_attributes_with_always_on_top(void*, _Bool);
extern void winit_window_attributes_with_full_size(void*, _Bool);
extern void* winit_create_window(void*, void*);
extern void winit_window_drop(void*);
extern void winit_window_set_ime_allowed(void*, _Bool);
extern void winit_window_set_inner_size(void*, unsigned int, unsigned int);
extern void winit_window_focus_window(void*);
extern _Bool winit_window_is_maximized(void*);
extern void* winit_event_queue_pop(void);
extern void winit_event_drop(void*);
*/
import "C"

import "unsafe"
import native "github.com/kanryu/go-winit/internal/libwinit"

func eventLoopNewNative() *byte { return (*byte)(C.winit_event_loop_new()) }
func eventLoopDropNative(ptr *byte) { C.winit_event_loop_drop(unsafe.Pointer(ptr)) }
func eventLoopRunReturnNative(ptr *byte) { C.winit_event_loop_run_return(unsafe.Pointer(ptr)) }
func windowAttributesNewNative() *byte { return (*byte)(C.winit_window_attributes_new()) }
func windowAttributesDropNative(ptr *byte) { C.winit_window_attributes_drop(unsafe.Pointer(ptr)) }
func windowAttributesWithTitleNative(ptr, title *byte) { C.winit_window_attributes_with_title(unsafe.Pointer(ptr), unsafe.Pointer(title)) }
func windowAttributesWithDecorationsNative(ptr *byte, value bool) { C.winit_window_attributes_with_decorations(unsafe.Pointer(ptr), C.bool(value)) }
func windowAttributesWithTransparencyNative(ptr *byte, value bool) { C.winit_window_attributes_with_transparency(unsafe.Pointer(ptr), C.bool(value)) }
func windowAttributesWithResizableNative(ptr *byte, value bool) { C.winit_window_attributes_with_resizable(unsafe.Pointer(ptr), C.bool(value)) }
func windowAttributesWithDimensionsNative(ptr *byte, width, height float64) { C.winit_window_attributes_with_dimensions(unsafe.Pointer(ptr), C.double(width), C.double(height)) }
func windowAttributesWithMaximizedNative(ptr *byte, value bool) { C.winit_window_attributes_with_maximized(unsafe.Pointer(ptr), C.bool(value)) }
func windowAttributesWithVisibilityNative(ptr *byte, value bool) { C.winit_window_attributes_with_visibility(unsafe.Pointer(ptr), C.bool(value)) }
func windowAttributesWithAlwaysOnTopNative(ptr *byte, value bool) { C.winit_window_attributes_with_always_on_top(unsafe.Pointer(ptr), C.bool(value)) }
func windowAttributesWithFullSizeNative(ptr *byte, value bool) { C.winit_window_attributes_with_full_size(unsafe.Pointer(ptr), C.bool(value)) }
func createWindowNative(eventLoop, attributes *byte) *byte { return (*byte)(C.winit_create_window(unsafe.Pointer(eventLoop), unsafe.Pointer(attributes))) }
func windowDropNative(ptr *byte) { C.winit_window_drop(unsafe.Pointer(ptr)) }
func windowSetIMEAllowedNative(ptr *byte, value bool) { C.winit_window_set_ime_allowed(unsafe.Pointer(ptr), C.bool(value)) }
func windowSetInnerSizeNative(ptr *byte, width, height uint32) { C.winit_window_set_inner_size(unsafe.Pointer(ptr), C.uint(width), C.uint(height)) }
func windowFocusWindowNative(ptr *byte) { C.winit_window_focus_window(unsafe.Pointer(ptr)) }
func windowIsMaximizedNative(ptr *byte) bool { return bool(C.winit_window_is_maximized(unsafe.Pointer(ptr))) }
func eventQueuePopNative() *native.Event { return (*native.Event)(C.winit_event_queue_pop()) }
func eventDropNative(ptr *byte) { C.winit_event_drop(unsafe.Pointer(ptr)) }
