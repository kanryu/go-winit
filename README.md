# go-winit

`go-winit` is a Go binding for [libwinit](https://github.com/kanryu/libwinit) on Windows.

In addition to exposing the libwinit C API, it provides Go callbacks, typed events, IME input, and user-event delivery from goroutines.

## About Winit

[Winit](https://github.com/rust-windowing/winit) is a cross-platform Rust library for window creation and management. It creates native windows and delivers events such as window resizing, keyboard input, mouse movement, focus changes, and window state changes.

Winit is intentionally a low-level building block. It does not provide a widget toolkit, a layout engine, or a renderer. Applications normally combine it with a graphics or UI library and use Winit's platform-specific handles when they need to draw into a window.

This repository exposes the functionality of the bundled [libwinit](https://github.com/kanryu/libwinit) C API to Go. The `winit` package keeps the native window and event-loop responsibilities while adding Go-oriented features such as typed event values, Go callbacks, IME-aware keyboard handling, and application-owned user events.

For the original Rust API and platform-specific details, see the [Winit documentation](https://docs.rs/winit) and the [`winit::platform` documentation](https://docs.rs/winit/latest/winit/platform/index.html).

## Receiving Events by Event ID

The low-level `Run` API dispatches events using libwinit event IDs.

```go
eventLoop.Run(func(event winit.Event) winit.ControlFlow {
	switch event.Type {
	case winit.WindowEventResized:
		log.Printf("resized: %d x %d", event.ResizedWidth, event.ResizedHeight)
	case winit.WindowEventKeyboardInput:
		log.Printf("key: %s", event.LogicalKey)
	case winit.WindowEventReceivedCharacter:
		log.Printf("character: %q", event.ReceivedCharacterString())
	case winit.WindowEventCloseRequested:
		window.Close()
		return winit.Exit
	}
	return winit.Wait
})
```

Common event IDs include:

- `WindowEventResized`
- `WindowEventMoved`
- `WindowEventCloseRequested`
- `WindowEventKeyboardInput`
- `WindowEventCursorMoved`
- `WindowEventReceivedCharacter`
- `WindowEventScaleFactorChanged`
- `WindowEventOccluded`
- `NewEvents`
- `MainEventsCleared`
- `UserEvent`

## Handling Events with Type Switches

`RunTyped` dispatches typed event structures, so application code does not need to branch on event IDs directly.

```go
eventLoop.RunTyped(func(event winit.EventValue) winit.ControlFlow {
	switch event := event.(type) {
	case winit.WindowResizedEvent:
		log.Printf("resized: %d x %d", event.Width, event.Height)
	case winit.WindowKeyboardInputEvent:
		log.Printf("key=%q state=%v", event.LogicalKey, event.State)
	case winit.WindowReceivedCharacterEvent:
		log.Printf("character: %q", event.Character)
	case winit.WindowOccludedEvent:
		log.Printf("minimized/hidden: %t", event.Occluded)
	case winit.WindowCloseRequestedEvent:
		window.Close()
		return winit.Exit
	case winit.UserEventValue:
		log.Printf("user event: %#v", event.Value)
	}
	return winit.Wait
})
```

Unknown event IDs are delivered as `UnknownEvent` values.

## Key Press and Key Release Events

`WindowKeyboardInputEvent` and the low-level `Event` contain the following keyboard information:

- `State`: `InputElementStatePressed` or `InputElementStateReleased`
- `LogicalKey`: a logical key name such as `Named(Enter)` or `Character("a")`
- `PhysicalKeyCode`: the physical key code
- `HasPhysicalKey`: whether a physical key code is available
- `IsSynthetic`: whether the event was synthesized

```go
case winit.WindowKeyboardInputEvent:
	switch event.State {
	case winit.InputElementStatePressed:
		log.Printf("pressed: %q", event.LogicalKey)
	case winit.InputElementStateReleased:
		log.Printf("released: %q", event.LogicalKey)
	}
```

## IME Input

IME input can be enabled for each window.

```go
window.SetIMEAllowed(true)
window.Focus()
```

Committed multibyte text is delivered as `WindowReceivedCharacterEvent` values.

```go
case winit.WindowReceivedCharacterEvent:
	log.Printf("IME commit: %q", event.Character)
```

The low-level API provides the same UTF-8 text through `ReceivedCharacterString`:

```go
case winit.WindowEventReceivedCharacter:
	log.Printf("IME commit: %q", event.ReceivedCharacterString())
```

While IME is enabled, ordinary `KeyboardInput` events generated during IME composition are automatically suppressed by libwinit. This prevents composition keystrokes from being interpreted as regular application commands; applications receive the committed IME text instead.

The actual IME on/off mode is controlled by the Windows input method. `SetIMEAllowed(true)` enables IME support for the target window.

## User Events

`EventLoopProxy` allows arbitrary Go values to be sent to the event loop from another goroutine.

```go
proxy := eventLoop.CreateProxy()
defer proxy.Close()

go func() {
	proxy.Send(struct {
		Message string
	}{Message: "event from goroutine"})
}()

eventLoop.RunTyped(func(event winit.EventValue) winit.ControlFlow {
	switch event := event.(type) {
	case winit.UserEventValue:
		log.Printf("user event: %#v", event.Value)
	}
	return winit.Wait
})
```

`Send` queues the value on the Go side and wakes the libwinit event loop. The user-event callback runs on the event-loop thread.

## Window State

Minimize and restore transitions are delivered as `WindowOccludedEvent` values. Maximization and restoration are delivered as resize events, so use `Window.IsMaximized()` to determine the current state.

```go
case winit.WindowResizedEvent:
	if window.IsMaximized() {
		log.Println("window maximized")
	}
case winit.WindowOccludedEvent:
	log.Printf("window occluded: %t", event.Occluded)
```

## Examples

- `examples/simple`: event-ID dispatch, IME, and keyboard input
- `examples/event_types`: typed event switches, IME, keyboard input, and user events

```powershell
go run ./examples/simple
go run ./examples/event_types
```

## Verification

```powershell
go test ./...
```
