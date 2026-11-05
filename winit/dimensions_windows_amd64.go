//go:build windows && amd64

package winit

import native "github.com/kanryu/go-winit/internal/libwinit"

func windowAttributesWithDimensionsNative(windowAttributes *byte, width, height float64) {
	native.WindowAttributesWithDimensionsFloat64(windowAttributes, width, height)
}
