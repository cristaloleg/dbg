package dbg_test

import (
	"bytes"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cristalhq/dbg"
)

var testBuf bytes.Buffer

func init() {
	dbg.SetOutput(&testBuf)
}

func ExampleWatch() {
	defer pleaseIgnoreThisFuncCall()

	defer dbg.Watch()()
	func() {
		defer dbg.Watch()()

		time.Sleep(time.Second)
	}()

	output := testBuf.String()
	mustContain(output, "dbg_test.ExampleWatch")
	mustContain(output, "dbg_test.ExampleWatch.func1")

	// Output:
}

func ExampleHit() {
	defer pleaseIgnoreThisFuncCall()
	defer dbg.PrintHits()

	for i := 0; i < 10; i++ {
		dbg.Hit()
		if i%2 == 0 {
			dbg.Hit()
		}
	}

	dbg.PrintHits()

	output := testBuf.String()
	mustContain(output, "example_test.go:40 10")
	mustContain(output, "example_test.go:42 5")

	// Output:
}

func ExamplePrintOnce() {
	defer pleaseIgnoreThisFuncCall()

	for i := 0; i < 10; i++ {
		dbg.PrintOnce("debuging")

		go func() {
			_ = "noop"
		}()
	}

	fmt.Println(testBuf.String())

	// Output:
	// debuging
}

func mustContain(s, substr string) {
	if !strings.Contains(s, substr) {
		panic("does not contain")
	}
}

func pleaseIgnoreThisFuncCall() {
	testBuf.Reset()
}
