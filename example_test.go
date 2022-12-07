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
	if !strings.Contains(output, "dbg_test.ExampleWatch") {
		panic("do not have func")
	}
	if !strings.Contains(output, "dbg_test.ExampleWatch.func1") {
		panic("do not have anon-func")
	}

	// Output:
}

func ExamplePrintOnce() {
	defer pleaseIgnoreThisFuncCall()

	var count int32
	for i := 0; i < 10; i++ {
		dbg.PrintOnce("debuging")

		go func() {
			if i == 10 {
				atomic.AddInt32(&count, 1)
			}
		}()
	}

	fmt.Println(testBuf.String())

	// Output:
	// debuging
}

func pleaseIgnoreThisFuncCall() {
	testBuf.Reset()
}
