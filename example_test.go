package dbg_test

import (
	"bytes"
	"fmt"
	"strings"
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

	for i := 0; i < 10; i++ {
		dbg.Hit()
		if i%2 == 0 {
			dbg.Hit()
		}
	}

	dbg.PrintHits()

	output := testBuf.String()
	mustContain(output, "example_test.go")
	mustContain(output, "example_test.go")

	// Output:
}

func ExampleOnce() {
	for i := 0; i < 10; i++ {
		dbg.Once(func() { fmt.Println("in loop") })
	}

	// Output:
	// in loop
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

func ExampleOnceButTwice() {
	fn := func() {
		fmt.Println("I'm printed twice!")
	}

	dbg.Once(fn)
	dbg.Once(fn)

	// Output:
	// I'm printed twice!
	// I'm printed twice!
}

func ExampleRarely() {
	var counter int
	for i := 0; i < 1000; i++ {
		dbg.Rarely(0.1, func(count int64) { counter++ })
	}

	fmt.Println(counter < 150)

	// Output:
	// true
}

func ExampleEvery() {
	var counter int
	for i := 0; i < 1000; i++ {
		dbg.Every(10, func(count int64) { counter++ })
	}

	fmt.Println(counter)

	// Output:
	// 100
}

func mustContain(s, substr string) {
	if !strings.Contains(s, substr) {
		panic("does not contain")
	}
}

func pleaseIgnoreThisFuncCall() {
	testBuf.Reset()
}
