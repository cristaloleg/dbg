//go:build !nodebug

package dbg

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

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
func Watch(labels ...any) func() time.Duration {
	// TODO: labels
	caller := Caller(2)
	start := time.Now()

	return func() time.Duration {
		took := time.Since(start)
		// TODO: add histogram
		fmt.Fprintln(output, caller, "took:", took.String())
		return took
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

	if x := rand.Float64(); x < prob {
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

var intervalMap sync.Map

// Interval run fn no often than the given interval.
func Interval(x time.Duration, fn func(last int64)) {
	last := get(&intervalMap, Location(2), new(int64))
	now := time.Now().UTC().UnixNano()

	lastTs := *last
	if lastTs == 0 || lastTs+x.Nanoseconds() <= now {
		if atomic.CompareAndSwapInt64(last, lastTs, now) {
			fn(*last)
		}
	}
}

func get[K any, V any](m *sync.Map, key K, def V) V {
	val, _ := m.LoadOrStore(key, def)
	return val.(V)
}
