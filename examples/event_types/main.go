package main

import (
	"log"
	"time"

	"github.com/kanryu/go-winit/winit"
)

func main() {
	eventLoop := winit.NewEventLoop()
	if eventLoop == nil {
		log.Fatal("winit: failed to create event loop")
	}
	defer eventLoop.Close()

	attributes := winit.NewWindowAttributes().
		WithTitle("Typed Events").
		WithDimensions(600, 400)
	window := eventLoop.CreateWindow(attributes)
	if window == nil {
		log.Fatal("winit: failed to create window")
	}
	defer window.Close()
	window.SetIMEAllowed(true)
	window.Focus()
	maximized := window.IsMaximized()
	proxy := eventLoop.CreateProxy()
	defer proxy.Close()
	go func() {
		time.Sleep(500 * time.Millisecond)
		proxy.Send(struct {
			Message string
		}{Message: "application event from a goroutine"})
	}()

	eventLoop.RunTyped(func(event winit.EventValue) winit.ControlFlow {
		switch event := event.(type) {
		case winit.WindowResizedEvent:
			log.Printf("window resized: %d x %d", event.Width, event.Height)
			currentMaximized := window.IsMaximized()
			if currentMaximized != maximized {
				log.Printf("window maximized: %t", currentMaximized)
				maximized = currentMaximized
			}
		case winit.WindowOccludedEvent:
			log.Printf("window occluded (minimized/hidden): %t", event.Occluded)
		case winit.WindowKeyboardInputEvent:
			state := "Unknown"
			switch event.State {
			case winit.InputElementStatePressed:
				state = "Pressed"
			case winit.InputElementStateReleased:
				state = "Released"
			}
			log.Printf("keyboard %s: key=%q physical=%t code=%d synthetic=%t",
				state, event.LogicalKey, event.HasPhysicalKey, event.PhysicalKeyCode, event.IsSynthetic)
		case winit.WindowReceivedCharacterEvent:
			log.Printf("received character: %q", event.Character)
		case winit.WindowCloseRequestedEvent:
			log.Println("window close requested")
			window.Close()
			return winit.Exit
		case winit.UserEventValue:
			log.Printf("user event: %#v", event.Value)
		}
		return winit.Wait
	})
}
