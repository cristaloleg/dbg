# dbg

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]
[![version-img]][version-url]

Debug helpers for Go.

## Acknowledge

Thanks to [Bohdan Storozhuk](https://github.com/storozhukbm) for [storozhukBM/dump](https://github.com/storozhukBM/dump) which inspired `dbg.Dump`.

## Features

* Simple.
* Useful.
* Really helps.

See [these docs][pkg-url] or [GUIDE.md](GUIDE.md) for more details.

## Install

Go version 1.19+

```
go get github.com/cristaloleg/dbg
```

## Example

```go
func main() {
	dbg.Hit() // count how many times this line was executed

	dbg.PrintOnce("debuging") // print once per program run

	for i := 0; i < 1000; i++ {
		// +1 every 10 calls
		dbg.Every(10, func(count int64) { counter++ })
	}

	x := 123
	str := "striiiing"
	dbg.Dump("so far we have ", x, str)
	// will be printed:
	// so far we have x: `123`; str: `striiiing`
}
```

See examples: [example_test.go](example_test.go).

## License

[MIT License](LICENSE).

[build-img]: https://github.com/cristaloleg/dbg/workflows/build/badge.svg
[build-url]: https://github.com/cristaloleg/dbg/actions
[pkg-img]: https://pkg.go.dev/badge/cristaloleg/dbg
[pkg-url]: https://pkg.go.dev/github.com/cristaloleg/dbg
[reportcard-img]: https://goreportcard.com/badge/cristaloleg/dbg
[reportcard-url]: https://goreportcard.com/report/cristaloleg/dbg
[coverage-img]: https://codecov.io/gh/cristaloleg/dbg/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/cristaloleg/dbg
[version-img]: https://img.shields.io/github/v/release/cristaloleg/dbg
[version-url]: https://github.com/cristaloleg/dbg/releases
