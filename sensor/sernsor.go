package sensor

type SensorData struct {
	Sequence           uint8   `desc:"Sequence Number"`
	Temperature        float32 `desc:"Temperature" unit:"°C"`
	Humidity           float32 `desc:"Relative Humidity" unit:"%RH"`
	AmbientLight       int16   `desc:"Ambient Light" unit:"lx"`
	BarometricPressure float32 `desc:"Barometric Pressure" unit:"hPa"`
	SoundNoise         float32 `desc:"Sound Noise" unit:"dB"`
	Tvoc               int16   `desc:"eTVOC" unit:"ppb"`
	Co2                int16   `desc:"eCO2" unit:"ppm"`
	DiscomfortIndex    float32 `desc:"Discomfort Index"`
	HeatStroke         float32 `desc:"Heat Stroke" unit:"°C"`
	Vibration          uint8   `desc:"Vibration Information"`
	SI                 float32 `desc:"Spectral Intensity" unit:"kine"`
	PGA                float32 `desc:"Peak Ground Acceleration" unit:"gal"`
	SeismicIntensity   float32 `desc:"Seismic Intensity"`
}

type DeviceInfo struct {
	Model       string `desc:"Model Number"`
	Serial      string `desc:"Serial Number"`
	Firmware    string `desc:"Firmware Revision"`
	Hardware    string `desc:"Hardware Revision"`
	Manufacture string `desc:"Manufacture Name"`
}

type Mode struct {
	Mode uint8 `desc:"Mode"`
}
