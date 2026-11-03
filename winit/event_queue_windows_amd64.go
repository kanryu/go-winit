//go:build windows && amd64

package winit

import (
	"runtime"
	"sync"
	"unsafe"

	native "github.com/kanryu/go-winit/internal/libwinit"
)

// userEventQueue is entirely Go-managed. The native proxy only wakes the
// event loop; application values never cross a C callback boundary.
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

func rawEventFromNative(raw *native.Event) Event {
	if raw == nil {
		return Event{}
	}
	e := Event{
		Type:           EventType(raw.EventType),
		WindowID:       [2]uint64{raw.WindowID.High, raw.WindowID.Low},
		ResizedWidth:   raw.WindowResized.Width,
		ResizedHeight:  raw.WindowResized.Height,
		CursorX:        raw.CursorMoved.X,
		CursorY:        raw.CursorMoved.Y,
		ScaleFactor:    raw.ScaleFactor.ScaleFactor,
		ScaleWidth:     raw.ScaleFactor.Width,
		ScaleHeight:    raw.ScaleFactor.Height,
		ReceivedLength: raw.ReceivedCharacter.Length,
		ReceivedBytes: [4]byte{
			raw.ReceivedCharacter.Byte1,
			raw.ReceivedCharacter.Byte2,
			raw.ReceivedCharacter.Byte3,
			raw.ReceivedCharacter.Byte4,
		},
		KeyboardState:   InputElementState(raw.KeyboardInput.State),
		HasPhysicalKey:  raw.KeyboardInput.HasPhysicalKeycode,
		PhysicalKeyCode: raw.KeyboardInput.PhysicalKeycode,
		IsSynthetic:     raw.KeyboardInput.IsSynthetic,
		Occluded:        raw.WindowOccluded.IsOccluded,
	}

	logicalKey := [...]byte{
		raw.KeyboardInput.LogicalKeyByte01, raw.KeyboardInput.LogicalKeyByte02,
		raw.KeyboardInput.LogicalKeyByte03, raw.KeyboardInput.LogicalKeyByte04,
		raw.KeyboardInput.LogicalKeyByte05, raw.KeyboardInput.LogicalKeyByte06,
		raw.KeyboardInput.LogicalKeyByte07, raw.KeyboardInput.LogicalKeyByte08,
		raw.KeyboardInput.LogicalKeyByte09, raw.KeyboardInput.LogicalKeyByte10,
		raw.KeyboardInput.LogicalKeyByte11, raw.KeyboardInput.LogicalKeyByte12,
		raw.KeyboardInput.LogicalKeyByte13, raw.KeyboardInput.LogicalKeyByte14,
		raw.KeyboardInput.LogicalKeyByte15, raw.KeyboardInput.LogicalKeyByte16,
		raw.KeyboardInput.LogicalKeyByte17, raw.KeyboardInput.LogicalKeyByte18,
		raw.KeyboardInput.LogicalKeyByte19, raw.KeyboardInput.LogicalKeyByte20,
		raw.KeyboardInput.LogicalKeyByte21, raw.KeyboardInput.LogicalKeyByte22,
		raw.KeyboardInput.LogicalKeyByte23, raw.KeyboardInput.LogicalKeyByte24,
		raw.KeyboardInput.LogicalKeyByte25, raw.KeyboardInput.LogicalKeyByte26,
		raw.KeyboardInput.LogicalKeyByte27, raw.KeyboardInput.LogicalKeyByte28,
		raw.KeyboardInput.LogicalKeyByte29, raw.KeyboardInput.LogicalKeyByte30,
		raw.KeyboardInput.LogicalKeyByte31, raw.KeyboardInput.LogicalKeyByte32,
		raw.KeyboardInput.LogicalKeyByte33, raw.KeyboardInput.LogicalKeyByte34,
		raw.KeyboardInput.LogicalKeyByte35, raw.KeyboardInput.LogicalKeyByte36,
		raw.KeyboardInput.LogicalKeyByte37, raw.KeyboardInput.LogicalKeyByte38,
		raw.KeyboardInput.LogicalKeyByte39, raw.KeyboardInput.LogicalKeyByte40,
		raw.KeyboardInput.LogicalKeyByte41, raw.KeyboardInput.LogicalKeyByte42,
		raw.KeyboardInput.LogicalKeyByte43, raw.KeyboardInput.LogicalKeyByte44,
		raw.KeyboardInput.LogicalKeyByte45, raw.KeyboardInput.LogicalKeyByte46,
		raw.KeyboardInput.LogicalKeyByte47, raw.KeyboardInput.LogicalKeyByte48,
		raw.KeyboardInput.LogicalKeyByte49, raw.KeyboardInput.LogicalKeyByte50,
		raw.KeyboardInput.LogicalKeyByte51, raw.KeyboardInput.LogicalKeyByte52,
		raw.KeyboardInput.LogicalKeyByte53, raw.KeyboardInput.LogicalKeyByte54,
		raw.KeyboardInput.LogicalKeyByte55, raw.KeyboardInput.LogicalKeyByte56,
		raw.KeyboardInput.LogicalKeyByte57, raw.KeyboardInput.LogicalKeyByte58,
		raw.KeyboardInput.LogicalKeyByte59, raw.KeyboardInput.LogicalKeyByte60,
		raw.KeyboardInput.LogicalKeyByte61, raw.KeyboardInput.LogicalKeyByte62,
		raw.KeyboardInput.LogicalKeyByte63, raw.KeyboardInput.LogicalKeyByte64,
	}
	length := raw.KeyboardInput.LogicalKeyLength
	if length > uint64(len(logicalKey)) {
		length = uint64(len(logicalKey))
	}
	e.LogicalKey = string(logicalKey[:length])
	return e
}

func (e *EventLoop) runGoEvents(callback Callback, typed TypedCallback) {
	for e != nil && e.ptr != nil {
		native.EventLoopRunReturn(e.ptr)
		for {
			raw := native.EventQueuePop()
			if raw == nil {
				break
			}
			event := rawEventFromNative(raw)
			native.EventDrop((*byte)(unsafe.Pointer(raw)))
			runtime.KeepAlive(raw)

			flow := Wait
			if typed != nil {
				if event.Type == UserEvent && e.userEvents != nil {
					for _, value := range e.userEvents.drain() {
						if typed(UserEventValue{Value: value}) == Exit {
							return
						}
					}
					continue
				}
				flow = typed(eventValueFromRaw(event))
			} else if callback != nil {
				flow = callback(event)
			}
			if flow == Exit {
				return
			}
		}
	}
}
