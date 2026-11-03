//go:build windows && amd64

package libwinit

// WindowAttributesWithDimensionsFloat64 is the ABI-correct float64 entry
// point. The generated Hike stub currently lowers float64 arguments as
// integer words, while the Rust C ABI expects them in XMM registers.
func WindowAttributesWithDimensionsFloat64(windowAttributes *byte, width, height float64)
