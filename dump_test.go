//go:build !nodebug

package dbg_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cristaloleg/dbg"
)

func TestDump(t *testing.T) {
	var testBuf bytes.Buffer
	dbg.SetTestOutput(&testBuf)

	idx := 1
	str := "some data"
	dbg.Dump("init have ", idx, str)

	kv := map[string]float64{"x": 5.6, "y": 4.5}
	sli := []bool{true, false, false}
	dbg.Dump(kv, sli)

	type Some struct {
		Data         string
		privateValue []map[string]string
	}
	structVal := Some{
		Data:         "data string",
		privateValue: []map[string]string{{"k": "v"}, {"a": "b"}},
	}
	dbg.Dump(structVal, structVal.privateValue[0]["k"], structVal.Data)

	dbg.Dump("other structure, goes here", idx)

	dbg.Dump()

	dbg.Dump(
		idx,
	)

	for i := 0; i < 3; i++ {
		dbg.Dump("repeated. ", i)
	}

	output := testBuf.String()
	if len(output) == 0 {
		t.Fatal()
	}
	mustContain2(t, output, "[DEBUG] ")
	mustContain2(t, output, "dump_test.go:17: init have idx: `1`; str: `some data`")
	mustContain2(t, output, "dump_test.go:21: kv: `map[x:5.6 y:4.5]`; sli: `[true false false]`")
	mustContain2(t, output, "structVal: `{Data:data string privateValue:[map[k:v] map[a:b]]}`; structVal.privateValue[0][\"k\"]: `v`; structVal.Data: `data string`")
	mustContain2(t, output, "\"other structure, goes here\", idx: `other structure, goes here`; `1`")
	mustContain2(t, output, "error parsing file ")
	mustContain2(t, output, "dbg/dump_test.go:37 to dump 1 vars: target line is invalid. Dump should start with `dbg.Dump(` and end with `41`: 37")
	mustContain2(t, output, "repeated. i: `0`")
	mustContain2(t, output, "repeated. i: `1`")
	mustContain2(t, output, "repeated. i: `2`")
}

func mustContain2(tb testing.TB, have, want string) {
	tb.Helper()
	if !strings.Contains(have, want) {
		tb.Errorf("\nhave: %s\nwant: %s", have, want)
	}
}
