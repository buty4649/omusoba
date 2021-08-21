package usb

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestCrc16(t *testing.T) {
	data := []byte{0x52, 0x42, 0x05, 0x00, 0x01, 0x0a, 0x18}
	got := crc16(data)
	if got != 0x8dfc {
		t.Errorf("crc16(%v) = %04x; want 8dfc", data, got)
	}
}

func TestMakeDataFrame(t *testing.T) {
	got := makeDataFrame(cmdRead, 0x180a, nil)
	expect := []byte{0x52, 0x42, 0x05, 0x00, 0x01, 0x0a, 0x18, 0xfc, 0x8d}
	if !bytes.Equal(got, expect) {
		t.Errorf("makeDataFrame(%v) = %v; want %v", got, expect, expect)
	}
}

func TestExtractResponse(t *testing.T) {
	test := makeDataFrame(cmdRead, 0x0000, []byte("OK"))
	got, _ := extractResponse(test, len(test))
	if !bytes.Equal(got, []byte("OK")) {
		t.Errorf("extractResponse(%v) = %v; want %v", test, got, "OK")
	}

	test[len(test)-2] = 0xff
	test[len(test)-1] = 0xff
	_, err := extractResponse(test, len(test))
	if err == nil || !strings.Contains(err.Error(), "Missmatch CRC:") {
		t.Errorf("extractResponse unexpected CRC check")
	}

	for kind, value := range map[string]uint8{"Read": 0x81, "Write": 0x82, "Unknown": 0xff} {
		test = makeDataFrame(value, 0x0000, []byte("0x01"))
		_, err := extractResponse(test, len(test))
		if err == nil || !strings.Contains(err.Error(), fmt.Sprintf("%s error:", kind)) {
			t.Errorf("extractResponse unexpted error response: %s", kind)
		}
	}

	for kind, value := range map[string]byte{
		"CRC": 0x01, "Command": 0x02, "Address": 0x03,
		"Length": 0x04, "Data": 0x05, "Busy": 0x06,
	} {
		test = makeDataFrame(0xff, 0x0000, []byte{value})
		_, err := extractResponse(test, len(test))
		if err == nil || !strings.Contains(err.Error(), fmt.Sprintf("Unknown error: %s", kind)) {
			t.Errorf("extractResponse unexpted error code: %s", kind)
		}
	}
}
