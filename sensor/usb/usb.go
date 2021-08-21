package usb

import (
	"buty4649/omusoba/sensor"
)

const (
	cmdRead       = 0x01
	cmdWrite      = 0x02
	cmdReadError  = 0x81
	cmdWriteError = 0x82
	cmdUnknown    = 0xff
)

type UsbSensor struct {
	serial SerialPort
}

func New(path string) UsbSensor {
	return UsbSensor{serial: &serialPort{path: path}}
}

func (us *UsbSensor) FetchLatestData() (*sensor.SensorData, error) {
	data, err := us.serial.Send(cmdRead, 0x5021, nil)
	if err != nil {
		return nil, err
	}

	return &sensor.SensorData{
		Sequence:           uint8(data[0]),
		Temperature:        float32(bytesToInt16(data[1:3])) * 0.01,
		Humidity:           float32(bytesToInt16(data[3:5])) * 0.01,
		AmbientLight:       bytesToInt16(data[5:7]),
		BarometricPressure: float32(bytesToInt32(data[7:11])) * 0.001,
		SoundNoise:         float32(bytesToInt16(data[11:13])) * 0.01,
		Tvoc:               bytesToInt16(data[13:15]),
		Co2:                bytesToInt16(data[15:17]),
		DiscomfortIndex:    float32(bytesToInt16(data[17:19])) * 0.01,
		HeatStroke:         float32(bytesToInt16(data[19:21])) * 0.01,
		Vibration:          uint8(data[21]),
		SI:                 float32(bytesToUint16(data[22:24])) * 0.1,
		PGA:                float32(bytesToUint16(data[24:26])) * 0.1,
		SeismicIntensity:   float32(bytesToUint16(data[26:28])) * 0.001,
	}, nil
}

func (us *UsbSensor) Info() (*sensor.DeviceInfo, error) {
	data, err := us.serial.Send(cmdRead, 0x180a, nil)
	if err != nil {
		return nil, err
	}

	return &sensor.DeviceInfo{
		Model:       string(data[0:10]),
		Serial:      string(data[10:20]),
		Firmware:    string(data[20:25]),
		Hardware:    string(data[25:30]),
		Manufacture: string(data[30:]),
	}, nil
}

func (us *UsbSensor) Mode() (*sensor.Mode, error) {
	data, err := us.serial.Send(cmdRead, 0x5117, nil)
	if err != nil {
		return nil, err
	}

	return &sensor.Mode{Mode: uint8(data[0])}, nil
}

func (us *UsbSensor) RawSend(cmd, addr int, data []byte) ([]byte, error) {
	return us.serial.Send(cmd, addr, data)
}
