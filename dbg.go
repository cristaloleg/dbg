package dbg

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	output     io.Writer = os.Stdout
	outputOnce sync.Once
)

// SetOutput for the dbg package. Can be set once.
func SetOutput(w io.Writer) {
	outputOnce.Do(func() { output = w })
}

// Watch the function timing.
// The most popular usage is:
//
//	defer dbg.Watch(...)()
//
// However, can be like that:
//
//	watch := dbg.Watch(...)
//	...
//	watch()
//
// to call in a specific place.
func Watch(labels ...any) func() {
	// TODO: labels
	caller := Caller(2)
	start := time.Now()

	return func() {
		took := time.Since(start)
		// TODO: add histogram
		fmt.Fprintln(output, caller, "took:", took.String())
	}
}

var hitMap sync.Map

// Hit increases counter for the given line.
// See PrintHits to print collected hits.
func Hit(labels ...any) {
	// TODO: labels
	loc := Location(2)

	counter, _ := hitMap.LoadOrStore(loc, new(int32))
	atomic.AddInt32(counter.(*int32), 1)
}

// PrintHits collected at the moment of call.
func PrintHits() {
	hitMap.Range(func(key, value any) bool {
		fmt.Fprintln(output, key, *value.(*int32))
		return true
	})
}

var printOnceMap sync.Map

// PrintOnce the given string.
func PrintOnce(s string) {
	loc := Location(2)

	_, ok := printOnceMap.LoadOrStore(loc, 1)
	if !ok {
		fmt.Fprintln(output, s)
	}
}

// Caller of the function but with a skipped callers in-between.
// If caller cannot be detected - Location(skip) is returned.
func Caller(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	details := runtime.FuncForPC(pc)
	if !ok || details == nil {
		return Location(skip + 1)
	}

	name := details.Name()
	idx := strings.LastIndexByte(name, '/')
	if idx != -1 {
		name = name[idx+1:]
	}
	return name

}

// Location of the function caller but with a skipped callers in-between.
func Location(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "<UNKNOWN:0>"
	}
	return fmt.Sprintf("%s:%d", file, line)
}
