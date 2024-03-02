//go:build !nodebug

package dbg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

// Dump variables with their names without specifying them.
func Dump(vars ...any) {
	if isDisabled() {
		return
	}
	if len(vars) == 0 {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		debug("error capturing stack to dump %d vars", len(vars))
		return
	}

	val, ok := lineCache.Load(cacheKey{file: file, line: line})
	if ok {
		info := val.(lineInfo)
		dumpLine(file, line, info.targetLine, info.names, vars)
		return
	}

	info, err := parseLine(file, line)
	if err != nil {
		debug("error parsing file %s:%d to dump %d vars: %s", file, line, len(vars), err)
		return
	}
	lineCache.Store(cacheKey{file: file, line: line}, info)

	dumpLine(file, line, info.targetLine, info.names, vars)
}

type lineInfo struct {
	targetLine string
	names      []string
}

var lineCache sync.Map

type cacheKey struct {
	file string
	line int
}

func parseLine(file string, line int) (lineInfo, error) {
	codeline, err := loadLine(file, line)
	if err != nil {
		return lineInfo{}, err
	}

	const prefix = "dbg.Dump("
	const suffix = ')'

	codeline = strings.Trim(codeline, ` 	`)
	startIdx := strings.Index(codeline, prefix)
	endIdx := strings.LastIndexByte(codeline, suffix)

	if startIdx < 0 || endIdx < 0 {
		return lineInfo{}, fmt.Errorf(
			"target line is invalid. Dump should start with `%s` and end with `%v`: %v\n",
			prefix, suffix, line,
		)
	}
	codeline = codeline[startIdx+len(prefix) : endIdx]

	return lineInfo{
		targetLine: codeline,
		names:      strings.Split(codeline, ", "),
	}, nil
}

func loadLine(file string, line int) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("open file: %v", file)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for i := 1; scanner.Scan(); i++ {
		if i != line {
			continue
		}
		return scanner.Text(), nil
	}
	return "", errors.New("no such line")
}

func dumpLine(file string, line int, targetLine string, names []string, vars []any) {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "[DEBUG] %v:%v: ", file, line)
	if len(names) != len(vars) {
		fmt.Fprintf(buf, "%v: ", targetLine)
		for i, val := range vars {
			fmt.Fprintf(buf, "`%+v`", val)
			if i < len(vars)-1 {
				fmt.Fprintf(buf, "; ")
			}
		}
	} else {
		for i, v := range names {
			if isStrLit(v) {
				fmt.Fprintf(buf, "%v", v[1:len(v)-1])
			} else {
				fmt.Fprintf(buf, "%v: `%+v`", v, vars[i])
				if i < len(vars)-1 {
					fmt.Fprintf(buf, "; ")
				}
			}
		}
	}
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(output, buf.String())
}

func isStrLit(s string) bool {
	return strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) ||
		(strings.HasPrefix(s, "`") && strings.HasSuffix(s, "`"))
}
