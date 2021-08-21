package usb

import (
	"fmt"
	"time"

	"go.bug.st/serial"
)

type SerialPort interface {
	Send(int, int, []byte) ([]byte, error)
}

type serialPort struct {
	path string
}

func (sp *serialPort) Send(command int, address int, data []byte) ([]byte, error) {
	port, err := serial.Open(sp.path, &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		return nil, err
	}
	defer port.Close()
	port.SetReadTimeout(100 * time.Millisecond)

	message := makeDataFrame(uint8(command), uint16(address), data)
	_, err = port.Write(message)
	if err != nil {
		return nil, err
	}
	time.Sleep(25 * time.Millisecond)

	buffer := make([]byte, 128)
	n, err := port.Read(buffer)
	if err != nil {
		return nil, err
	}

	result, err := extractResponse(buffer, n)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func makeDataFrame(command uint8, address uint16, data []byte) []byte {
	length := 5 // command(1byte) + address(2byte) + crc16(2byte)
	if data != nil {
		length += len(data)
	}

	result := []byte{
		0x52, 0x42, // "R" "B"
		byte(length & 0xff), byte(length >> 8),
		command,
		byte(address & 0xff), byte(address >> 8),
	}
	if data != nil {
		result = append(result, data...)
	}

	crc := crc16(result)
	result = append(result, byte(crc&0xff), byte(crc>>8))

	return result
}

func crc16(data []byte) uint16 {
	crc := 0xffff

	for _, v := range data {
		crc ^= int(v)
		for i := 0; i < 8; i++ {
			flag := crc & 0x01
			crc >>= 1
			if flag == 1 {
				crc ^= 0xA001
			}
		}
	}

	return uint16(crc)
}

func extractResponse(response []byte, length int) ([]byte, error) {
	crc := uint16(response[length-2]) + (uint16(response[length-1]) << 8)
	expected := crc16(response[:length-2])

	if crc != expected {
		return nil, fmt.Errorf("Missmatch CRC: %04x != %04x", crc, expected)
	}

	command := response[4]
	data := response[7 : length-2]

	if command == 0x81 || command == 0x82 || command == 0xff {
		kind := map[byte]string{0x81: "Read", 0x82: "Write", 0xff: "Unknown"}[command]
		reason := map[byte]string{
			0x01: "CRC", 0x02: "Command", 0x03: "Address",
			0x04: "Length", 0x05: "Data", 0x06: "Busy",
		}[data[0]]
		return nil, fmt.Errorf("%s error: %s", kind, reason)
	}

	return data, nil
}
