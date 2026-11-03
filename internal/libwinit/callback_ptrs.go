package libwinit

// Export the callback addresses emitted by hikec's Windows ABI stubs.
func EventLoopRunDataCallbackPtr() uintptr   { return _hike_EventLoopRunDataCallbackPtr() }
func EventLoopRunReturnCallbackPtr() uintptr { return _hike_EventLoopRunReturnCallbackPtr() }
func PollingSemaphoreCallbackPtr() uintptr   { return _hike_PollingSemaphoreCallbackPtr() }
func PollingMainEventsCallbackPtr() uintptr  { return _hike_PollingMainEventsCallbackPtr() }
func PollingResizeCallbackPtr() uintptr      { return _hike_PollingResizeCallbackPtr() }
func PollingRedrawCallbackPtr() uintptr      { return _hike_PollingRedrawCallbackPtr() }
