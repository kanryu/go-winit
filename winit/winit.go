// Package winit provides the Go-facing handles for libwinit.
//
// Native calls are made through the HikeC-generated internal/libwinit package;
// the internal/libwinit_native package supplies the linked .syso object.
package winit

import (
	"runtime"
	"unsafe"

	native "github.com/kanryu/go-winit/internal/libwinit"
	_ "github.com/kanryu/go-winit/internal/libwinit_native"
)

type EventLoopType uint8

const (
	EventLoopWindows EventLoopType = iota
	EventLoopMacOS
	EventLoopX11
	EventLoopWayland
	EventLoopUnknown
)

type ControlFlow uint32

const (
	Poll ControlFlow = iota
	Wait
	Exit
)

type EventType uint32

const (
	EventUnknown EventType = iota
	WindowEventResized
	WindowEventMoved
	WindowEventCloseRequested
	WindowEventDestroyed
	WindowEventDroppedFile
	WindowEventHoveredFile
	WindowEventHoveredFileCancelled
	WindowEventReceivedCharacter
	WindowEventFocused
	WindowEventKeyboardInput
	WindowEventCursorMoved
	WindowEventCursorEntered
	WindowEventCursorLeft
	WindowEventMouseWheel
	WindowEventMouseInput
	WindowEventTouchpadPressure
	WindowEventAxisMotion
	WindowEventTouch
	WindowEventScaleFactorChanged
	NewEvents
	MainEventsCleared
	LoopDestroyed
	Suspended
	Resumed
	RedrawRequested
	RedrawEventsCleared
	ModifiersChanged
	UserEvent
	WindowEventOccluded
)

type Event struct {
	Type            EventType
	WindowID        [2]uint64
	ResizedWidth    uint32
	ResizedHeight   uint32
	CursorX         float64
	CursorY         float64
	ScaleFactor     float64
	ScaleWidth      uint32
	ScaleHeight     uint32
	ReceivedLength  uint64
	ReceivedBytes   [4]byte
	KeyboardState   InputElementState
	HasPhysicalKey  bool
	PhysicalKeyCode uint32
	IsSynthetic     bool
	LogicalKey      string
	Occluded        bool
}

// ReceivedCharacterString returns the UTF-8 character carried by an IME
// commit or ReceivedCharacter event.
func (e Event) ReceivedCharacterString() string {
	length := e.ReceivedLength
	if length > uint64(len(e.ReceivedBytes)) {
		length = uint64(len(e.ReceivedBytes))
	}
	return string(e.ReceivedBytes[:length])
}

type Callback func(Event) ControlFlow

// TypedCallback receives the public event value corresponding to each native
// event. Use RunTyped when dispatching by Go event type instead of Event.Type.
type TypedCallback func(EventValue) ControlFlow

// StringBox is the native string-box handle used by libwinit's title APIs.
type StringBox struct{ ptr *byte }

// BoxerStringFromByteString creates a native StringBox from UTF-8 bytes.
func BoxerStringFromByteString(value []byte) *StringBox {
	var data *byte
	if len(value) != 0 {
		data = (*byte)(unsafe.Pointer(&value[0]))
	}
	s := &StringBox{ptr: native.BoxerStringFromByteString(data, uint64(len(value)))}
	runtime.KeepAlive(value)
	return s
}

// Drop releases the native StringBox.
func (s *StringBox) Drop() {
	if s == nil || s.ptr == nil {
		return
	}
	native.BoxerStringDrop(s.ptr)
	s.ptr = nil
}

// BoxerStringDrop releases a native StringBox.
func BoxerStringDrop(s *StringBox) {
	if s != nil {
		s.Drop()
	}
}

type EventLoop struct {
	ptr        *byte
	userEvents *userEventQueue
}

func NewEventLoop() *EventLoop {
	e := &EventLoop{ptr: native.EventLoopNew()}
	runtime.SetFinalizer(e, (*EventLoop).Close)
	return e
}

func (e *EventLoop) Close() {
	if e == nil || e.ptr == nil {
		return
	}
	native.EventLoopDrop(e.ptr)
	e.ptr = nil
	runtime.SetFinalizer(e, nil)
}

// EventLoopProxy wakes the event loop from another goroutine and delivers
// application-owned values as UserEventValue values to RunTyped.
type EventLoopProxy struct {
	ptr   *byte
	queue *userEventQueue
}

// CreateProxy creates a thread-safe event-loop proxy.
func (e *EventLoop) CreateProxy() *EventLoopProxy {
	if e == nil || e.ptr == nil {
		return nil
	}
	proxy := &EventLoopProxy{ptr: native.EventLoopCreateProxy(e.ptr), queue: &userEventQueue{}}
	if proxy.ptr == nil {
		return nil
	}
	e.userEvents = proxy.queue
	runtime.SetFinalizer(proxy, (*EventLoopProxy).Close)
	return proxy
}

// Close releases the native event-loop proxy.
func (p *EventLoopProxy) Close() {
	if p == nil || p.ptr == nil {
		return
	}
	native.EventLoopDropProxy(p.ptr)
	p.ptr = nil
	runtime.SetFinalizer(p, nil)
}

// Send queues an application-owned value and wakes the event loop. The value
// is delivered on the event-loop goroutine as UserEventValue.Value.
func (p *EventLoopProxy) Send(value any) bool {
	if p == nil || p.ptr == nil || p.queue == nil {
		return false
	}
	p.queue.push(value)
	if native.EventLoopWake(p.ptr) {
		return true
	}
	return false
}

// Wake wakes the event loop without attaching an application value.
func (p *EventLoopProxy) Wake() bool {
	if p == nil || p.ptr == nil {
		return false
	}
	return native.EventLoopWake(p.ptr)
}

// Exit requests the event loop to exit through this proxy.
func (p *EventLoopProxy) Exit() bool {
	if p == nil || p.ptr == nil {
		return false
	}
	return native.EventLoopExit(p.ptr)
}

func (e *EventLoop) Type() EventLoopType {
	if e == nil || e.ptr == nil {
		return EventLoopUnknown
	}
	return EventLoopType(native.EventLoopGetType(e.ptr))
}

// Run enters the native event loop and dispatches each event to callback in Go.
// Custom native callbacks can be supplied through RunNative.
func (e *EventLoop) Run(callback Callback) {
	if e == nil || e.ptr == nil {
		return
	}
	runGoCallback(e, callback)
}

// RunTyped enters the native event loop and dispatches typed event values.
func (e *EventLoop) RunTyped(callback TypedCallback) {
	if e == nil || e.ptr == nil {
		return
	}
	runGoTypedCallback(e, callback)
}

// RunNative enters the event loop with a native callback pointer obtained from
// the Hike binding. data is passed unchanged to that callback.
func (e *EventLoop) RunNative(data unsafe.Pointer, callback uintptr) {
	if e == nil || e.ptr == nil || callback == 0 {
		return
	}
	native.EventLoopRunData(e.ptr, (*byte)(data), (*byte)(unsafe.Pointer(callback)))
}

// RunWindow runs the event loop with the native window as callback data. The
// Hike callback releases that window when a close request arrives.
func (e *EventLoop) RunWindow(window *Window) {
	if e == nil || e.ptr == nil || !window.valid() {
		return
	}
	e.RunNative(unsafe.Pointer(window.ptr), native.EventLoopRunDataCallbackPtr())
	window.ptr = nil
	runtime.SetFinalizer(window, nil)
}

type WindowAttributes struct{ ptr *byte }

func NewWindowAttributes() *WindowAttributes {
	a := &WindowAttributes{ptr: native.WindowAttributesNew()}
	runtime.SetFinalizer(a, (*WindowAttributes).Close)
	return a
}
func (a *WindowAttributes) Close() {
	if a == nil || a.ptr == nil {
		return
	}
	native.WindowAttributesDrop(a.ptr)
	a.ptr = nil
	runtime.SetFinalizer(a, nil)
}
func (a *WindowAttributes) valid() bool { return a != nil && a.ptr != nil }
func (a *WindowAttributes) WithDecorations(v bool) *WindowAttributes {
	if a.valid() {
		native.WindowAttributesWithDecorations(a.ptr, v)
	}
	return a
}
func (a *WindowAttributes) WithTitle(title string) *WindowAttributes {
	if a.valid() {
		box := BoxerStringFromByteString([]byte(title))
		if box != nil {
			native.WindowAttributesWithTitle(a.ptr, box.ptr)
			box.Drop()
		}
	}
	return a
}
func (a *WindowAttributes) WithTransparency(v bool) *WindowAttributes {
	if a.valid() {
		native.WindowAttributesWithTransparency(a.ptr, v)
	}
	return a
}
func (a *WindowAttributes) WithResizable(v bool) *WindowAttributes {
	if a.valid() {
		native.WindowAttributesWithResizable(a.ptr, v)
	}
	return a
}
func (a *WindowAttributes) WithDimensions(width, height float64) *WindowAttributes {
	if a.valid() {
		native.WindowAttributesWithDimensionsFloat64(a.ptr, width, height)
	}
	return a
}
func (a *WindowAttributes) WithMaximized(v bool) *WindowAttributes {
	if a.valid() {
		native.WindowAttributesWithMaximized(a.ptr, v)
	}
	return a
}
func (a *WindowAttributes) WithVisibility(v bool) *WindowAttributes {
	if a.valid() {
		native.WindowAttributesWithVisibility(a.ptr, v)
	}
	return a
}
func (a *WindowAttributes) WithAlwaysOnTop(v bool) *WindowAttributes {
	if a.valid() {
		native.WindowAttributesWithAlwaysOnTop(a.ptr, v)
	}
	return a
}
func (a *WindowAttributes) WithFullSize(v bool) *WindowAttributes {
	if a.valid() {
		native.WindowAttributesWithFullSize(a.ptr, v)
	}
	return a
}

type Window struct{ ptr *byte }

func (e *EventLoop) CreateWindow(a *WindowAttributes) *Window {
	if e == nil || e.ptr == nil || !a.valid() {
		return nil
	}
	w := &Window{ptr: native.CreateWindow(e.ptr, a.ptr)}
	// libwinit consumes WindowAttributes during window creation.
	a.ptr = nil
	runtime.SetFinalizer(a, nil)
	if w.ptr != nil {
		runtime.SetFinalizer(w, (*Window).Close)
	}
	return w
}
func (w *Window) Close() {
	if w == nil || w.ptr == nil {
		return
	}
	native.WindowDrop(w.ptr)
	w.ptr = nil
	runtime.SetFinalizer(w, nil)
}
func (w *Window) valid() bool { return w != nil && w.ptr != nil }
func (w *Window) RequestRedraw() {
	if w.valid() {
		native.WindowRequestRedraw(w.ptr)
	}
}
func (w *Window) ScaleFactor() float64 {
	if !w.valid() {
		return 0
	}
	return native.WindowGetScaleFactor(w.ptr)
}
func (w *Window) SetSize(width, height uint32) {
	if w.valid() {
		native.WindowSetInnerSize(w.ptr, width, height)
	}
}
func (w *Window) SetPosition(x, y int32) {
	if w.valid() {
		native.WindowSetPosition(w.ptr, x, y)
	}
}
func (w *Window) SetMaximized(v bool) {
	if w.valid() {
		native.WindowSetMaximized(w.ptr, v)
	}
}

// IsMaximized reports whether the window is currently maximized.
func (w *Window) IsMaximized() bool {
	if !w.valid() {
		return false
	}
	return native.WindowIsMaximized(w.ptr)
}
func (w *Window) Focus() {
	if w.valid() {
		native.WindowFocusWindow(w.ptr)
	}
}

// SetIMEAllowed enables or disables IME composition for the window.
func (w *Window) SetIMEAllowed(allowed bool) {
	if w.valid() {
		native.WindowSetIMEAllowed(w.ptr, allowed)
	}
}
func (w *Window) HWND() unsafe.Pointer {
	if !w.valid() {
		return nil
	}
	return unsafe.Pointer(native.WindowGetHWND(w.ptr))
}

func (w *Window) ID() (high, low uint64) {
	if !w.valid() {
		return
	}
	value := native.U128{}
	box := native.U128Box{Tag: 0, Value: uintptr(unsafe.Pointer(&value))}
	native.WindowGetID(w.ptr, (*byte)(unsafe.Pointer(&box)))
	return value.High, value.Low
}

func (w *Window) Size() (width, height uint32) {
	if !w.valid() {
		return
	}
	value := native.SizeU32{}
	box := native.ValueBoxSizeBoxU32{Tag: 0, Value: uintptr(unsafe.Pointer(&value))}
	native.WindowGetInnerSize(w.ptr, (*byte)(unsafe.Pointer(&box)))
	return value.Width, value.Height
}

func (w *Window) Position() (x, y int32) {
	if !w.valid() {
		return
	}
	value := native.PointBoxI32{}
	box := native.ValueBoxPointBoxI32{Tag: 0, Value: uintptr(unsafe.Pointer(&value))}
	native.WindowGetPosition(w.ptr, (*byte)(unsafe.Pointer(&box)))
	return value.X, value.Y
}
