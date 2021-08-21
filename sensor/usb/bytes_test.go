package usb

import "testing"

func TestBytesToUint16(t *testing.T) {

	expect := uint16(513)
	got := bytesToUint16([]byte{0x01, 0x02})
	if got != expect {
		t.Errorf("got: %v; expect: %v", got, expect)
	}

	expect = uint16(65022)
	got = bytesToUint16([]byte{0xfe, 0xfd})
	if got != expect {
		t.Errorf("got: %v; expect: %v", got, expect)
	}
}

func TestBytesToInt16(t *testing.T) {

	expect := int16(513)
	got := bytesToInt16([]byte{0x01, 0x02})
	if got != expect {
		t.Errorf("got: %v; expect: %v", got, expect)
	}

	expect = int16(-514)
	got = bytesToInt16([]byte{0xfe, 0xfd})
	if got != expect {
		t.Errorf("got: %v; expect: %v", got, expect)
	}
}

func TestBytesToUint32(t *testing.T) {

	expect := uint32(67305985)
	got := bytesToUint32([]byte{0x01, 0x02, 0x03, 0x04})
	if got != expect {
		t.Errorf("got: %v; expect: %v", got, expect)
	}

	expect = uint32(4244504319)
	got = bytesToUint32([]byte{0xff, 0xfe, 0xfd, 0xfc})
	if got != expect {
		t.Errorf("got: %v; expect: %v", got, expect)
	}
}

func TestBytesToInt32(t *testing.T) {

	expect := int32(67305985)
	got := bytesToInt32([]byte{0x01, 0x02, 0x03, 0x04})
	if got != expect {
		t.Errorf("got: %v; expect: %v", got, expect)
	}

	expect = int32(-50462977)
	got = bytesToInt32([]byte{0xff, 0xfe, 0xfd, 0xfc})
	if got != expect {
		t.Errorf("got: %v; expect: %v", got, expect)
	}
}
