package dbg

import (
	"fmt"
	"io"
	"math/rand"
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

	// NOTE: not needed after Go 1.20
	rnd = rand.New(rand.NewSource(time.Now().Unix()))
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

// Hit increment counter for the given line.
// See PrintHits to print collected hits.
func Hit() {
	counter := get(&hitMap, Location(2), new(int64))
	atomic.AddInt64(counter, 1)
}

// PrintHits collected at the moment of call.
func PrintHits() {
	hitMap.Range(func(key, value any) bool {
		fmt.Fprintln(output, key, atomic.LoadInt64(value.(*int64)))
		return true
	})
}

var onceMap sync.Map

// Once will run the given fn once on the line of call.
func Once(fn func()) {
	once := get(&onceMap, Location(2), new(sync.Once))
	once.Do(fn)
}

// PrintOnce the given string.
func PrintOnce(s string) {
	once := get(&onceMap, Location(2), new(sync.Once))
	once.Do(func() {
		fmt.Fprintln(output, s)
	})
}

var rarelyMap sync.Map

// Rarely run fn with a given probability.
func Rarely(prob float64, fn func(count int64)) {
	counter := get(&rarelyMap, Location(2), new(int64))
	done := atomic.AddInt64(counter, 1)

	if x := rnd.Float64(); x < prob {
		fn(done)
	}
}

var everyMap sync.Map

// Every x calls run fn.
func Every(x int64, fn func(count int64)) {
	counter := get(&everyMap, Location(2), new(int64))
	done := atomic.AddInt64(counter, 1)

	if done > 0 && done%x == 0 {
		fn(done)
	}
}

// Sink any value like it's used.
// Treat it as `_ = x` or `_, _, ... = x, y, ...`.
func Sink(values ...any) {}

// Caller of the function but with a skipped callers in-between.
// If caller cannot be detected - Location(skip) is returned.
func Caller(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	details := runtime.FuncForPC(pc)
	if !ok || details == nil {
		return "<UNKNOWN:0>"
	}

	name := details.Name()
	idx := strings.LastIndexByte(name, '/')
	if idx != -1 {
		name = name[idx+1:]
	}
	return name
}

// Callers returns stack of callers no deeper than a given value.
func Callers(depth int) []string {
	fpcs := make([]uintptr, depth)
	if runtime.Callers(2, fpcs) == 0 {
		return nil
	}

	callers := make([]string, 0, len(fpcs))
	for _, p := range fpcs {
		caller := runtime.FuncForPC(p - 1)
		if caller != nil {
			callers = append(callers, caller.Name())
		}
	}
	return callers
}

// Location of the function caller but with a skipped callers in-between.
func Location(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "<UNKNOWN:0>"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func get[K any, V any](m *sync.Map, key K, def V) V {
	val, _ := m.LoadOrStore(key, def)
	return val.(V)
}
