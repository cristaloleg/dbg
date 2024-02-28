package dbg

// Sink any value like it's used.
// Treat it as `_ = x` or `_, _, ... = x, y, ...`.
func Sink(values ...any) {}
