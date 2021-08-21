package cmdutil

import "testing"

func TestMarshal(t *testing.T) {
	test := struct {
		A string `desc:"a" unit:"unit"`
	}{A: "bar"}

	got, err := marshal("label", test)
	if err != nil {
		t.Errorf("unexpected error: %v; got: %v", err, got)
	}

	expect := "a: bar unit\n"
	if got != expect {
		t.Errorf("unexpected marshal: label; got: %v, expected: %v", got, expect)
	}

	got, _ = marshal("yaml", test)
	expect = "a: bar\n"
	if got != expect {
		t.Errorf("unexpected marshal: yaml; got: %v, expected: %v", got, expect)
	}

	got, _ = marshal("json", test)
	expect = `{
  "A": "bar"
}
`
	if got != expect {
		t.Errorf("unexpected marshal: json; got: %v, expected: %v", got, expect)
	}

	uint_test := struct {
		Uint   uint   `desc:"uint"`
		Uint8  uint8  `desc:"uint8"`
		Uint16 uint16 `desc:"uint16"`
		Uint32 uint32 `desc:"uint32"`
	}{
		Uint: uint(100), Uint8: uint8(100), Uint16: uint16(100), Uint32: uint32(100),
	}
	got, _ = marshal("label", uint_test)
	expect = `  uint: 100
 uint8: 100
uint16: 100
uint32: 100
`
	if got != expect {
		t.Errorf("unexpected marshal: label; got: %v, expected: %v", got, expect)
	}

	int_test := struct {
		Int   int   `desc:"int"`
		Int8  int8  `desc:"int8"`
		Int16 int16 `desc:"int16"`
		Int32 int32 `desc:"int32"`
	}{
		Int: int(100), Int8: int8(100), Int16: int16(100), Int32: int32(100),
	}
	got, _ = marshal("label", int_test)
	expect = `  int: 100
 int8: 100
int16: 100
int32: 100
`
	if got != expect {
		t.Errorf("unexpected marshal: label; got: %v, expected: %v", got, expect)
	}

	float_test := struct {
		A float32 `desc:"a"`
		B float32 `desc:"b"`
	}{
		A: float32(0.1), B: float32(1.0),
	}
	got, _ = marshal("label", float_test)
	expect = "a: 0.1\nb: 1\n"
	if got != expect {
		t.Errorf("unexpected marshal: label; got: %v, expected: %v", got, expect)
	}

	slice_test := struct {
		A []byte `desc:"a"`
	}{[]byte{0x01, 0x02, 0x03}}
	got, _ = marshal("label", slice_test)
	expect = "a: 01 02 03\n"
	if got != expect {
		t.Errorf("unexpected marshal: label; got: %v, expected: %v", got, expect)
	}
}
