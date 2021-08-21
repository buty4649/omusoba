package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/spf13/cobra"
)

type EnvSensor struct {
	client      ble.Client
	Address     string
	Model       string
	Serial      string
	Firmware    string
	Hardware    string
	Manufacture string
}

var (
	ErrNotEnvSensor = errors.New("Not env sensor")
	cfgAdapter      uint8
	scanCmd         = &cobra.Command{
		Use: "scan",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScanCmd()
		},
	}
)

var (
	modelUUID       = ble.UUID16(0x2A24) // Model Number
	serialUUID      = ble.UUID16(0x2A25) // Serial Number
	firmwareUUID    = ble.UUID16(0x2A26) // Firmware Revision
	harwareUUID     = ble.UUID16(0x2A27) // Hardware Revision
	manufactureUUID = ble.UUID16(0x2A29) // Manufacture Name
)

func init() {
	scanCmd.Flags().Uint8VarP(&cfgAdapter, "adapter", "A", 0, "Adapter Number")
	rootCmd.AddCommand(scanCmd)
}

func runScanCmd() error {
	adapter, err := linux.NewDevice(ble.OptDeviceID(int(cfgAdapter)))
	if err != nil {
		return err
	}
	ble.SetDefaultDevice(adapter)
	fmt.Printf("adapter: hci%d (%s)\n", cfgAdapter, adapter.Address())

	fmt.Println("Started scanning for devices (Ctrl+C to stop)")
	ctx := context.Background()
	defer ctx.Done()

	go func() {
		var checked_address []ble.Addr
		checked := func(a ble.Addr) bool {
			for _, addr := range checked_address {
				if a.String() == addr.String() {
					return true
				}
			}
			return false
		}

		ble.Scan(ctx, false, func(a ble.Advertisement) {
			addr := a.Addr()
			if checked(addr) {
				return
			}

			s, err := NewEnvSensor(ctx, a)
			if err != nil {
				if err != ErrNotEnvSensor {
					//fmt.Printf("Error: %s\n", err.Error())
					return
				}
			} else {
				fmt.Printf("[%s] Model: %s Serial: %s\n", s.Address, s.Model, s.Serial)
				s.Disconnect()
			}

			checked_address = append(checked_address, addr)
		}, nil)
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

	return nil
}

func NewEnvSensor(ctx context.Context, a ble.Advertisement) (*EnvSensor, error) {
	name := a.LocalName()
	if name != "Env" && name != "Rbt" {
		return nil, ErrNotEnvSensor
	}

	addr := a.Addr()
	client, err := ble.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}

	services, err := client.DiscoverServices([]ble.UUID{ble.DeviceInfoUUID})
	if err != nil {
		client.CancelConnection()
		return nil, err
	}

	filter := []ble.UUID{
		modelUUID, serialUUID, firmwareUUID, harwareUUID, manufactureUUID,
	}
	char, err := client.DiscoverCharacteristics(filter, services[0])
	if err != nil {
		client.CancelConnection()
		return nil, err
	}

	result := &EnvSensor{client: client, Address: addr.String()}

	for _, c := range char {
		v, err := client.ReadCharacteristic(c)
		if err != nil {
			client.CancelConnection()
			return nil, err
		}

		switch string(c.UUID) {
		case string(modelUUID):
			result.Model = string(v)
		case string(serialUUID):
			result.Serial = string(v)
		case string(firmwareUUID):
			result.Firmware = string(v)
		case string(harwareUUID):
			result.Hardware = string(v)
		case string(manufactureUUID):
			result.Manufacture = string(v)
		}
	}

	return result, nil
}

func (s *EnvSensor) Disconnect() error {
	return s.client.CancelConnection()
}
