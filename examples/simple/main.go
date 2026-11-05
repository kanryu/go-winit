package main

import (
	"log"
	"os"

	"github.com/kanryu/go-winit/winit"
)

func main() {
	if os.Getenv("WINIT_UNIX_BACKEND") != "wayland" {
		_ = os.Setenv("WINIT_UNIX_BACKEND", "x11")
	}
	eventLoop := winit.NewEventLoop()
	if eventLoop == nil {
		log.Fatal("winit: failed to create event loop")
	}
	defer eventLoop.Close()

	// Window creation consumes the native attributes object, matching the
	// updated libwinit C API example.
	attributes := winit.NewWindowAttributes().
		WithTitle("Hello World").
		WithDimensions(600, 400)
	if attributes == nil {
		log.Fatal("winit: failed to create window attributes")
	}

	window := eventLoop.CreateWindow(attributes)
	if window == nil {
		log.Fatal("winit: failed to create window")
	}
	defer window.Close()
	window.SetSize(600, 400)
	window.SetIMEAllowed(true)
	window.Focus()
	maximized := window.IsMaximized()

	eventLoop.Run(func(event winit.Event) winit.ControlFlow {
		switch event.Type {
		case winit.WindowEventCloseRequested:
			log.Println("Window close requested")
			window.Close()
			return winit.Exit
		case winit.WindowEventResized:
			log.Printf("Window resized to: (%d x %d)\n", event.ResizedWidth, event.ResizedHeight)
			currentMaximized := window.IsMaximized()
			if currentMaximized != maximized {
				log.Printf("Window maximized: %t\n", currentMaximized)
				maximized = currentMaximized
			}
		case winit.WindowEventCursorMoved:
			log.Printf("Cursor moved to: (%f @ %f)\n", event.CursorX, event.CursorY)
		case winit.WindowEventReceivedCharacter:
			log.Printf("Received character: %q\n", event.ReceivedCharacterString())
		case winit.WindowEventKeyboardInput:
			state := "Unknown"
			switch event.KeyboardState {
			case winit.InputElementStatePressed:
				state = "Pressed"
			case winit.InputElementStateReleased:
				state = "Released"
			}
			log.Printf("Keyboard %s: key=%q physical=%t code=%d synthetic=%t\n",
				state, event.LogicalKey, event.HasPhysicalKey, event.PhysicalKeyCode, event.IsSynthetic)
		case winit.WindowEventScaleFactorChanged:
			log.Printf("Scale factor changed to: %f. New physical size: (%d x %d)\n",
				event.ScaleFactor, event.ScaleWidth, event.ScaleHeight)
		case winit.WindowEventOccluded:
			log.Printf("Window occluded (minimized/hidden): %t\n", event.Occluded)
		}
		return winit.Wait
	})
}
