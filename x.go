//go:build !nodebug

package dbg

// X will print value and return it. Read `X` as `X-ray`.
// Can be used in any expression:
//
//	req, err := http.Post(dbg.X(url), "localhost:8080/ping", http.NoBody)
func X[T any](value T) T {
	if isDisabled() {
		return value
	}
	loc := Location(1)
	debug("%s: %v", loc, value)
	return value
}
