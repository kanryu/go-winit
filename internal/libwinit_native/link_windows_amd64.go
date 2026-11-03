//go:build windows && amd64

package libwinit_native

// #cgo windows,amd64 LDFLAGS: -lntdll -ladvapi32 -luser32 -lgdi32 -ldwmapi -luxtheme -luserenv -lkernel32 -lole32 -loleaut32 -lshell32 -lshlwapi -lcomdlg32 -limm32 -lws2_32 -lruntimeobject -lucrt -lucrtbase -lm -lgcc -lgcc_eh -lmingwex -lmingw32
import "C"

// libwinit is built for the Windows GNU target. Its archive leaves the
// Windows, Winsock, COM and MinGW runtime libraries to the final external
// link. This cgo-only helper carries those flags; the public Go binding does
// not expose cgo.
