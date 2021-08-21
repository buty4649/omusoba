package usb

import (
	"reflect"
	"testing"

	mock "buty4649/omusoba/sensor/usb/mock"

	"github.com/golang/mock/gomock"
)

type testData struct {
	name   string
	expect interface{}
}

func TestFetchLastData(t *testing.T) {
	c, us := setup(t, func(e *mock.MockSerialPortMockRecorder) {
		v := []byte{
			0x01, 0x89, 0x09, 0x7c, 0x17, 0x16, 0x01, 0x71,
			0x7b, 0x0f, 0x00, 0x82, 0x1c, 0x5b, 0x00, 0xe9,
			0x03, 0x20, 0x1c, 0x9d, 0x08, 0x02, 0x4c, 0x00,
			0x3a, 0x07, 0x3c, 0x10, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00,
		}
		e.Send(cmdRead, 0x5021, nil).Return(v, nil).Times(1)
	})
	defer c.Finish()

	resp, err := us.FetchLatestData()

	tests := []testData{
		{"Sequence", uint8(0x01)},
		{"Temperature", float32(24.41)},
		{"Humidity", float32(60.12)},
		{"AmbientLight", int16(278)},
		{"BarometricPressure", float32(1014.64105)},
		{"SoundNoise", float32(72.979996)},
		{"Tvoc", int16(91)},
		{"Co2", int16(1001)},
		{"DiscomfortIndex", float32(72)},
		{"HeatStroke", float32(22.05)},
		{"Vibration", uint8(2)},
		{"SI", float32(7.6)},
		{"PGA", float32(185)},
		{"SeismicIntensity", float32(4.156)},
	}
	check(t, *resp, err, tests)
}

func TestInfo(t *testing.T) {
	c, us := setup(t, func(e *mock.MockSerialPortMockRecorder) {
		v := []byte("2JCIE-BU010000MY000000.0099.99OMRON")
		e.Send(cmdRead, 0x180a, nil).Return(v, nil).Times(1)
	})
	defer c.Finish()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resp, err := us.Info()

	tests := []testData{
		{"Model", "2JCIE-BU01"},
		{"Serial", "0000MY0000"},
		{"Firmware", "00.00"},
		{"Hardware", "99.99"},
		{"Manufacture", "OMRON"},
	}
	check(t, *resp, err, tests)
}

func TestMode(t *testing.T) {
	c, us := setup(t, func(e *mock.MockSerialPortMockRecorder) {
		v := []byte{0x00}
		e.Send(cmdRead, 0x5117, nil).Return(v, nil).Times(1)
	})
	defer c.Finish()

	resp, err := us.Mode()
	tests := []testData{
		{"Mode", uint8(0x00)},
	}
	check(t, *resp, err, tests)
}

func TestRawSend(t *testing.T) {
	c, us := setup(t, func(e *mock.MockSerialPortMockRecorder) {
		e.Send(0x01, 0x5117, nil).Return([]byte{0x00}, nil).Times(1)
	})
	defer c.Finish()

	d, err := us.RawSend(0x01, 0x5117, nil)
	if err != nil {
		t.Errorf("err is always not nil; got %v", err)
	}

	expect := []byte{0x00}
	if got := d; !reflect.DeepEqual(got, expect) {
		t.Errorf("got: %v; expect %v", got, expect)
	}
}

func setup(t *testing.T, fn func(*mock.MockSerialPortMockRecorder)) (*gomock.Controller, *UsbSensor) {
	ctrl := gomock.NewController(t)
	mock := mock.NewMockSerialPort(ctrl)
	fn(mock.EXPECT())

	us := &UsbSensor{}
	us.serial = mock

	return ctrl, us
}

func check(t *testing.T, value interface{}, err error, tests []testData) {

	if err != nil {
		t.Errorf("err is always not nil; got %v", err)
	}

	v := reflect.ValueOf(value)
	for _, tt := range tests {
		name := tt.name
		expect := tt.expect
		got := v.FieldByName(name).Interface()
		if got != expect {
			t.Errorf("unexpected %v; got: %v; expect %v", name, got, expect)
		}
	}

}
