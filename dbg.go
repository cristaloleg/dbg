package dbg

import (
	"io"
	"os"
	"sync"
)

var (
	output     io.Writer = os.Stdout
	outputOnce sync.Once
)

// SetOutput for the dbg package. Can be set once.
func SetOutput(w io.Writer) {
	outputOnce.Do(func() { output = w })
}
