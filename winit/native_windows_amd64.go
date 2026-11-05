//go:build windows && amd64

package winit

import native "github.com/kanryu/go-winit/internal/libwinit"

func eventLoopNewNative() *byte                        { return native.EventLoopNew() }
func eventLoopDropNative(ptr *byte)                    { native.EventLoopDrop(ptr) }
func eventLoopRunReturnNative(ptr *byte)               { native.EventLoopRunReturn(ptr) }
func windowAttributesNewNative() *byte                 { return native.WindowAttributesNew() }
func windowAttributesDropNative(ptr *byte)             { native.WindowAttributesDrop(ptr) }
func windowAttributesWithTitleNative(ptr, title *byte) { native.WindowAttributesWithTitle(ptr, title) }
func windowAttributesWithDecorationsNative(ptr *byte, value bool) {
	native.WindowAttributesWithDecorations(ptr, value)
}
func windowAttributesWithTransparencyNative(ptr *byte, value bool) {
	native.WindowAttributesWithTransparency(ptr, value)
}
func windowAttributesWithResizableNative(ptr *byte, value bool) {
	native.WindowAttributesWithResizable(ptr, value)
}
func windowAttributesWithMaximizedNative(ptr *byte, value bool) {
	native.WindowAttributesWithMaximized(ptr, value)
}
func windowAttributesWithVisibilityNative(ptr *byte, value bool) {
	native.WindowAttributesWithVisibility(ptr, value)
}
func windowAttributesWithAlwaysOnTopNative(ptr *byte, value bool) {
	native.WindowAttributesWithAlwaysOnTop(ptr, value)
}
func windowAttributesWithFullSizeNative(ptr *byte, value bool) {
	native.WindowAttributesWithFullSize(ptr, value)
}
func createWindowNative(eventLoop, attributes *byte) *byte {
	return native.CreateWindow(eventLoop, attributes)
}
func windowDropNative(ptr *byte)                      { native.WindowDrop(ptr) }
func windowSetIMEAllowedNative(ptr *byte, value bool) { native.WindowSetIMEAllowed(ptr, value) }
func windowSetInnerSizeNative(ptr *byte, width, height uint32) {
	native.WindowSetInnerSize(ptr, width, height)
}
func windowFocusWindowNative(ptr *byte)      { native.WindowFocusWindow(ptr) }
func windowIsMaximizedNative(ptr *byte) bool { return native.WindowIsMaximized(ptr) }
func eventQueuePopNative() *native.Event     { return native.EventQueuePop() }
func eventDropNative(ptr *byte)              { native.EventDrop(ptr) }
