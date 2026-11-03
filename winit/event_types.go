package winit

// EventValue is the typed representation of one libwinit event.  The low-level
// Event record is kept separately because it mirrors the C ABI; these values
// are the types used by the public event dispatcher.
type EventValue interface {
	EventType() EventType
}

func eventValueFromRaw(raw Event) EventValue {
	window := WindowEvent{WindowID: WindowID{High: raw.WindowID[0], Low: raw.WindowID[1]}}
	switch raw.Type {
	case WindowEventResized:
		return WindowResizedEvent{WindowEvent: window, Width: raw.ResizedWidth, Height: raw.ResizedHeight}
	case WindowEventCloseRequested:
		return WindowCloseRequestedEvent{WindowEvent: window}
	case WindowEventCursorMoved:
		return WindowCursorMovedEvent{WindowEvent: window, X: raw.CursorX, Y: raw.CursorY}
	case WindowEventReceivedCharacter:
		return WindowReceivedCharacterEvent{WindowEvent: window, Character: raw.ReceivedCharacterString()}
	case WindowEventKeyboardInput:
		return WindowKeyboardInputEvent{
			WindowEvent:     window,
			State:           raw.KeyboardState,
			HasPhysicalKey:  raw.HasPhysicalKey,
			PhysicalKeyCode: raw.PhysicalKeyCode,
			IsSynthetic:     raw.IsSynthetic,
			LogicalKey:      raw.LogicalKey,
		}
	case WindowEventOccluded:
		return WindowOccludedEvent{WindowEvent: window, Occluded: raw.Occluded}
	case WindowEventScaleFactorChanged:
		return WindowScaleFactorChangedEvent{WindowEvent: window, ScaleFactor: raw.ScaleFactor, Width: raw.ScaleWidth, Height: raw.ScaleHeight}
	case NewEvents:
		return NewEventsEvent{}
	case MainEventsCleared:
		return MainEventsClearedEvent{}
	case LoopDestroyed:
		return LoopDestroyedEvent{}
	case Suspended:
		return SuspendedEvent{}
	case Resumed:
		return ResumedEvent{}
	case RedrawRequested:
		return RedrawRequestedEvent{WindowEvent: window}
	case RedrawEventsCleared:
		return RedrawEventsClearedEvent{}
	case ModifiersChanged:
		return ModifiersChangedEvent{WindowEvent: window}
	case UserEvent:
		return UserEventValue{}
	default:
		return UnknownEvent{Type: raw.Type, WindowID: window.WindowID}
	}
}

// WindowID identifies the window associated with an event.
type WindowID struct {
	High uint64
	Low  uint64
}

// WindowEvent is embedded by window-specific event values.
type WindowEvent struct {
	WindowID WindowID
}

type WindowResizedEvent struct {
	WindowEvent
	Width  uint32
	Height uint32
}

func (WindowResizedEvent) EventType() EventType { return WindowEventResized }

type WindowMovedEvent struct {
	WindowEvent
	X int32
	Y int32
}

func (WindowMovedEvent) EventType() EventType { return WindowEventMoved }

type WindowCloseRequestedEvent struct{ WindowEvent }

func (WindowCloseRequestedEvent) EventType() EventType { return WindowEventCloseRequested }

type WindowDestroyedEvent struct{ WindowEvent }

func (WindowDestroyedEvent) EventType() EventType { return WindowEventDestroyed }

type WindowDroppedFileEvent struct{ WindowEvent }

func (WindowDroppedFileEvent) EventType() EventType { return WindowEventDroppedFile }

type WindowHoveredFileEvent struct{ WindowEvent }

func (WindowHoveredFileEvent) EventType() EventType { return WindowEventHoveredFile }

type WindowHoveredFileCancelledEvent struct{ WindowEvent }

func (WindowHoveredFileCancelledEvent) EventType() EventType {
	return WindowEventHoveredFileCancelled
}

type WindowReceivedCharacterEvent struct {
	WindowEvent
	Character string
}

func (WindowReceivedCharacterEvent) EventType() EventType {
	return WindowEventReceivedCharacter
}

type WindowFocusedEvent struct {
	WindowEvent
	Focused bool
}

func (WindowFocusedEvent) EventType() EventType { return WindowEventFocused }

type WindowOccludedEvent struct {
	WindowEvent
	Occluded bool
}

func (WindowOccludedEvent) EventType() EventType { return WindowEventOccluded }

type WindowKeyboardInputEvent struct {
	WindowEvent
	DeviceID        int64
	State           InputElementState
	HasPhysicalKey  bool
	PhysicalKeyCode uint32
	IsSynthetic     bool
	LogicalKey      string
}

func (WindowKeyboardInputEvent) EventType() EventType { return WindowEventKeyboardInput }

type WindowCursorMovedEvent struct {
	WindowEvent
	DeviceID int64
	X        float64
	Y        float64
}

func (WindowCursorMovedEvent) EventType() EventType { return WindowEventCursorMoved }

type WindowCursorEnteredEvent struct{ WindowEvent }

func (WindowCursorEnteredEvent) EventType() EventType { return WindowEventCursorEntered }

type WindowCursorLeftEvent struct{ WindowEvent }

func (WindowCursorLeftEvent) EventType() EventType { return WindowEventCursorLeft }

type WindowMouseWheelEvent struct {
	WindowEvent
	DeviceID int64
	Phase    TouchPhase
	Delta    MouseScrollDelta
}

func (WindowMouseWheelEvent) EventType() EventType { return WindowEventMouseWheel }

type WindowMouseInputEvent struct {
	WindowEvent
	DeviceID int64
	State    InputElementState
	Button   MouseButton
}

func (WindowMouseInputEvent) EventType() EventType { return WindowEventMouseInput }

type WindowTouchpadPressureEvent struct{ WindowEvent }

func (WindowTouchpadPressureEvent) EventType() EventType {
	return WindowEventTouchpadPressure
}

type WindowAxisMotionEvent struct{ WindowEvent }

func (WindowAxisMotionEvent) EventType() EventType { return WindowEventAxisMotion }

type WindowTouchEvent struct {
	WindowEvent
	DeviceID int64
	Phase    TouchPhase
	X        float64
	Y        float64
	ID       uint64
}

func (WindowTouchEvent) EventType() EventType { return WindowEventTouch }

type WindowScaleFactorChangedEvent struct {
	WindowEvent
	ScaleFactor float64
	Width       uint32
	Height      uint32
}

func (WindowScaleFactorChangedEvent) EventType() EventType {
	return WindowEventScaleFactorChanged
}

type NewEventsEvent struct{}

func (NewEventsEvent) EventType() EventType { return NewEvents }

type MainEventsClearedEvent struct{}

func (MainEventsClearedEvent) EventType() EventType { return MainEventsCleared }

type LoopDestroyedEvent struct{}

func (LoopDestroyedEvent) EventType() EventType { return LoopDestroyed }

type SuspendedEvent struct{}

func (SuspendedEvent) EventType() EventType { return Suspended }

type ResumedEvent struct{}

func (ResumedEvent) EventType() EventType { return Resumed }

type RedrawRequestedEvent struct{ WindowEvent }

func (RedrawRequestedEvent) EventType() EventType { return RedrawRequested }

type RedrawEventsClearedEvent struct{}

func (RedrawEventsClearedEvent) EventType() EventType { return RedrawEventsCleared }

type ModifiersChangedEvent struct {
	WindowEvent
	Modifiers ModifiersState
}

func (ModifiersChangedEvent) EventType() EventType { return ModifiersChanged }

type UserEventValue struct {
	Kind  UserEventKind
	Value any
}

func (UserEventValue) EventType() EventType { return UserEvent }

// UnknownEvent preserves an event ID not known by this version of the binding.
// Raw is reserved for a future ABI-preserving payload copy.
type UnknownEvent struct {
	Type     EventType
	WindowID WindowID
	Raw      []byte
}

func (e UnknownEvent) EventType() EventType { return e.Type }

type TouchPhase uint32

const (
	TouchPhaseUnknown TouchPhase = iota
	TouchPhaseStarted
	TouchPhaseMoved
	TouchPhaseEnded
	TouchPhaseCancelled
)

type MouseScrollDeltaType uint32

const (
	MouseScrollDeltaUnknown MouseScrollDeltaType = iota
	MouseScrollDeltaLine
	MouseScrollDeltaPixel
)

type MouseScrollDelta struct {
	Type MouseScrollDeltaType
	X    float64
	Y    float64
}

type InputElementState uint32

const (
	InputElementStateUnknown InputElementState = iota
	InputElementStatePressed
	InputElementStateReleased
)

type MouseButtonType uint32

const (
	MouseButtonUnknown MouseButtonType = iota
	MouseButtonLeft
	MouseButtonRight
	MouseButtonMiddle
	MouseButtonOther
)

type MouseButton struct {
	Type MouseButtonType
	Code uint16
}

type ModifiersState struct {
	Shift bool
	Ctrl  bool
	Alt   bool
	Logo  bool
}

type UserEventKind uint32

const (
	UserEventWakeUp UserEventKind = iota
	UserEventExit
)
