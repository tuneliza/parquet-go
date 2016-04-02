package bitpacking

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func bin(values ...string) (out []byte) {
	for _, v := range values {
		v = strings.Replace(v, " ", "", -1)

		b, err := strconv.ParseUint(v, 2, 8)
		if err != nil {
			panic(err)
		}

		out = append(out, byte(b))
	}

	return
}

var testcases = []struct {
	bitWidth uint
	input    []int32
	output   []byte
}{
	// with one bit you can encode 2 values
	{1, []int32{1}, bin("1")},
	{1, []int32{1, 1}, bin("11")},
	{1, []int32{1, 1, 1}, bin("111")},
	{1, []int32{0, 1, 1, 1}, bin("1110")},
	{1, []int32{1, 0, 1, 1, 1}, bin("11101")},

	{1, []int32{1, 1, 1, 1,
		1, 1, 1, 1, 1}, bin("1111 1111", "1")},

	// with two bit you can encode 4 values
	{2, []int32{0, 1, 2, 3}, bin("11 10 01 00")},
	{2, []int32{0, 1, 2, 3,
		0, 3, 3, 3}, bin("11 10 01 00", "11 11 11 00")},

	// sample documentation case
	{3, []int32{0, 1, 2, 3, 4, 5, 6, 7},
		bin("10001000", "11000110", "11111010")},

	{8, []int32{0, 1, 2, 4, 8, 16, 32, 64, 128},
		bin("0", "1", "10", "100", "1000",
			"1 0000", "10 0000", "100 0000", "1000 0000")},
}

func TestEncoding(t *testing.T) {
	for idx, tc := range testcases {
		var w bytes.Buffer
		enc := NewEncoder(tc.bitWidth, RLE)

		if _, err := enc.Write(&w, tc.input); err != nil {
			t.Fatalf("write: %s", err)
		}

		if bytes.Equal(w.Bytes(), tc.output) == false {
			t.Fatalf("%d: %#v != %#v", idx, w.Bytes(), tc.output)
		}
	}
}

func TestDecoding(t *testing.T) {
	for idx, tc := range testcases {
		dec := NewDecoder(tc.bitWidth)
		out := make([]int32, 8)
		if err := dec.Read(bytes.NewReader(tc.output), out); err != nil {
			t.Errorf("%d: %s", idx, err)
		}

		if !reflect.DeepEqual(out, tc.input) {
			t.Logf("%v != %v", out, tc.input)
		}
	}
}
