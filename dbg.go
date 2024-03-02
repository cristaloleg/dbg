package dbg

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

var enabled int32 = 1

// IsEnabled reports whether debug is enabled.
func IsEnabled() bool { return !isDisabled() }

func isDisabled() bool { return atomic.LoadInt32(&enabled) == 0 }

// Enable debugging. Pass `true` to enable, pass `false` to disable.
func Enable(flag bool) {
	var v int32
	if flag {
		v = 1
	}
	atomic.StoreInt32(&enabled, v)
}

var (
	output     io.Writer = os.Stdout
	outputOnce sync.Once
)

// SetOutput for the dbg package. Can be set once.
func SetOutput(w io.Writer) {
	outputOnce.Do(func() { output = w })
}

// Sink any value like it's used.
// Treat it as `_ = x` or `_, _, ... = x, y, ...`.
func Sink(values ...any) {}

func debug(format string, args ...any) {
	fmt.Fprintf(output, "[DEBUG] "+format+"\n", args...)
}
