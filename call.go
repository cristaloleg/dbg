package dbg

import (
	"fmt"
	"runtime"
	"strings"
)

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

// Location of the function caller but with a skipped callers in-between.
func Location(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "<UNKNOWN:0>"
	}
	return fmt.Sprintf("%s:%d", file, line)
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

// DumpGoroutines returns stacktrace for the current goroutine or all of them.
func DumpGoroutines(all bool) string {
	buf := make([]byte, 4096)
	for {
		n := runtime.Stack(buf, all)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}
