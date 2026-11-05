//go:build linux && amd64

package libwinit_native

// #cgo linux,amd64 LDFLAGS: -lX11 -lwayland-client -lxkbcommon -ldl -lpthread -lm
import "C"

// The Linux libwinit archive is static and leaves the platform libraries to
// the final executable link. Keep these flags in a cgo-only file so the public
// Go API does not expose cgo requirements.
