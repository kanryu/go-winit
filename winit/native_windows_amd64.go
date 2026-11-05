//go:build windows && amd64

package winit

import native "github.com/kanryu/go-winit/internal/libwinit"

func eventLoopNewNative() *byte { return native.EventLoopNew() }
func eventLoopDropNative(ptr *byte) { native.EventLoopDrop(ptr) }
func eventLoopRunReturnNative(ptr *byte) { native.EventLoopRunReturn(ptr) }
func windowSetInnerSizeNative(ptr *byte, width, height uint32) { native.WindowSetInnerSize(ptr, width, height) }
