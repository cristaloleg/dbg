# Guide for dbg

## Watch how long function takes

For that we have `dbg.Watch` function. In most cases it will be used in defer statement like:

```go
func Debug1() {
	defer dbg.Watch()() // NOTE: double () !!!
	// ... code to debug
}

func Debug1() {
	// ... preparation code

	watch := dbg.Watch()
	// ... code to debug
	watch()

	// ... some other code
}
```

## Collect how many times line was executed

Function `dbg.Hit` is a simple counter exactly for the line where it's called. Can be called from different goroutines. 

Use `dbg.PrintHits` to print the result of all hits.

```go
func Debug() {
	for i := 0; i < 1000; i++ {
		go func() {
			if i % 3 == 0 {
				dbg.Hit()
			}
		}()
	}

	dbg.PrintHits()
}
```

## Call function once on the line

During debugging quite often you might need to define `sync.Once` and use it once (obviously). But defining it every time might be cumbersome.

To make this simpler we have `dbg.Once`, see:

```go
var parallel = os.Getenv("PARALLEL") == "true"

func Debug() {
	if parallel {
		dbg.Once(func() { /* do something once */ })
		// ...
		return
	}
	// ...
}
```

And quite often the only thing you need once is to print something, `dbg.PrintOnce` does exactly that:

```go
var parallel = os.Getenv("PARALLEL") == "true"

func Debug() {
	if parallel {
		dbg.PrintOnce("yes, parallel is enabled")
		// ...
		return
	}
	// ...
}
```

## Rare or repeating actions

There are cases when you want to do an action but not very often, `dbg.Rarely` can call your function with a given probability:

```go
func StringLen(input string) int {
	// 0.1 which is approximately 10% of all StringLen calls
	dbg.Rarely(0.1, func(count int64) {
		println("got input len:", len(input))
	})
	return len(input)
}
```

Similar function `dbg.Every` does the same but every X calls.

```go
func StringLen(input string) int {
	// every 12 StringLen calls do func
	dbg.Every(12, func(count int64) {
		println("got input len:", len(input))
	})
	return len(input)
}
```
