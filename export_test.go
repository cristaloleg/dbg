package dbg

import "io"

func SetTestOutput(w io.Writer) {
	output = w
}
