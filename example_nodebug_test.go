//go:build nodebug

package dbg_test

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/cristaloleg/dbg"
)

var testBuf bytes.Buffer

func init() {
	dbg.SetOutput(&testBuf)
}

func ExampleSink() {
	x, y := 123, 456
	// TODO: commented for now
	// x = y / 0

	dbg.Sink(x, y)

	// Output:
}

func ExampleWatch() {
	defer dbg.Watch()()
	func() {
		defer dbg.Watch()()

		time.Sleep(time.Second)
	}()

	fmt.Println(testBuf.String())

	// Output:
}

func ExampleHit() {
	for i := 0; i < 10; i++ {
		dbg.Hit()
		if i%2 == 0 {
			dbg.Hit()
		}
	}

	dbg.PrintHits()

	fmt.Println(testBuf.String())

	// Output:
}

func ExampleOnce() {
	for i := 0; i < 10; i++ {
		dbg.Once(func() { fmt.Println("in loop") })
	}

	// Output:
}

func ExamplePrintOnce() {
	for i := 0; i < 10; i++ {
		dbg.PrintOnce("debuging")

		go func() {
			_ = "noop"
		}()
	}

	fmt.Println(testBuf.String())

	// Output:
}

func Example_onceButTwice() {
	fn := func() {
		fmt.Println("I'm printed twice!")
	}

	dbg.Once(fn)
	dbg.Once(fn)

	// Output:
}

func ExampleRarely() {
	var counter int
	for i := 0; i < 1000; i++ {
		dbg.Rarely(0.1, func(count int64) { counter++ })
	}

	fmt.Println(counter == 0)

	// Output:
	// true
}

func ExampleEvery() {
	var counter int
	for i := 0; i < 1000; i++ {
		dbg.Every(10, func(count int64) { counter++ })
	}

	fmt.Println(counter == 0)

	// Output:
	// true
}

func ExampleCallers() {
	var callers []string
	f1 := func() {
		callers = dbg.Callers(20)
	}
	f2 := func() { f1() }
	f3 := func() { f2() }
	f4 := func() { f3() }

	f4()

	fmt.Println(strings.Join(callers, "\n"))

	// Output:
	// github.com/cristaloleg/dbg_test.ExampleCallers.func1
	// github.com/cristaloleg/dbg_test.ExampleCallers.func2
	// github.com/cristaloleg/dbg_test.ExampleCallers.func3
	// github.com/cristaloleg/dbg_test.ExampleCallers.func4
	// github.com/cristaloleg/dbg_test.ExampleCallers
	// testing.runExample
	// testing.runExamples
	// testing.(*M).Run
	// main.main
	// runtime.main
	// runtime.goexit
}

func ExampleX() {
	foo := func(a ...any) any {
		return a[0]
	}

	res := foo(dbg.X(123), 456, 789)
	fmt.Println(res)
	fmt.Println(testBuf.String() == "")

	// Output:
	// 123
	// true
}

func ExampleDump() {
	offset := struct {
		TxName   string
		idx      uint64
		deadline uint64
	}{
		TxName:   "Final",
		idx:      34,
		deadline: 16000000000,
	}
	body := "txBody%1"
	hashCode := uint64(9487746)
	codeIsValid := false

	dbg.Dump("Tx commit ", offset, body, hashCode, codeIsValid)

	fmt.Println(testBuf.String() == "")

	// Output:
}
