//go:build nodebug

package dbg

import (
	"time"
)

func Do(fn func())                            {}
func When(x bool, fn func())                  {}
func Want(cond bool, format string, a ...any) {}
func Watch(labels ...any) func() time.Duration {
	return func() time.Duration { return 0 }
}
func Hit()                                          {}
func PrintHits()                                    {}
func Once(fn func())                                {}
func PrintOnce(s string)                            {}
func First(x int64, fn func(count int64))           {}
func Rarely(prob float64, fn func(count int64))     {}
func Every(x int64, fn func(count int64))           {}
func Interval(x time.Duration, fn func(last int64)) {}
