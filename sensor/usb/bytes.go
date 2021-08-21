package usb

import "encoding/binary"

func bytesToUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func bytesToInt16(b []byte) int16 {
	return int16(bytesToUint16(b))
}

func bytesToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func bytesToInt32(b []byte) int32 {
	return int32(bytesToUint32(b))
}
